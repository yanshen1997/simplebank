package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yanshen1997/simplebank/api"
	db "github.com/yanshen1997/simplebank/db/sqlc"
)

const (
	sqlDriver     = "postgres"
	sqlSourceName = "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable"
	address       = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(sqlDriver, sqlSourceName)
	if err != nil {
		log.Fatal("can not open db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(*store)

	if err = server.Start(address); err != nil {
		log.Fatal("can not start server:", err)
	}
}
