package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	sqlDriver     = "postgres"
	sqlSourceName = "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(sqlDriver, sqlSourceName)
	if err != nil {
		log.Fatal("can not open db:", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
