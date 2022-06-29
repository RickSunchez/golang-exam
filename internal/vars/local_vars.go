package vars

import (
	"path/filepath"
	"runtime"
)

func WorkDirectory() string {
	_, path, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(path), "../..")
}

var SmsDataFile = filepath.Join(WorkDirectory(), "src/simulator/sms.data")
var VoiceDataFile = filepath.Join(WorkDirectory(), "src/simulator/voice.data")
var EmailDataFile = filepath.Join(WorkDirectory(), "src/simulator/email.data")
var BillingDataFile = filepath.Join(WorkDirectory(), "src/simulator/billing.data")

var Alpha2CodesFile = filepath.Join(WorkDirectory(), "src/data/alpha2.json")

var SmsMmsProviders = []string{"Topolo", "Rond", "Kildy"}
var VoiceCallProviders = []string{"TransparentCalls", "E-Voice", "JustPhone"}
var EmailProviders = []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast", "AOL", "Live", "RediffMail", "GMX", "Protonmail",
	"Yandex", "Mail.ru"}
