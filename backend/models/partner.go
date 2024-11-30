package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/shopspring/decimal"
)

type Partner struct {
	gorm.Model
	ID              int             `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"` // รหัส Partner
	Name            string          `gorm:"type:varchar(255);column:name;NOT NULL"`      // ชื่อ Partner
	Email           string          `gorm:"type:varchar(255);column:email;unique;NOT NULL"` // อีเมล
	Phone           string          `gorm:"type:varchar(50);column:phone"`              // เบอร์โทร
	AffiliateKey    string          `gorm:"type:varchar(50);unique;NOT NULL"`           // รหัส Affiliate ที่ใช้แนบในลิงค์
	CommissionRate  decimal.Decimal `gorm:"type:decimal(5,2);column:commission_rate"`   // เปอร์เซ็นต์ค่าคอมมิชชั่น เช่น 10.00%
	TotalCommission decimal.Decimal `gorm:"type:decimal(10,2);column:total_commission"` // ยอดค่าคอมมิชชั่นสะสม
	Status          string          `gorm:"type:varchar(50);column:status;default:'active'"` // สถานะ (active/inactive)
	CreatedAt       time.Time       `gorm:"column:created_at;NOT NULL"`                 // วันที่สร้าง
	UpdatedAt       time.Time       `gorm:"column:updated_at;NOT NULL"`                 // วันที่อัปเดตล่าสุด
	Affiliates []Affiliate `gorm:"foreignKey:PartnerID"` // ความสัมพันธ์กับ Affiliate
	TotalEarnings   decimal.Decimal `gorm:"type:decimal(10,2);default:0"` // ค่าคอมมิชชั่นสะสม

}

type PartnerResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	AffiliateKey string `json:"affiliatekey"`
	CommissionRate decimal.Decimal `json:"commissionrate"`
	TotalCommission decimal.Decimal `json:"totalcommission"`
	Status string `json:"status"`
}

func (Partner) TableName() string {
	return "Partners"
}
