package webhook

import (
	"Peregrine/alerter"
	"Peregrine/log"
	"Peregrine/stru"
	"bytes"
	"encoding/json"
	"net/http"
)

type webhook struct {
	URL string
}

var hook webhook

var webhookAlert chan stru.AlarmContext

func Init(way stru.Way) {
	hook.URL = way.WebhookURL
	webhookAlert = make(chan stru.AlarmContext)
	alerter.Alerters[way.Name] = webhookAlert
	go hook.wrok()
}

func (w *webhook) wrok() {
	for {
		select {
		case alert := <-webhookAlert:
			msg, _ := json.Marshal(alert.Entry)
			if len(alert.Target.To) <= 0 {
				return
			}
			code, body, e := w.send(msg)

			if e != nil {
				log.Error(alert.Way+"告警器发送错误", e.Error())
			}
			if code != 200 {
				log.Error(alert.Way+"告警器发送状态码错误", body)
			}
			log.Debug(alert.Way+" 告警器发送成功", body)
		}
	}
}
func (w *webhook) send(msg []byte) (statusCode int, bodyData string, err error) {
	resp, err := http.Post(w.URL, "application/json", bytes.NewBuffer(msg))
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return statusCode, "", err
	}
	bodyData = buf.String()

	return statusCode, bodyData, nil
}
