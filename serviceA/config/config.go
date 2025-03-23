package config

import "github.com/spf13/viper"

type AppSettings struct {
	Port        string
	UrlConsulta string
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

	viper.SetDefault("PORT", "")
	viper.SetDefault("URL_CONSULTA", "")
	viper.SetDefault("URL_ZIPKIN", "")
	viper.SetDefault("SERVICE_NAME", "")

	appConfig := &AppSettings{
		Port:        viper.GetString("PORT"),
		UrlConsulta: viper.GetString("URL_CONSULTA"),
		UrlZipkin:   viper.GetString("URL_ZIPKIN"),
		ServiceName: viper.GetString("SERVICE_NAME"),
	}

	return appConfig, nil
}
