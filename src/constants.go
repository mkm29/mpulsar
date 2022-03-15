package main

const (
	pulsarURL         = "localhost"
	pulsarPort        = 6650
	topicName         = "locations"
	subscriptionName  = "geocode-sub"
	peliasURL         = "http://localhost:4000/v1/search?text=%s&size=1&layers=address"
	cassandraIP       = "localhost"
	cassandraPort     = 9042
	cassandraKeyspace = "mortal_mint"
	cassandraCL       = "Quorum"
)
