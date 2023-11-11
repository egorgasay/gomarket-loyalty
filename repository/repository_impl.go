package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (repository *repositoryImpl) SetUser(ctx context.Context, user model.User) error {
	_, err := repository.db.Collection("users").InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return exception.ErrAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (repository *repositoryImpl) AddMechanic(ctx context.Context, bonus model.Mechanic) error {
	_, err := repository.db.Collection("mechanics").InsertOne(ctx, bonus)
	if mongo.IsDuplicateKeyError(err) {
		return exception.ErrAlreadyExists
	}

	if err != nil {
		return err
	}
	return nil

}

func (repository *repositoryImpl) UpdateBonusUser(ctx context.Context, id string, bonus int) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"bonus": bonus}}
	_, err := repository.db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repository *repositoryImpl) CreateOrder(ctx context.Context, order model.Order) error {
	_, err := repository.db.Collection("order").InsertOne(ctx, order)
	if mongo.IsDuplicateKeyError(err) {
		return exception.ErrAlreadyExists
	}

	if err != nil {
		return err
	}
	return nil
}

func (repository *repositoryImpl) GetAllMechanics(ctx context.Context) ([]model.Mechanic, error) {
	var mechanics []model.Mechanic
	cursor, err := repository.db.Collection("mechanics").Find(ctx, bson.M{})
	if err != nil {
		return mechanics, err
	}

	err = cursor.All(ctx, &mechanics)
	if err != nil {
		return mechanics, err
	}

	return mechanics, nil
}
