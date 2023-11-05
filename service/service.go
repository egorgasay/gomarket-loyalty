package service

import (
	"gomarket-loyalty/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Service
type Service interface {
	Base() string
	Register(request model.RegisterRequest) (string, error)

	HashPassword(password string) string
	ValidateDataRegister(user model.RegisterRequest) error
}
