package mytypes

import (
	"sync"
	"time"
)

// Data types

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

type SMSDataFlow struct {
	Error error
	Data  *[][]SMSData
}

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type MMSDataFlow struct {
	Error error
	Data  *[][]MMSData
}

type VoiceCallData struct {
	Country             string  `json:"country"`
	Bandwidth           string  `json:"bandwidth"`
	ResponseTime        string  `json:"response_time"`
	Provider            string  `json:"provider"`
	ConnectionStability float32 `json:"connection_stability"`
	TTFB                int     `json:"ttfb"`
	VoicePurity         int     `json:"voice_purity"`
	MedianOfCallsTime   int     `json:"median_of_call_time"`
}

type VCDataFlow struct {
	Error error
	Data  *[]VoiceCallData
}

type EmailData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
}

type EmailDataFlow struct {
	Error error
	Data  *map[string][][]EmailData
}

type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"recurring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}

type BillingDataFlow struct {
	Error error
	Data  *BillingData
}

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

type SupportDataFlow struct {
	Error error
	Data  *[]int
}

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы: active и closed
}

type IncidentDataFlow struct {
	Error error
	Data  *[]IncidentData
}

// Server types

type ResultT struct {
	Status bool       `json:"status"`
	Data   ResultSetT `json:"data"`
	Error  string     `json:"error"`
}

type ResultSetT struct {
	SMS       [][]SMSData              `json:"sms"`
	MMS       [][]MMSData              `json:"mms"`
	VoiceCall []VoiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []IncidentData           `json:"incident"`
}

type ResultSetTFlow struct {
	Error error
	Data  *ResultSetT
}

// Service types

type Alpha2Row struct {
	Country string `json:"country"`
	Alpha3  string `json:"alpha3"`
	ISO1    string `json:"ISO1"`
}
type Alpha2Codes map[string]Alpha2Row

type Synchronizer struct {
	LastResponseTime time.Time
	sync.Mutex
}

func (s *Synchronizer) NewLoop(newTime time.Time) {
	s.Lock()
	s.LastResponseTime = newTime
	s.Unlock()
}
