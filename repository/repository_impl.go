package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gomarket-loyalty/config"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
)

func NewRepository(database *mongo.Database) Repository {
	return &repositoryImpl{
		db: database,
	}
}

type repositoryImpl struct {
	db *mongo.Database
}

func (repository *repositoryImpl) SetUser(user model.User) error {
	ctx, cancel := config.NewMongoContext()
	defer cancel()
	_, err := repository.db.Collection("users").InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return exception.ErrLoginAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}
