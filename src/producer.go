package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/gin-gonic/gin"
	logger "github.com/mkm29/mpulsar/pkg/logging"
)

// func publish_http(w http.ResponseWriter, r *http.Request) {

// }

func publish(c *gin.Context) {
	// Log request
	logger.Log("INFO", logger.WithRequest(c.Request))
	// Declare a new User struct.
	var m Message

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	de := json.NewDecoder(c.Request.Body).Decode(&m)
	if de != nil {
		logger.Log("ERROR", logger.WithRequest(c.Request), de)
		c.JSON(http.StatusBadRequest, gin.H{"error": de.Error()})
		return
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: fmt.Sprint("pulsar://", pulsarURL, ":", pulsarPort),
	})

	defer client.Close()

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topicName,
	})

	// Decode Message object to JSON
	js, err := json.Marshal(m)
	if err != nil {
		logger.Log("ERROR", logger.WithRequest(c.Request), de)
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
