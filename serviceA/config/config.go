package config

import "github.com/spf13/viper"

type AppSettings struct {
	Port   string
	UrlCep string
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

	viper.SetDefault("URL_CEP", "")
	viper.SetDefault("PORT", "")

	appConfig := &AppSettings{
		Port:   viper.GetString("PORT"),
		UrlCep: viper.GetString("URL_CEP"),
	}

	return appConfig, nil
}
