package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type EntityWallet struct {
	ID          string `gorm:"type:varchar(255);primaryKey;"`
	UserID      string `gorm:"type:varchar(255);not null"`
	Balance     int64  `gorm:"type:int"`
	Active      bool   `gorm:"type:bool;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	EntityUsers EntityUsers `gorm:"foreignKey:UserID"`
}

func (entity *EntityWallet) BeforeCreate(db *gorm.DB) error {
	entity.ID = uuid.New().String()
	entity.CreatedAt = time.Now().Local()
	entity.Active = true
	return nil
}

func (entity *EntityWallet) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
