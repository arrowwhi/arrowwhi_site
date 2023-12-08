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
		//go db.processQueue()
	}
	return db
}

func (dbe *DbEngine) ConnectToDB(dsn string) {
	var err error
	dbe.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Connect to database error! \n", err.Error())
	}
}

func (dbe *DbEngine) CreateClientsTable() {
	if err := dbe.db.AutoMigrate(&Client{}); err != nil {
		log.Error(err.Error())
	}
	if err := dbe.db.AutoMigrate(&Message{}); err != nil {
		log.Error(err.Error())
	}
	if err := dbe.db.AutoMigrate(&FeedbackType{}); err != nil {
		log.Error(err.Error())
	}
	if err := dbe.db.AutoMigrate(&Feedback{}); err != nil {
		log.Error(err.Error())
	}

	feedbackTypes := []FeedbackType{
		{ID: 1, Type: "revision"},
		{ID: 2, Type: "bug"},
	}
	if out := dbe.db.Create(&feedbackTypes); out.Error != nil {
		log.Print(out.Error)
	}
}
