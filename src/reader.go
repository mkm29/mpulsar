package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
)

func read(w http.ResponseWriter, r *http.Request) {
	var users []string
	client, err := pulsar.NewClient(pulsar.ClientOptions{URL: fmt.Sprintf("pulsar://%s:%d", PULSAR_URL, PULSAR_PORT)})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	reader, err := client.CreateReader(ro)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	for reader.HasNext() {
		msg, err := reader.Next(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// create User object from message
		js, err := json.Marshal(msg.Payload())
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, string(js))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
