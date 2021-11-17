package database

import (
	"log"

	"github.com/gocql/gocql"
)

type DBConnection struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

var conn DBConnection

func SetupDBConn() {
	conn.cluster = gocql.NewCluster("127.0.0.1")
	conn.cluster.Consistency = gocql.Quorum
	conn.cluster.Keyspace = "std"
	conn.session, _ = conn.cluster.CreateSession()
}

func ExecuteQ(query string, values ...interface{}) error {

	err := conn.session.Query(query).Bind(values...).Exec()

	if err != nil {
		log.Fatal(err)
	}
	return err

}

func ProceedQ(query string, values ...interface{}) gocql.Scanner {

	Scanner := conn.session.Query(query).Bind(values...).Iter().Scanner()

	return Scanner

}
