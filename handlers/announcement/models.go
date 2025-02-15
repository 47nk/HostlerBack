package announcement

import "time"

type Announcement struct {
	ID          uint                     `gorm:"primaryKey"`
	Title       string                   `gorm:"size:255;not null" json:"title"`
	Type        string                   `gorm:"size:100;not null" json:"type"`
	Description string                   `gorm:"type:text" json:"description"`
	ChannelId   int                      `gorm:"not null" json:"channel_id"`
	CreatedBy   uint                     `json:"created_by,omitempty"`
	UpdatedBy   uint                     `json:"updated_by,omitempty"`
	CreatedAt   time.Time                `json:"created_at,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at,omitempty"`
	Creator     User                     `gorm:"foreignKey:CreatedBy;references:ID"`
	Attachments []AnnouncementAttachment `gorm:"foreignKey:AnnouncementID;references:ID" json:"attachments"`
}

func (Announcement) TableName() string {
	return "announcements"
}

type AnnouncementAttachment struct {
	ID             uint      `gorm:"primaryKey"`
	AnnouncementID uint      `gorm:"index" json:"announcement_id"`
	FilePath       string    `gorm:"size:255" json:"file_path"`
	FileType       string    `gorm:"size:50" json:"file_type"`
	FileSize       int64     `json:"file_size"`
	UploadedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"uploaded_at"`
	Active         bool      `gorm:"default:true" json:"active"`
}

func (AnnouncementAttachment) TableName() string {
	return "announcement_attachments"
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
	ID          uint      `gorm:"primaryKey"`
	EntityID    uint      `gorm:"not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Type        string    `gorm:"type:varchar(50);not null"`
	Description string    `gorm:"type:text"`
	Active      bool      `gorm:"default:true" json:"active"`
	CreatedBy   uint      `gorm:"index;default:NULL"`
	UpdatedBy   uint      `gorm:"index;default:NULL"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Channel) TableName() string {
	return "channels"
}

type User struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	FirstName    string `gorm:"type:text;not null"`
	LastName     string `gorm:"type:text;not null"`
	Username     string `gorm:"type:text;not null"`
	RoleId       int64  `gorm:"type:numeric;not null"`
	MobileNumber string `gorm:"type:text;not null"`
}

func (User) TableName() string {
	return "users"
}
