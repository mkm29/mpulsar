package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
)

// declare constants for Pulsar
const (
	server = "54.211.97.53"
	port   = 6650
	topic  = "users"
)

var (
	client   pulsar.Client
	producer pulsar.Producer
	consumer pulsar.Consumer
	reader   pulsar.Reader
)

// create struct to hold message
type Message struct {
	ID   pulsar.MessageID
	Text string
}

type User struct {
	ID        int
	FirstName string
	LastName  string
}

func main() {
	// initialize Pulsar
	initPulsar()

	mux := http.NewServeMux()
	mux.HandleFunc("/user/create", userCreate)
	mux.HandleFunc("/user/list", listUsers)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}

// initialize all Pulsar variables
func initPulsar() {
	// create client
	client, err := createClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	// create producer
	producer, err := createProducer(client, "my-topic")
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	// create consumer
	consumer, err := createConsumer(client, "my-topic", "my-sub")
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	// create reader
	// create ReaderOptions
	ro := pulsar.ReaderOptions{
		Topic: topic,
		Name:  "my-reader",
		Properties: map[string]string{
			"key1": "value1",
		},
		StartMessageID:          pulsar.EarliestMessageID(),
		StartMessageIDInclusive: true,
		// MessageChannel:          make(chan pulsar.ConsumerMessage),
		ReceiverQueueSize: 10,
	}
	reader, err := createReader(client, ro)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
}

func userCreate(w http.ResponseWriter, r *http.Request) {
	// Declare a new Person struct.
	var m Message

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Produce message to Pulsar
	msgId, err := sendMessage(producer, m.Text)
	_ = msgId
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "User: %+v", m)

}

func createClient() (pulsar.Client, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: fmt.Sprintf("pulsar://%s:%d", server, port),
	})

	return client, err
}

func createConsumer(client pulsar.Client, t string, s string) (pulsar.Consumer, error) {
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            t,
		SubscriptionName: s,
		Type:             pulsar.Shared,
	})

	return consumer, err
}

func createProducer(client pulsar.Client, t string) (pulsar.Producer, error) {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: t,
	})
	return producer, err
}

func sendMessage(producer pulsar.Producer, msg string) (pulsar.MessageID, error) {
	ctx := context.Background()

	msgId, err := producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: []byte(msg),
	})

	if err != nil {
		log.Fatal(err)
	}

	return msgId, err
}

func consumeMessage(consumer pulsar.Consumer) pulsar.Message {
	ctx := context.Background()

	msg, err := consumer.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}

	consumer.Ack(msg)

	return msg
}

func createReader(client pulsar.Client, ro pulsar.ReaderOptions) (pulsar.Reader, error) {
	/*
	 * Create a Reader with ReaderOptions
	 *
	 * ReaderOptions are optional and can be used to configure the Reader
	 *
	 * For example, you can set the Reader's name, schema, and schema type
	 *
	 * For more information on ReaderOptions, see:
	 * https://godoc.org/github.com/apache/pulsar-client-go/pulsar#ReaderOptions
	 */
	reader, err := client.CreateReader(pulsar.ReaderOptions{
		Topic:                   ro.Topic,
		Name:                    ro.Name,
		Properties:              ro.Properties,
		StartMessageID:          ro.StartMessageID,
		StartMessageIDInclusive: ro.StartMessageIDInclusive,
		MessageChannel:          ro.MessageChannel,
		ReceiverQueueSize:       ro.ReceiverQueueSize,
		SubscriptionRolePrefix:  ro.SubscriptionRolePrefix,
		Decryption:              ro.Decryption,
	})
	return reader, err
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	// set limit of 10
	limit := 10
	// create list of byte strings to hold users
	var users []string
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
		// check if we have reached the limit
		if len(users) == limit {
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
