package actor

import (
	"gaffa/com"
	"gaffa/types"
	"github.com/google/uuid"
	"testing"
)

type TestStruct struct {
	str string
}

func TestBehaviour_Call(t *testing.T) {
	testStr := "Hello! 123 Test --"
	bhv := NewBehaviour[TestStruct]("TestTopic", func(self *State, message com.Message[TestStruct]) {
		if message.Data.str != testStr {
			t.Error("Mismatch in Test")
		}
	})

	test := TestStruct{str: testStr}
	msg := com.NewMessage("TestTopic", (types.ActorID)(uuid.New()), &test)
	bhv.Call(nil, msg.ToPtrValue())
}
