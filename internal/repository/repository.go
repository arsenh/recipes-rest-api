package repository

import (
	"context"

	"github.com/arsenh/recipes-api/internal/models"
)

type RecipeRepository interface {
	// List returns all recipes (TODO: add pagination later if needed)
	List(ctx context.Context) ([]models.Recipe, error)

	// Create stores a new recipe and returns it with generated ID
	Create(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error)

	// GetByID returns a single recipe or error if not found
	GetByID(ctx context.Context, id string) (*models.Recipe, error)

	// Update replaces a recipe by ID
	Update(ctx context.Context, id string, recipe *models.Recipe) (*models.Recipe, error)

	// Delete removes a recipe by ID
	Delete(ctx context.Context, id string) error

	// SearchByTag returns recipes containing the given tag (case-insensitive)
	SearchByTag(ctx context.Context, tag string) ([]models.Recipe, error)
}
