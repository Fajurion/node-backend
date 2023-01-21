package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/entities/app"
	"node-backend/entities/app/projects"
	"node-backend/entities/node"
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
	db.AutoMigrate(&properties.AccountSetting{})

	// Migrate node related tables
	db.AutoMigrate(&node.Node{})

	// Migrate app related tables
	db.AutoMigrate(&app.App{})
	db.AutoMigrate(&app.AppNode{})
	db.AutoMigrate(&app.AppSetting{})

	// Migrate project related tables
	db.AutoMigrate(&projects.Project{})
	db.AutoMigrate(&projects.Container{})
	db.AutoMigrate(&projects.Event{})
	db.AutoMigrate(&projects.Member{})

	// Assign the database to the global variable
	DBConn = db
}
