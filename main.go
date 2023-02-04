package main

import (
	"database/sql"
	"log"

	"github.com/cyhe50/simple_bank/api"
	db "github.com/cyhe50/simple_bank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress  = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
