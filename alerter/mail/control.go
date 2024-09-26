package mail

import (
	"Peregrine/alerter"
	"Peregrine/stru"
	"crypto/tls"
	"errors"
	"fmt"

	"gopkg.in/gomail.v2"
)

type smtp struct {
	Form string
}

var mail smtp

var dialer *gomail.Dialer

var mailAlert chan alerter.Alert

func Init(way stru.Way) {
	mail.Form = way.SMTPForm
	dialer = gomail.NewDialer(way.SMTPHost, way.SMTPPort, way.SMTPUsername, way.SMTPPasswd)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: way.SMTPTLS,
	}
	mailAlert = make(chan alerter.Alert)
	alerter.Alerters[way.Name] = mailAlert
	go mail.work()
}

func (s *smtp) work() error {
	for {
		alert := <-mailAlert
		err := s.send(alert)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (s *smtp) send(alert alerter.Alert) error {

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.Form)
	var to []string
	if len(alert.Target) <= 0 {
		return errors.New("没有目标邮箱")
	}
	for _, v := range alert.Target {
		to = append(to, v.To...)
	}
	msg.SetHeader("To", to...)

	msg.SetHeader("Subject", "告警")
	msg.SetBody("text/plain", fmt.Sprintf("Level: %s\nDescription: %s\nExpr: %s", alert.Entry.Level, alert.Entry.Description, alert.Entry.Expr))
	fmt.Println(msg)
	if err := dialer.DialAndSend(msg); err != nil {
		panic(err)
	}
	fmt.Println("发送成功")
	return nil
}
