package login

import (
	"time"
)

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	FirstName    string    `gorm:"type:text;not null"`
	LastName     string    `gorm:"type:text;not null"`
	RollNumber   string    `gorm:"type:text;not null"`
	RoleId       int64     `gorm:"type:numeric;not null"`
	MobileNumber string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamptz"`
	Password     string    `gorm:"type:text;not null"`
	Active       bool      `gorm:"type:boolean"`
}

func (User) TableName() string {
	return "users"
}
