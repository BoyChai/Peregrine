package dingding

import (
	"Peregrine/alerter"
	"Peregrine/log"
	"Peregrine/stru"
	"encoding/json"
	"fmt"

	"github.com/BoyChai/CoralBot/action"
)

type dingding struct {
	handler *action.DingDingHandle
}

var dingdingAlert chan stru.AlarmContext
var ding dingding

func Init(way stru.Way) {
	ding.handler = action.NewDingDingHandleByWebHook(way.DingdingWebhook)

	dingdingAlert = make(chan stru.AlarmContext)
	alerter.Alerters[way.Name] = dingdingAlert
	go ding.work()
}

type CardMsg struct {
	Msgtype    string `json:"msgtype"`
	ActionCard struct {
		Title       string `json:"title"`
		Text        string `json:"text"`
		SingleTitle string `json:"singleTitle"`
		SingleURL   string `json:"singleURL"`
	} `json:"actionCard"`
}

func (d *dingding) work() {
	for {
		select {
		case alert := <-dingdingAlert:

			msg, _ := json.Marshal(CardMsg{
				Msgtype: "actionCard",
				ActionCard: struct {
					Title       string `json:"title"`
					Text        string `json:"text"`
					SingleTitle string `json:"singleTitle"`
					SingleURL   string `json:"singleURL"`
				}{
					Title: "监控告警",
					Text:  fmt.Sprintf("Level: %s  \nDescription: %s  \nExpr: %s  \n", alert.Entry.Level, alert.Entry.Description, alert.Entry.Expr),
				},
			})
			code, body, e := d.handler.SendGroupMessageByWebhook(string(msg))
			if e != nil {
				log.Error(e.Error())
			}
			if code != 200 {
				log.Error(alert.Way+" 告警器告警失败,状态码错误", body)
			}
			log.Debug(alert.Way+" 告警器告警成功", body)
		}
	}
}
