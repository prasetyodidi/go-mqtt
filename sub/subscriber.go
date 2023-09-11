package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	strMessage := string(message.Payload())
	// fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())
	fmt.Printf("Received message: %s\n", message.Payload())
	if strMessage == "n" {
		client.Disconnect(1)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected!")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println("Connect lost: %v", err)
}

func getClient() mqtt.Client {
	broker := "broker.emqx.io"
	port := 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_sub")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func main() {
	client := getClient()
	sub(client)

	time.Sleep(time.Minute)
	client.Disconnect(uint(time.Minute))
}

func sub(client mqtt.Client) {
	topic := "didi-topic"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subcribed to topic: %s\n", topic)
}
