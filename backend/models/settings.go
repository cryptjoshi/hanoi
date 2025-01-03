package models

import "gorm.io/gorm"

type Settings struct {
	gorm.Model
	// CompanyName string `gorm:"size:255;column:companyCode" json:"companyNamee"`
	// BaseCurrency string `gorm:"size:255;column:baseCurrency" json:"baseCurrency"`
	// CustomerCurrency string `gorm:"size:255;column:customerCurrency" json:"customerCurrency"`	
	// BaseRate float64 `gorm:"column:baseRate" json:"baseRate"`
	// CustomerRate float64 `gorm:"column:customerRate" json:"customerRate"`
	// Active int `gorm:"column:active" json:"active"`
	// Remark string `gorm:"size:255;column:remark" json:"remark"`
	Key string `gorm:"size:255;column:key" json:"key"`
	Value string `gorm:"size:255;column:value" json:"value"`
}


func (Settings) TableName() string {
	return "Settings"
}
