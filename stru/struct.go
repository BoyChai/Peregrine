package stru

type Asset struct {
	Name  string   `yaml:"name"`
	Hosts []string `yaml:"hosts"`
}

type Alerter struct {
	SMTPHost     string   `yaml:"smtp_host"`
	SMTPUsername string   `yaml:"smtp_username"`
	SMTPPasswd   string   `yaml:"smtp_passwd"`
	SMTPPort     int      `yaml:"smtp_port"`
	SMTPForm     string   `yaml:"smtp_form"`
	Target       []Target `yaml:"target"`
}

type Target struct {
	Name      string   `yaml:"name"`
	EmailTo   []string `yaml:"email_to,omitempty"`
	WebhookTo []string `yaml:"webhook_to,omitempty"`
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
	TriggerCount  int         `yaml:"trigger_count"`
	ProbeInterval int         `yaml:"probe_interval"`
	For           string      `yaml:"for"`
}

type Config struct {
	Asset   []Asset `yaml:"asset"`
	Alerter Alerter `yaml:"alerter"`
	Rule    []Rule  `yaml:"rule"`
}
