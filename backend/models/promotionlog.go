package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)
 

type PromotionLog struct {
	gorm.Model
	ID            int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`
	CreatedOn     string `gorm:"size:255;not null;" json:"CreatedOn"`
	Promotioncode   string `gorm:"size:255;not null;" binding:"required" json:"promotioncode"`
	Promotionname   string `gorm:"size:255;not null;" binding:"required" json:"promotionname"`
	Promoamount decimal.Decimal   `gorm:"type:numeric"  json:"Promoamount"`
	UserID      int  `gorm:"column:userid;"gorm:"not null;" json:"userid"`
	StatementID int  `gorm:"column:statementid;"json:"statementid"`
	WalletID    int  `gorm:"column:walletid;"gorm:"not null;" json:"walletid"`
	Transactionamount decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:transactionamount;NOT NULL"`
	Beforebalance     decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:beforebalance"`
	Proamount         decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:proamount"`
	AddOnamount       decimal.Decimal    `gorm:"column:addonamount;"gorm:"type:numeric(10,2);gorm:"column:addonamount"`
	Balance           decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:balance"`
	Status        int  `json:"status"`
	Example string `gorm:"size:255;not null;" json:"example"`
	Uid               string    `gorm:"type:varchar(255);column:uid;"` 
 
}


func (m *PromotionLog) TableName() string {
	return "PromotionLog"
}
