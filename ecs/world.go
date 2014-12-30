package ecs

import "sync/atomic"

type World struct {
	tags       map[Entity]*tag
	lastEntity int64
}

func NewWorld() *World {
	return &World{
		tags: make(map[Entity]*tag),
	}
}

func (w *World) NewEntity() Entity {
	return Entity(atomic.AddInt64(&w.lastEntity, 1))
}

func (w *World) AddTag(e Entity, key interface{}, val interface{}) {
	w.tags[e] = &tag{w.tags[e], key, val}
}

func (w *World) Tag(e Entity, key interface{}) interface{} {
	return w.tags[e].value(key)
}

type tag struct {
	next     *tag
	key, val interface{}
}

func (t *tag) value(key interface{}) interface{} {
	if t == nil {
		return nil
	}

	if t.key == key {
		return t.val
	}

	return t.next.value(key)
}
