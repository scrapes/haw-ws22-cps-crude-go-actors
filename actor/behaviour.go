package actor

import (
	"gaffa/com"
	"github.com/google/uuid"
	"reflect"
)

type Callback[T any] func(self *State, message com.Message[T])

type Behaviour struct {
	callback reflect.Value
	topic    string
	typ      reflect.Type
	id       uuid.UUID
	JSON     bool
}

func NewBehaviour[T any](topic string, callback Callback[T]) *Behaviour {
	return _NewBehaviour[T](topic, callback, false)
}

func NewBehaviourJson[T any](topic string, callback Callback[T]) *Behaviour {
	return _NewBehaviour[T](topic, callback, true)
}

func _NewBehaviour[T any](topic string, callback Callback[T], json bool) *Behaviour {
	genType := com.Message[T]{}
	bhv := Behaviour{
		callback: reflect.ValueOf(callback),
		topic:    topic,
		typ:      reflect.TypeOf(genType),
		id:       uuid.New(),
		JSON:     json,
	}
	return &bhv
}

func (bhv *Behaviour) Call(self *State, messagePtr reflect.Value) {
	// dereference pointer
	message := reflect.Indirect(messagePtr)

	bhv.callback.Call([]reflect.Value{
		reflect.ValueOf(self),
		message,
	})
}

func (bhv *Behaviour) GetTopic() string {
	return bhv.topic
}

func (bhv *Behaviour) GetTyp() reflect.Type {
	return bhv.typ
}

func (bhv *Behaviour) GetID() uuid.UUID {
	return bhv.id
}
