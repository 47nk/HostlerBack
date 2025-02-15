package dashboard

import (
	"time"
)

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	FirstName    string    `gorm:"type:text;not null"`
	LastName     string    `gorm:"type:text;not null"`
	Username     string    `gorm:"type:text;not null"`
	MobileNumber string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"` // Matches 'timestamp with time zone' type
	UpdatedAt    time.Time `gorm:"type:timestamptz"`
}

func (User) TableName() string {
	return "users"
}

type Bill struct {
	ID            int64   `gorm:"primaryKey;autoIncrement"`
	UserId        int64   `gorm:"numeric;not null"`
	Amount        float64 `gorm:"numeric;not null"`
	BillType      string  `gorm:"type:text;not null"`
	BillingMonth  string  `gorm:"type:text;not null"`
	PaymentStatus string  `gorm:"type:text;not null"`
}

func (Bill) TableName() string {
	return "bills"
}

type Transaction struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	CreatedAt       time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	BillId          int64     `gorm:"numeric;not null"`
	Price           float64   `gorm:"numeric;not null"`
	Items           int64     `gorm:"numeric;not null"`
	ExtraPrice      float64   `gorm:"numeric;not null"`
	ExtraItems      int64     `gorm:"numeric;not null"`
	Description     string    `gorm:"type:text;not null"`
	TransactionType string    `gorm:"size:50;not null"`
}

func (Transaction) TableName() string {
	return "transactions"
}
