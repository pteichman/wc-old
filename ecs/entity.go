package ecs

import "fmt"

type Entity int64

var cur Entity

func New() Entity {
	cur++
	return Entity(cur)
}

func (e Entity) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%x\"", e)
	return []byte(s), nil
}
