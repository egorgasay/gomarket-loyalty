package repository

import "gomarket-loyalty/model"

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name=Repository
type Repository interface {
	SetUser(user model.User) error
}
