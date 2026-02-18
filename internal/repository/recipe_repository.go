package repository

import "go.mongodb.org/mongo-driver/mongo"

type RecipeRepository struct {
	collection *mongo.Collection
}

func NewRecipeRepository(db *mongo.Database) *RecipeRepository {
	return &RecipeRepository{
		collection: db.Collection("recipes"),
	}
}
