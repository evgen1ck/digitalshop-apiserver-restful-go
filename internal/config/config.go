package config

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
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
	osName := runtime.GOOS
	if osName == "linux" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.NewError("Failed to get project directory", err)
		}
		envFilePath := filepath.Join(dir, ".env")
		if err := godotenv.Load(envFilePath); err != nil {
			logger.NewError("Failed to load .env file", err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			logger.NewError("No .env file found", err)
		}
	}

	// General settings
	flag.StringVar(&cfg.App.ServiceUrl, "service-url", GetEnv("APP_SERVICE_URL", logger), "service url")
	flag.StringVar(&cfg.App.Host, "app-host", GetEnv("APP_HOST", logger), "server host")
	flag.StringVar(&cfg.App.Port, "app-port", GetEnv("APP_PORT", logger), "server port")
	flag.BoolVar(&cfg.App.DebugMode, "debug-mode", GetEnvAsBool("APP_DEBUG_MODE", logger), "debug mode")
	flag.StringVar(&cfg.App.JwtSecret, "jwt-secret", GetEnv("APP_JWT_SECRET", logger), "jwt secret")

	// STMP for noreply mail
	flag.StringVar(&cfg.Smtp1.Username, "mailer-username", GetEnv("STMP1_USERNAME", logger), "mailer username")
	flag.StringVar(&cfg.Smtp1.Password, "mailer-password", GetEnv("STMP1_PASSWORD", logger), "mailer password")
	flag.StringVar(&cfg.Smtp1.Host, "mailer-host", GetEnv("STMP1_HOST", logger), "mailer host")
	flag.IntVar(&cfg.Smtp1.Port, "mailer-port", GetEnvAsInt("STMP1_PORT", logger), "mailer port")
	flag.StringVar(&cfg.Smtp1.From, "mailer-from", GetEnv("STMP1_FROM", logger), "mailer sender")

	// Postgres DSN
	flag.StringVar(&cfg.Postgres.User, "database-user", GetEnv("POSTGRES_USER", logger), "username for database")
	flag.StringVar(&cfg.Postgres.Password, "database-password", GetEnv("POSTGRES_PASSWORD", logger), "password for database")
	flag.StringVar(&cfg.Postgres.Ip, "database-ip", GetEnv("POSTGRES_IP", logger), "hostname/address for database")
	flag.StringVar(&cfg.Postgres.Port, "database-port", GetEnv("POSTGRES_PORT", logger), "port for database")
	flag.StringVar(&cfg.Postgres.Database, "database-database", GetEnv("POSTGRES_DATABASE", logger), "maintenance database for database")

	// TLS files
	flag.StringVar(&cfg.Tls.CertFile, "tls-cert-file", GetEnv("TLS_CERTFILE", logger), "tls certificate file")
	flag.StringVar(&cfg.Tls.KeyFile, "tls-key-file", GetEnv("TLS_KEYFILE", logger), "tls key file")

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
