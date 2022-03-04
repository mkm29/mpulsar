package main

import (
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
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
	ro = pulsar.ReaderOptions{
		Topic:          TOPIC_NAME,
		StartMessageID: pulsar.EarliestMessageID(),
		// MessageChannel:    make(chan pulsar.ConsumerMessage),
		ReceiverQueueSize: 10,
	}
}
