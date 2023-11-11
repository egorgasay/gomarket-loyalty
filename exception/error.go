package exception

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrEnabledData   = errors.New("enabled data")
	ErrNotFound      = errors.New("not found")
)

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}

}
