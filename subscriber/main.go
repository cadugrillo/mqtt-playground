package main

// Connect to the broker, subscribe, and write messages received to a file

import (
	"fmt"
	Mqttbuffer "mqtt-playground/mqttbuffer"
	"mqtt-playground/subscriber/config"
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

// handler is a simple struct that provides a function to be called when a message is received. The message is parsed
// and the count followed by the raw message is written to the file (this makes it easier to sort the file)
type handler struct {
	f *os.File
}

func NewHandler() *handler {
	var f *os.File
	if ConfigFile.Logs.WriteToDisk {
		var err error
		f, err = os.Create(ConfigFile.Logs.OutputFile)
		if err != nil {
			panic(err)
		}
	}
	return &handler{f: f}
}

// Close closes the file
func (o *handler) Close() {
	if o.f != nil {
		if err := o.f.Close(); err != nil {
			fmt.Printf("ERROR closing file: %s", err)
		}
		o.f = nil
	}
}

// handle is called when a message is received
func (o *handler) handle(_ mqtt.Client, msg mqtt.Message) {

	var recmsg Mqttbuffer.Message
	recmsg.Duplicate = msg.Duplicate()
	recmsg.Qos = msg.Qos()
	recmsg.Retained = msg.Retained()
	recmsg.Topic = msg.Topic()
	recmsg.MessageID = msg.MessageID()
	recmsg.Payload = string(msg.Payload())

	b = b.AddMessage(recmsg)

	if o.f != nil {
		if _, err := o.f.WriteString(recmsg.Payload); err != nil {
			fmt.Printf("ERROR writing to file: %s", err)
		}
	}
}

func main() {

	ConfigFile = config.ReadConfig()

	// Enable logging by uncommenting the below
	// mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	// mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	// mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	// Create a handler that will deal with incoming messages
	h := NewHandler()
	defer h.Close()

	// Now we establish the connection to the mqtt broker
	opts := mqtt.NewClientOptions()
	opts.AddBroker(ConfigFile.Client.ServerAddress)
	opts.SetClientID(ConfigFile.Client.ClientId)

	opts.SetOrderMatters(ConfigFile.Client.OrderMaters)                                      // Allow out of order messages (use this option unless in order delivery is essential)
	opts.ConnectTimeout = (time.Duration(ConfigFile.Client.ConnectionTimeout) * time.Second) // Minimal delays on connect
	opts.WriteTimeout = (time.Duration(ConfigFile.Client.WriteTimeout) * time.Second)        // Minimal delays on writes
	opts.KeepAlive = int64(ConfigFile.Client.KeepAlive)                                      // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = (time.Duration(ConfigFile.Client.PingTimeout) * time.Second)          // local broker so response should be quick

	// Automate connection management (will keep trying to connect and will reconnect if network drops)
	opts.ConnectRetry = ConfigFile.Client.ConnectRetry
	opts.AutoReconnect = ConfigFile.Client.AutoConnect

	// If using QOS2 and CleanSession = FALSE then it is possible that we will receive messages on topics that we
	// have not subscribed to here (if they were previously subscribed to they are part of the session and survive
	// disconnect/reconnect). Adding a DefaultPublishHandler lets us detect this.
	opts.DefaultPublishHandler = func(_ mqtt.Client, msg mqtt.Message) {
		fmt.Printf("UNEXPECTED MESSAGE: %s\n", msg)
	}

	// Log events
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		fmt.Println("connection lost")
	}

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("connection established")

		// Establish the subscription - doing this here means that it will happen every time a connection is established
		// (useful if opts.CleanSession is TRUE or the broker does not reliably store session data)
		for i := 0; i < len(ConfigFile.Topics.Topic); i++ {
			t := c.Subscribe(ConfigFile.Topics.Topic[i], byte(ConfigFile.Client.Qos), h.handle)
			id := i

			// the connection handler is called in a goroutine so blocking here would not cause an issue. However as blocking
			// in other handlers does cause problems its best to just assume we should not block
			go func() {
				_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
				if t.Error() != nil {
					fmt.Printf("ERROR SUBSCRIBING: %s\n", t.Error())
				} else {
					fmt.Println("subscribed to: ", ConfigFile.Topics.Topic[id])
				}
			}()
		}
	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		fmt.Println("attempting to reconnect")
	}

	//
	// Connect to the broker
	//
	client := mqtt.NewClient(opts)

	// If using QOS2 and CleanSession = FALSE then messages may be transmitted to us before the subscribe completes.
	// Adding routes prior to connecting is a way of ensuring that these messages are processed
	for i := 0; i < len(ConfigFile.Topics.Topic); i++ {
		client.AddRoute(ConfigFile.Topics.Topic[i], h.handle)
	}

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connection is up")

	go func() {
		for {
			if ConfigFile.Logs.WriteToLog {
				if b.NewMessage() {
					msg, err := b.ReadMessage(b.GetReadPointer())
					if err != nil {
						panic(err.Error())
					}
					fmt.Println(msg.Payload)
					b = b.NextMessage()
				}

			}
		}
	}()

	// Messages will be delivered asynchronously so we just need to wait for a signal to shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	client.Disconnect(1000)
	fmt.Println("shutdown complete")
}
