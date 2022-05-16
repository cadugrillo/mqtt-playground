package Mqttbuffer

import (
	"errors"
	"fmt"
)

type Mqttbuffer struct {
	buffer       [5]Message
	readPointer  int
	writePointer int
}

// Message
type Message struct {
	Duplicate bool
	Qos       byte
	Retained  bool
	Topic     string
	MessageID uint16
	Payload   string
	Ack       bool
}

func NewMqttbuffer() Mqttbuffer {
	b := Mqttbuffer{}
	return b
}

func (b Mqttbuffer) GetReadPointer() int {
	return b.readPointer
}

func (b Mqttbuffer) GetWritePointer() int {
	return b.writePointer
}

func (b Mqttbuffer) AddMessage(message Message) Mqttbuffer {
	if b.writePointer == len(b.buffer)-1 {
		b.buffer[b.writePointer] = message
		b.writePointer = 0
		return b
	}
	b.buffer[b.writePointer] = message
	b.writePointer++
	return b
}

func (b Mqttbuffer) ReadMessage(index int) (Message, error) {
	if index < len(b.buffer) {
		return b.buffer[index], nil
	}
	msg := Message{}
	return msg, errors.New(fmt.Sprintf("Index %d greater then buffer size [%d]", index, len(b.buffer)))
}

func (b Mqttbuffer) NextMessage() Mqttbuffer {
	if b.readPointer == len(b.buffer)-1 {
		b.readPointer = 0
		return b
	}
	if b.readPointer != b.writePointer {
		b.readPointer++
		return b
	}
	fmt.Println("No new messages on the buffer")
	return b
}

func (b Mqttbuffer) NewMessage() bool {
	return b.writePointer != b.readPointer
}
