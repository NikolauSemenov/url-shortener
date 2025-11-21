package errorsApp

import "errors"

var (
	ErrUrlAlreadyExists = errors.New("url already exists")
	ErrCacheMiss        = errors.New("cache: key not found")
)
