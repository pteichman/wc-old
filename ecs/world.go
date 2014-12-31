package ecs

import "sync/atomic"

type World struct {
	tags       map[Entity]*Tag
	lastEntity int64
}

func NewWorld() *World {
	return &World{
		tags: make(map[Entity]*Tag),
	}
}

func (w *World) NewEntity() Entity {
	return Entity(atomic.AddInt64(&w.lastEntity, 1))
}

func (w *World) AddTag(e Entity, key interface{}, val interface{}) {
	w.tags[e] = &Tag{w.tags[e], key, val}
}

func (w *World) Tag(e Entity, key interface{}) interface{} {
	return w.tags[e].Value(key)
}

func (w *World) AllTags(e Entity) *Tag {
	return w.tags[e]
}

func (w *World) ForTagged(tag interface{}, f func(Entity, *Tag)) {
	for entity, tags := range w.tags {
		if tags.Value(tag) != nil {
			f(entity, tags)
		}
	}
}

type Tag struct {
	next     *Tag
	key, val interface{}
}

func (t *Tag) Value(key interface{}) interface{} {
	if t == nil {
		return nil
	}

	if t.key == key {
		return t.val
	}

	return t.next.Value(key)
}
