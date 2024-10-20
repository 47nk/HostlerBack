package dashboard

import (
	"time"
)

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"` // Matches 'bigint' type in PostgreSQL
	FirstName    string    `gorm:"type:text;not null"`       // Matches 'text' type in PostgreSQL
	LastName     string    `gorm:"type:text;not null"`
	RollNumber   string    `gorm:"type:text;not null"`
	MobileNumber string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"` // Matches 'timestamp with time zone' type
	UpdatedAt    time.Time `gorm:"type:timestamptz"`                           // Matches 'timestamp with time zone' type                   // Matches 'timestamp with time zone' type
}

func (User) TableName() string {
	return "users"
}

type Bill struct {
	ID           int64   `gorm:"primaryKey;autoIncrement"` // Matches 'bigint' type in PostgreSQL
	UserId       int64   `gorm:"numeric;not null"`
	Amount       float64 `gorm:"numeric;not null"`
	BillType     string  `gorm:"type:text;not null"`
	BillingMonth string  `gorm:"type:text;not null"`
}

func (Bill) TableName() string {
	return "bills"
}

type Transaction struct {
	ID              int64   `gorm:"primaryKey;autoIncrement"` // Matches 'bigint' type in PostgreSQL
	BillId          int64   `gorm:"numeric;not null"`
	Price           float64 `gorm:"numeric;not null"`
	Items           int64   `gorm:"numeric;not null"`
	ExtraPrice      float64 `gorm:"numeric;not null"`
	ExtraItems      int64   `gorm:"numeric;not null"`
	Description     string  `gorm:"type:text;not null"`
	TransactionType string  `gorm:"size:50;not null"`
}

func (Transaction) TableName() string {
	return "transactions"
}
