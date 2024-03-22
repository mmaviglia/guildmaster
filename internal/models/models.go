package models

import (
	"guildmaster/internal/config"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var namespace uuid.UUID = uuid.Must(uuid.FromString("19bd48dc-b100-450b-b6ed-d0022e4f23d5"))

// Establish a connection to the database and run all migrations.
func SetupDB() error {

	// Attempt to establish the database connection, retrying if necessary
	attempts := 0
	for {
		attempts++

		db, err := gorm.Open(postgres.Open(config.DB_DSN), &gorm.Config{})
		if err != nil {
			if attempts < 10 {
				time.Sleep(time.Second)
				continue
			}
			return err
		}

		DB = db
		break
	}

	if err := migrateDatabase(); err != nil {
		return err
	}

	return nil
}

// Run migrations for all models.
func migrateDatabase() error {
	// Add all database models to this call
	return DB.AutoMigrate(GuildActivity{})
}

// Calculate the offset using the limit and page.
func offset(limit, page int) int {
	return (page - 1) * limit
}
