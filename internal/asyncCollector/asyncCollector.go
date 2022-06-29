package asyncCollector

import (
	"fmt"
	"last_lesson/internal/alpha2"
	"last_lesson/internal/getters"
	"last_lesson/internal/mytypes"
	"sort"
	"strings"
	"sync"

	"github.com/jinzhu/copier"
)

func GetResult() (*mytypes.ResultSetT, error) {
	result, err := awaitResult()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}

func awaitResult() (mytypes.ResultSetT, error) {
	var wg sync.WaitGroup
	var result mytypes.ResultSetT
	var resultError error = nil

	wg.Add(7)
	go func() {
		defer wg.Done()

		sms := <-smsStream(&wg)
		if sms.Error != nil && resultError == nil {
			resultError = sms.Error
		} else {
			result.SMS = (*sms.Data)
		}
	}()

	go func() {
		defer func() {
			wg.Done()
		}()

		mms := <-mmsStream(&wg)

		if mms.Error != nil && resultError == nil {
			resultError = mms.Error
		} else {
			result.MMS = (*mms.Data)
		}
	}()

	go func() {
		defer wg.Done()

		vc := <-vcStream(&wg)
		if vc.Error != nil && resultError == nil {
			resultError = vc.Error
		} else {
			result.VoiceCall = (*vc.Data)
		}
	}()

	go func() {
		defer wg.Done()

		email := <-emailStream(&wg)
		if email.Error != nil && resultError == nil {
			resultError = email.Error
		} else {
			result.Email = (*email.Data)
		}
	}()

	go func() {
		defer wg.Done()

		billing := <-billingStream(&wg)
		if billing.Error != nil && resultError == nil {
			resultError = billing.Error
		} else {
			result.Billing = (*billing.Data)
		}
	}()

	go func() {
		defer wg.Done()

		support := <-supportStream(&wg)
		if support.Error != nil && resultError == nil {
			resultError = support.Error
		} else {
			result.Support = (*support.Data)
		}
	}()

	go func() {
		defer wg.Done()

		incidents := <-incidentsStream(&wg)
		if incidents.Error != nil && resultError == nil {
			resultError = incidents.Error
		} else {
			result.Incidents = (*incidents.Data)
		}
	}()

	wg.Wait()

	return result, resultError
}

func smsStream(wg *sync.WaitGroup) chan mytypes.SMSDataFlow {
	out := make(chan mytypes.SMSDataFlow)
	data := mytypes.SMSDataFlow{}

	data.Data = nil

	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer func() {
			close(out)
		}()

		smsData, err := getters.SMS()
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		err = alpha2.Alpha2ToCountrySMS(&smsData)
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		var smsDataCopied []mytypes.SMSData
		err = copier.Copy(&smsDataCopied, &smsData)
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		sort.SliceStable(smsData, func(i, j int) bool {
			return strings.Compare(smsData[i].Provider, smsData[j].Provider) == -1
		})
		sort.SliceStable(smsDataCopied, func(i, j int) bool {
			return strings.Compare(smsDataCopied[i].Country, smsDataCopied[j].Country) == -1
		})

		var result [][]mytypes.SMSData

		result = append(result, smsData)
		result = append(result, smsDataCopied)

		data.Data = &result

		out <- data
	}()

	return out
}

func mmsStream(wg *sync.WaitGroup) chan mytypes.MMSDataFlow {
	out := make(chan mytypes.MMSDataFlow)
	data := mytypes.MMSDataFlow{}

	data.Data = nil

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(out)

		// MMS
		mmsData, err := getters.MMS()
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		err = alpha2.Alpha2ToCountryMMS(&mmsData)
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		var mmsDataCopied []mytypes.MMSData
		err = copier.Copy(&mmsDataCopied, &mmsData)
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		sort.SliceStable(mmsData, func(i, j int) bool {
			return strings.Compare(mmsData[i].Provider, mmsData[j].Provider) == -1
		})
		sort.SliceStable(mmsDataCopied, func(i, j int) bool {
			return strings.Compare(mmsDataCopied[i].Country, mmsDataCopied[j].Country) == -1
		})

		var result [][]mytypes.MMSData

		result = append(result, mmsData)
		result = append(result, mmsDataCopied)

		data.Data = &result

		out <- data
	}()

	return out
}

func vcStream(wg *sync.WaitGroup) chan mytypes.VCDataFlow {
	out := make(chan mytypes.VCDataFlow)
	data := mytypes.VCDataFlow{}

	data.Data = nil

	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer func() {
			close(out)
		}()

		// VoiceCall
		voiceCallData, err := getters.VoiceCall()
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		var result [][]mytypes.VoiceCallData
		result = append(result, voiceCallData)

		data.Data = &result

		out <- data
	}()

	return out
}

func emailStream(wg *sync.WaitGroup) chan mytypes.EmailDataFlow {
	out := make(chan mytypes.EmailDataFlow)
	data := mytypes.EmailDataFlow{}

	data.Data = nil

	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer func() {
			close(out)
		}()

		// Email
		emailData, err := getters.Email()
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		sort.SliceStable(emailData, func(i, j int) bool {
			switch strings.Compare(emailData[i].Country, emailData[j].Country) {
			case -1:
				return true
			case 1:
				return false
			}

			return emailData[i].DeliveryTime < emailData[j].DeliveryTime
		})

		previousCode := emailData[0].Country
		beginIndex := 0
		result := make(map[string][][]mytypes.EmailData)

		for i, row := range emailData {
			if row.Country != previousCode {
				var min3 []mytypes.EmailData
				var max3 []mytypes.EmailData

				for n := beginIndex; n < beginIndex+3; n++ {
					min3 = append(min3, emailData[n])
				}

				for m := i - 3; m < i; m++ {
					max3 = append(max3, emailData[m])
				}

				result[previousCode] = append(result[previousCode], min3)
				result[previousCode] = append(result[previousCode], max3)

				previousCode = row.Country
				beginIndex = i
			}
		}

		data.Data = &result

		out <- data
	}()

	return out
}

func billingStream(wg *sync.WaitGroup) chan mytypes.BillingDataFlow {
	out := make(chan mytypes.BillingDataFlow)
	data := mytypes.BillingDataFlow{}

	data.Data = nil

	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer func() {
			close(out)
		}()
		// Billing
		billingData, err := getters.Billing()
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		data.Data = &billingData

		out <- data
	}()

	return out
}

func supportStream(wg *sync.WaitGroup) chan mytypes.SupportDataFlow {
	out := make(chan mytypes.SupportDataFlow)
	data := mytypes.SupportDataFlow{}

	data.Data = nil

	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer func() {
			close(out)
		}()
		// Support
		// Что значит "Профильные задачи?? Какие из них профильные??"

		supportData, err := getters.Support()
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		var sumOfTasks int
		var result []int
		for _, row := range supportData {
			sumOfTasks += row.ActiveTickets
		}

		if sumOfTasks < 9 {
			result = append(result, 1)
		} else if (sumOfTasks >= 9) && (sumOfTasks <= 16) {
			result = append(result, 2)
		} else {
			result = append(result, 3)
		}

		mpt := float32(60) / float32(18) // minutes per task
		result = append(result, int(mpt*float32(sumOfTasks)))

		data.Data = &result

		out <- data
	}()

	return out
}

func incidentsStream(wg *sync.WaitGroup) chan mytypes.IncidentDataFlow {
	out := make(chan mytypes.IncidentDataFlow)
	data := mytypes.IncidentDataFlow{}

	data.Data = nil

	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer func() {
			close(out)
		}()
		incidentData, err := getters.Incidents()
		if err != nil {
			data.Error = err
			out <- data
			return
		}

		sort.SliceStable(incidentData, func(i, j int) bool {
			return strings.Compare(incidentData[i].Status, incidentData[j].Status) == -1
		})

		data.Data = &incidentData

		out <- data
	}()

	return out
}

// func <sample>(wg *sync.WaitGroup) chan <type> {
// 	out := make(chan <type>)
// 	data := <type>{}

// 	data.Data = nil

// 	wg.Add(1)

// 	go func(){
// 		defer func() {
// 			wg.Done()
// 		}()
// 		defer func() {
// 			close(out)
// 		}()
// 		for {
// if err != nil {
// 	data.Error = err
// 	out <- data
// 	continue
// }
// out <- data
// 		}
// 	}()

// 	return out
// }
