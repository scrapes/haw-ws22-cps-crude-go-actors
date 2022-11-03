package com

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"log"
	url2 "net/url"
	"os"
	"reflect"
	"runtime"
	"sync"
)

type SubCallback struct {
	Callback func(msg reflect.Value)
	ID       uuid.UUID
}

type MqttClient struct {
	client         mqtt.Client
	verbose        bool
	qos            byte
	serializer     *Serializer
	mutexCallbacks sync.Mutex
	callbacks      map[string][]SubCallback
	types          map[string]reflect.Type
	publishMutex   sync.Mutex
	clientUUID     uuid.UUID
}

func mqttWait(t mqtt.Token) error {
	_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
	if t.Error() != nil {
		return t.Error()
	}
	return nil
}

func MqttConfigure(verbose bool) {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)

	if verbose {
		mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
		mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
		mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	}

}

func NewMqttClient(mqttHost string, verbose bool) *MqttClient {
	MqttConfigure(verbose)
	mqttUrl, err := url2.Parse(mqttHost)
	hosts := []*url2.URL{mqttUrl}

	clientID := uuid.New()

	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return nil
		}
		os.Exit(1)
	}

	options := mqtt.ClientOptions{
		Servers:  hosts,
		ClientID: clientID.String(),
	}

	options.SetOrderMatters(true)
	options.SetCleanSession(true)
	client := mqtt.NewClient(&options)
	ret := MqttClient{
		client:     client,
		verbose:    verbose,
		qos:        1,
		serializer: NewSerializer(),
		callbacks:  make(map[string][]SubCallback),
		types:      make(map[string]reflect.Type),
		clientUUID: clientID,
	}
	return &ret
}

func (_client *MqttClient) Unsubscribe(topic string, id uuid.UUID) {
	_client.mutexCallbacks.Lock()
	defer _client.mutexCallbacks.Unlock()

	index := -1
	for i, callback := range _client.callbacks[topic] {
		if callback.ID == id {
			index = i
			break
		}
	}
	if index > 0 {
		_client.callbacks[topic] = append(_client.callbacks[topic][:index], _client.callbacks[topic][index+1:]...)
		if len(_client.callbacks[topic]) == 0 {
			_client.types[topic] = nil
		}
	}
}

func (_client *MqttClient) Subscribe(topic string, typ reflect.Type, callback SubCallback) error {
	return _client.subscribe(false, topic, typ, callback)
}

func (_client *MqttClient) SubscribeJson(topic string, typ reflect.Type, callback SubCallback) error {
	return _client.subscribe(true, topic, typ, callback)
}

func (_client *MqttClient) subscribe(json bool, topic string, typ reflect.Type, callback SubCallback) error {
	_client.mutexCallbacks.Lock()
	defer _client.mutexCallbacks.Unlock()

	if msgType, ok := _client.types[topic]; ok {
		if msgType != typ {
			return errors.New("channel Types do not match")
		}
		//do something here
	} else {

		_client.types[topic] = typ
	}
	_client.callbacks[topic] = append(_client.callbacks[topic], callback)

	isJson := json
	return mqttWait(_client.client.Subscribe(topic, _client.qos, func(client mqtt.Client, message mqtt.Message) {
		var err error
		var objectPointer reflect.Value
		if isJson {
			objectPointer, err = _client.serializer.DecodeJson(_client.types[topic], message.Payload())
		} else {
			objectPointer, err = _client.serializer.Decode(_client.types[topic], message.Payload())
		}

		if err != nil {
			fmt.Println(err)
			fmt.Println(topic)
			runtime.Breakpoint()
			return
		}

		_client.mutexCallbacks.Lock()
		for _, handler := range _client.callbacks[topic] {
			go handler.Callback(objectPointer)
		}
		_client.mutexCallbacks.Unlock()
	}))
}

func (_client *MqttClient) ConnectSync() error {
	return mqttWait(_client.Connect())
}

func (_client *MqttClient) Connect() mqtt.Token {
	return _client.client.Connect()
}

func (_client *MqttClient) GetSerializer() *Serializer {
	return _client.serializer
}

func (_client *MqttClient) PublishValue(topic string, obj reflect.Value) error {
	_client.publishMutex.Lock()

	buffer, err := _client.serializer.EncodeValue(obj)
	if err != nil {
		return err
	}
	if _client.verbose {
		fmt.Println(buffer)
	}
	token := _client.client.Publish(topic, _client.qos, false, buffer)
	_client.publishMutex.Unlock()
	return mqttWait(token)
}

func (_client *MqttClient) Publish(topic string, obj any) error {
	_client.publishMutex.Lock()

	buffer, err := _client.serializer.Encode(obj)
	if err != nil {
		return err
	}
	if _client.verbose {
		fmt.Println(buffer)
	}
	token := _client.client.Publish(topic, _client.qos, false, buffer)
	_client.publishMutex.Unlock()
	return mqttWait(token)
}

func (_client *MqttClient) PublishValueJson(topic string, obj reflect.Value) error {
	return _client.PublishJson(topic, obj.Interface())
}

func (_client *MqttClient) PublishJson(topic string, obj any) error {
	_client.publishMutex.Lock()

	buffer, err := _client.serializer.EncodeJson(obj)
	if err != nil {
		return err
	}
	if _client.verbose {
		fmt.Println(string(buffer))
	}
	token := _client.client.Publish(topic, _client.qos, false, buffer)
	_client.publishMutex.Unlock()
	return mqttWait(token)
}
