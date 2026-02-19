package repository

import (
	"context"

	"github.com/arsenh/recipes-api/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type recipeRepositoryMongo struct {
	collection *mongo.Collection
}

func NewRecipeRepository(db *mongo.Database) RecipeRepository {
	return &recipeRepositoryMongo{
		collection: db.Collection("recipes"),
	}
}

func (r *recipeRepositoryMongo) List(ctx context.Context) ([]models.Recipe, error) {
	return nil, nil
}

func (r *recipeRepositoryMongo) Create(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error) {
	return nil, nil
}

func (r *recipeRepositoryMongo) GetByID(ctx context.Context, id string) (*models.Recipe, error) {
	return nil, nil
}

func (r *recipeRepositoryMongo) Update(ctx context.Context, id string, recipe *models.Recipe) (*models.Recipe, error) {
	return nil, nil
}

func (r *recipeRepositoryMongo) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *recipeRepositoryMongo) SearchByTag(ctx context.Context, tag string) ([]models.Recipe, error) {
	return nil, nil
}
