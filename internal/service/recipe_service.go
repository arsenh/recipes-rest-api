package service

import "github.com/arsenh/recipes-api/internal/repository"

type RecipeService struct {
	repo repository.RecipeRepository
}

func NewRecipeService(repo repository.RecipeRepository) *RecipeService {
	return &RecipeService{repo: repo}
}
