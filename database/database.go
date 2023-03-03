package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/entities/account/settings"
	"node-backend/entities/app"
	"node-backend/entities/node"
)

var DBConn *gorm.DB

func Connect() {
	url := "host=" + DB_HOST + " user=" + DB_USERNAME + " password=" + DB_PASSWORD + " dbname=" + DB_DATABASE + " port=" + DB_PORT

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatal("Something went wrong during the connection with the database.", err)
	}

	log.Println("Successfully connected to the database.")

	// Configure the database driver
	driver, _ := db.DB()

	driver.SetMaxIdleConns(10)
	driver.SetMaxOpenConns(100)
	driver.SetConnMaxLifetime(time.Hour)

	// Migrate the schema

	// Migrate account related tables
	db.AutoMigrate(&account.Account{})
	db.AutoMigrate(&account.Authentication{})
	db.AutoMigrate(&account.Session{})
	db.AutoMigrate(&account.Subscription{})
	db.AutoMigrate(&account.Rank{})
	db.AutoMigrate(&properties.Friend{})
	db.AutoMigrate(&settings.AccountSetting{})

	// Migrate node related tables
	db.AutoMigrate(&node.Cluster{})
	db.AutoMigrate(&node.Node{})
	db.AutoMigrate(&node.NodeCreation{})

	// Migrate app related tables
	db.AutoMigrate(&app.App{})
	db.AutoMigrate(&app.AppSetting{})

	// Assign the database to the global variable
	DBConn = db
}
