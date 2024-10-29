package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	ID            int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`
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
	Name string `gorm:"size:255;not null;" json:"name"`
	PercentDiscount decimal.Decimal `gorm:"type:numeric" json:"percentDiscount"`
	MaxDiscount decimal.Decimal `gorm:"type:numeric" json:"maxDiscount"`
	Unit string  `gorm:"size:50;not null;" json:"unit"`
	UsageLimit int `json:"usageLimit"`
	SpecificTime string `gorm:"size:255;not null;" json:"specificTime"`
	PaymentMethod string `gorm:"size:255;not null;" json:"paymentMethod"`
	MinDept decimal.Decimal `gorm:"type:numeric" json:"minDept"`
	MinSpend decimal.Decimal `gorm:"type:numeric" json:"minSpend"`
	MaxSpend decimal.Decimal `gorm:"type:numeric" json:"maxSpend"`
	TermsAndConditions string `gorm:"size:255;not null;" json:"termsAndConditions"`
	Description string `gorm:"size:255;not null;" json:"description"`
	StartDate string `gorm:"size:255;not null;" json:"startDate"`
	EndDate string `gorm:"size:255;not null;" json:"endDate"`
	Status        int  `json:"status"`
	Example string `gorm:"size:255;not null;" json:"example"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}


 
