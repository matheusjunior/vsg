package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique;not null;index" json:"email"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}

func (r *User) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return nil
}
