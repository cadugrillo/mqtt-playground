package main

import (
	"fmt"
	"log"
	"mqtt-playground/connector/config"
	Mqttbuffer "mqtt-playground/mqttbuffer"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	once       sync.Once
	ConfigFile config.Config
	b          Mqttbuffer.Mqttbuffer
)

func init() {
	once.Do(initialise)
}

func initialise() {
	b = Mqttbuffer.NewMqttbuffer()
}

type handler struct {
	f bool
}

func NewHandler() *handler {
	var f bool
	return &handler{f: f}
}

func (o *handler) handle(_ mqtt.Client, msg mqtt.Message) {

	var recmsg Mqttbuffer.Message
	recmsg.Duplicate = msg.Duplicate()
	recmsg.Qos = msg.Qos()
	recmsg.Retained = msg.Retained()
	recmsg.MessageID = msg.MessageID()
	recmsg.Payload = string(msg.Payload())

	for i := 0; i < len(ConfigFile.TopicsSub.Topic); i++ {
		if ConfigFile.TopicsSub.Topic[i] == msg.Topic() {
			recmsg.Topic = ConfigFile.TopicsPub.Topic[i]
			break
		}
	}

	b = b.AddMessage(recmsg)
}

func main() {

	ConfigFile = config.ReadConfig()

	//logs
	if ConfigFile.Logs.Error {
		mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	}
	if ConfigFile.Logs.Critical {
		mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	}
	if ConfigFile.Logs.Warning {
		mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	}
	if ConfigFile.Logs.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	}

	h := NewHandler()

	optsSub := mqtt.NewClientOptions()
	optsSub.AddBroker(ConfigFile.ClientSub.ServerAddress)
	optsSub.SetClientID(ConfigFile.ClientSub.ClientId)
	optsSub.SetOrderMatters(ConfigFile.ClientSub.OrderMaters)                                      // Allow out of order messages (use this option unless in order delivery is essential)
	optsSub.ConnectTimeout = (time.Duration(ConfigFile.ClientSub.ConnectionTimeout) * time.Second) // Minimal delays on connect
	optsSub.WriteTimeout = (time.Duration(ConfigFile.ClientSub.WriteTimeout) * time.Second)        // Minimal delays on writes
	optsSub.KeepAlive = int64(ConfigFile.ClientSub.KeepAlive)                                      // Keepalive every 10 seconds so we quickly detect network outages
	optsSub.PingTimeout = (time.Duration(ConfigFile.ClientSub.PingTimeout) * time.Second)          // local broker so response should be quick
	optsSub.ConnectRetry = ConfigFile.ClientSub.ConnectRetry                                       // Automate connection management (will keep trying to connect and will reconnect if network drops)
	optsSub.AutoReconnect = ConfigFile.ClientSub.AutoConnect
	optsSub.DefaultPublishHandler = func(_ mqtt.Client, msg mqtt.Message) { fmt.Printf("SUB BROKER - UNEXPECTED : %s\n", msg) }
	optsSub.OnConnectionLost = func(cl mqtt.Client, err error) { fmt.Println("SUB BROKER - CONNECTION LOST") } // Log events

	optsSub.OnConnect = func(c mqtt.Client) {
		fmt.Println("SUB BROKER - CONNECTION STABLISHED")

		// Establish the subscription - doing this here means that it will happen every time a connection is established
		// (useful if opts.CleanSession is TRUE or the broker does not reliably store session data)
		for i := 0; i < len(ConfigFile.TopicsSub.Topic); i++ {
			t := c.Subscribe(ConfigFile.TopicsSub.Topic[i], byte(ConfigFile.ClientSub.Qos), h.handle)
			id := i

			// the connection handler is called in a goroutine so blocking here would not cause an issue. However as blocking
			// in other handlers does cause problems its best to just assume we should not block
			go func() {
				_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
				if t.Error() != nil {
					fmt.Printf("SUB BROKER - ERROR SUBSCRIBING TO : %s\n", t.Error())
				} else {
					fmt.Println("SUB BROKER - SUBSCRIBED TO : ", ConfigFile.TopicsSub.Topic[id])
				}
			}()
		}
	}

	optsSub.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) { fmt.Println("SUB BROKER - ATTEMPTING TO RECONNECT") }

	/////opts for Pub Broker
	optsPub := mqtt.NewClientOptions()
	optsPub.AddBroker(ConfigFile.ClientPub.ServerAddress)
	optsPub.SetClientID(ConfigFile.ClientPub.ClientId)
	optsPub.SetOrderMatters(ConfigFile.ClientPub.OrderMaters)                                      // Allow out of order messages (use this option unless in order delivery is essential)
	optsPub.ConnectTimeout = (time.Duration(ConfigFile.ClientPub.ConnectionTimeout) * time.Second) // Minimal delays on connect
	optsPub.WriteTimeout = (time.Duration(ConfigFile.ClientPub.WriteTimeout) * time.Second)        // Minimal delays on writes
	optsPub.KeepAlive = int64(ConfigFile.ClientPub.KeepAlive)                                      // Keepalive every 10 seconds so we quickly detect network outages
	optsPub.PingTimeout = (time.Duration(ConfigFile.ClientPub.PingTimeout) * time.Second)          // local broker so response should be quick
	optsPub.ConnectRetry = ConfigFile.ClientPub.ConnectRetry                                       // Automate connection management (will keep trying to connect and will reconnect if network drops)
	optsPub.AutoReconnect = ConfigFile.ClientPub.AutoConnect
	optsPub.DefaultPublishHandler = func(_ mqtt.Client, msg mqtt.Message) { fmt.Printf("PUB BROKER - UNEXPECTED : %s\n", msg) }
	optsPub.OnConnectionLost = func(cl mqtt.Client, err error) { fmt.Println("PUB BROKER - CONNECTION LOST") } // Log events
	optsPub.OnConnect = func(c mqtt.Client) { fmt.Println("PUB BROKER - CONNECTION STABLISHED") }
	optsPub.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) { fmt.Println("PUB BROKER - ATTEMPTING TO RECONNECT") }

	//
	// Connect to the SUB broker
	//
	clientSub := mqtt.NewClient(optsSub)

	// If using QOS2 and CleanSession = FALSE then messages may be transmitted to us before the subscribe completes.
	// Adding routes prior to connecting is a way of ensuring that these messages are processed
	for i := 0; i < len(ConfigFile.TopicsSub.Topic); i++ {
		clientSub.AddRoute(ConfigFile.TopicsSub.Topic[i], h.handle)
	}

	if tokenSub := clientSub.Connect(); tokenSub.Wait() && tokenSub.Error() != nil {
		panic(tokenSub.Error())
	}
	fmt.Println("SUB BROKER  - CONNECTION IS UP")

	//
	//connect to PUB broker
	//
	clientPub := mqtt.NewClient(optsPub)

	if tokenPub := clientPub.Connect(); tokenPub.Wait() && tokenPub.Error() != nil {
		panic(tokenPub.Error())
	}
	fmt.Println("PUB BROKER  - CONNECTION IS UP")

	go func() {
		for {
			if b.NewMessage() {
				msg, err := b.ReadMessage(b.GetReadPointer())
				if err != nil {
					panic(err.Error())
				}
				if ConfigFile.Logs.SubPayload {
					fmt.Println(msg.Payload)
				}

				clientPub.Publish(msg.Topic, msg.Qos, msg.Retained, msg.Payload)
				b = b.NextMessage()
			}

		}
	}()

	// Messages will be delivered asynchronously so we just need to wait for a signal to shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	clientSub.Disconnect(1000)
	clientPub.Disconnect(1000)
	fmt.Println("shutdown complete")
}
