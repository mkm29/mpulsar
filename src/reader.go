package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	logger "../pkg/log"
	"github.com/apache/pulsar-client-go/pulsar"
)

func read(w http.ResponseWriter, r *http.Request) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{URL: fmt.Sprintf("pulsar://%s:%d", PULSAR_URL, PULSAR_PORT)})
	if err != nil {
		logger.Error(err)
	}

	defer client.Close()

	reader, err := client.CreateReader(ro)
	if err != nil {
		logger.WithRequest(r).Error(err)
	}
	defer reader.Close()

	for reader.HasNext() {
		msg, err := reader.Next(context.Background())
		if err != nil {
			logger.WithRequest(r).Error(err)
		}
		fmt.Printf("Received message: %v\n", string(msg.Payload()))
		var m Message
		err = json.Unmarshal(msg.Payload(), &m)
		if err != nil {
			logger.WithRequest(r).Error(err)
		}
		// print Message object
		fmt.Printf("Unmarshalled to Message struct: %+v", m)
		// send Message struct to channel for processing
		readChan <- pulsar.ReaderMessage{Message: msg}
	}
}
