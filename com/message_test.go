package com

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewMessage(t *testing.T) {
	bhvName := "BehaviourName"
	actorID := uuid.New()
	testString := "DeadBeef"

	message := NewDirectMessage[string](bhvName, actorID, &testString)
	if message.BehaviourName != bhvName || message.Sender != actorID || message.Data != testString {
		t.Error("Mismatched Data")
	}
}
