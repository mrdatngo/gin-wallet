package model

import (
	"time"

	"github.com/google/uuid"
	util "github.com/mrdatngo/gin-wallet/utils"
	"gorm.io/gorm"
)

type EntityUsers struct {
	ID        string `gorm:"type:varchar(255);primaryKey;"`
	Fullname  string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(255);unique;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Active    bool   `gorm:"type:bool;default:false"`
	RoleID    int64  `gorm:"type:uint;"`
	MetaData  string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (entity *EntityUsers) BeforeCreate(db *gorm.DB) error {
	entity.ID = uuid.New().String()
	entity.Password = util.HashPassword(entity.Password)
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityUsers) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
