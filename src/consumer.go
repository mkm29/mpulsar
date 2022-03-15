package main

import (
	"context"
	"fmt"

	logger "github.com/mkm29/mpulsar/pkg/logging"

	"github.com/apache/pulsar-client-go/pulsar"
)

func subscribe() {
	logger.Log("INFO", "Subscribing to topic")
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: fmt.Sprintf("pulsar://%s:%d", pulsarURL, pulsarPort),
	})

	defer client.Close()

	consumer, err := client.Subscribe(co)

	defer consumer.Close()

	msg, err := consumer.Receive(context.Background())
	if err != nil {
		logger.Log("ERROR", err)
		return
	}

	logger.Log("INFO", fmt.Sprintf("Received message: %s", msg.Payload()))
	// Push message onto consumeChan
	consumeChan <- pulsar.ConsumerMessage{Message: msg}
	// Acknowledge message
	consumer.Ack(msg)
}
