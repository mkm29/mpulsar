package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
	logger "github.com/mkm29/mpulsar/pkg/logging"
)

// func publish_http(w http.ResponseWriter, r *http.Request) {

// }

func publish(w http.ResponseWriter, r *http.Request) {
	// Log request
	logger.Log("INFO", logger.WithRequest(r))
	// Declare a new User struct.
	var m Message

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	de := json.NewDecoder(r.Body).Decode(&m)
	if de != nil {
		http.Error(w, de.Error(), http.StatusBadRequest)
		logger.Log("ERROR", logger.WithRequest(r), de)
		return
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: fmt.Sprint("pulsar://", PULSAR_URL, ":", PULSAR_PORT),
	})

	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: TOPIC_NAME,
	})

	// Decode User object to JSON
	js, err := json.Marshal(m)
	if err != nil {
		logger.Log("ERROR", logger.WithRequest(r), de)
		return
	}
	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: js,
	})

	defer producer.Close()

	if err != nil {
		fmt.Println("Failed to publish message", err)
		logger.Log("ERROR", err)
	} else {
		fmt.Println("Published message")
	}
}
