package myErrors

import "errors"

var (
	ErrPermissonAlreadyExists = errors.New("permisson already exists")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrDBInsert               = errors.New("failed to insert user into database")
)
