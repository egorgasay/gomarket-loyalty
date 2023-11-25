package service

import (
	"context"
	"gomarket-loyalty/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Service
type Service interface {
	Base() string
	Create(ctx context.Context, request model.RegisterRequest) error
	ValidateDataRegister(user model.RegisterRequest) error
	AddMechanic(ctx context.Context, bonus model.Mechanic) error
	AddBonus(mechanic model.Mechanic, item model.Item) int
	CreateOrder(ctx context.Context, clientID string, orderID string, order model.Items) error
	GetInfoOrders(ctx context.Context, clientID string) ([]model.Order, error)
}
