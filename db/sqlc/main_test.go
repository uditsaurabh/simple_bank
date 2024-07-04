package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/uditsaurabh/simple_bank/util"
)

var testQueries *Queries
var testdb *sql.DB
var Config util.Config

func TestMain(m *testing.M) {
	var err error
	//loading config
	Config, err = util.LoadConfig("../..")
	if err != nil {
		log.Println("cannot load config")
	}

	if testdb, err = sql.Open(Config.DBDriver, Config.DBSource); err != nil {
		log.Fatal("cannot connect to the database")
	}
	testQueries = New(testdb)
	os.Exit(m.Run())
}
