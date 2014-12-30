package wc

import (
	"math/rand"

	"yasty.org/peter/wc/ecs"
)

type taxFiller struct {
	min int
	max int
}

func (tf taxFiller) Fill(world *ecs.World, f []ecs.Entity, w, h int) {
	for _, s := range f {
		world.AddTag(s, taxTag, randRange(tf.min, tf.max))
	}
}

// randRange produces a random integer N in the range a <= N < b
func randRange(a, b int) int {
	return a + rand.Intn(b-a)
}
