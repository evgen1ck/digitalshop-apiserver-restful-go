package config

import (
	"flag"
	"os"
	"path/filepath"
	tl "test-server-go/internal/tools"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Service struct {
			Name string `yaml:"name"`
			Url  struct {
				Client string `yaml:"client"`
				Server string `yaml:"server"`
			} `yaml:"url"`
		} `yaml:"service"`
		Port  int    `yaml:"port"`
		Debug bool   `yaml:"debug"`
		Jwt   string `yaml:"jwt"`
	} `yaml:"app"`
	Prometheus struct {
		Port int `yaml:"port"`
	} `yaml:"prometheus"`
	MailNoreply struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		From     string `yaml:"from"`
	} `yaml:"mailNoreply"`
	MailSupport struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		From     string `yaml:"from"`
	} `yaml:"mailSupport"`
	Postgres struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Ip       string `yaml:"ip"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`
	Redis struct {
		Password string `yaml:"password"`
		Ip       string `yaml:"ip"`
		Port     int    `yaml:"port"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
	Tls struct {
		CertFile string `yaml:"certfile"`
		KeyFile  string `yaml:"keyfile"`
	} `yaml:"tls"`
	Payments struct {
		Freekassa struct {
			ShopId           int    `yaml:"shopId"`
			ApiKey           string `yaml:"apiKey"`
			FirstSecretWord  string `yaml:"firstSecretWord"`
			SecondSecretWord string `yaml:"secondSecretWord"`
		} `yaml:"freekassa"`
	} `yaml:"payments"`
}

func SetupYaml() (*Config, error) {
	var cfg Config
	path, err := tl.GetExecutablePath()
	if err != nil {
		return nil, err
	}

	configPath := flag.String("config", filepath.Join(path, "server.yaml"), "Path to the YAML configuration file")
	yamlFile, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}

	// General settings
	flag.StringVar(&cfg.App.Service.Name, "app-service-name", cfg.App.Service.Name, "service name")
	flag.StringVar(&cfg.App.Service.Url.Client, "app-service-url-client", cfg.App.Service.Url.Client, "service url for client app")
	flag.StringVar(&cfg.App.Service.Url.Server, "app-service-url-server", cfg.App.Service.Url.Server, "service url for server app")
	flag.IntVar(&cfg.App.Port, "app-port", cfg.App.Port, "server port")
	flag.BoolVar(&cfg.App.Debug, "app-debug", cfg.App.Debug, "debug mode")
	flag.StringVar(&cfg.App.Jwt, "app-jwt", cfg.App.Jwt, "jwt secret")

	// Prometheus settings
	flag.IntVar(&cfg.Prometheus.Port, "prometheus-port", cfg.Prometheus.Port, "prometheus port")

	// E-mail for noreply@example.com
	flag.StringVar(&cfg.MailNoreply.Username, "mail-noreply-username", cfg.MailNoreply.Username, "noreply mail username")
	flag.StringVar(&cfg.MailNoreply.Password, "mail-noreply-password", cfg.MailNoreply.Password, "noreply mail password")
	flag.StringVar(&cfg.MailNoreply.Host, "mail-noreply-host", cfg.MailNoreply.Host, "noreply mail host")
	flag.IntVar(&cfg.MailNoreply.Port, "mail-noreply-port", cfg.MailNoreply.Port, "noreply mail port")
	flag.StringVar(&cfg.MailNoreply.From, "mail-noreply-from", cfg.MailNoreply.From, "noreply mail sender")

	// E-mail for support@example.com
	flag.StringVar(&cfg.MailSupport.Username, "mail-support-username", cfg.MailSupport.Username, "support mail username")
	flag.StringVar(&cfg.MailSupport.Password, "mail-support-password", cfg.MailSupport.Password, "support mail password")
	flag.StringVar(&cfg.MailSupport.Host, "mail-support-host", cfg.MailSupport.Host, "support mail host")
	flag.IntVar(&cfg.MailSupport.Port, "mail-support-port", cfg.MailSupport.Port, "support mail port")
	flag.StringVar(&cfg.MailSupport.From, "mail-support-from", cfg.MailSupport.From, "support mail sender")

	// Postgres
	flag.StringVar(&cfg.Postgres.User, "postgres-user", cfg.Postgres.User, "username for postgres")
	flag.StringVar(&cfg.Postgres.Password, "postgres-password", cfg.Postgres.Password, "password for postgres password")
	flag.StringVar(&cfg.Postgres.Ip, "postgres-ip", cfg.Postgres.Ip, "hostname/address for postgres")
	flag.IntVar(&cfg.Postgres.Port, "postgres-port", cfg.Postgres.Port, "port for postgres")
	flag.StringVar(&cfg.Postgres.Database, "postgres-database", cfg.Postgres.Database, "maintenance database for postgres")

	// Redis
	flag.StringVar(&cfg.Redis.Password, "redis-password", cfg.Redis.Password, "password for redis password")
	flag.StringVar(&cfg.Redis.Ip, "redis-ip", cfg.Redis.Ip, "hostname/address for redis")
	flag.IntVar(&cfg.Redis.Port, "redis-port", cfg.Redis.Port, "port for redis")
	flag.IntVar(&cfg.Redis.Database, "redis-database", cfg.Redis.Database, "maintenance database for redis")

	// TLS
	flag.StringVar(&cfg.Tls.CertFile, "tls-certfile", cfg.Tls.CertFile, "tls certificate file")
	flag.StringVar(&cfg.Tls.KeyFile, "tls-keyfile", cfg.Tls.KeyFile, "tls key file")

	// Payments
	flag.IntVar(&cfg.Payments.Freekassa.ShopId, "payments-freekassa-shopId", cfg.Payments.Freekassa.ShopId, "payments shop id for freekassa")
	flag.StringVar(&cfg.Payments.Freekassa.ApiKey, "payments-freekassa-apiKey", cfg.Payments.Freekassa.ApiKey, "payments api key for freekassa")
	flag.StringVar(&cfg.Payments.Freekassa.FirstSecretWord, "payments-freekassa-firstSecretWord", cfg.Payments.Freekassa.FirstSecretWord, "payments first secret word for freekassa")
	flag.StringVar(&cfg.Payments.Freekassa.SecondSecretWord, "payments-freekassa-secondSecretWord", cfg.Payments.Freekassa.SecondSecretWord, "payments second secret word for freekassa")

	flag.Parse()

	return &cfg, nil
}
