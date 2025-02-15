package login

import (
	"time"
)

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	FirstName    string    `gorm:"type:text;not null"`
	LastName     string    `gorm:"type:text;not null"`
	Username     string    `gorm:"type:text;not null"`
	RoleId       int64     `gorm:"type:numeric;not null"`
	MobileNumber string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamptz"`
	Password     string    `gorm:"type:text;not null"`
	Active       bool      `gorm:"type:boolean"`
	UserRole     Role      `gorm:"foreignKey:RoleId;references:ID"`
}

func (User) TableName() string {
	return "users"
}

type Role struct {
	ID          uint      `gorm:"primaryKey"`
	Role        string    `gorm:"type:varchar(100);not null"`
	Active      bool      `gorm:"default:true"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp"`
	CreatedBy   *uint     `gorm:"index;constraint:OnDelete:Set NULL;"`
	UpdatedBy   *uint     `gorm:"index;constraint:OnDelete:Set NULL;"`
}

func (Role) TableName() string {
	return "roles"
}
