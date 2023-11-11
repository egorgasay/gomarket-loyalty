package repository

import (
	"context"
	"gomarket-loyalty/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	SetUser(ctx context.Context, user model.User) error
	AddMechanic(ctx context.Context, bonus model.Mechanic) error
	UpdateBonusUser(ctx context.Context, id string, bonus int) error
	CreateOrder(ctx context.Context, order model.Order) error
	GetAllMechanics(ctx context.Context) ([]model.Mechanic, error)
	GetInfoOrders(ctx context.Context, clientID string) ([]model.Order, error)
}
