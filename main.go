package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yanshen1997/simplebank/api"
	db "github.com/yanshen1997/simplebank/db/sqlc"
	"github.com/yanshen1997/simplebank/util"
)

func main() {
	config, err := util.Load(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not open db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("can not start server:", err)
	}
}
