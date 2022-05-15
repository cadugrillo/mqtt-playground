package main

import (
	"fmt"
	"mqtt-playground/mqttbuffer"
)

func main() {

	b := mqttbuffer.Newmqttbuffer()

	fmt.Println("Write pointer start value: ", b.GetWritePointer())
	fmt.Println("Read pointer start value: ", b.GetReadPointer())

	b = b.AddMessage("first message")
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())
	msg := b.ReadMessage(b.GetReadPointer())
	fmt.Println("Received message: ", msg)
	b = b.ReadNextMessage()
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())

}
