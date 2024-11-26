package models
import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)
type Referral struct {
	gorm.Model
	ID         int             `gorm:"primaryKey"`
	UserID     int             `gorm:"not null"`       // ผู้แนะนำ (User ID)
	RefereeID  int             `gorm:"not null"`       // ผู้ถูกแนะนำ (User ID)
	Turnover   decimal.Decimal `gorm:"type:decimal(10,2);default:0"` // ยอดเล่นของ Referee
	Commission decimal.Decimal `gorm:"type:decimal(10,2);default:0"` // ค่าคอมมิชชั่นที่ได้
}

func (m *Referral) TableName() string {
	return "Referrals"
}