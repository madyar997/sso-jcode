package main

import (
	"log"

	"github.com/madyar997/sso-jcode/config"
	"github.com/madyar997/sso-jcode/internal/app"
)

func main() {
	cfg, err := config.NewViperConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
