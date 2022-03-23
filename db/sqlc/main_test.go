package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/yanshen1997/simplebank/util"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.Load("../..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not open db:", err)
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
