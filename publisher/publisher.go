package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())
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
	opts.SetClientID("go_mqtt_pub")
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

	client := getClient()

	publish(client, scanner)

	client.Disconnect(250)
}

func publish(client mqtt.Client, scanner *bufio.Scanner) {
	isAlive := true
	for isAlive {
		fmt.Print("Enter some text: ")

		scanner.Scan()
		userInput := scanner.Text()

		if userInput == "n" {
			isAlive = false
		}
		
		token := client.Publish("didi-topic", 0, false, userInput)
		token.Wait()
		time.Sleep(time.Second)
	}
}
