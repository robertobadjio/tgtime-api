package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config appConfig

type appConfig struct {
	Env string
	HostName string
	HostPort int
	UserName string
	Password string
	DataBaseName string
	SslMode string
	BotToken string
	RouterAddress string
	RouterUserName string
	RouterPassword string
	WebHookPath string
	TelegramBot string
}

func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("officetime")
	v.AutomaticEnv()
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}
	return v.Unmarshal(&Config)
}