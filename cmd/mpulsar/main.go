package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
)

// declare constants for Pulsar
const (
	server     = "54.211.97.53"
	port       = 6650
	topicName1 = "my-topic-1"
	topicName2 = "my-topic-2"
)

func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Example Pulsar Client in Go!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	log.Println("Listing for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func createClient() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: fmt.Sprintf("pulsar://%s:%d", server, port),
	})

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
}
