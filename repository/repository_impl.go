package repository

import (
	"go.mongodb.org/mongo-driver/bson"
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
		return exception.ErrAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (repository *repositoryImpl) AddMechanic(bonus model.Mechanic) error {
	ctx, cancel := config.NewMongoContext()
	defer cancel()
	_, err := repository.db.Collection("mechanics").InsertOne(ctx, bonus)
	if mongo.IsDuplicateKeyError(err) {
		return exception.ErrAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil

}

func (repository *repositoryImpl) UpdateBonusUser(id string, bonus int) error {
	ctx, cancel := config.NewMongoContext()
	defer cancel()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"bonus": bonus}}
	_, err := repository.db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repository *repositoryImpl) CreateOrder(order model.Order) error {
	ctx, cancel := config.NewMongoContext()
	defer cancel()
	_, err := repository.db.Collection("order").InsertOne(ctx, order)
	if mongo.IsDuplicateKeyError(err) {
		return exception.ErrAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (repository *repositoryImpl) GetBonus(id int) (model.Mechanic, error) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()
	var mechanic model.Mechanic
	err := repository.db.Collection("mechanics").FindOne(ctx, bson.M{"_id": id}).Decode(&mechanic)
	if err != nil {
		return mechanic, err
	}
	return mechanic, nil
}
