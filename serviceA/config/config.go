package config

import "github.com/spf13/viper"

type AppSettings struct {
	Port      string
	UrlCep    string
	UrlTempo  string
	UrlZipkin string
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

	viper.SetDefault("PORT", "")
	viper.SetDefault("URL_CEP", "")
	viper.SetDefault("URL_TEMPO", "")
	viper.SetDefault("URL_ZIPKIN", "")

	appConfig := &AppSettings{
		Port:     viper.GetString("PORT"),
		UrlCep:   viper.GetString("URL_CEP"),
		UrlTempo: viper.GetString("URL_TEMPO"),
		UrlZipkin: viper.GetString("URL_ZIPKIN"),
	}

	return appConfig, nil
}
