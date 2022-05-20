/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

/*
To run this sample, The following certificates
must be created:
  rootCA-crt.pem - root certificate authority that is used
                   to sign and verify the client and server
                   certificates.
  rootCA-key.pem - keyfile for the rootCA.
  server-crt.pem - server certificate signed by the CA.
  server-key.pem - keyfile for the server certificate.
  client-crt.pem - client certificate signed by the CA.
  client-key.pem - keyfile for the client certificate.
  CAfile.pem     - file containing concatenated CA certificates
                   if there is more than 1 in the chain.
                   (e.g. root CA -> intermediate CA -> server cert)
  Instead of creating CAfile.pem, rootCA-crt.pem can be added
  to the default openssl CA certificate bundle. To find the
  default CA bundle used, check:
  $GO_ROOT/src/pks/crypto/x509/root_unix.go
  To use this CA bundle, just set tls.Config.RootCAs = nil.
*/

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("samplecerts/CAfile.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair("samplecerts/client-crt.pem", "samplecerts/client-key.pem")
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type Message struct {
	Message string
}

func main() {
	//tlsconfig := NewTLSConfig()

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://test.mosquitto.org:1883")
	opts.SetClientID("cg-playground") //.SetTLSConfig(tlsconfig)
	opts.SetDefaultPublishHandler(f)

	// Start the connection
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Connection is up")

	i := 1
<<<<<<< HEAD:publishers/pub1/main.go
	for range time.Tick(time.Duration(500) * time.Millisecond) {
		if i == 1001 {
=======
	for range time.Tick(time.Duration(100) * time.Millisecond) {
		if i == 5001 {
>>>>>>> db369cd50a38d3b8838f3cafabeac7c15910afb3:publisher/main.go
			break
		}
		text := fmt.Sprintf("this is msg #%d!", i)
		msg, err := json.Marshal(Message{Message: text})
		if err != nil {
			panic(err)
		}

<<<<<<< HEAD:publishers/pub1/main.go
		c.Publish("cg-connector/sample/1", 0, false, msg)
=======
		c.Publish("/cg-connector/sample/1", 0, false, msg)
>>>>>>> db369cd50a38d3b8838f3cafabeac7c15910afb3:publisher/main.go
		fmt.Println(msg)
		i++
	}

	c.Disconnect(1000)
	fmt.Println("shutdown complete")
}