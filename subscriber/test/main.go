package main

import (
	"fmt"
	Mqttbuffer "mqtt-playground/mqttbuffer"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	b := Mqttbuffer.NewMqttbuffer()

	///read messages
	go func() {
		for {
			if b.NewMessage() {
				msg, err := b.ReadMessage(b.GetReadPointer())
				if err != nil {
					panic(err.Error())
				}
				fmt.Println(msg)
				b = b.NextMessage()
				fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())
				//time.Sleep(time.Second)
			}
		}
	}()

	fmt.Println("Write pointer start value: ", b.GetWritePointer())
	fmt.Println("Read pointer start value: ", b.GetReadPointer())

	///new message from type Message
	msg := Mqttbuffer.Message{}

	///first message
	msg.Payload = "first message"
	msg.Topic = "device/1"
	b = b.AddMessage(msg)
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	time.Sleep(time.Second)

	///second message
	msg.Payload = "second message"
	msg.Topic = "device/2"
	b = b.AddMessage(msg)
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	time.Sleep(time.Second)

	///third message
	msg.Payload = "third message"
	msg.Topic = "device/3"
	b = b.AddMessage(msg)
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	time.Sleep(time.Second)

	///fourth message
	msg.Payload = "fourth message"
	msg.Topic = "device/4"
	b = b.AddMessage(msg)
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	time.Sleep(time.Second)

	///fifth message
	msg.Payload = "fifth message"
	msg.Topic = "device/5"
	b = b.AddMessage(msg)
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	time.Sleep(time.Second)

	///sixth message
	msg.Payload = "sixth message"
	msg.Topic = "device/6"
	b = b.AddMessage(msg)
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	time.Sleep(time.Second)

	///seventh message
	msg.Payload = "seventh message"
	msg.Topic = "device/7"
	b = b.AddMessage(msg)
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())

	fmt.Println(b)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	fmt.Println("shutdown complete")

}
