package models

import (
	"github.com/shopspring/decimal"
	//"gorm.io/gorm"
	"time"
)

type Users struct {
	//gorm.Model
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
	MinTurnover      decimal.Decimal   `gorm:"type:decimal(10,2);gorm:"column:minturnover;default:0"`
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
   