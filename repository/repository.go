package repository

import "gomarket-loyalty/model"

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	SetUser(user model.User) error
	AddMechanic(bonus model.Mechanic) error
	UpdateBonusUser(id string, bonus int) error
	CreateOrder(order model.Order) error
	GetBonus(id int) (model.Mechanic, error)
}
