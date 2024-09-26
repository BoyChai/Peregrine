package control

import (
	"Peregrine/alerter/dingding"
	"Peregrine/alerter/mail"
	"Peregrine/alerter/webhook"
	"Peregrine/asset"
	"Peregrine/config"
	"Peregrine/monitor"
	"Peregrine/stru"
	"fmt"
	"log"
)

var TriggerChan chan stru.AlarmTrigger

// 资产

func InitPeregrine() {
	cfg := config.ReadConfig()
	// 注册资产
	asset.Init(cfg.Asset)

	// 注册告警器
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
	// 注册监视器
	TriggerChan = make(chan stru.AlarmTrigger)
	monitor.Run(TriggerChan, cfg.Rule)

	// 拿数据
	for {
		trigger := <-TriggerChan
		// 发送告警
		fmt.Println(trigger)
	}
}
