package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

func subscribe() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: fmt.Sprintf("pulsar://%s:%d", PULSAR_URL, PULSAR_PORT),
	})

	defer client.Close()

	consumer, err := client.Subscribe(co)

	defer consumer.Close()

	msg, err := consumer.Receive(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
		msg.ID(), string(msg.Payload()))
	// Push message onto consumeChan
	consumeChan <- pulsar.ConsumerMessage{Message: msg}
	// Acknowledge message
	consumer.Ack(msg)
}
