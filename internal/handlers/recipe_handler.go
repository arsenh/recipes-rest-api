// @title Recipes API
// @version 1.0
// @description Simple REST API for managing recipes
// @BasePath /
// @schemes http https
// @host localhost:8080
// @contact.name Arsen
// @contact.email your.email@example.com
// @license.name MIT

package handlers

import (
	"log"
	"net/http"

	"github.com/arsenh/recipes-api/dto"
	apperrors "github.com/arsenh/recipes-api/internal/errors"
	"github.com/arsenh/recipes-api/internal/service"
	"github.com/gin-gonic/gin"
)

// RecipeHandler handles all recipe-related HTTP requests
type RecipeHandler struct {
	service *service.RecipeService
}

func NewRecipeHander(service *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: service}
}

// ListRecipesHandler godoc
// @Summary List all recipes
// @Description Returns the complete list of recipes currently in the system
// @Tags recipes
// @Produce json
// @Success 200 {array} models.Recipe
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /recipes [get]
func (h *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	recipes, err := h.service.ListRecipes(c.Request.Context())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't get list of recipes",
		})
		return
	}
	c.JSON(http.StatusOK, recipes)
}

// GetRecipeByIdHandler godoc
// @Summary Get recipe by ID
// @Description Returns a single recipe by its unique ID
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID" example(ckx123abc456)
// @Success 200 {object} models.Recipe
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Recipe not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /recipes/{id} [get]
func (h *RecipeHandler) GetRecipeByIdHandler(c *gin.Context) {
	id := c.Param("id")
	recipe, err := h.service.GetRecipeById(c.Request.Context(), id)

	if err != nil {
		switch err {
		case apperrors.ErrBadID:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		case apperrors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
	c.JSON(http.StatusOK, recipe)
}

// NewRecipeHandler godoc
// @Summary Create a new recipe
// @Description Creates a new recipe and returns it with generated ID and timestamp
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body models.Recipe true "Recipe data"
// @Success 201 {object} models.Recipe
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /recipes [post]
func (h *RecipeHandler) NewRecipeHandler(c *gin.Context) {
	var recipe dto.CreateRecipeRequest
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	createdRecipe, err := h.service.NewRecipe(c.Request.Context(), recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, createdRecipe)
}

// UpdateRecipeHandler godoc
// @Summary Update an existing recipe
// @Description Updates a recipe by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param recipe body models.Recipe true "Updated recipe data"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Recipe not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /recipes/{id} [put]
func (h *RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe dto.UpdateRecipeRequest
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	updatedRecipe, err := h.service.UpdateRecipeById(c.Request.Context(), id, recipe)
	if err != nil {
		switch err {
		case apperrors.ErrBadID:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		case apperrors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
	c.JSON(http.StatusOK, updatedRecipe)
}

// DeleteRecipeHandler godoc
// @Summary Delete a recipe
// @Description Deletes a recipe by ID
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} map[string]string "Recipe has been deleted"
// @Failure 404 {object} map[string]string "Recipe not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /recipes/{id} [delete]
func (h *RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteRecipeById(c.Request.Context(), id)
	if err != nil {
		switch err {
		case apperrors.ErrBadID:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		case apperrors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

// SearchRecipeHandler godoc
// @Summary Search recipes by tag
// @Description Returns all recipes that contain the given tag (case-insensitive)
// @Tags recipes
// @Produce json
// @Param tag query string true "Recipe tag (e.g. 'vegan', 'dessert')"
// @Success 200 {array} models.Recipe
// @Failure 400 {object} map[string]string "Tag parameter is required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /recipes/search [get]
func (h *RecipeHandler) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")

	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tag parameter is required"})
	}

	recipes, err := h.service.SearchByTag(c.Request.Context(), tag)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, recipes)
}
