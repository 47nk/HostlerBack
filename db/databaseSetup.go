package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDB() (*gorm.DB, error) {
	dsn := "host=127.0.0.1 user=postgres password=Naresh@007 dbname=manage port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable detailed logging
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")
	return db, nil
}

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	FirstName string         `gorm:"unique;not null"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
