package main

import (
	"database/sql"

	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/uditsaurabh/simple_bank/api"
	db "github.com/uditsaurabh/simple_bank/db/sqlc"
	"github.com/uditsaurabh/simple_bank/util"
)

func main() {
	var conn *sql.DB
	var err error

	//loading config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Println("cannot log config")
	}
	//connect to database
	if conn, err = sql.Open(config.DBDriver, config.DBSource); err != nil {
		log.Fatal("cannot connect to the database")
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(store, &config)
	if err != nil {
		log.Fatal("cannot create server", err.Error(), len(config.EncryptionKey))
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Println("server exited...")
		os.Exit(-1)
	}
}
