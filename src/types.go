package main

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

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
}
