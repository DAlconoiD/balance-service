package main

import (
	"flag"
	"github.com/dalconoid/balance-service/server"
	"github.com/dalconoid/balance-service/storage"
	"github.com/dalconoid/balance-service/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to application config file")
	flag.Parse()

	config, err := utils.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	db := &storage.Database{ConnString: config.DBConnectionString, PaginationNum: config.PaginationNumber}
	err = db.Open()
	if err != nil {
		log.Fatal(err)
	}
	s := server.New()
	s.ConfigureRouter(db)
	log.Infof("Starting server on %s", config.ServerAddress)
	log.Fatal(s.Start(config.ServerAddress))
}
