package control

import (
	"Peregrine/alerter"
	"Peregrine/alerter/dingding"
	"Peregrine/alerter/mail"
	"Peregrine/alerter/webhook"
	"Peregrine/asset"
	"Peregrine/config"
	"Peregrine/monitor"
	"Peregrine/stru"
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
		// fmt.Println("告警告警3")
		// fmt.Println("告警资产:", trigger.AlerterTarget)
		// fmt.Println("告警通知对象:", trigger.AlerterWay)
		// fmt.Println("告警的注解:", trigger.Entry.Description)
		// fmt.Println("告警的规则:", trigger.Entry.Expr)
		// fmt.Println("告警的等级:", trigger.Entry.Level)
		// 告警内容组装
		alarmContext := stru.AlarmContext{
			Asset:    trigger.AssetName,
			Way:      trigger.AlerterWay,
			Entry:    trigger.Entry,
			Instance: trigger.Instance,
			Value:    trigger.Value,
		}
		for _, v := range cfg.Alerter.Target {
			if v.Name == trigger.AlerterTarget {
				alarmContext.Target = stru.Target{
					Name: trigger.AlerterTarget,
					To:   v.To,
				}
			}
			break
		}
		alerter.Alerters[trigger.AlerterWay] <- alarmContext
	}
}
