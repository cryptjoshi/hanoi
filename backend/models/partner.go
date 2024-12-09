package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/shopspring/decimal"
)

type Partner struct {
	gorm.Model
	ID              int             `gorm:"column:id;primaryKey;autoIncrement;NOT NULL" json:"id"` // รหัส Partner
	Name            string          `gorm:"type:varchar(255);column:name;NOT NULL" json:"name"`      // ชื่อ Partner
	Username        string          `gorm:"index:idx_username,unique;type:varchar(255);column:username;NOT NULL" json:"username"`
	Preferredname    string    `gorm:"type:varchar(250);column:preferredname" json:"preferredname"`
	Password        string          `gorm:"type:text;column:password;NOT NULL" json:"password"`
	Email           string          `gorm:"type:varchar(255);column:email;unique;NOT NULL" json:"email"` // อีเมล
	Phone           string          `gorm:"type:varchar(50);column:phone" json:"phone"`              // เบอร์โทร
	AffiliateKey    string          `gorm:"type:varchar(50);column:affiliatekey;unique;NOT NULL" json:"affiliatekey"`           // รหัส Affiliate ที่ใช้แนบในลิงค์
	CommissionRate  decimal.Decimal  `gorm:"type:decimal(5,2);column:commissionrate" json:"commissionrate"`   // เปอร์เซ็นต์ค่าคอมมิชชั่น เช่น 10.00%
	TotalCommission decimal.Decimal  `gorm:"type:decimal(10,2);column:totalcommission" json:"totalcommission"` // ยอดค่าคอมมิชชั่นสะสม
	Status          string          `gorm:"type:varchar(50);column:status;default:'active'" json:"status"` // สถานะ (active/inactive)
	Prefix          string          `gorm:"type:varchar(50);column:prefix;NOT NULL" json:"prefix"`    
	CreatedAt       time.Time       `gorm:"column:created_at;NOT NULL" json:"created_at"`                 // วันที่สร้าง
	UpdatedAt       time.Time       `gorm:"column:updated_at;NOT NULL" json:"updated_at"`                 // วันที่อัปเดตล่าสุด
	Affiliates      []Affiliate     `gorm:"foreignKey:PartnerID" json:"affiliates"` // ความสัมพันธ์กับ Affiliate
	TotalEarnings   decimal.Decimal  `gorm:"type:decimal(10,2);column:totalearnings;default:0" json:"totalearnings"` // ค่าคอมมิชชั่นสะสม
	Token           string          `gorm:"type:text;column:token" json:"token"`
	Bankname        string          `gorm:"type:varchar(250);column:bankname" json:"bankname"`
	Banknumber      string          `gorm:"type:varchar(50);column:banknumber;NOT NULL" json:"banknumber"`
	Balance         decimal.Decimal  `gorm:"type:numeric(8,2);column:balance;default:0" json:"balance"`
}

// type Partner struct {
// 	gorm.Model
// 	ID              int             `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"` // รหัส Partner
// 	Name            string          `gorm:"type:varchar(255);column:name;NOT NULL"`      // ชื่อ Partner
// 	Username         string    `gorm:"index:idx_username,unique;type:varchar(255);column:username;NOT NULL"`
// 	Password         string    `gorm:"type:text;column:password;NOT NULL"`
// 	Email           string          `gorm:"type:varchar(255);column:email;unique;NOT NULL"` // อีเมล
// 	Phone           string          `gorm:"type:varchar(50);column:phone"`              // เบอร์โทร
// 	AffiliateKey    string          `gorm:"type:varchar(50);column:affiliatekey;unique;NOT NULL"`           // รหัส Affiliate ที่ใช้แนบในลิงค์
// 	CommissionRate  decimal.Decimal `gorm:"type:decimal(5,2);column:commission_rate"`   // เปอร์เซ็นต์ค่าคอมมิชชั่น เช่น 10.00%
// 	TotalCommission decimal.Decimal `gorm:"type:decimal(10,2);column:total_commission"` // ยอดค่าคอมมิชชั่นสะสม
// 	Status          string          `gorm:"type:varchar(50);column:status;default:'active'"` // สถานะ (active/inactive)
// 	Prefix          string    		`gorm:"type:varchar(50);column:prefix;NOT NULL"`    
// 	CreatedAt       time.Time       `gorm:"column:created_at;NOT NULL"`                 // วันที่สร้าง
// 	UpdatedAt       time.Time       `gorm:"column:updated_at;NOT NULL"`                 // วันที่อัปเดตล่าสุด
// 	Affiliates []Affiliate `gorm:"foreignKey:PartnerID"` // ความสัมพันธ์กับ Affiliate
// 	TotalEarnings   decimal.Decimal `gorm:"type:decimal(10,2);default:0"` // ค่าคอมมิชชั่นสะสม
// 	Token            string    `gorm:"type:text";gorm:"column:token"`
// 	Bankname         string    `gorm:"type:varchar(250);column:bankname"`
// 	Banknumber       string    `gorm:"type:varchar(50)";column:banknumber;NOT NULL"`
// 	Balance          decimal.Decimal   `gorm:"type:numeric(8,2);column:balance;default:0"`


// }

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
