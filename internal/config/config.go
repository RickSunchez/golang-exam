package config

import (
	"github.com/spf13/viper"
)

type EnviromentT struct {
	DataFolder     string
	Alpha2File     string
	BillingFile    string
	EmailFile      string
	SmsFile        string
	VoiceFile      string
	AccendentsFile string
}

func Enviroment() (EnviromentT, error) {
	viper.SetConfigFile("config/.env")

	err := viper.ReadInConfig()
	if err != nil {
		return EnviromentT{}, err
	}

	return EnviromentT{
		DataFolder:     viper.GetString("DATA_FOLDER"),
		Alpha2File:     viper.GetString("ALPHA2_CODES_FILE"),
		BillingFile:    viper.GetString("BILLING_FILE"),
		EmailFile:      viper.GetString("EMAIL_FILE"),
		SmsFile:        viper.GetString("SMS_FILE"),
		VoiceFile:      viper.GetString("VOICE_FILE"),
		AccendentsFile: viper.GetString("ACCENDENTS_FILE"),
	}, nil

}
