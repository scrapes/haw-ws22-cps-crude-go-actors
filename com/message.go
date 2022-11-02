package com

import (
	"github.com/google/uuid"
	"reflect"
	"strings"
)

type ReflectedMessage reflect.Value
type ReflectedMessagePtr reflect.Value

type Message[T any] struct {
	isGroup       bool
	BehaviourName string
	Group         uuid.UUID
	Sender        uuid.UUID
	Receiver      uuid.UUID
	Data          T
}

func NewDirectMessage[T any](BehaviourName string, receiver uuid.UUID, data *T) Message[T] {
	return Message[T]{
		isGroup:       false,
		BehaviourName: BehaviourName,
		Receiver:      receiver,
		Group:         uuid.Nil,
		Sender:        uuid.Nil,
		Data:          *data,
	}
}

func NewGroupMessage[T any](BehaviourName string, group uuid.UUID, data *T) Message[T] {
	return Message[T]{
		isGroup:       true,
		BehaviourName: BehaviourName,
		Group:         group,
		Receiver:      uuid.Nil,
		Sender:        uuid.Nil,
		Data:          *data,
	}
}
func (msg *Message[T]) SetSender(sender uuid.UUID) {
	msg.Sender = sender
}

func (msg Message[T]) ToValue() reflect.Value {
	return reflect.ValueOf(msg)
}

func (msg *Message[T]) ToPtrValue() reflect.Value {
	return reflect.ValueOf(msg)
}

func (msg *Message[T]) GetTopic() string {
	str := ""
	if msg.isGroup {
		str = "group/by-id/" + msg.Group.String() + "/bhv/by-name/" + msg.BehaviourName
	} else {
		str = "actor/by-id/" + msg.Receiver.String() + "/bhv/by-name/" + msg.BehaviourName
	}
	return strings.ToValidUTF8(str, "--")
}

func (msg *Message[T]) Send(mqttClient *MqttClient) error {
	return mqttClient.Publish(msg.GetTopic(), *msg)
}

func (msg *Message[T]) SendJson(mqttClient *MqttClient) error {
	return mqttClient.PublishJson(msg.GetTopic(), *msg)
}

func (msg ReflectedMessagePtr) ToMessage() ReflectedMessage {
	return ReflectedMessage(reflect.Indirect(reflect.Value(msg)))
}
