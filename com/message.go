package com

import (
	"gitlab.com/anwski/crude-go-actors/types"
	"reflect"
)

type ReflectedMessage reflect.Value
type ReflectedMessagePtr reflect.Value

type Message[T any] struct {
	Topic  string
	Sender types.ActorID
	Data   T
}

func NewMessage[T any](topic string, sender types.ActorID, data *T) Message[T] {
	return Message[T]{
		Topic:  topic,
		Sender: sender,
		Data:   *data,
	}
}

func (msg Message[T]) ToValue() reflect.Value {
	return reflect.ValueOf(msg)
}

func (msg *Message[T]) ToPtrValue() reflect.Value {
	return reflect.ValueOf(msg)
}

func (msg ReflectedMessagePtr) ToMessage() ReflectedMessage {
	return ReflectedMessage(reflect.Indirect(reflect.Value(msg)))
}
