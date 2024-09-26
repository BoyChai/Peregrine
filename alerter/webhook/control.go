package webhook

import (
	"Peregrine/alerter"
	"Peregrine/stru"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type webhook struct {
	URL string
}

var hook webhook

var webhookAlert chan alerter.Alert

func Init(way stru.Way) {
	fmt.Println(way.WebhookURL)
	hook.URL = way.WebhookURL
	webhookAlert = make(chan alerter.Alert)
	alerter.Alerters[way.Name] = webhookAlert
	go hook.wrok()
}

func (w *webhook) wrok() {
	for {
		select {
		case alert := <-webhookAlert:
			msg, _ := json.Marshal(alert.Entry)
			if len(alert.Target) <= 0 {
				return
			}
			for _, v := range alert.Target {
				d, d2, e := w.send(msg, v)
				if e != nil {
					log.Println(e)
				}
				log.Println(d, d2)
			}
		}
	}
}
func (w *webhook) send(msg []byte, target stru.Target) (statusCode int, bodyData string, err error) {
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
