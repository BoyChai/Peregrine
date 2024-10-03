package stru

import "time"

type Log struct {
	File  bool   `yaml:"file"`
	Path  string `yaml:"path"`
	Json  bool   `yaml:"json"`
	Level int    `yaml:"level"`
}

type Asset struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}
type Alerter struct {
	Way    []Way    `yaml:"way"`
	Target []Target `yaml:"target"`
}
type Way struct {
	Name            string `yaml:"name"`
	Type            string `yaml:"type"`
	SMTPHost        string `yaml:"smtp_host,omitempty"`
	SMTPUsername    string `yaml:"smtp_username,omitempty"`
	SMTPPasswd      string `yaml:"smtp_passwd,omitempty"`
	SMTPPort        int    `yaml:"smtp_port,omitempty"`
	SMTPForm        string `yaml:"smtp_form,omitempty"`
	SMTPTLS         bool   `yaml:"smtp_tls,omitempty"`
	DingdingWebhook string `yaml:"dingding_webhook,omitempty"`
	WebhookURL      string `yaml:"webhook_url,omitempty"`
}
type Target struct {
	Name string   `yaml:"name"`
	To   []string `yaml:"to"`
}

type RuleEntry struct {
	Expr        string `yaml:"expr"`
	Description string `yaml:"description"`
	Level       string `yaml:"level"`
}
type Rule struct {
	Entry         []RuleEntry `yaml:"entry"`
	AssetName     string      `yaml:"asset_name"`
	AlerterTarget string      `yaml:"alerter_target"`
	AlerterWay    string      `yaml:"alerter_way"` // 添加 alerter_way
	TriggerCount  int         `yaml:"trigger_count"`
	ProbeInterval int         `yaml:"probe_interval"`
	For           int         `yaml:"for"`
}
type Config struct {
	Log     Log     `yaml:"log"`
	Asset   []Asset `yaml:"asset"`
	Alerter Alerter `yaml:"alerter"`
	Rule    []Rule  `yaml:"rule"`
}

type AlarmTrigger struct {
	AssetName     string
	AlerterTarget string
	AlerterWay    string
	Entry         RuleEntry
	ID            int
	Time          time.Time
	Instance      []string
	Value         []string
}

type PrometheusResp struct {
	Status string             `json:"status"`
	Data   PrometheusRespData `json:"data"`
}

type PrometheusRespData struct {
	ResultType string                     `json:"resultType"`
	Result     []PrometheusRespDataResult `json:"result"`
}

type PrometheusRespDataResult struct {
	Metric PrometheusRespDataResultMetric `json:"metric"`
	Value  []interface{}                  `json:"value"`
}

type PrometheusRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type AlarmContext struct {
	// 资产名字
	Asset string
	// 告警器名字
	Way string
	// 告警目标
	Target Target
	// 告警相关规则信息
	Entry RuleEntry
	// 触发主机
	Instance []string
	// 出发值
	Value []string
}
