package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Users struct {
	gorm.Model
	ID               int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`
	Walletid         int       `gorm:"column:walletid;NOT NULL"`
	Username         string    `gorm:"index:idx_username,unique";gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:username;NOT NULL"`
	Password         string    `gorm:"type:text";gorm:"column:password;NOT NULL"`
	ProviderPassword string    `gorm:"type:varchar(100)";gorm:"column:provider_password"`
	Fullname         string    `gorm:"type:text";gorm:"column:fullname"`
	Preferredname    string    `gorm:"type:varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:preferredname"`
	Bankname         string    `gorm:"type:varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:bankname"`
	Banknumber       string    `gorm:"type:varchar(50)";gorm:"column:banknumber;NOT NULL"`
	Balance          decimal.Decimal   `gorm:"type:numeric(8,2);gorm:"column:balance;default:0"`
	Beforebalance    decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:beforebalance;default:0;NOT NULL"`
	Token            string    `gorm:"type:text";gorm:"column:token"`
	Role             string    `gorm:"type:varchar(50)";gorm:"column:role"`
	Salt             string    `gorm:"type:varchar(150)";gorm:"column:salt;NOT NULL"`
	CreatedAt        time.Time `gorm:"column:createdAt;NOT NULL"`
	UpdatedAt        time.Time `gorm:"column:updatedAt;NOT NULL"`
	DeletionAt    time.Time `gorm:"default:current_timestamp(3)";gorm:"column:deletionAt;NOT NULL"`
	Status           int       `gorm:"column:status;default:1"`
	Betamount        decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:betamount;default:0"`
	Win              decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:win;default:0"`
	Lose             decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:lose;default:0"`
	Turnover         decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:turnover;default:0"`
	MinTurnover      string    `gorm:"column:minturnover;type:varchar(50);default:0"`
	MinTurnoverDef   string    `gorm:"column:minturnoverdef;"gorm:"type:varchar(50);"gorm:"default:'10%'"`
	BetLimit         string    `gorm:"column:betlimit;"gorm:"type:varchar(50);"`
	Currency         string    `gorm:"type:varchar(50)";gorm:"column:currency"`
	ProID            string    `gorm:"type:varchar(50)";gorm:"column:pro_id"`
	PartnersKey      string    `gorm:"type:varchar(50)";gorm:"column:partners_key"`
	ProStatus        string    `gorm:"type:varchar(50)";gorm:"column:pro_status;default:none"`
	ProBalance       decimal.Decimal   `gorm:"type:decimal(10,2);column:probalance;default:0"`
	Firstname        string    `gorm:"type:varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:firstname"`
	Lastname         string    `gorm:"type:varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:lastname"`
	Deposit          decimal.Decimal   `gorm:"type:decimal(10,2);column:deposit;default:0"`
	Withdraw         decimal.Decimal   `gorm:"type:decimal(10,2);column:withdraw;default:0"`
	Credit           decimal.Decimal   `gorm:"type:decimal(10,2);column:credit;default:0"`
	Prefix           string    `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:prefix;NOT NULL"`
	Actived          *time.Time `gorm:"default:current_timestamp(3)";gorm:"column:actived;NOT NULL"`
	TempSecret       string    `gorm:"type:varchar(50)";gorm:"column:temp_secret"`
	Secret           string    `gorm:"type:text";gorm:"column:secret"`
	OtpAuthUrl       string    `gorm:"type:text";gorm:"column:otpAuthUrl"`
	LastProamount      decimal.Decimal   `gorm:"type:decimal(10,2);column:lastproamount;default:0"`
	LastDeposit      decimal.Decimal   `gorm:"type:decimal(10,2);column:lastdeposit;default:0"`
	LastWithdraw     decimal.Decimal   `gorm:"type:decimal(10,2);column:lastwithdraw;default:0"`
	TotalTurnover     decimal.Decimal `gorm:"type:decimal(15,2);column:total_turnover;default:0"` // Turnover รวมที่เกิดจากสมาชิกและ Affiliated
	CommissionEarned  decimal.Decimal `gorm:"type:decimal(15,2);column:commission_earned;default:0"` // ค่าคอมมิชชันสะสมจาก Affiliated
	PartnerID int `gorm:"column:partner_id"` // ใช้บันทึก ID ของ partner ที่เชื่อมโยง
	AffiliateLink string `gorm:"type:varchar(255);column:affiliate_link"` // ใช้บันทึกลิงค์ที่ใช้
	ReferralCode  string          `gorm:"column:referral_code;unique;not null"` // รหัสแนะนำ
	ReferredBy    string          `gorm:"column:referred_by"` // ผู้แนะนำ (User ID)
	TotalEarnings decimal.Decimal `gorm:"column:total_earnings;type:decimal(10,2);default:0"` // ค่าคอมมิชชั่นสะสม
}

func (m *Users) TableName() string {
	return "Users"
}

type LoginResponse struct {
	Token string `json:"token"`
   }

type MyResponse struct {
	Status bool `json:"status"`
	Data   string  `json:"data"`
	Message string `json:"message"` 
}

type SwaggerBody struct {
	Prefix string `json:"prefix"`
}

type SwaggerUser struct {
	ID               int       `gorm:"column:id;NOT NULL"`
	Walletid         int       `gorm:"column:walletid;NOT NULL"`
	Username         string    `gorm:"index:idx_username,unique";gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:username;NOT NULL"`
	Password         string    `gorm:"type:text";gorm:"column:password;NOT NULL"`
	ProviderPassword string    `gorm:"type:varchar(100)";gorm:"column:provider_password"`
	Fullname         string    `gorm:"type:text";gorm:"column:fullname"`
	Preferredname    string    `gorm:"type:varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:preferredname"`
	Bankname         string    `gorm:"type:varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:bankname"`
	Banknumber       string    `gorm:"type:varchar(50)";gorm:"column:banknumber;NOT NULL"`
	Balance          decimal.Decimal   `gorm:"type:numeric(8,2);gorm:"column:balance;default:0"`
	Beforebalance    decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:beforebalance;default:0;NOT NULL"`
	Token            string    `gorm:"type:text";gorm:"column:token"`
	Role             string    `gorm:"type:varchar(50)";gorm:"column:role"`
	Salt             string    `gorm:"type:varchar(150)";gorm:"column:salt;NOT NULL"`
	CreatedAt        time.Time `gorm:"column:createdAt;NOT NULL"`
	UpdatedAt        time.Time `gorm:"column:updatedAt;NOT NULL"`
	DeletionAt    time.Time `gorm:"default:current_timestamp(3)";gorm:"column:deletionAt;NOT NULL"`
	Status           int       `gorm:"column:status;default:1"`
	Betamount        decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:betamount;default:0"`
	Win              decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:win;default:0"`
	Lose             decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:lose;default:0"`
	Turnover         decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:turnover;default:0"`
	ProID            string    `gorm:"type:varchar(50)";gorm:"column:pro_id"`
	PartnersKey      string    `gorm:"type:varchar(50)";gorm:"column:partners_key"`
	ProStatus        string    `gorm:"type:varchar(50)";gorm:"column:pro_status;default:none"`
	Firstname        string    `gorm:"type:varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:firstname"`
	Lastname         string    `gorm:"type:varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:lastname"`
	Deposit          decimal.Decimal   `gorm:"type:decimal(10,2);column:deposit;default:0"`
	Withdraw         decimal.Decimal   `gorm:"type:decimal(10,2);column:withdraw;default:0"`
	Credit           decimal.Decimal   `gorm:"type:decimal(10,2);column:credit;default:0"`
	Prefix           string    `gorm:"type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci";gorm:"column:prefix;NOT NULL"`
	Actived          time.Time `gorm:"default:current_timestamp(3)";gorm:"column:actived;NOT NULL"`
	TempSecret       string    `gorm:"type:varchar(50)";gorm:"column:temp_secret"`
	Secret           string    `gorm:"type:text";gorm:"column:secret"`
	OtpAuthUrl       string    `gorm:"type:text";gorm:"column:otpAuthUrl"`
}
   