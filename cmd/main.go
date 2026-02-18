// @title Recipes API
// @version 1.0
// @description Simple REST API for managing recipes
// @BasePath /

package main

import (
	"github.com/joho/godotenv"

	"github.com/arsenh/recipes-api/internal/app"
)

func main() {
	godotenv.Load()
	app := app.New()
	app.Router.Run(":8080")
}
