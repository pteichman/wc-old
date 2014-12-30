package ecs

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleWorld() {
	type key int
	const (
		nameKey key = iota
	)

	w := NewWorld()
	e := w.NewEntity()

	w.AddTag(e, nameKey, "Alice")
	fmt.Println(w.Tag(e, nameKey))

	// Output: Alice
}

func TestNewWorld(t *testing.T) {
	w := NewWorld()

	e1 := w.NewEntity()
	e2 := w.NewEntity()

	if e1 == e2 {
		t.Fatalf("Unexpected e1 == e2 (%v, %v)", e1, e2)
	}
}

func TestTag(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity()

	w.AddTag(e, 1, "foo")
	w.AddTag(e, 2, "bar")

	if !reflect.DeepEqual(w.Tag(e, 1), "foo") {
		t.Fatalf("Unexpected tag: %s (expected \"foo\")", w.Tag(e, 1))
	}
}
