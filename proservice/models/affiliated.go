package models
import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Affiliate struct {
	gorm.Model
	ID          int             `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`  // รหัส Affiliate
	PartnerID   int             `gorm:"column:partner_id;NOT NULL"`                   // รหัส Partner (FK)
	UserID      int             `gorm:"column:user_id;NOT NULL"`                      // รหัส User (FK)
	Turnover    decimal.Decimal `gorm:"type:decimal(10,2);column:turnover;default:0"` // ยอด Turnover ของผู้ใช้
	Commission  decimal.Decimal `gorm:"type:decimal(10,2);column:commission;default:0"` // ค่าคอมมิชชั่นที่ได้จาก Turnover
	ReferralURL string          `gorm:"type:varchar(255);column:referral_url"`        // ลิงค์ที่ใช้ (เช่น https://example.com?partner_id=123)
	Status      string          `gorm:"type:varchar(50);column:status;default:'active'"` // สถานะการเชื่อมโยง (active/inactive)
	CreatedAt   time.Time       `gorm:"column:created_at;NOT NULL"`                   // วันที่สร้าง
	UpdatedAt   time.Time       `gorm:"column:updated_at;NOT NULL"`                   // วันที่อัปเดตล่าสุด
	Partner Partner `gorm:"foreignKey:PartnerID;references:ID"` // ความสัมพันธ์กับ Partner
	User    Users   `gorm:"foreignKey:UserID;references:ID"`    // ความสัมพันธ์กับ Users
}



func (m *Affiliate) TableName() string {
	return "Affiliates"
}