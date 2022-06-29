package getters

import (
	"encoding/json"
	"io"
	"last_lesson/internal/alpha2"
	"last_lesson/internal/mytypes"
	"last_lesson/internal/sub"
	"last_lesson/internal/vars"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
)

func MMS() ([]mytypes.MMSData, error) {
	resp, err := http.Get("http://127.0.0.1:8383/mms")
	if err != nil {
		return []mytypes.MMSData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []mytypes.MMSData{}, http.ErrNoLocation
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []mytypes.MMSData{}, err
	}

	var mmsData []mytypes.MMSData

	if json.Unmarshal(body, &mmsData) != nil {
		return []mytypes.MMSData{}, err
	}

	var filteredData []mytypes.MMSData
	alpha2codes, err := alpha2.GetAlpha2Codes(false)
	if err != nil {
		return []mytypes.MMSData{}, err
	}

	for _, row := range mmsData {
		bandWidth, bwErr := strconv.Atoi(row.Bandwidth)
		_, rtErr := strconv.Atoi(row.ResponseTime)

		alpha2Err := !sub.StringInSlice(row.Country, alpha2codes)
		providerErr := !sub.StringInSlice(row.Provider, vars.SmsMmsProviders)
		bandWidthErr := (bandWidth < 0) || (bandWidth > 100)

		if (bwErr != nil) || (rtErr != nil) || alpha2Err || providerErr || bandWidthErr {
			continue
		}

		filteredData = append(filteredData, row)
	}

	return filteredData, nil
}

func SMS() ([]mytypes.SMSData, error) {
	smsData, err := sub.ReadCSVFromFile(vars.SmsDataFile, ";", 4)
	if err != nil {
		return []mytypes.SMSData{}, err
	}

	alpha2codes, err := alpha2.GetAlpha2Codes(false) // false - read from file; true - sync data from wiki
	if err != nil {
		return []mytypes.SMSData{}, err
	}

	result := []mytypes.SMSData{}
	for _, row := range smsData {
		intData, err := sub.CheckInts([]int{1, 2}, row)
		if err != nil {
			continue
		}

		alpha2Err := !sub.StringInSlice(row[0], alpha2codes)
		providerErr := !sub.StringInSlice(row[3], vars.SmsMmsProviders)
		bandWidthErr := (intData[0] < 0) || (intData[0] > 100)

		if (err != nil) || alpha2Err || providerErr || bandWidthErr {
			continue
		}

		result = append(result, mytypes.SMSData{
			Country:      row[0],
			Bandwidth:    row[1],
			ResponseTime: row[2],
			Provider:     row[3],
		})
	}

	return result, nil
}

func VoiceCall() ([]mytypes.VoiceCallData, error) {
	voiceData, err := sub.ReadCSVFromFile(vars.VoiceDataFile, ";", 8)
	if err != nil {
		return []mytypes.VoiceCallData{}, err
	}

	alpha2codes, err := alpha2.GetAlpha2Codes(false) // false - read from file; true - sync data from wiki
	if err != nil {
		return []mytypes.VoiceCallData{}, err
	}

	var result []mytypes.VoiceCallData
	for _, row := range voiceData {
		intData, err := sub.CheckInts([]int{1, 2, 5, 6, 7}, row)
		if err != nil {
			continue
		}

		alpha2Err := !sub.StringInSlice(row[0], alpha2codes)
		providerErr := !sub.StringInSlice(row[3], vars.VoiceCallProviders)
		connectionStability, csErr := strconv.ParseFloat(row[4], 32)
		bandWidthErr := (intData[0] < 0) || (intData[0] > 100)

		if (csErr != nil) || alpha2Err || providerErr || bandWidthErr {
			continue
		}

		result = append(result, mytypes.VoiceCallData{
			Country:             row[0],
			Bandwidth:           row[1],
			ResponseTime:        row[2],
			Provider:            row[3],
			ConnectionStability: float32(connectionStability),
			TTFB:                intData[2],
			VoicePurity:         intData[3],
			MedianOfCallsTime:   intData[4],
		})
	}

	return result, nil
}

func Email() ([]mytypes.EmailData, error) {
	emailData, err := sub.ReadCSVFromFile(vars.EmailDataFile, ";", 3)
	if err != nil {
		return []mytypes.EmailData{}, err
	}

	alpha2codes, err := alpha2.GetAlpha2Codes(false) // false - read from file; true - sync data from wiki
	if err != nil {
		return []mytypes.EmailData{}, err
	}

	var result []mytypes.EmailData
	for _, row := range emailData {
		alpha2Err := !sub.StringInSlice(row[0], alpha2codes)
		providerErr := !sub.StringInSlice(row[1], vars.EmailProviders)
		deliveryTime, err := strconv.Atoi(row[2])

		if (err != nil) || alpha2Err || providerErr {
			continue
		}

		result = append(result, mytypes.EmailData{
			Country:      row[0],
			Provider:     row[1],
			DeliveryTime: deliveryTime,
		})
	}

	return result, err
}

func Billing() (mytypes.BillingData, error) {
	billingData, err := sub.ReadFromFile(vars.BillingDataFile)
	if err != nil {
		return mytypes.BillingData{}, err
	}

	var bitMask uint8

	for _, c := range billingData {
		bitMask <<= 1

		if strings.Compare(string(c), "1") == 0 {
			bitMask |= 1
		}
	}

	flags := []bool{false, false, false, false, false, false}
	for i := range flags {
		flag := bitMask & 1
		flags[i] = flag == 1

		bitMask >>= 1
	}

	result := mytypes.BillingData{
		CreateCustomer: flags[0],
		Purchase:       flags[1],
		Payout:         flags[2],
		Recurring:      flags[3],
		FraudControl:   flags[4],
		CheckoutPage:   flags[5],
	}

	return result, err
}

func Support() ([]mytypes.SupportData, error) {
	resp, err := http.Get("http://127.0.0.1:8383/support")
	if err != nil {
		return []mytypes.SupportData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []mytypes.SupportData{}, http.ErrNoLocation
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []mytypes.SupportData{}, err
	}

	var supportData []mytypes.SupportData

	if json.Unmarshal(body, &supportData) != nil {
		return []mytypes.SupportData{}, err
	}

	return supportData, err
}

func Incidents() ([]mytypes.IncidentData, error) {
	resp, err := http.Get("http://127.0.0.1:8383/accendent")
	if err != nil {
		return []mytypes.IncidentData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []mytypes.IncidentData{}, http.ErrNoLocation
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []mytypes.IncidentData{}, err
	}

	var incidentData []mytypes.IncidentData
	if json.Unmarshal(body, &incidentData) != nil {
		return []mytypes.IncidentData{}, err
	}

	var filteredData []mytypes.IncidentData
	for _, incident := range incidentData {
		notActive := strings.Compare(incident.Status, "active") != 0
		notClosed := strings.Compare(incident.Status, "closed") != 0
		if notActive && notClosed {
			continue
		}

		filteredData = append(filteredData, incident)
	}

	return filteredData, err
}

func GetResultData() (mytypes.ResultSetT, error) {
	var result mytypes.ResultSetT

	// SMS
	smsData, err := SMS()
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	err = alpha2.Alpha2ToCountrySMS(&smsData)
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	var smsDataCopied []mytypes.SMSData
	err = copier.Copy(&smsDataCopied, &smsData)
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	sort.SliceStable(smsData, func(i, j int) bool {
		return strings.Compare(smsData[i].Provider, smsData[j].Provider) == -1
	})
	sort.SliceStable(smsDataCopied, func(i, j int) bool {
		return strings.Compare(smsDataCopied[i].Country, smsDataCopied[j].Country) == -1
	})

	result.SMS = append(result.SMS, smsData)
	result.SMS = append(result.SMS, smsDataCopied)

	// MMS
	mmsData, err := MMS()
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	err = alpha2.Alpha2ToCountryMMS(&mmsData)
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	var mmsDataCopied []mytypes.MMSData
	err = copier.Copy(&mmsDataCopied, &mmsData)
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	sort.SliceStable(mmsData, func(i, j int) bool {
		return strings.Compare(mmsData[i].Provider, mmsData[j].Provider) == -1
	})
	sort.SliceStable(mmsDataCopied, func(i, j int) bool {
		return strings.Compare(mmsDataCopied[i].Country, mmsDataCopied[j].Country) == -1
	})

	result.MMS = append(result.MMS, mmsData)
	result.MMS = append(result.MMS, mmsDataCopied)

	// VoiceCall
	voiceCallData, err := VoiceCall()
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	result.VoiceCall = append(result.VoiceCall, voiceCallData)

	// Email
	emailData, err := Email()
	if err != nil {
		return mytypes.ResultSetT{}, err
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
	result.Email = make(map[string][][]mytypes.EmailData)

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

			result.Email[previousCode] = append(result.Email[previousCode], min3)
			result.Email[previousCode] = append(result.Email[previousCode], max3)

			previousCode = row.Country
			beginIndex = i
		}
	}

	// Billing
	billingData, err := Billing()
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	result.Billing = billingData

	// Support
	// Что значит "Профильные задачи?? Какие из них профильные??"

	supportData, err := Support()
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	var sumOfTasks int
	for _, row := range supportData {
		sumOfTasks += row.ActiveTickets
	}

	if sumOfTasks < 9 {
		result.Support = append(result.Support, 1)
	} else if (sumOfTasks >= 9) && (sumOfTasks <= 16) {
		result.Support = append(result.Support, 2)
	} else {
		result.Support = append(result.Support, 3)
	}

	mpt := float32(60) / float32(18) // minutes per task
	result.Support = append(result.Support, int(mpt*float32(sumOfTasks)))

	// Incidents

	incidentData, err := Incidents()
	if err != nil {
		return mytypes.ResultSetT{}, err
	}

	sort.SliceStable(incidentData, func(i, j int) bool {
		return strings.Compare(incidentData[i].Status, incidentData[j].Status) == -1
	})

	result.Incidents = incidentData

	return result, err
}
