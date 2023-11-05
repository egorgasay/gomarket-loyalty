package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func NewProductRepository(database *mongo.Database) ProductRepository {
	return &productRepositoryImpl{
		Collection: database.Collection("TEST"),
	}
}

type productRepositoryImpl struct {
	Collection *mongo.Collection
}
