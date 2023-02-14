package main

import (
	"fmt"
	"log"
	"test-server-go/internal/config"
	"test-server-go/internal/mailer"
	"test-server-go/internal/postgres"
	"test-server-go/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		//logger.Fatal(err)
		log.Fatal(err)
	}

	pdb, err := postgres.New(cfg.GetPostgresDSN())
	if err != nil {
		//logger.Fatal(err)
		log.Fatal(err)
	}
	defer pdb.Close()

	smtpMailer := mailer.NewSmtp(*cfg)

	app := &server.Application{
		Config:   cfg,
		Postgres: pdb,
		Mailer:   smtpMailer,
		//logger:       logger,
	}

	fmt.Println(app.Config)
	fmt.Println("OK")

	err = app.ServerRun()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("server is stopped")
}
