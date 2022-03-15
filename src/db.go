package main

/*
	This should actually be a seperate
	microservice, when geocoded metadata is added
	to location, push to Pulsar topic, and have
	a consumer pull from the topic and insert message
	into Cassandra.
*/

import (
	"errors"
	"os"

	"github.com/gocql/gocql"
	logger "github.com/mkm29/mpulsar/pkg/logging"
)

func connect(useAuth bool) (*gocql.Session, error) {
	// use gocql to connect to Cassandra cluster
	logger.Log("INFO", "Connecting to Cassandra cluster")

	// create a cluster object
	cluster := gocql.NewCluster(cassandraIP)
	cluster.Keyspace = cassandraKeyspace
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	// if useAuth is true, set the credentials
	if useAuth {
		logger.Log("INFO", "Using authentication")
		missingCreds := false
		// Get username and password from environment variables
		password, ok := os.LookupEnv("CASSANDRA_PASSWORD")
		if !ok {
			logger.Log("ERROR", "CASSANDRA_PASSWORD not set")
			missingCreds = true
		}
		username, ok := os.LookupEnv("CASSANDRA_USERNAME")
		if !ok {
			logger.Log("ERROR", "CASSANDRA_USERNAME not set")
			missingCreds = true
		}
		if !missingCreds {
			return nil, errors.New("CASSANDRA_USERNAME/PASSWORD not set")
		}

		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: password,
			Password: username,
		}
	}

	// create session
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Log("ERROR", err)
		return nil, err
	}
	defer session.Close()
	return session, nil

}

func createTable(session *gocql.Session) (bool, error) {
	// create a locations table in Cassandra
	// with the following columns:
	// id: UUID
	// point: POINTTYPE
	// metadata: TEXT
	// PointType is only supported in DSE 6.7+
	isDse := true
	var sql string
	if isDse {
		sql = "CREATE TABLE IF NOT EXISTS locations (id uuid, point 'PointType', metadata text, PRIMARY KEY((id), point));"
	} else {
		sql = "CREATE TABLE IF NOT EXISTS locations (id uuid, latitude float, longitude float, metadata text, PRIMARY KEY((id), latitude, longitude));"
	}
	// execute the query
	err := session.Query(sql).Exec()
	if err != nil {
		logger.Log("ERROR", err)
		return false, err
	}
	return true, nil
}

func insertData(session *gocql.Session, l Location) (bool, error) {
	// insert a location into the locations table
	isDse := true
	var sql string
	if isDse {
		sql = "INSERT INTO locations (id, point, metadata) VALUES (uuid(), POINT(? ?), ?);"
	} else {
		sql = "INSERT INTO locations (id, latitude, longitude, metadata) VALUES (uuid, ?, ?, ?);"
	}
	// execute the query
	err := session.Query(sql, l.Latitude, l.Longitude, l.Label).Exec()

	if err != nil {
		logger.Log("ERROR", err)
		return false, err
	}
	return true, nil
	// return the id of the inserted location

}
