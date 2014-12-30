package ecs

import "fmt"

type Entity int64

func (e Entity) String() string {
	return fmt.Sprintf("%x", int64(e))
}

func (e Entity) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%x\"", int64(e))
	return []byte(s), nil
}
