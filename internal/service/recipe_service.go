package service

import (
	"context"
	"fmt"

	"github.com/arsenh/recipes-api/internal/models"
	"github.com/arsenh/recipes-api/internal/repository"
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
