package service

import (
	"gomarket-loyalty/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Service
type Service interface {
	Base() string
	Create(request model.RegisterRequest) error
	ValidateDataRegister(user model.RegisterRequest) error
	AddMechanic(bonus model.Mechanic) error
}
