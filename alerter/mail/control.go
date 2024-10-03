package mail

import (
	"Peregrine/alerter"
	"Peregrine/alerter/template"
	"Peregrine/log"
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

var mailAlert chan stru.AlarmContext

func Init(way stru.Way) {
	mail.Form = way.SMTPForm
	dialer = gomail.NewDialer(way.SMTPHost, way.SMTPPort, way.SMTPUsername, way.SMTPPasswd)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: way.SMTPTLS,
	}
	mailAlert = make(chan stru.AlarmContext)
	alerter.Alerters[way.Name] = mailAlert
	go mail.work()
}

func (s *smtp) work() error {
	for {
		alert := <-mailAlert
		err := s.send(alert)
		if err != nil {
			log.Error(alert.Way+" 发送告警时出现错误", err.Error())
		}
	}
}
func (s *smtp) send(alert stru.AlarmContext) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.Form)
	var to []string
	if len(alert.Target.To) <= 0 {
		return errors.New("没有目标邮箱")
	}
	for _, v := range alert.Target.To {
		to = append(to, v)
	}
	msg.SetHeader("To", to...)

	msg.SetHeader("Subject", "告警")
	msg.SetBody("text/plain", template.GetMailText(alert))
	fmt.Println(msg)
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}
	log.Info(alert.Way + "告警器发送成")
	return nil
}
