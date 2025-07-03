package smtpclient

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

var config SMTPConfig

var defaultConfig = SMTPConfig{
	Host:     "smtp.example.com",
	Port:     25,
	Username: "",
	Password: "",
	UseTLS:   false,
}

func init() {
	config = defaultConfig
}

func Setup(cfg SMTPConfig) {
	config = cfg
}
