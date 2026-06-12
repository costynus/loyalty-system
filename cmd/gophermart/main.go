package main

import (
	"log"

	"github.com/costynus/loyalty-system/config"
	"github.com/costynus/loyalty-system/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
        log.Fatalf("Config error: %s", err)
	}

    app.Run(cfg)
}
