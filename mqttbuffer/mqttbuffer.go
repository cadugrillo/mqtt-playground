package mqttbuffer

import (
	"errors"
	"fmt"
)

type mqttbuffer struct {
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

func Newmqttbuffer() mqttbuffer {
	b := mqttbuffer{}
	return b
}

func (b mqttbuffer) GetReadPointer() int {
	return b.readPointer
}

func (b mqttbuffer) GetWritePointer() int {
	return b.writePointer
}

func (b mqttbuffer) AddMessage(message Message) mqttbuffer {
	if b.writePointer == len(b.buffer) {
		b.writePointer = 0
	}
	b.buffer[b.writePointer] = message
	b.writePointer++
	return b
}

func (b mqttbuffer) ReadMessage(index int) (Message, error) {
	if index < len(b.buffer) {
		return b.buffer[index], nil
	}
	msg := Message{}
	return msg, errors.New(fmt.Sprintf("Index %d greater then buffer size [%d]", index, len(b.buffer)))
}

func (b mqttbuffer) ReadNextMessage() mqttbuffer {
	if b.readPointer == len(b.buffer) {
		b.readPointer = 0
	}
	if b.readPointer < b.writePointer {
		b.readPointer++
		return b
	}
	fmt.Println("No new messages on the buffer")
	return b
}
