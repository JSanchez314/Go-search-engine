package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConn holds the connection to the database and is shared across the application.
var DBConn *gorm.DB

// InitDB initializes the database connection, sets up required extensions, and performs auto-migration.
func InitDB() {
	// Get the database URL from the environment variable.
	dburl := os.Getenv("DATABASE_URL")

	var err error
	// Open a connection to the PostgreSQL database using GORM.
	DBConn, err = gorm.Open(postgres.Open(dburl), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		panic("failed to connect to database")
	}

	// Install the 'uuid-ossp' extension if it's not already installed.
	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		fmt.Println("Cannot install uuid extension:", err)
		panic(err)
	}

	// Automatically migrate the database schema for the specified models.
	err = DBConn.AutoMigrate(&User{}, &SearchSetting{}, CrawledUrl{}, SearchIndex{})
	if err != nil {
		panic(err)
	}
}

// GetDB returns the current database connection instance.
func GetDB() *gorm.DB {
	return DBConn
}
