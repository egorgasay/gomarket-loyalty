package exception

import "errors"

var (
	ErrLoginAlreadyExists = errors.New("login is already exists")
	ErrEnabledData        = errors.New("enabled data")
)

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}

}
