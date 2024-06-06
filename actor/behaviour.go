package actor

import (
	"github.com/google/uuid"
	"gitlab.com/anwski/crude-go-actors/com"
	"reflect"
	"strings"
)

type Callback[T any] func(self *Actor, message com.Message[T])

type Behaviour struct {
	callback reflect.Value
	Name     string
	typ      reflect.Type
	id       uuid.UUID
	JSON     bool
}

func NewBehaviour[T any](Name string, callback Callback[T]) *Behaviour {
	return _NewBehaviour[T](Name, callback, false)
}

func NewBehaviourJson[T any](Name string, callback Callback[T]) *Behaviour {
	return _NewBehaviour[T](Name, callback, true)
}

func _NewBehaviour[T any](Name string, callback Callback[T], json bool) *Behaviour {
	genType := com.Message[T]{}
	bhv := Behaviour{
		callback: reflect.ValueOf(callback),
		Name:     strings.Replace(Name, " ", "_", -1), //replace whitespace to prevent mqtt errors
		typ:      reflect.TypeOf(genType),
		id:       uuid.New(),
		JSON:     json,
	}
	return &bhv
}

func (bhv *Behaviour) Call(self *Actor, messagePtr reflect.Value) {
	// dereference pointer
	message := reflect.Indirect(messagePtr)

	println(self.Type, "Calling behaviour")
	self.lock.Lock()
	bhv.callback.Call([]reflect.Value{
		reflect.ValueOf(self),
		message,
	})
	self.lock.Unlock()
	println(self.Type, "Done calling behaviour")
}

func (bhv *Behaviour) GetName() string {
	return bhv.Name
}

func (bhv *Behaviour) GetTyp() reflect.Type {
	return bhv.typ
}

func (bhv *Behaviour) GetID() uuid.UUID {
	return bhv.id
}
