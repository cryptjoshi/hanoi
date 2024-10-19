package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID            int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`
	Name         string          `gorm:"size:255;not null;" json:"name" binding:"required"`
	SerialNumber string          `gorm:"size:255;not null;unique" json:"serialNumber" binding:"required"`
	Quantity     uint            `json:"quantity" binding:"required"`
	Price        decimal.Decimal `gorm:"type:numeric" json:"price" binding:"required"`
	Description  string          `gorm:"type:text" json:"description" binding:"required"`
}