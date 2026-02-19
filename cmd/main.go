package main

import (
	"log"

	"github.com/arsenh/recipes-api/internal/app"
	"github.com/arsenh/recipes-api/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	app := app.New(cfg)
	defer app.Close()

	app.Router.Run(":8080")
}
