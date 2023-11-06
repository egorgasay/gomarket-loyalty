package service

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (service *serviceImpl) Register(request model.RegisterRequest) error {
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

//func (service *serviceImpl) HashPassword(password string) string {
//	hash := md5.New()
//	hash.Write([]byte(password))
//	hashString := hex.EncodeToString(hash.Sum(nil))
//	return hashString
//}

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
