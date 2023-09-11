package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("\n%s ", message.Payload())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected!")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println("Connect lost: %v", err)
}

func getClient(name string) mqtt.Client {
	broker := "broker.emqx.io"
	port := 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(name)
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
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("enter name: ")
	scanner.Scan()
	name := scanner.Text()

	client := getClient(name)

	sub(client)
	publish(client, scanner, name)

	client.Disconnect(250)
}

func publish(client mqtt.Client, scanner *bufio.Scanner, name string) {
	isAlive := true
	for isAlive {
		fmt.Print("\nEnter some text: ")

		scanner.Scan()
		userInput := scanner.Text()

		if userInput == "n" {
			isAlive = false
		}
		
		message := name + ": " + userInput
		
		token := client.Publish("didi-topic", 0, false, message)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "didi-topic"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subcribed to topic: %s\n", topic)
}
