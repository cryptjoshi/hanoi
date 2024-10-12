package models
import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)


type BanksAccount struct {
	gorm.Model
	ID            int       `gorm:"column:id;NOT NULL"`
	Accountname   string    `gorm:"type:varchar(255)";gorm:"column:accountname;NOT NULL"`
	Accountnumber string    `gorm:"type:varchar(20)";gorm:"column:accountnumber;NOT NULL"`
	Bankname      string    `gorm:"type:varchar(150)";gorm:"column:bankname"`
	Bankid        string    `gorm:"type:varchar(4)";gorm:"column:bankid"`
	Balance       decimal.Decimal   `gorm:"column:balance;default:0"`
	Level         int       `gorm:"column:level"`
	Status        int       `gorm:"column:status;NOT NULL"`
	Active        int       `gorm:"column:active;NOT NULL"`
	CreatedOn     time.Time `gorm:"column:CreatedOn;NOT NULL"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;NOT NULL"`
	DeletedAt     time.Time `gorm:"column:deletedAt"`
	Deviceid      string    `gorm:"type:varchar(255)";gorm:"column:deviceid"`
	Pin           string    `gorm:"type:varchar(6)";gorm:"column:pin"`
}

func (m *BanksAccount) TableName() string {
	return "BanksAccount"
}
