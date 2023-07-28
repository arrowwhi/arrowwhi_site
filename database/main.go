package database

import (
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbEngine struct {
	db *gorm.DB
}

var db *dbEngine

func GetDb() *dbEngine {
	if db == nil {
		db = &dbEngine{}
	}
	return db
}

func (dbe *dbEngine) ConnectToDB(dsn string) {
	var err error
	//dsn := "user=golang password=golang dbname=file_library host=localhost port=5432 sslmode=disable TimeZone=Europe/Moscow"
	dbe.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Connect to database error! \n", err.Error())
	}
}

func (dbe *dbEngine) CreateClientsTable() {
	// AutoMigrate создает таблицу "clients" согласно структуре Client.
	err := dbe.db.AutoMigrate(&Client{})
	if err != nil {
		log.Error(err.Error())
	}
}
