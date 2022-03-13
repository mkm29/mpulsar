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

	logger "../pkg/log"
	"github.com/gocql/gocql"
)

func connect(use_auth bool) (error, *gocql.Session) {
	// use gocql to connect to Cassandra cluster

	// create a cluster object
	cluster := gocql.NewCluster(CASSANDRA_IP)
	cluster.Keyspace = CASSANDRA_KEYSPACE
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	// if use_auth is true, set the credentials
	if use_auth {
		missing_creds := false
		// Get username and password from environment variables
		password, ok := os.LookupEnv("CASSANDRA_PASSWORD")
		if !ok {
			logger.Error("CASSANDRA_PASSWORD not set")
			missing_creds = true
		}
		username, ok := os.LookupEnv("CASSANDRA_USERNAME")
		if !ok {
			logger.Error("CASSANDRA_USERNAME not set")
			missing_creds = true
		}
		if !missing_creds {
			return errors.New("CASSANDRA_USERNAME/PASSWORD not set"), nil
		}

		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: password,
			Password: username,
		}
	}

	// create session
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Error(err)
		return err, nil
	}
	defer session.Close()
	return nil, session

}

func create_table(session *gocql.Session) (error, bool) {
	// create a locations table in Cassandra
	// with the following columns:
	// id: UUID
	// point: POINTTYPE
	// metadata: TEXT
	// PointType is only supported in DSE 6.7+
	is_dse := true
	var sql string
	if is_dse {
		sql = "CREATE TABLE IF NOT EXISTS locations (id uuid, point 'PointType', metadata text, PRIMARY KEY((id), point));"
	} else {
		sql = "CREATE TABLE IF NOT EXISTS locations (id uuid, latitude float, longitude float, metadata text, PRIMARY KEY((id), latitude, longitude));"
	}
	// execute the query
	err := session.Query(sql).Exec()
	if err != nil {
		logger.Error(err)
		return err, false
	}
	return nil, true
}

func insert_data(session *gocql.Session, l Location) (error, bool) {
	// insert a location into the locations table
	is_dse := true
	var sql string
	if is_dse {
		sql = "INSERT INTO locations (id, point, metadata) VALUES (uuid(), POINT(? ?), ?);"
	} else {
		sql = "INSERT INTO locations (id, latitude, longitude, metadata) VALUES (uuid, ?, ?, ?);"
	}
	// execute the query
	err := session.Query(sql, l.Latitude, l.Longitude, l.Label).Exec()

	if err != nil {
		logger.Error(err)
		return err, false
	}
	return nil, true
	// return the id of the inserted location

}
