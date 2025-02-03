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
}

func (Announcement) TableName() string {
	return "announcements"
}

type Entity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Address   string    `gorm:"type:text" json:"address"`
	CreatedBy uint      `gorm:"foreignKey:CreatedBy;references:ID" json:"created_by"`
	UpdatedBy uint      `gorm:"foreignKey:UpdatedBy;references:ID" json:"updated_by"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Active    bool      `gorm:"default:true" json:"active"`
	Channels  []Channel `gorm:"foreignKey:EntityID;references:ID"`
}

func (Entity) TableName() string {
	return "entities"
}

type Channel struct {
	ID          uint      `gorm:"primaryKey"`                 // Primary key
	EntityID    uint      `gorm:"not null"`                   // The entity (hostel, office, etc.) this channel belongs to
	Name        string    `gorm:"type:varchar(255);not null"` // Name of the channel
	Type        string    `gorm:"type:varchar(50);not null"`  // Type of channel (e.g., 'general', 'announcements', 'events')
	Description string    `gorm:"type:text"`                  // Optional description of the channel
	CreatedBy   uint      `gorm:"index;default:NULL"`         // User who created the channel (nullable, references users table)
	UpdatedBy   uint      `gorm:"index;default:NULL"`         // User who last updated the channel (references users table)
	CreatedAt   time.Time `gorm:"autoCreateTime"`             // Timestamp when the channel was created
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`             // Last updated timestamp

}

func (Channel) TableName() string {
	return "channels"
}
