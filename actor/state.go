package actor

import (
	"gaffa/com"
	"gaffa/types"
	"github.com/google/uuid"
	"reflect"
)

type State struct {
	ID         types.ActorID
	Type       string
	mqttClient *com.MqttClient
	behaviours map[string]*Behaviour
}

func NewActor(com *com.MqttClient, typ string) *State {
	act := State{
		ID:         (types.ActorID)(uuid.New()),
		Type:       typ,
		mqttClient: com,
		behaviours: make(map[string]*Behaviour),
	}

	return &act
}

func (actor *State) AddBehaviour(bhv *Behaviour) error {
	actor.behaviours[bhv.GetTopic()] = bhv
	if bhv.JSON {
		return actor.mqttClient.SubscribeJson(bhv.GetTopic(), bhv.GetTyp(), actor.GetSubCallback(bhv))
	} else {
		return actor.mqttClient.Subscribe(bhv.GetTopic(), bhv.GetTyp(), actor.GetSubCallback(bhv))
	}
}
func (actor *State) GetCallback(bhv *Behaviour) func(msg reflect.Value) {
	return func(msg reflect.Value) {
		bhv.Call(actor, msg)
	}
}
func (actor *State) GetSubCallback(bhv *Behaviour) com.SubCallback {
	return com.SubCallback{
		Callback: actor.GetCallback(bhv),
		ID:       bhv.GetID(),
	}
}

func (actor *State) RemoveBehaviour(topic string) {
	actor.mqttClient.Unsubscribe(topic, actor.behaviours[topic].GetID())
	actor.behaviours[topic] = nil
}

func SendMessage[T any](actor *State, message com.Message[T]) error {
	return actor.mqttClient.Publish(message.Topic, message)
}

func SendMessageJson[T any](actor *State, message com.Message[T]) error {
	return actor.mqttClient.PublishJson(message.Topic, message)
}
