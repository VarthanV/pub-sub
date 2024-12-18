package main

import (
	"context"

	"github.com/VarthanV/pub-sub/broker"
	"github.com/VarthanV/pub-sub/pkg/config"
	"github.com/VarthanV/pub-sub/pkg/database"
	"github.com/VarthanV/pub-sub/server"
	"github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatal("unable to load config ", err)
	}

	db, err := database.Init(cfg.Database.Name,
		cfg.Database.DoMigrations)
	if err != nil {
		logrus.Fatal("error in initializing db ", err)
	}

	// Init broker
	b := broker.New(db, cfg)

	controller := server.NewController(db, b)

	logrus.Info("Starting broker.....")
	err = b.Start(context.Background())
	if err != nil {
		logrus.Fatal("error in starting broker ", err)
	}

	logrus.Info(`
  _       _      __      _  
 |_) | | |_) __ (_  | | |_) 
 |   |_| |_)    __) |_| |_) 
                            
`)
	server.InitServer(cfg.Server.HTTPPort, controller)

}
