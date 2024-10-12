package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	CreatedOn     string `gorm:"size:255;not null;" json:"CreatedOn"`
	Promotioncode   string `gorm:"size:255;not null;" binding:"required" json:"promotioncode"`
	Promotionname   string `gorm:"size:255;not null;" binding:"required" json:"promotionname"`
	Active        int  `json:"active"`
	Freecredit      string  `gorm:"size:255;not null;" json:"freecredit"`
	Freecreditmax   string   `gorm:"size:255;not null;" json:"freecreditmax"`
	Maxtime    decimal.Decimal   `gorm:"type:numeric"  json:"maxtime"`
	Typetime   string   `gorm:"size:255;not null;" json:"typetime"`
	Days   string   `gorm:"size:255;not null;" json:"days"`
	Includegames   string   `gorm:"size:255;not null;" json:"includegames"`
	Excludegames   string   `gorm:"size:255;not null;" json:"excludegames"`
	Turnamount   string   `gorm:"size:255;not null;" json:"turnamount"`
	Widthdrawmax       decimal.Decimal   `gorm:"type:numeric" json:"Widthdrawmax"`
    Formular   string   `gorm:"size:255;not null;" json:"formular"`
	Options   string   `gorm:"size:255;not null;" json:"options"`
	// Bankid        string `gorm:"size:255;not null;" json:"bankid"`
	// Bankname      string `gorm:"size:255;not null;" json:"bankname"`
	DeletedAt     string `gorm:"size:255;not null;" json:"deletedAt"`
	// Deviceid      string `gorm:"type:text"json:"deviceid"`
	// ID            uint  `json:"id"`
	// Level         int  `json:"level"`
	// Pin           string `gorm:"type:text" json:"pin"`
	Status        int  `json:"status"`
}
