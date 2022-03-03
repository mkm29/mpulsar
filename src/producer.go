package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
)

func publish(w http.ResponseWriter, r *http.Request) {
	// Declare a new User struct.
	var u User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	de := json.NewDecoder(r.Body).Decode(&u)
	if de != nil {
		http.Error(w, de.Error(), http.StatusBadRequest)
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
	js, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: js,
	})

	defer producer.Close()

	if err != nil {
		fmt.Println("Failed to publish message", err)
	} else {
		fmt.Println("Published message")
	}
}
