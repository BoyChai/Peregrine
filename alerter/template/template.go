package template

import (
	"Peregrine/stru"
	"bytes"
	"html/template"
	"log"
)

var mailTemplateText = `告警资产: {{.Asset}}
告警方式: {{.Way}}
告警目标: {{.Target.Name}} ({{range .Target.To}}{{.}} {{end}})
告警规则: {{.Entry.Description}} (表达式: {{.Entry.Expr}}, 等级: {{.Entry.Level}})
触发主机: {{range .Instance}}{{.}} {{end}}
当前值: {{range .Value}}{{.}} {{end}}
`
var mailTempalte *template.Template

func init() {
	// 解析模板
	var err error
	mailTempalte, err = template.New("mail").Parse(mailTemplateText)
	if err != nil {
		log.Fatalln("Error parsing template:", err)
		return
	}
}

func GetMailText(context stru.AlarmContext) string {
	var output bytes.Buffer

	err := mailTempalte.Execute(&output, context)
	if err != nil {
		log.Println("渲染失败:", err)
		return ""
	}
	return output.String()
}
