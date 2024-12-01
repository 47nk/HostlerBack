package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	dsn := "postgresql://postgres.ixuwoesbfkzrwlzxtsic:Naresh@007@aws-0-ap-south-1.pooler.supabase.com:6543/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	fmt.Println("Successfully connected to the PostgreSQL database!")
	return db, nil
}

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	FirstName string         `gorm:"unique;not null"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
