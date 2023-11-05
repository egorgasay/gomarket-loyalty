package service

import (
	"crypto/md5"
	"encoding/hex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/repository"
	"regexp"
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

func (service *serviceImpl) Register(request model.RegisterRequest) (string, error) {
	err := service.ValidateDataRegister(request)
	if err != nil {
		return "", err
	}

	cookie := uuid.New().String()
	request.Password = service.HashPassword(request.Password)

	user := model.User{
		Login:    request.Login,
		Password: request.Password,
		Cookie:   cookie,
	}

	err = service.repository.SetUser(user)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

func (service *serviceImpl) HashPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	hashString := hex.EncodeToString(hash.Sum(nil))
	return hashString
}

// ValidateDataRegister validates the user data for registration.
// It checks if the login is between 3 and 15 characters long,
// if the password is between 8 and 20 characters long,
// and if the password only contains alphanumeric characters.
// If any validation fails, it returns the ErrEnabledData error.
func (service *serviceImpl) ValidateDataRegister(user model.RegisterRequest) error {
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Login, validation.Required, validation.Length(3, 15)),
		validation.Field(&user.Password, validation.Required, validation.Length(8, 20), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]*$"))),
	)

	if err != nil {
		return exception.ErrEnabledData
	}
	return nil
}
