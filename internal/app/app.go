package app

import (
	"github.com/arsenh/recipes-api/docs"
	"github.com/arsenh/recipes-api/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	Router *gin.Engine
}

func New() *App {
	// load config
	// connect to database
	// create repository to give database url
	// create service that will use this repository
	// create handler which will use service
	handler := handlers.NewRecipeHander(nil)

	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/recipes", handler.ListRecipesHandler)
	router.POST("/recipes", handler.NewRecipeHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.GET("/recipes/search", handler.SearchRecipeHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &App{Router: router}
}
