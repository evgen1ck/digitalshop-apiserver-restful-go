package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	App struct {
		Service struct {
			Name string `yaml:"name"`
			Url  struct {
				App string `yaml:"app"`
				Api string `yaml:"api"`
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
	Tls struct {
		CertFile string `yaml:"certfile"`
		KeyFile  string `yaml:"keyfile"`
	} `yaml:"tls"`
}

func SetupYaml() (*Config, error) {
	var cfg Config

	path, _ := os.Getwd()
	configPath := flag.String("config", filepath.Join(path, "server.yaml"), "Path to the YAML configuration file")
	yamlFile, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	// General settings
	flag.StringVar(&cfg.App.Service.Name, "app-service-name", cfg.App.Service.Name, "service name")
	flag.StringVar(&cfg.App.Service.Url.App, "app-service-url-app", cfg.App.Service.Url.App, "service url for app")
	flag.StringVar(&cfg.App.Service.Url.Api, "app-service-url-api", cfg.App.Service.Url.Api, "service url for api")
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
	flag.StringVar(&cfg.Postgres.Password, "postgres-password", cfg.Postgres.Password, "password for postgres username")
	flag.StringVar(&cfg.Postgres.Ip, "postgres-ip", cfg.Postgres.Ip, "hostname/address for postgres")
	flag.IntVar(&cfg.Postgres.Port, "postgres-port", cfg.Postgres.Port, "port for postgres")
	flag.StringVar(&cfg.Postgres.Database, "postgres-database", cfg.Postgres.Database, "maintenance database for postgres")

	// TLS
	flag.StringVar(&cfg.Tls.CertFile, "tls-certfile", cfg.Tls.CertFile, "tls certificate file")
	flag.StringVar(&cfg.Tls.KeyFile, "tls-keyfile", cfg.Tls.KeyFile, "tls key file")

	flag.Parse()

	return &cfg, nil
}

func (cfg *Config) GetPostgresDSN() string {
	return cfg.Postgres.User + ":" +
		cfg.Postgres.Password + "@" +
		cfg.Postgres.Ip + ":" +
		strconv.Itoa(cfg.Postgres.Port) + "/" +
		cfg.Postgres.Database
}
