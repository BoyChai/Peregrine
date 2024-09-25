package mail

import (
	"Peregrine/alerter"
	"Peregrine/stru"
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

type smtp struct {
	Form string
}

var mail smtp

var dialer *gomail.Dialer

var mailAlert chan alerter.Alert

func Init(alter stru.Alerter) {
	mail.Form = alter.SMTPForm
	dialer = gomail.NewDialer(alter.SMTPHost, alter.SMTPPort, alter.SMTPUsername, alter.SMTPPasswd)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: alter.SMTPTLS,
	}
	mailAlert = make(chan alerter.Alert)
	alerter.Alerters["mail"] = mailAlert
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
	msg.SetHeader("To", alert.Target...)

	msg.SetHeader("Subject", "告警")
	msg.SetBody("text/plain", fmt.Sprintf("Level: %s\nDescription: %s\nExpr: %s", alert.Entry.Level, alert.Entry.Description, alert.Entry.Expr))
	fmt.Println(msg)
	if err := dialer.DialAndSend(msg); err != nil {
		panic(err)
	} // 发送成功
	fmt.Println("发送成功")
	return nil
}