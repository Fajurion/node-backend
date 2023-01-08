package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"node-backend/entities/account"
)

var DBConn *gorm.DB

func Connect() {
	url := "host=localhost user=postgres password=deinemutter123 dbname=chat port=5432"

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Something went wrong during the connection with the database.", err)
	}

	log.Println("Successfully connected to the database.")
	db.Logger = logger.Default.LogMode(logger.Info)

	// Migrate the schema
	db.AutoMigrate(&account.Account{})

	// Assign the database to the global variable
	DBConn = db
}
