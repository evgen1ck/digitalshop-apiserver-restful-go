package config

import (
	"flag"
	"github.com/joho/godotenv"
	"test-server-go/internal/env"
	"test-server-go/internal/logger"
)

type Config struct {
	App struct {
		ServiceUrl string // example: example.com
		Host       string // example: localhost
		Port       string // example: 3000
		DebugMode  bool   // example: true
		JwtSecret  string // example: 3d93fe63ff2d2ebc9c0c814f7364fbba1653eddeb63a006b59cf7b2985545242
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

func New(logger *logger.Logger) (*Config, error) {
	var cfg Config

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		logger.NewError("No .env file found", err)
	}

	// General settings
	flag.StringVar(&cfg.App.ServiceUrl, "service-url", env.GetEnv("APP_SERVICE_URL", logger), "service url")
	flag.StringVar(&cfg.App.Host, "app-host", env.GetEnv("APP_HOST", logger), "server host")
	flag.StringVar(&cfg.App.Port, "app-port", env.GetEnv("APP_PORT", logger), "server port")
	flag.BoolVar(&cfg.App.DebugMode, "debug-mode", env.GetEnvAsBool("APP_DEBUG_MODE", logger), "debug mode")
	flag.StringVar(&cfg.App.JwtSecret, "jwt-secret", env.GetEnv("APP_JWT_SECRET", logger), "jwt secret")

	// STMP for noreply mail
	flag.StringVar(&cfg.Smtp1.Username, "mailer-username", env.GetEnv("STMP1_USERNAME", logger), "mailer username")
	flag.StringVar(&cfg.Smtp1.Password, "mailer-password", env.GetEnv("STMP1_PASSWORD", logger), "mailer password")
	flag.StringVar(&cfg.Smtp1.Host, "mailer-host", env.GetEnv("STMP1_HOST", logger), "mailer host")
	flag.IntVar(&cfg.Smtp1.Port, "mailer-port", env.GetEnvAsInt("STMP1_PORT", logger), "mailer port")
	flag.StringVar(&cfg.Smtp1.From, "mailer-from", env.GetEnv("STMP1_FROM", logger), "mailer sender")

	// Postgres DSN
	flag.StringVar(&cfg.Postgres.User, "postgres-user", env.GetEnv("POSTGRES_USER", logger), "username for postgres")
	flag.StringVar(&cfg.Postgres.Password, "postgres-password", env.GetEnv("POSTGRES_PASSWORD", logger), "password for postgres")
	flag.StringVar(&cfg.Postgres.Ip, "postgres-ip", env.GetEnv("POSTGRES_IP", logger), "hostname/address for postgres")
	flag.StringVar(&cfg.Postgres.Port, "postgres-port", env.GetEnv("POSTGRES_PORT", logger), "port for postgres")
	flag.StringVar(&cfg.Postgres.Database, "postgres-database", env.GetEnv("POSTGRES_DATABASE", logger), "maintenance database for postgres")

	// TLS files
	flag.StringVar(&cfg.Tls.CertFile, "tls-cert-file", env.GetEnv("TLS_CERTFILE", logger), "tls certificate file")
	flag.StringVar(&cfg.Tls.KeyFile, "tls-key-file", env.GetEnv("TLS_KEYFILE", logger), "tls key file")

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
