package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testdb *sql.DB

const (
	drivername     = "postgres"
	dataSourceName = "postgresql://udit:root@127.0.0.1:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {

	var err error

	if testdb, err = sql.Open(drivername, dataSourceName); err != nil {
		log.Fatal("cannot connect to the database")
	}
	testQueries = New(testdb)
	os.Exit(m.Run())
}
