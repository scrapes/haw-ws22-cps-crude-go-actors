package com

import (
	"github.com/google/uuid"
	"reflect"
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

func (msg *Message[T]) GetTopic(name string) string {
	var Topic string

	if msg.isGroup {
		Topic = "group/"
	} else {
		Topic = "actor/"
	}

	if name == "" {
		Topic += "by-id/"
		if msg.isGroup {
			Topic += msg.Group.String()
		} else {
			Topic += msg.Receiver.String()
		}
	} else {
		Topic += "by-name/"
		Topic += name
	}
	Topic += "/bhv/by-name/" + msg.BehaviourName

	return Topic
}

func (msg *Message[T]) Send(mqttClient *MqttClient) error {
	return mqttClient.Publish(msg.GetTopic(""), *msg)
}

func (msg *Message[T]) SendJson(mqttClient *MqttClient) error {
	return mqttClient.PublishJson(msg.GetTopic(""), *msg)
}

func (msg *Message[T]) SendAsName(mqttClient *MqttClient, name string) error {
	return mqttClient.Publish(msg.GetTopic(name), *msg)
}

func (msg *Message[T]) SendJsonAsName(mqttClient *MqttClient, name string) error {
	return mqttClient.PublishJson(msg.GetTopic(name), *msg)
}

func (msg ReflectedMessagePtr) ToMessage() ReflectedMessage {
	return ReflectedMessage(reflect.Indirect(reflect.Value(msg)))
}
