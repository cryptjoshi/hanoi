package models

import (
	//"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type TsxAdmin struct {
	gorm.Model
	ID            int       `gorm:"column:id;NOT NULL"`
	Username      string    `gorm:"type:varchar(100)";gorm:"column:username;NOT NULL"`
	Password      string    `gorm:"type:varchar(100)";gorm:"column:password;NOT NULL"`
	Preferredname string    `gorm:"type:varchar(255)";gorm:"column:preferredname"`
	Token         string    `gorm:"type:text";gorm:"column:token"`
	Role          string    `gorm:"type:varchar(20)";gorm:"column:role"`
	Salt          string    `gorm:"type:text";gorm:"column:salt;NOT NULL"`
	CreatedAt     time.Time `gorm:"column:createdAt;NOT NULL"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;NOT NULL"`
	DeletedAt     time.Time `gorm:"column:deletedAt"`
	Prefix        string    `gorm:"type:varchar(50)";gorm:"column:prefix"`
	Secret        string    `gorm:"type:varchar(255)";gorm:"column:secret"`
	Status        int       `gorm:"column:status;default:0"`
}

func (m *TsxAdmin) TableName() string {
	return "TsxAdmin"
}