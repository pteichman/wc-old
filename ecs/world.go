package ecs

type World struct {
	Entities map[Entity]Entity
	Tags     map[TagType]interface{}
}

func (w *World) AddEntity(e Entity) {
	w.Entities[e] = e
}

func (w *World) AddTag(e Entity, tt TagType, t interface{}) {
}

func (w *World) Tag(e Entity, tt TagType) interface{} {
	return nil
}
