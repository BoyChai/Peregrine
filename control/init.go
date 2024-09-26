package control

import (
	"Peregrine/alerter/dingding"
	"Peregrine/alerter/mail"
	"Peregrine/alerter/webhook"
	"Peregrine/config"
	"log"
)

func InitPeregrine() {
	cfg := config.ReadConfig()
	if len(cfg.Alerter.Way) > 0 {
		for _, way := range cfg.Alerter.Way {
			switch way.Type {
			case "smtp":
				mail.Init(way)
			case "dingding":
				dingding.Init(way)
			case "webhook":
				webhook.Init(way)
			default:
				log.Fatal("Unknown alarm type：", way.Type)
			}
		}
	}
}
