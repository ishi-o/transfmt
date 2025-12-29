// Package config
package config

import "github.com/spf13/viper"

type FanyiConfig struct {
	SourceLang string `mapstructure:"default_source_lang"`
	TargetLang string `mapstructure:"default_target_lang"`
	Backend    string `mapstructure:"default_backend"`
}

type UserSetting struct {
	Fanyi FanyiConfig `mapstructure:"fanyi"`
}

type Config struct {
	App struct {
		Name      string `mapstructure:"name"`
		Port      int    `mapstructure:"port"`
		SecretKey string `mapstructure:"secret_key"`
	}
	UserSetting UserSetting `mapstructure:"user_setting"`
	Redis       struct {
		Host       string `mapstructure:"host"`
		Port       int    `mapstructure:"port"`
		Password   string `mapstructure:"password"`
		DB         int    `mapstructure:"db"`
		PoolSize   int    `mapstructure:"pool_size"`
		SessionTTL int    `mapstructure:"session_ttl"`
	}
}

var C Config

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
