package models

import "gorm.io/gorm"

type Games struct {
	gorm.Model
	ID            int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`
	ProductCode string `gorm:"size:255;column:productCode" json:"productCode"`
	Product string `gorm:"size:255;column:product" json:"product"`
	GameType string `gorm:"size:255;column:gameType" json:"gameType"`
	Active int `gorm:"column:active" json:"active"`
	Remark string `gorm:"size:255;column:remark" json:"remark"`
	Position string `gorm:"size:255;column:position" json:"position"`
	Urlimage string `gorm:"size:255;column:urlimage" json:"urlimage"`
	Name string `gorm:"size:255;column:name" json:"name"`
	Status string `gorm:"size:255;column:status" json:"status"`
	
}

func (Games) TableName() string {
	return "Games"
}
