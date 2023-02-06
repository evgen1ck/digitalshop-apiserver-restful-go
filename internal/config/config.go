package config

import (
	"flag"
	"test-server-go/internal/env"
)

type Config struct {
	App struct {
		Host      string
		Port      string
		DebugMode bool
	}
	Smtp1 struct {
		Username string
		Password string
		Host     string
		Port     int
		From     string
	}
	Postgres struct {
		User     string
		Password string
		Ip       string
		Port     string
		Database string
	}
	Tls struct {
		CertFile string
		KeyFile  string
	}
}

func New() (*Config, error) {
	var cfg Config

	// General settings
	flag.StringVar(&cfg.App.Host, "app-host", env.GetEnv("APP_HOST"), "server host")
	flag.StringVar(&cfg.App.Port, "app-port", env.GetEnv("APP_PORT"), "server port")
	flag.BoolVar(&cfg.App.DebugMode, "debug-mode", env.GetEnvAsBool("APP_DEBUG_MODE"), "debug mode")

	// STMP for noreply mail
	flag.StringVar(&cfg.Smtp1.Username, "mailer-username", env.GetEnv("STMP1_USERNAME"), "mailer username")
	flag.StringVar(&cfg.Smtp1.Password, "mailer-password", env.GetEnv("STMP1_PASSWORD"), "mailer password")
	flag.StringVar(&cfg.Smtp1.Host, "mailer-host", env.GetEnv("STMP1_HOST"), "mailer host")
	flag.IntVar(&cfg.Smtp1.Port, "mailer-port", env.GetEnvAsInt("STMP1_PORT"), "mailer port")
	flag.StringVar(&cfg.Smtp1.From, "mailer-from", env.GetEnv("STMP1_FROM"), "mailer sender")

	// Postgres DSN
	flag.StringVar(&cfg.Postgres.User, "postgres-user", env.GetEnv("POSTGRES_USER"), "username for postgres")
	flag.StringVar(&cfg.Postgres.Password, "postgres-password", env.GetEnv("POSTGRES_PASSWORD"), "password for postgres")
	flag.StringVar(&cfg.Postgres.Ip, "postgres-ip", env.GetEnv("POSTGRES_IP"), "hostname/address for postgres")
	flag.StringVar(&cfg.Postgres.Port, "postgres-port", env.GetEnv("POSTGRES_PORT"), "port for postgres")
	flag.StringVar(&cfg.Postgres.Database, "postgres-database", env.GetEnv("POSTGRES_DATABASE"), "maintenance database for postgres")

	// TLS files
	flag.StringVar(&cfg.Tls.CertFile, "tls-cert-file", env.GetEnv("TLS_CERTFILE"), "tls certificate file")
	flag.StringVar(&cfg.Tls.KeyFile, "tls-key-file", env.GetEnv("TLS_KEYFILE"), "tls key file")

	flag.Parse()

	return &cfg, nil
}

func (cfg *Config) GetPostgresDSN() string {
	return cfg.Postgres.User + ":" +
		cfg.Postgres.Password + "@" +
		cfg.Postgres.Ip + ":" +
		cfg.Postgres.Port + "/" +
		cfg.Postgres.Database
}

func (cfg *Config) GetURL() string {
	return cfg.App.Host + ":" + cfg.App.Port
}
