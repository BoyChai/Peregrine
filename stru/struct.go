package stru

type Asset struct {
	Name  string `yaml:"name"`
	Hosts string `yaml:"hosts"`
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
	AssetGroup    string      `yaml:"asset_group"`
	AlerterTarget string      `yaml:"alerter_target"`
	AlerterWay    string      `yaml:"alerter_way"` // 添加 alerter_way
	TriggerCount  int         `yaml:"trigger_count"`
	ProbeInterval int         `yaml:"probe_interval"`
	For           string      `yaml:"for"`
}
type Config struct {
	Asset   []Asset `yaml:"asset"`
	Alerter Alerter `yaml:"alerter"`
	Rule    []Rule  `yaml:"rule"`
}
