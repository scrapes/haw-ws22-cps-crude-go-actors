package com

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

type MessageDecoder func(serializer *Serializer, buffer []byte) any
type MessageEncoder func(serializer *Serializer, obj any) ([]byte, error)

type Serializer struct {
	Encoder *gob.Encoder
	eBuffer bytes.Buffer
	eMutex  sync.Mutex
	Decoder *gob.Decoder
	dBuffer bytes.Buffer
	dMutex  sync.Mutex
}

func NewSerializer() *Serializer {
	serializer := Serializer{}
	serializer.Init()
	return &serializer
}

func (obj *Serializer) Init() {
	obj.Encoder = gob.NewEncoder(&obj.eBuffer)
	obj.Decoder = gob.NewDecoder(&obj.dBuffer)
}

func (obj *Serializer) EncodeJson(msg any) ([]byte, error) {
	return json.Marshal(msg)
}

func (obj *Serializer) Encode(msg any) ([]byte, error) {
	if reflect.ValueOf(msg).Kind() == reflect.Ptr {
		return nil, errors.New("message cannot be pointer")
	}

	obj.eMutex.Lock()
	defer obj.eMutex.Unlock()

	obj.eBuffer.Reset()
	err := obj.Encoder.Encode(msg)
	if err != nil {
		return nil, err
	}
	b := obj.eBuffer.Bytes()
	buffer := make([]byte, len(b))
	copy(buffer, b)
	return buffer, nil
}

func (obj *Serializer) EncodeValue(msg reflect.Value) ([]byte, error) {
	if reflect.ValueOf(msg).Kind() == reflect.Ptr {
		return nil, errors.New("message cannot be pointer")
	}

	obj.eMutex.Lock()
	defer obj.eMutex.Unlock()

	obj.eBuffer.Reset()
	err := obj.Encoder.EncodeValue(msg)
	if err != nil {
		return nil, err
	}
	b := obj.eBuffer.Bytes()
	buffer := make([]byte, len(b))
	copy(buffer, b)
	return buffer, nil
}

func (obj *Serializer) DecodeBufferValue(msg []byte, decoded reflect.Value) error {
	obj.dMutex.Lock()
	defer obj.dMutex.Unlock()

	obj.dBuffer.Reset()
	obj.dBuffer.Write(msg)

	err := obj.Decoder.DecodeValue(decoded)
	if err != nil {
		return err
	}

	return nil
}

func (obj *Serializer) DecodeValue(msg []byte, typ reflect.Type) (reflect.Value, error) {
	obj.dMutex.Lock()
	defer obj.dMutex.Unlock()

	buffer := reflect.New(typ)
	obj.dBuffer.Reset()
	obj.dBuffer.Write(msg)

	err := obj.Decoder.DecodeValue(buffer)

	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.Indirect(buffer), nil
}

func (obj *Serializer) Decode(typ reflect.Type, msg []byte) (reflect.Value, error) {
	buff := reflect.New(typ)
	err := obj.DecodeBufferValue(msg, buff)
	if err != nil {
		return reflect.Value{}, err
	}
	return buff, nil
}

func (obj *Serializer) DecodeJson(typ reflect.Type, msg []byte) (reflect.Value, error) {
	buff := reflect.New(typ)
	err := json.Unmarshal(msg, buff.Interface())
	if err != nil {
		return reflect.Value{}, err
	}
	return buff, nil
}
