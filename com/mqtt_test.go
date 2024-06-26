package com

import (
	"fmt"
	"github.com/google/uuid"
	crude_go_actors "github.com/scrapes/haw-ws22-cps-crude-go-actors"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestMqttClient_Subscribe(t *testing.T) {
	client := NewMqttClient("mqtt://127.0.0.1:1883", true, 2)
	err := client.ConnectSync()
	if err != nil {
		t.Error(err)
	}
	crude_go_actors.Logger.Info("Connected")
	sum := 1

	wg := sync.WaitGroup{}
	wg.Add(sum)

	topic := "/test/mqtt/sub/%d"

	msg := Message[string]{
		BehaviourName: topic,
		Data:          "lil",
		Sender:        uuid.New(),
	}
	for i := 0; i < sum; i++ {
		tpc := fmt.Sprintf(topic, i)
		uid := uuid.New()
		err22 := client.SubscribeJson(tpc, reflect.TypeOf(msg), SubCallback{Callback: func(msgPtr reflect.Value) {
			check := tpc
			msg := reflect.Indirect(msgPtr)
			reflect.ValueOf(func(msg Message[string]) {
				if check != msg.Data {
					runtime.Breakpoint()
				}
				wg.Done()
				client.Unsubscribe(tpc, uid)
			}).Call([]reflect.Value{msg})
		}, ID: uid})

		if err22 != nil {
			t.Error(err22)
		}
	}

	for i := 0; i < sum; i++ {
		tpc := fmt.Sprintf(topic, i)
		msg.Data = tpc
		msg.BehaviourName = tpc
		err3 := client.PublishJson(tpc, msg)
		if err3 != nil {
			t.Error(err3)
		}
	}
	wg.Wait()
}
