package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

// TODO move const to env file
const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
