package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	CreatedOn     string `gorm:"size:255;not null;" json:"CreatedOn"`
	Accountname   string `gorm:"size:255;not null;" binding:"required" json:"accountname"`
	Accountnumber string `gorm:"size:255;not null;" json:"accountnumber"`
	Active        int  `json:"active"`
	Balance       decimal.Decimal   `gorm:"type:numeric" json:"balance"`
	Bankid        string `gorm:"size:255;not null;" json:"bankid"`
	Bankname      string `gorm:"size:255;not null;" json:"bankname"`
	DeletedAt     string `gorm:"size:255;not null;" json:"deletedAt"`
	Deviceid      string `gorm:"type:text"json:"deviceid"`
	ID            uint  `json:"id"`
	Level         int  `json:"level"`
	Pin           string `gorm:"type:text" json:"pin"`
	Status        int  `json:"status"`
}
