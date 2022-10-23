package actor

import (
	"fmt"
	"gaffa/com"
	"sync"
	"testing"
)

func TestState_AddBehaviour(t *testing.T) {
	client := com.NewMqttClient("mqtt://127.0.0.1:1883")
	err := client.ConnectSync()
	if err != nil {
		return
	}
	fmt.Println("Connected")

	wg := sync.WaitGroup{}

	topic := "/test/actor/sub/"
	sum := 10
	add := 10

	for i := 0; i < sum*add; i++ {
		actor := NewActor(client, "TestActor")
		err := actor.AddBehaviour(NewBehaviourJson[int](fmt.Sprintf(topic+"%d", i), func(act *State, msg com.Message[int]) {
			msg.Data += add
			msg.Topic = fmt.Sprintf(topic+"%d", msg.Data)
			err := act.mqttClient.PublishJson(msg.Topic, msg)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}))
		if err != nil {
			t.Fatal(err)
		}
	}

	for {
		wg.Add(sum * add)
		fmt.Println("Starting round...")
		for i := 0; i < add; i++ {
			tpc := fmt.Sprintf(topic+"%d", i)
			toSend := com.Message[int]{
				Data:  i,
				Topic: tpc,
			}
			err2 := client.PublishJson(toSend.Topic, toSend)
			if err2 != nil {
				t.Fatal(err2)
			}
		}

		wg.Wait()

		fmt.Println("Done!")
	}

}

func TestSendMessageAndJson(t *testing.T) {
	mqttClient := com.NewMqttClient("mqtt://127.0.0.1:1883")
	sendingActor := NewActor(mqttClient, "SendingActor")
	type args struct {
		actor   *State
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
			testMessage := com.NewMessage(tt.args.topic, tt.args.actor.ID, &tt.args.message)
			if err := SendMessage(tt.args.actor, testMessage); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testMessage := com.NewMessage(tt.args.topic, tt.args.actor.ID, &tt.args.message)
			if err := SendMessageJson(tt.args.actor, testMessage); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
