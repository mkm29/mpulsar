package main

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

type Message struct {
	Text string
}
