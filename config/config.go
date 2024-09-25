package config

import (
	"log"
	"os"

	"Peregrine/stru"

	"gopkg.in/yaml.v3"
)

// 读取配置
func ReadConfig() stru.Config {

	data, err := os.ReadFile("peregrine.yaml")
	if err != nil {
		log.Fatalln(err)

	}
	var cfg stru.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return cfg
}
