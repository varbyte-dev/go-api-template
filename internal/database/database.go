package database

import (
	"log"

	"go-api-template/internal/config"
	"go-api-template/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	logLevel := logger.Silent
	if config.App.AppEnv == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(sqlite.Open(config.App.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("Database connected")

	migrate()
}

func migrate() {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database migrated")
}
