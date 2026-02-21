package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arsenh/recipes-api/dto"
	"github.com/arsenh/recipes-api/internal/models"
	"github.com/arsenh/recipes-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeService struct {
	repo repository.RecipeRepository
}

func NewRecipeService(repo repository.RecipeRepository) *RecipeService {
	return &RecipeService{repo: repo}
}

func (r *RecipeService) ListRecipes(ctx context.Context) ([]models.Recipe, error) {
	recipes, err := r.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error on getting list of recipes from repository")
	}
	return recipes, nil
}

func (r *RecipeService) GetRecipeById(ctx context.Context, id string) (*models.Recipe, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *RecipeService) NewRecipe(ctx context.Context, recipe dto.CreateRecipeRequest) (*models.Recipe, error) {
	newRecipe := models.Recipe{
		ID:           primitive.NewObjectID(),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  time.Now(),
	}

	return r.repo.Create(ctx, &newRecipe)
}

func (r *RecipeService) UpdateRecipeById(ctx context.Context, id string, recipe dto.UpdateRecipeRequest) (*models.Recipe, error) {

	updatedRecipe := models.Recipe{
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
	}

	err := r.repo.Update(ctx, id, &updatedRecipe)
	if err != nil {
		return nil, err
	}

	return r.repo.GetByID(ctx, id)
}

func (r *RecipeService) DeleteRecipeById(ctx context.Context, id string) error {
	return r.repo.Delete(ctx, id)
}

func (r *RecipeService) SearchByTag(ctx context.Context, tag string) ([]models.Recipe, error) {
	return r.repo.SearchByTag(ctx, tag)
}
