package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type EntityTransaction struct {
	ID           string `gorm:"primaryKey;"`
	Description  string `gorm:"type:varchar(255);not null"`
	Status       int64  `gorm:"type:varchar(255);"`
	Amount       int64  `gorm:"type:int;"`
	UserID       string `gorm:"type:varchar(255);not null"`
	WalletID     string `gorm:"type:varchar(255);not null"`
	Type         string `gorm:"type:varchar(255);"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	EntityWallet EntityWallet `gorm:"foreignKey:WalletID"`
	EntityUsers  EntityUsers  `gorm:"foreignKey:UserID"`
}

func (entity *EntityTransaction) BeforeCreate(db *gorm.DB) error {
	entity.ID = uuid.New().String()
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityTransaction) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
