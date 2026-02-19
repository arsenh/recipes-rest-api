package app

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/arsenh/recipes-api/docs"
	"github.com/arsenh/recipes-api/internal/config"
	"github.com/arsenh/recipes-api/internal/database"
	"github.com/arsenh/recipes-api/internal/handlers"
	"github.com/arsenh/recipes-api/internal/models"
	"github.com/arsenh/recipes-api/internal/repository"
	"github.com/arsenh/recipes-api/internal/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
)

type App struct {
	Router  *gin.Engine
	MongoDB *database.MongoDatabase
}

func addDummyDataToMongoDB(mongoDb *database.MongoDatabase) {

	var recipes []models.Recipe

	bytes, err := os.ReadFile("DB.json")
	if err != nil {
		log.Println("Cannot open DB.json file")
		os.Exit(-1)
	}
	if err = json.Unmarshal(bytes, &recipes); err != nil {
		log.Println("Error on parsing json DB data")
	}

	ctx := context.Background()

	collection := mongoDb.DB.Collection("recipes")
	// DELETE all existing data
	_, err = collection.DeleteMany(ctx, bson.M{})

	if err != nil {
		log.Fatal("Error clearing recipes collection:", err)
	}

	var listOfRecipes []interface{}
	for _, r := range recipes {
		recipe := bson.M{
			"name":         r.Name,
			"tags":         r.Tags,
			"ingredients":  r.Ingredients,
			"instructions": r.Instructions,
			"publishedAt":  r.PublishedAt,
		}
		listOfRecipes = append(listOfRecipes, recipe)
	}

	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal("Bulk write exception:", err)
	}
	log.Println("Inserted recipes:", len(insertManyResult.InsertedIDs))
}

func New(config *config.Config) *App {
	// connect to database
	mongoDb, err := database.ConnectMongo(config.DatabaseURL, config.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}

	// Add some data in database
	// DELETE: for testing only
	addDummyDataToMongoDB(mongoDb)

	// create repositories
	recipeRepo := repository.NewRecipeRepository(mongoDb.DB)

	// create service that will use this repository
	recipeService := service.NewRecipeService(recipeRepo)
	// create handler which will use service
	handler := handlers.NewRecipeHander(recipeService)

	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/recipes", handler.ListRecipesHandler)
	router.GET("/recipes/:id", handler.GetRecipeByIdHandler)
	router.POST("/recipes", handler.NewRecipeHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.GET("/recipes/search", handler.SearchRecipeHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &App{Router: router}
}

func (a *App) Close() {
	if a.MongoDB != nil {
		_ = a.MongoDB.Close(context.Background())
	}
}
