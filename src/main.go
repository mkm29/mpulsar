package main

import (
	"net/http"

	logger "github.com/mkm29/mpulsar/pkg/log"

	"github.com/apache/pulsar-client-go/pulsar"
)

var (
	ro          pulsar.ReaderOptions
	co          pulsar.ConsumerOptions
	readChan    chan pulsar.ReaderMessage
	consumeChan chan pulsar.ConsumerMessage
)

func main() {
	logger.Info("Starting smigPulsar Go service")

	// initialize Pulsar variables
	initializeVars()
	// create HTTP server
	mux := http.NewServeMux()
	// Add handler for /user/list
	mux.HandleFunc("/locations/list", read)
	// Add handler for publishing message
	mux.HandleFunc("/locations/create", publish)
	// Create handler for geocode
	mux.HandleFunc("/locations/geocode", geocode_http)

	// Listen on port 8080
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}

func initializeVars() {
	// create a channel for reading messages
	readChan = make(chan pulsar.ReaderMessage)
	// create a channel for consuming messages
	consumeChan = make(chan pulsar.ConsumerMessage)
	ro = pulsar.ReaderOptions{
		Topic:             TOPIC_NAME,
		StartMessageID:    pulsar.EarliestMessageID(),
		MessageChannel:    readChan,
		ReceiverQueueSize: 10,
	}
	// log ReaderOptions object
	logger.Info("ReaderOptions: %+v", ro)
	co = pulsar.ConsumerOptions{
		Topic:            TOPIC_NAME,
		SubscriptionName: SUBSCRIPTION_NAME,
		Type:             pulsar.Shared,
		MessageChannel:   consumeChan,
	}
	// log ConsumerOptions object
	logger.Info("ConsumerOptions: %+v", co)
}
