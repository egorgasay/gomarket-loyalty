package exception

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrEnabledData   = errors.New("enabled data")
)

func PanicIfNeeded(err interface{}) {
	if err != nil {
		panic(err)
	}

}
