// @title Recipes API
// @version 1.0
// @description Simple REST API for managing recipes
// @BasePath /

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	docs "github.com/arsenh/recipes-api/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Tags         []string  `json:"tags" binding:"required"`
	Ingredients  []string  `json:"ingredients" binding:"required"`
	Instructions []string  `json:"instructions" binding:"required"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

var ctx context.Context
var err error
var client *mongo.Client

func init() {
	bytes, err := os.ReadFile("DB.json")
	if err != nil {
		fmt.Println("Cannot open DB.json file")
		os.Exit(-1)
	}
	if err = json.Unmarshal(bytes, &recipes); err != nil {
		fmt.Println("Error on parsing json DB data")
	}

	ctx = context.Background()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println("Cannot connect to MongoDB")
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Iserted recipes: ", len(insertManyResult.InsertedIDs))
}

// ListRecipesHandler godoc
// @Summary List all recipes
// @Description Get all recipes
// @Tags recipes
// @Produce json
// @Success 200 {array} Recipe
// @Router /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// NewRecipeHandler godoc
// @Summary Create a new recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body Recipe true "Recipe data"
// @Success 201 {object} Recipe
// @Failure 400
// @Router /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, recipe)
}

// UpdateRecipeHandler godoc
// @Summary Update a recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param recipe body Recipe true "Recipe data"
// @Success 200 {object} Recipe
// @Failure 404
// @Router /recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipe.ID = id
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipeHandler godoc
// @Summary Delete a recipe
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200
// @Failure 404
// @Router /recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

// SearchRecipeHandler godoc
// @Summary Search recipes by tag
// @Tags recipes
// @Produce json
// @Param tag query string true "Recipe tag"
// @Success 200 {array} Recipe
// @Router /recipes/search [get]
func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}

		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	router.GET("/recipes", ListRecipesHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run()
}
