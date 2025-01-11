package announcement

import "time"

type Announcement struct {
	ID          uint      `gorm:"primaryKey" json:"id,omitempty"`
	Title       string    `gorm:"size:255;not null" json:"title,omitempty"`
	Type        string    `gorm:"size:100;not null" json:"type,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedBy   uint      `json:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"active,omitempty"`
}

func (Announcement) TableName() string {
	return "announcements"
}
