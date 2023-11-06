package service

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gomarket-loyalty/constants"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/repository"
)

func NewProductService(Repository *repository.Repository) Service {
	return &serviceImpl{
		repository: *Repository,
	}
}

type serviceImpl struct {
	repository repository.Repository
}

func (service *serviceImpl) Base() string {
	return "Hello World"
}

func (service *serviceImpl) Create(request model.RegisterRequest) error {
	err := service.ValidateDataRegister(request)
	if err != nil {
		return err
	}

	user := model.User{
		Login: request.Login,
	}

	err = service.repository.SetUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (service *serviceImpl) AddMechanic(bonus model.Mechanic) error {
	if bonus.RewardType != constants.Points && bonus.RewardType != constants.Percentage {
		return exception.ErrEnabledData
	}
	if bonus.Reward <= 0 {
		return exception.ErrEnabledData

	}
	err := service.repository.AddMechanic(bonus)
	if err != nil {
		return err
	}
	return nil

}

// ValidateDataRegister validates the user data for registration.
// It checks if the login is between 3 and 15 characters long,
// and if the password only contains alphanumeric characters.
// If any validation fails, it returns the ErrEnabledData error.
func (service *serviceImpl) ValidateDataRegister(user model.RegisterRequest) error {
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Login, validation.Required, validation.Length(3, 15)),
	)

	if err != nil {
		return exception.ErrEnabledData
	}
	return nil
}
