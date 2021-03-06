package ecs

import "errors"

var (
	ErrExists   = errors.New("ecs: object already exists")
	ErrHasID    = errors.New("ecs: create() with nonzero id")
	ErrNotFound = errors.New("ecs: object not found")
)
