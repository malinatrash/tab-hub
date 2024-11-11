package myErrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrDBInsert          = errors.New("failed to insert user into database")
)
