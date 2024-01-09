package main 

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"Simple-Bank/api"
	db "Simple-Bank/db/sqlc"
	"Simple-Bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to open database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err!= nil {
		log.Fatal("Failed to start server", err)
	}

}