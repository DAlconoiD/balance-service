package main

import (
	"flag"
	"github.com/DAlconoiD/balance-service/server"
	"github.com/DAlconoiD/balance-service/storage"
	"github.com/DAlconoiD/balance-service/utils"
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
	log.Fatal(s.Start(config.ServerAddress))
}
