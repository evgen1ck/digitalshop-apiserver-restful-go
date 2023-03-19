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
		ServiceName   string // example: flower shop
		ServiceUrl    string // example: https://example.com
		ApiServiceUrl string // example: https://api.example.com
		Host          string // example: localhost
		Port          string // example: 9990
		DebugMode     bool   // example: true
		JwtSecret     string // example: 3d93fe63ff2d2ebc9c0c814f7364fbba1653eddeb63a006b59cf7b2985545242
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
			logger.NewErrorWithExit("Failed to get project directory", err)
		}
		envFilePath := filepath.Join(dir, ".env")
		if err := godotenv.Load(envFilePath); err != nil {
			logger.NewErrorWithExit("Failed to load .env file", err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			logger.NewErrorWithExit("No .env file found", err)
		}
	}

	// General settings
	flag.StringVar(&cfg.App.ServiceName, "service-name", getEnv("APP_SERVICE_NAME", logger), "service name")
	flag.StringVar(&cfg.App.ServiceUrl, "service-url", getEnv("APP_SERVICE_URL", logger), "service url")
	flag.StringVar(&cfg.App.ApiServiceUrl, "api_v1-service-url", getEnv("APP_API_SERVICE_URL", logger), "api_v1 service url")
	flag.StringVar(&cfg.App.Host, "app-host", getEnv("APP_HOST", logger), "server host")
	flag.StringVar(&cfg.App.Port, "app-port", getEnv("APP_PORT", logger), "server port")
	flag.BoolVar(&cfg.App.DebugMode, "debug-mode", getEnvAsBool("APP_DEBUG_MODE", logger), "debug mode")
	flag.StringVar(&cfg.App.JwtSecret, "jwt-secret", getEnv("APP_JWT_SECRET", logger), "jwt secret")

	// STMP for noreply mail
	flag.StringVar(&cfg.Smtp1.Username, "mailer-username", getEnv("STMP1_USERNAME", logger), "mailer username")
	flag.StringVar(&cfg.Smtp1.Password, "mailer-password", getEnv("STMP1_PASSWORD", logger), "mailer password")
	flag.StringVar(&cfg.Smtp1.Host, "mailer-host", getEnv("STMP1_HOST", logger), "mailer host")
	flag.IntVar(&cfg.Smtp1.Port, "mailer-port", getEnvAsInt("STMP1_PORT", logger), "mailer port")
	flag.StringVar(&cfg.Smtp1.From, "mailer-from", getEnv("STMP1_FROM", logger), "mailer sender")

	// Postgres DSN
	flag.StringVar(&cfg.Postgres.User, "database-user", getEnv("POSTGRES_USER", logger), "username for database")
	flag.StringVar(&cfg.Postgres.Password, "database-password", getEnv("POSTGRES_PASSWORD", logger), "password for database")
	flag.StringVar(&cfg.Postgres.Ip, "database-ip", getEnv("POSTGRES_IP", logger), "hostname/address for database")
	flag.StringVar(&cfg.Postgres.Port, "database-port", getEnv("POSTGRES_PORT", logger), "port for database")
	flag.StringVar(&cfg.Postgres.Database, "database-database", getEnv("POSTGRES_DATABASE", logger), "maintenance database for database")

	// TLS files
	flag.StringVar(&cfg.Tls.CertFile, "tls-cert-file", getEnv("TLS_CERTFILE", logger), "tls certificate file")
	flag.StringVar(&cfg.Tls.KeyFile, "tls-key-file", getEnv("TLS_KEYFILE", logger), "tls key file")

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
