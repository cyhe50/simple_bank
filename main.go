package main

import (
	"database/sql"
	"log"

	"github.com/cyhe50/simple_bank/api"
	db "github.com/cyhe50/simple_bank/db/sqlc"
	"github.com/cyhe50/simple_bank/util"

	_ "github.com/lib/pq"
)

func main() {
	envConfig, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env config: ", err)
	}

	conn, err := sql.Open(envConfig.DriverName, envConfig.DataSourceName)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(envConfig.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
