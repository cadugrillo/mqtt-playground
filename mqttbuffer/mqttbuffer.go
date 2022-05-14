package mqttbuffer

import (
	"fmt"
)

type mqttbuffer struct {
	buffer       [5]string
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

func (b mqttbuffer) AddMessage(message string) mqttbuffer {
	if b.writePointer == len(b.buffer) {
		b.writePointer = 0
	}
	b.buffer[b.writePointer] = message
	b.writePointer++
	return b
}

func (b mqttbuffer) ReadMessage(index int) string {
	if index < len(b.buffer) {
		return b.buffer[index]
	}
	fmt.Printf("Index %d greater then buffer size [%d]", index, len(b.buffer))
	return ""
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
