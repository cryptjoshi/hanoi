package models
import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)



type Provider struct {
	gorm.Model
	ID                int       `gorm:"column:id;NOT NULL"`
	Providername          string    `gorm:"type:varchar(255)";gorm:"column:providername"`
	Balance           decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:balance"`
	CreatedAt         time.Time `gorm:"column:createdAt;NOT NULL"`
	UpdatedAt         time.Time `gorm:"column:updatedAt;NOT NULL"`
	DeletionAt    time.Time `gorm:"default:current_timestamp(3)";gorm:"column:deletionAt;NOT NULL"`
	Status            string    `gorm:"column:status;"`
}

// func (m *Provider) TableName() string {
// 	return "Provider"
// }