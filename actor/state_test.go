package actor

import (
	"fmt"
	"gitlab.com/anwski/crude-go-actors/com"
	"sync"
	"testing"
)

func TestState_AddBehaviour(t *testing.T) {
	client := com.NewMqttClient("mqtt://127.0.0.1:1883", false, 2)
	err := client.ConnectSync()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Connected")

	wg := sync.WaitGroup{}

	sum := 99

	grp := NewGroup("AddNumberGroup")

	bhv := NewBehaviourJson[int]("AddNumberBhv", func(act *Actor, msg com.Message[int]) {
		if act.GetState() == msg.Data {
			fmt.Println(msg.Data)
			data := msg.Data
			data++

			pongMsg := com.NewGroupMessage[int]("AddNumberBhv", grp.ID, &data)
			err := ActorSendMessageJson(act, pongMsg)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}

	})

	for i := 0; i < sum; i++ {
		actor := NewActor(client, "AddNumberActor")
		actor.Become(i)
		actor.JoinGroup(grp)

		err := actor.AddBehaviour(bhv)
		if err != nil {
			t.Fatal(err)
		}
	}

	num := 0
	initialMsg := com.NewGroupMessage[int](bhv.Name, grp.ID, &num)
	wg.Add(sum)

	err2 := initialMsg.SendJson(client)

	if err2 != nil {
		fmt.Println(err2)
	}
	wg.Wait()

	fmt.Println("Done!")

}

func TestSendMessageAndJson(t *testing.T) {
	mqttClient := com.NewMqttClient("mqtt://127.0.0.1:1883", true, 2)
	err := mqttClient.ConnectSync()
	if err != nil {
		t.Error(err)
	}
	sendingActor := NewActor(mqttClient, "SendingActor")
	type args struct {
		actor   *Actor
		topic   string
		message any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestStringSend", args: args{sendingActor, "test/state/send/string", "HelloWorld"}, wantErr: false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testMessage := com.NewDirectMessage(tt.args.topic, tt.args.actor.ID, &tt.args.message)
			if err := ActorSendMessage(tt.args.actor, testMessage); (err != nil) != tt.wantErr {
				t.Errorf("ActorSendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testMessage := com.NewDirectMessage(tt.args.topic, tt.args.actor.ID, &tt.args.message)
			if err := ActorSendMessageJson(tt.args.actor, testMessage); (err != nil) != tt.wantErr {
				t.Errorf("ActorSendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
