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
	Name string `gorm:"size:255;not null;" json:"name"`
	PercentDiscount decimal.Decimal `gorm:"type:numeric" json:"percentDiscount"`
	MaxDiscount decimal.Decimal `gorm:"type:numeric" json:"maxDiscount"`
	UsageLimit string `gorm:"size:255;not null;" json:"usageLimit"`
	SpecificTime string `gorm:"size:255;not null;" json:"specificTime"`
	PaymentMethod string `gorm:"size:255;not null;" json:"paymentMethod"`
	MinSpend decimal.Decimal `gorm:"type:numeric" json:"minSpend"`
	MaxSpend decimal.Decimal `gorm:"type:numeric" json:"maxSpend"`
	TermsAndConditions string `gorm:"size:255;not null;" json:"termsAndConditions"`
	// Deviceid      string `gorm:"type:text"json:"deviceid"`
	// ID            uint  `json:"id"`
	// Level         int  `json:"level"`
	// Pin           string `gorm:"type:text" json:"pin"`
	Status        int  `json:"status"`
}
