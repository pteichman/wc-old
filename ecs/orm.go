package ecs

import "errors"

var (
	ErrExists = errors.New("ecs: object already exists")
	ErrHasId  = errors.New("ecs: create() with nonzero id")
)
