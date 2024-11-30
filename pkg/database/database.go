package database

import (
	"fmt"
	"log"
	"os"

	"github.com/VarthanV/pub-sub/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dbName string, doMigrations bool) (*gorm.DB, error) {
	l := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			Colorful: true,
			LogLevel: logger.Info,
		})

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", dbName)),
		&gorm.Config{
			Logger: l,
		})
	if err != nil {
		logrus.Error("unable to open db ", err)
		return nil, err
	}

	if doMigrations {
		err = db.AutoMigrate(
			&models.Exchange{},
			&models.Queue{},
			&models.Binding{})
		if err != nil {
			logrus.Error("unable to do migrations")
			return nil, err
		}
	}

	return db, nil
}
