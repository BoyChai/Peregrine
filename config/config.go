package config

import (
	"log"

	"Peregrine/stru"

	"github.com/spf13/viper"
)

// 初始化配置
func InitPeregrine() {
	viper.SetConfigName("peregrine")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	var config stru.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
