package control

import (
	"Peregrine/alerter"
	"Peregrine/alerter/dingding"
	"Peregrine/alerter/mail"
	"Peregrine/alerter/template"
	"Peregrine/alerter/webhook"
	"Peregrine/asset"
	"Peregrine/config"
	"Peregrine/log"
	"Peregrine/monitor"
	"Peregrine/stru"
)

var TriggerChan chan stru.AlarmTrigger

// 资产

func InitPeregrine() {
	cfg := config.ReadConfig()
	// 配置日志
	log.InitLogOut(cfg.Log)
	// 注册资产
	log.Debug("正在注册资产")
	asset.Init(cfg.Asset)
	// 读取告警模板
	log.Debug("正在读取告警模板")
	template.ReadAlerterTemplate()
	// 注册告警器
	if len(cfg.Alerter.Way) > 0 {
		for _, way := range cfg.Alerter.Way {
			log.Debug("注册" + way.Name + "告警器")
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
	log.Debug("正在注册监视器")
	// 注册监视器
	TriggerChan = make(chan stru.AlarmTrigger)
	monitor.Run(TriggerChan, cfg.Rule)

	// 接收告警
	log.Debug("开始接收告警")
	for {
		trigger := <-TriggerChan
		// 发送告警
		log.Info("接收到资产"+trigger.AlerterTarget+"的告警,正在通知告警器", trigger.Entry.Expr)
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
