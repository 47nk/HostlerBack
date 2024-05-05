package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	// Connection string
	dsn := "user=username password=password dbname=mydatabase sslmode=disable"

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
