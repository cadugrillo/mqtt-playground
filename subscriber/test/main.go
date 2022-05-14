package main

import (
	"fmt"
	"mqtt-playground/mqttbuffer"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	var msg mqtt.Message

	b := mqttbuffer.Newmqttbuffer()

	fmt.Println("Write pointer before adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer before adding new msg: ", b.GetReadPointer())

	b.AddMessage(msg)

	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())

	if b.GetReadPointer() < b.GetWritePointer() {
		fmt.Println("Received message: ", b.ReadNextMessage())
	}
	fmt.Println("Read pointer after reading next msg: ", b.GetReadPointer())
}
