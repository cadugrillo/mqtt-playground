package main

import (
	"fmt"
	"mqtt-playground/mqttbuffer"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	b := mqttbuffer.Newmqttbuffer()

	///read messages
	go func() {
		for {
			if b.NewMessage() {
				msg := b.ReadMessage(b.GetReadPointer())
				fmt.Println(msg)
				b = b.NextMessage()
			}
		}
	}()

	fmt.Println("Write pointer start value: ", b.GetWritePointer())
	fmt.Println("Read pointer start value: ", b.GetReadPointer())
	///first message
	b = b.AddMessage("first message")
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())
	///second message
	b = b.AddMessage("second message")
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())

	///third message
	b = b.AddMessage("third message")
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())

	///fourth message
	b = b.AddMessage("fourth message")
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())

	///fifth message
	b = b.AddMessage("fifth message")
	fmt.Println("Write pointer after adding new msg: ", b.GetWritePointer())
	fmt.Println("Read pointer after adding new msg: ", b.GetReadPointer())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	fmt.Println("shutdown complete")

}
