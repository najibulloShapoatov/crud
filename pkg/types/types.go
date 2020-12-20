package types

import (
	"errors"
)

var (
	//ErrNotFound ...
	ErrNotFound = errors.New("item not found")
	//ErrInternal ...
	ErrInternal = errors.New("internal error")
	//ErrTokenNotFound ...
	ErrTokenNotFound = errors.New("token not found")
	//ErrNoSuchUser ...
	ErrNoSuchUser = errors.New("no such user")
	//ErrInvalidPassword ..
	ErrInvalidPassword = errors.New("invalid password")
	//ErrPhoneUsed ...
	ErrPhoneUsed = errors.New("phone alredy registered")
	//ErrTokenExpired ...
	ErrTokenExpired = errors.New("token expired")
)


