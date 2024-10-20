package login

import (
	"time"
)

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"` // Matches 'bigint' type in PostgreSQL
	FirstName    string    `gorm:"type:text;not null"`       // Matches 'text' type in PostgreSQL
	LastName     string    `gorm:"type:text;not null"`
	RollNumber   string    `gorm:"type:text;not null"`
	RoleId       int64     `gorm:"type:numeric;not null"`
	MobileNumber string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"` // Matches 'timestamp with time zone' type
	UpdatedAt    time.Time `gorm:"type:timestamptz"`                           // Matches 'timestamp with time zone' type
	Password     string    `gorm:"type:text;not null"`
	Active       bool      `gorm:"type:boolean"`
}

func (User) TableName() string {
	return "users"
}
