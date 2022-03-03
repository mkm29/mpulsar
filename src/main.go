package main

import (
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
)

const (
	PULSAR_URL        = "54.211.97.53"
	PULSAR_PORT       = 6650
	TOPIC_NAME        = "users"
	SUBSCRIPTION_NAME = "test-sub"
)

var (
	ro pulsar.ReaderOptions
)

func main() {
	// initialize Pulsar variables
	initializeVars()
	// create HTTP server
	mux := http.NewServeMux()
	// Add handler for /user/list
	mux.HandleFunc("/user/list", read)
	// Add handler for publishing message
	mux.HandleFunc("/user/create", publish)

	// Listen on port 8080
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}

func initializeVars() {
	ro = pulsar.ReaderOptions{
		Topic:          TOPIC_NAME,
		StartMessageID: pulsar.EarliestMessageID(),
		// MessageChannel:    make(chan pulsar.ConsumerMessage),
		ReceiverQueueSize: 10,
	}
}
