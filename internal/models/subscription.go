package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ServiceName string         `gorm:"type:varchar(255);not null" json:"service_name"`
	Price       int            `gorm:"type:integer;not null" json:"price"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	StartDate   time.Time      `gorm:"type:date;not null;index" json:"start_date"`
	EndDate     *time.Time     `gorm:"type:date;index" json:"end_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
