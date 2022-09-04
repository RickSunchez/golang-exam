package config

import (
	"github.com/spf13/viper"
)

type EnviromentT struct {
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

	dataFolder := viper.GetString("DATA_FOLDER")

	return EnviromentT{
		Alpha2File:     dataFolder + viper.GetString("ALPHA2_CODES_FILE"),
		BillingFile:    dataFolder + viper.GetString("BILLING_FILE"),
		EmailFile:      dataFolder + viper.GetString("EMAIL_FILE"),
		SmsFile:        dataFolder + viper.GetString("SMS_FILE"),
		VoiceFile:      dataFolder + viper.GetString("VOICE_FILE"),
		AccendentsFile: dataFolder + viper.GetString("ACCENDENTS_FILE"),
	}, nil

}
