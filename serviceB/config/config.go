package config

import "github.com/spf13/viper"

type AppSettings struct {
	Port        string
	UrlCep      string
	UrlTempo    string
	TempoApiKey string
	UrlZipkin   string
	ServiceName string
}

func ProvideConfig() *AppSettings {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func LoadConfig() (*AppSettings, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.SetDefault("URL_TEMPO", "")
	viper.SetDefault("PORT", "")
	viper.SetDefault("API_KEY_TEMPO", "")
	viper.SetDefault("URL_ZIPKIN", "")
	viper.SetDefault("SERVICE_NAME", "")

	appConfig := &AppSettings{
		Port:        viper.GetString("PORT"),
		UrlTempo:    viper.GetString("URL_TEMPO"),
		TempoApiKey: viper.GetString("API_KEY_TEMPO"),
		UrlZipkin:   viper.GetString("URL_ZIPKIN"),
		ServiceName: viper.GetString("SERVICE_NAME"),
	}

	return appConfig, nil
}
