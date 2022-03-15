package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	logger "github.com/mkm29/mpulsar/pkg/logging"
	"github.com/mkm29/mpulsar/pkg/utils"

	"github.com/apache/pulsar-client-go/pulsar"
)

var (
	ro          pulsar.ReaderOptions
	co          pulsar.ConsumerOptions
	readChan    chan pulsar.ReaderMessage
	consumeChan chan pulsar.ConsumerMessage
	// LOGLEVEL = level to use for logging
	LOGLEVEL string = utils.GetEnv("LOGLEVEL", "INFO")
)

func main() {
	// Configure logging
	logger.Configure()
	logger.Log("INFO", "Starting smigPulsar Go service")

	// initialize Pulsar variables
	initializeVars()
	// create HTTP server with Gin
	r := gin.Default()

	// Add handler for /user/list
	r.GET("/locations/list", read)

	// Add handler for publishing message
	r.GET("/locations/create", publish)
	// Create handler for geocode
	r.GET("/locations/geocode", geocodeHTTP)

	// Listen on port 8080
	logger.Log("INFO", "Listening on port 8080")
	r.Run(":8080")
}

func initializeVars() {
	logger.Log("INFO", "Initializing Pulsar variables")
	// create a channel for reading messages
	readChan = make(chan pulsar.ReaderMessage)
	// create a channel for consuming messages
	consumeChan = make(chan pulsar.ConsumerMessage)
	ro = pulsar.ReaderOptions{
		Topic:             topicName,
		StartMessageID:    pulsar.EarliestMessageID(),
		MessageChannel:    readChan,
		ReceiverQueueSize: 10,
	}
	logger.Log("INFO", fmt.Sprintf("ReaderOptions: %+v", ro))
	co = pulsar.ConsumerOptions{
		Topic:            topicName,
		SubscriptionName: subscriptionName,
		Type:             pulsar.Shared,
		MessageChannel:   consumeChan,
	}
	logger.Log("INFO", fmt.Sprintf("ConsumerOptions: %+v", co))
}
