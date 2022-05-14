package mqttbuffer

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttbuffer struct {
	buffer       [5]mqtt.Message
	readPointer  int
	writePointer int
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

func (b mqttbuffer) AddMessage(message mqtt.Message) {
	if b.writePointer == (len(b.buffer) - 1) {
		b.writePointer = 0
	}
	b.buffer[b.writePointer] = message
	b.writePointer++
	fmt.Println("Write pointer after adding new msg: ", b.writePointer)
}

func (b mqttbuffer) ReadMessage(index int) mqtt.Message {
	if index < len(b.buffer) {
		return b.buffer[index]
	}
	fmt.Printf("Index %d greater then buffer size [%d]", index, (len(b.buffer) - 1))
	return nil
}

func (b mqttbuffer) ReadNextMessage() mqtt.Message {
	if b.readPointer == (len(b.buffer) - 1) {
		b.readPointer = 0
	}
	if b.readPointer < b.writePointer {
		msg := b.buffer[b.readPointer]
		b.readPointer = b.readPointer + 1
		return msg
	}
	fmt.Println("No new messages on the buffer")
	return nil
}
