package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	dsn := " REMOVED"
	//open connection to db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	err = db.Exec("DEALLOCATE ALL").Error
	if err != nil {
		return nil, fmt.Errorf("failed to deallocate prepared statements: %w", err)
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
