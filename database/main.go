package database

import (
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbEngine struct {
	db *gorm.DB
}

var db *DbEngine

func Get() *DbEngine {
	if db == nil {
		db = &DbEngine{}
		go db.processQueue()
	}
	return db
}

func (dbe *DbEngine) ConnectToDB(dsn string) {
	var err error
	//dsn := "user=golang password=golang dbname=file_library host=localhost port=5432 sslmode=disable TimeZone=Europe/Moscow"
	dbe.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Connect to database error! \n", err.Error())
	}
}

func (dbe *DbEngine) CreateClientsTable() {
	err := dbe.db.AutoMigrate(&Client{})
	if err != nil {
		log.Error(err.Error())
	}
	err = dbe.db.AutoMigrate(&Message{})
	if err != nil {
		log.Error(err.Error())
	}

}
