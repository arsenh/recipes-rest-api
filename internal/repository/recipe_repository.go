package repository

import (
	"context"
	"log"

	apperrors "github.com/arsenh/recipes-api/internal/errors"
	"github.com/arsenh/recipes-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "recipes"

type recipeRepositoryMongo struct {
	collection *mongo.Collection
}

func NewRecipeRepository(db *mongo.Database) RecipeRepository {
	return &recipeRepositoryMongo{
		collection: db.Collection(collectionName),
	}
}

func (r *recipeRepositoryMongo) List(ctx context.Context) ([]models.Recipe, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("error on geting cursor for all documents in recipes collection")
		return nil, err
	}
	defer cur.Close(ctx)

	recipes := []models.Recipe{}

	for cur.Next(ctx) {
		var recipe models.Recipe
		err := cur.Decode(&recipe)
		if err != nil {
			log.Println("error on decode recipe document")
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *recipeRepositoryMongo) Create(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error) {
	_, err := r.collection.InsertOne(ctx, recipe)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

func (r *recipeRepositoryMongo) GetByID(ctx context.Context, id string) (*models.Recipe, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return nil, apperrors.ErrBadID
	}

	filter := bson.M{"_id": objID}

	var recipe models.Recipe
	err = r.collection.FindOne(ctx, filter).Decode(&recipe)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &recipe, nil
}

func (r *recipeRepositoryMongo) Update(ctx context.Context, id string, recipe *models.Recipe) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return apperrors.ErrBadID
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"name":         recipe.Name,
			"tags":         recipe.Tags,
			"ingredients":  recipe.Ingredients,
			"instructions": recipe.Instructions,
			// ID and PublishedAt are NOT included → they stay unchanged
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("error on update recipe")
		return err
	}

	if result.MatchedCount == 0 && result.ModifiedCount == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *recipeRepositoryMongo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return apperrors.ErrBadID
	}

	filter := bson.M{"_id": objID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("error on delete recipe")
		return err
	}
	if result.DeletedCount == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *recipeRepositoryMongo) SearchByTag(ctx context.Context, tag string) ([]models.Recipe, error) {

	filter := bson.M{"tags": tag}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	recipes := []models.Recipe{}
	for cursor.Next(ctx) {
		var recipe models.Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}
