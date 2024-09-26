package dingding

import (
	"Peregrine/alerter"
	"Peregrine/stru"
	"encoding/json"
	"fmt"

	"github.com/BoyChai/CoralBot/action"
)

type dingding struct {
	handler *action.DingDingHandle
}

var dingdingAlert chan alerter.Alert
var ding dingding

func Init(alter stru.Alerter) {
	if alter.DingdingWebhook == "" {
		return
	}
	ding.handler = action.NewDingDingHandleByWebHook(alter.DingdingWebhook)

	dingdingAlert = make(chan alerter.Alert)
	alerter.Alerters["dingding"] = dingdingAlert
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
			d, d2, e := d.handler.SendGroupMessageByWebhook(string(msg))
			if e != nil {
				panic(e)
			}
			fmt.Println(d, d2)
		}
	}

}
