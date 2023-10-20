package main

import (
	"log"

	"github.com/madyar997/practice_7/config"
	"github.com/madyar997/practice_7/internal/app"
)

func main() {
	cfg, err := config.NewViperConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
