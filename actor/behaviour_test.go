package actor

import (
	"github.com/google/uuid"
	"github.com/scrapes/haw-ws22-cps-crude-go-actors/com"
	"testing"
)

type TestStruct struct {
	str string
}

func TestBehaviour_Call(t *testing.T) {
	testStr := "Hello! 123 Test --"
	bhv := NewBehaviour[TestStruct]("TestTopic", func(self *Actor, message com.Message[TestStruct]) {
		if message.Data.str != testStr {
			t.Error("Mismatch in Test")
		}
	})

	test := TestStruct{str: testStr}
	msg := com.NewDirectMessage("TestTopic", uuid.New(), &test)
	bhv.Call(nil, msg.ToPtrValue())
}
