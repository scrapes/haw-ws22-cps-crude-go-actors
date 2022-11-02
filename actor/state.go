package actor

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/anwski/crude-go-actors/com"
	"reflect"
	"strings"
)

type Actor struct {
	ID         uuid.UUID
	Type       string
	mqttClient *com.MqttClient
	behaviours map[uuid.UUID]*Behaviour
	State      any
	groups     map[uuid.UUID]*Group
}

func NewActor(com *com.MqttClient, typ string) *Actor {
	act := Actor{
		ID:         uuid.New(),
		Type:       typ,
		mqttClient: com,
		behaviours: make(map[uuid.UUID]*Behaviour),
		groups:     make(map[uuid.UUID]*Group),
	}

	return &act
}

func (actor *Actor) JoinGroup(grp *Group) {
	actor.groups[grp.ID] = grp
	for _, bhv := range actor.behaviours {
		topic := grp.GetTopic(bhv.GetName())
		callback := actor.GetSubCallback(bhv)
		if bhv.JSON {
			err := actor.mqttClient.SubscribeJson(topic, bhv.GetTyp(), callback)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := actor.mqttClient.Subscribe(topic, bhv.GetTyp(), callback)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (actor *Actor) LeaveGroup(grp *Group) {
	actor.groups[grp.ID] = nil
	for _, behaviour := range actor.behaviours {
		topic := grp.GetTopic(behaviour.GetName())
		actor.mqttClient.Unsubscribe(topic, behaviour.GetID())
	}
}

func (actor *Actor) GetGroup(name string) *Group {
	for _, group := range actor.groups {
		if group.Name == name {
			return group
		}
	}
	return nil
}

func (actor *Actor) GetTopic(name string) string {
	str := "actor/by-id/" + actor.ID.String() + "/bhv/by-name/" + name
	return strings.ToValidUTF8(str, "--")
}
func (actor *Actor) AddBehaviour(bhv *Behaviour) error {
	if actor.behaviours[bhv.GetID()] != nil {
		return fmt.Errorf("behaviour already added, please remove first")
	}

	actor.behaviours[bhv.GetID()] = bhv
	topicByname := actor.GetTopic(bhv.GetName())
	callback := actor.GetSubCallback(bhv)
	/*topic_byID := "/actor/by-id/" + actor.ID.String() + "/bhv/by-ID/" + bhv.GetID().String()*/
	if bhv.JSON {
		for _, group := range actor.groups {
			groupTopicname := group.GetTopic(bhv.GetName())
			err := actor.mqttClient.SubscribeJson(groupTopicname, bhv.GetTyp(), callback)
			if err != nil {
				fmt.Println(err)
			}
		}
		return actor.mqttClient.SubscribeJson(topicByname, bhv.GetTyp(), callback)
	} else {
		for _, group := range actor.groups {
			groupTopicname := group.GetTopic(bhv.GetName())
			err := actor.mqttClient.Subscribe(groupTopicname, bhv.GetTyp(), callback)
			if err != nil {
				fmt.Println(err)
			}
		}
		return actor.mqttClient.Subscribe(topicByname, bhv.GetTyp(), callback)
	}
}

func (actor *Actor) GetCallback(bhv *Behaviour) func(msg reflect.Value) {
	return func(msg reflect.Value) {
		bhv.Call(actor, msg)
	}
}
func (actor *Actor) GetSubCallback(bhv *Behaviour) com.SubCallback {
	return com.SubCallback{
		Callback: actor.GetCallback(bhv),
		ID:       bhv.GetID(),
	}
}

func (actor *Actor) RemoveBehaviour(bhv *Behaviour) {
	topic := actor.GetTopic(bhv.GetName())
	actor.mqttClient.Unsubscribe(topic, bhv.GetID())
	actor.behaviours[bhv.GetID()] = nil
}

func (actor *Actor) Become(state any) {
	actor.State = state
}

func (actor *Actor) UnBecome() {
	actor.State = ""
}

func (actor *Actor) GetState() any {
	return actor.State
}

func (actor *Actor) GetMqttClient() *com.MqttClient {
	return actor.mqttClient
}

func ActorSendMessage[T any](actor *Actor, message com.Message[T]) error {
	message.SetSender(actor.ID)
	return message.Send(actor.mqttClient)
}

func ActorSendMessageJson[T any](actor *Actor, message com.Message[T]) error {
	message.SetSender(actor.ID)
	return message.SendJson(actor.mqttClient)
}
