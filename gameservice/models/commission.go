package models
import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Commission struct {
	gorm.Model
	PartnerID int             `gorm:"column:partner_id;NOT NULL"`    // Partner ที่ได้รับค่าคอมมิชชั่น
	UserID    int             `gorm:"column:user_id;NOT NULL"`       // ผู้เล่นที่สร้าง transaction
	Amount    decimal.Decimal `gorm:"type:decimal(10,2);column:amount;NOT NULL"` // ค่าคอมมิชชั่น
	GameID    string          `gorm:"type:varchar(255);column:game_id"` // รหัสเกมที่เกี่ยวข้อง
}

func (m *Commission) TableName() string {
	return "Commission"
}