package models
import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type AffiliateLog struct {
	gorm.Model
	ID                int             `gorm:"column:id;primaryKey;autoIncrement"`
	AffiliateUserID   int             `gorm:"column:affiliate_user_id;NOT NULL"`           // ID ของสมาชิกที่แนะนำ
	ReferredUserID    int             `gorm:"column:referred_user_id;NOT NULL"`            // ID ของสมาชิกที่ถูกแนะนำ
	TransactionID     int             `gorm:"column:transaction_id;NOT NULL"`              // ID ของรายการธุรกรรม
	TurnoverAmount    decimal.Decimal `gorm:"type:decimal(15,2);column:turnover_amount;NOT NULL"` // Turnover ที่เกิดขึ้น
	CommissionAmount  decimal.Decimal `gorm:"type:decimal(15,2);column:commission_amount;NOT NULL"` // ค่าคอมมิชชันจาก Turnover
	CreatedAt         time.Time       `gorm:"column:created_at;NOT NULL"`
}

func (m *AffiliateLog) TableName() string {
	return "AffiliateLog"
}