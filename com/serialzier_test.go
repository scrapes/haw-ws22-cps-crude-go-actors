package com

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestSerializer_Encode(t *testing.T) {
	type test struct {
		Name string
		Num  int64
	}
	typeOf := reflect.TypeOf(test{})
	sr := NewSerializer()
	ch := make(chan []byte, 100)
	var counter int64 = 0

	fnc := reflect.ValueOf(func(t test) (string, int64) {
		return t.Name, t.Num
	})

	go func() {
		var cnt int64 = 0
		for {
			buff := <-ch
			obj, err := sr.DecodeValue(buff, typeOf)
			if err != nil {
				fmt.Println(err)
			}

			values := fnc.Call([]reflect.Value{obj})
			name := values[0].String()
			num := values[1].Int()

			if num != cnt {
				fmt.Println("err")
				fmt.Println(name)
				runtime.Breakpoint()
			}

			cnt++
		}
	}()
	for i := 0; i < 10000000; i++ {
		str := test{
			Num:  counter,
			Name: fmt.Sprintf("/ja/lol/ey/%d", counter),
		}

		encode, err := sr.Encode(str)
		if err != nil {
			runtime.Breakpoint()
			return
		}

		buffer := make([]byte, len(encode))
		copy(buffer, encode)
		ch <- buffer

		counter++
	}
}
func TestSerializer_Decode(t *testing.T) {
	type test struct {
		Name string
		Num  int64
	}
	sr := NewSerializer()
	testObject := test{Num: 0xDeadbeef, Name: "DeadGnu"}
	buffer, err := sr.Encode(testObject)
	if err != nil {
		t.Error(err)
	}

	message, err2 := sr.Decode(reflect.TypeOf(testObject), buffer)

	if err2 != nil {
		t.Error(err2)
	}

	name := message.Interface().(test).Name
	num := message.Interface().(test).Num

	if testObject.Name != name || testObject.Num != num {
		t.Error("Mismatch of Test Object")
	}
}
func TestSerializer_DecodeJson(t *testing.T) {
	type test struct {
		Num  int
		Name string
	}

	obj := test{Num: 12345, Name: "this is a name"}
	sr := NewSerializer()
	json, err := sr.EncodeJson(obj)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(json))

	ptr, err2 := sr.DecodeJson(reflect.TypeOf(obj), json)

	if err2 != nil {
		t.Fatal(err2)
	}
	fmt.Println(ptr)
}
