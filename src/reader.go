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
		fmt.Printf("Received message: %v\n", string(msg.Payload()))
		var u User
		err = json.Unmarshal(msg.Payload(), &u)
		if err != nil {
			log.Fatal(err)
		}
		// print User object
		fmt.Printf("Unmarshalled to User struct: %+v", u)
		// send User struct to channel for processing
		// users <- u
	}
}
