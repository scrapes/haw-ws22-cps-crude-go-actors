package com

import (
	"github.com/google/uuid"
	"gitlab.com/anwski/crude-go-actors/types"
	"testing"
)

func TestNewMessage(t *testing.T) {
	topic := "TestTopic"
	actorID := (types.ActorID)(uuid.New())
	testString := "DeadBeef"

	message := NewMessage[string](topic, actorID, &testString)
	if message.Topic != topic || message.Sender != actorID || message.Data != testString {
		t.Error("Mismatched Data")
	}

}
