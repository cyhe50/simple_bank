package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cyhe50/simple_bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB
var envConfig util.EnvConfig

func TestMain(m *testing.M) {
	var err error
	envConfig, err = util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load env config: ", err)
	}

	testDB, err = sql.Open(envConfig.DriverName, envConfig.DataSourceName)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
