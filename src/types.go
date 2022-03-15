package main

// Location struct
type Location struct {
	Confidence  float64
	Latitude    float64
	Longitude   float64
	City        string
	Address     string
	Region      string
	Country     string
	CountryCode string
	PostalCode  int
	Label       string
	Geohashes   []string
}

// Message type (used with Pub/Sub Pulsar)
type Message struct {
	Payload string
}
