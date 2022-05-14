package main

import (
	"fmt"
	"mqtt-playground/mqttbuffer"
)

func main() {

	b := mqttbuffer.Newmqttbuffer()

	fmt.Println("Write pointer before adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer before adding new msg: ", b.GetReadPointer())

	b = b.AddMessage("first message")
	fmt.Println(b)
	b = b.AddMessage("second message")
	fmt.Println(b)
	b = b.AddMessage("third message")
	fmt.Println(b)
	b = b.AddMessage("fourth message")
	fmt.Println(b)
	b = b.AddMessage("fifth message")
	fmt.Println(b)
	b = b.AddMessage("sixth message")
	fmt.Println(b)

	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())

	if b.GetReadPointer() < b.GetWritePointer() {
		msg := b.ReadMessage(b.GetReadPointer())
		fmt.Println("Received message: ", msg)
		b = b.ReadNextMessage()
	}
	fmt.Println("Read pointer after reading next msg: ", b.GetReadPointer())

}
