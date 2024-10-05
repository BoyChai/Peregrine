package template

import (
	"Peregrine/log"
	"Peregrine/stru"
	"bytes"
	"html/template"
)

var mailTempalte *template.Template
var dingdingTempalte *template.Template

func ReadAlerterTemplate() {
	// 解析模板
	var err error
	mailTempalte, err = template.ParseFiles("./template/mail.template")
	if err != nil {
		log.Fatal("读取mail模板时出现错误", err)
		return
	}
	dingdingTempalte, err = template.ParseFiles("./template/dingding.template")
	if err != nil {
		log.Fatal("读取dingding模板时出现错误", err)
		return
	}
}

func GetMailText(context stru.AlarmContext) string {
	var output bytes.Buffer

	err := mailTempalte.Execute(&output, context)
	if err != nil {
		log.Error("获取模板时出现错误")
		return ""
	}
	return output.String()
}

func GetDingDingText(context stru.AlarmContext) string {
	var output bytes.Buffer

	err := dingdingTempalte.Execute(&output, context)
	if err != nil {
		log.Error("获取模板时出现错误")
		return ""
	}
	return output.String()
}
