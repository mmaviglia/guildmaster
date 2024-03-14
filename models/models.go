package models

import (
	"guildmaster/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Establish a connection to the database and run all migrations.
func SetupDB() error {
	db, err := gorm.Open(postgres.Open(config.DB_DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	if err := migrateDatabase(); err != nil {
		return err
	}

	return nil
}

// Run migrations for all models.
func migrateDatabase() error {
	// Add all database models to this call
	return DB.AutoMigrate()
}

func offset(limit, page int) int {
	return (page - 1) * limit
}
