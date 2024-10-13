package login

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`                   // Matches 'bigint' type in PostgreSQL
	FirstName string         `gorm:"type:text;not null"`                         // Matches 'text' type in PostgreSQL
	CreatedAt time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"` // Matches 'timestamp with time zone' type
	UpdatedAt time.Time      `gorm:"type:timestamptz"`                           // Matches 'timestamp with time zone' type
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index"`                     // Matches 'timestamp with time zone' type
}

func (User) TableName() string {
	return "users"
}
