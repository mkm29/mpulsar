package main

const (
	PULSAR_URL         = "localhost"
	PULSAR_PORT        = 6650
	TOPIC_NAME         = "locations"
	SUBSCRIPTION_NAME  = "geocode-sub"
	PELIAS_URL         = "http://localhost:4000/v1/search?text=%s&size=1&layers=address"
	CASSANDRA_IP       = "localhost"
	CASSANDRA_PORT     = 9042
	CASSANDRA_KEYSPACE = "mortal_mint"
	CASSANDRA_CL       = "Quorum"
)
