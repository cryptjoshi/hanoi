package models
import (
	"github.com/shopspring/decimal"
	//"gorm.io/gorm"
	"time"
)



type BankStatement struct {
	//gorm.Model
	ID                int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`
	Uid               string    `gorm:"type:varchar(255);column:uid;"` 
	Accountno         string    `gorm:"type:varchar(255)";gorm:"column:accountno"`
	Bankname          string    `gorm:"type:varchar(255)";gorm:"column:bankname"`
	Userid            int       `gorm:"column:userid;NOT NULL"`
	Walletid          int       `gorm:"column:walletid;NOT NULL"`
	Bankcode          int       `gorm:"column:bankcode"`
	BetAmount 		  decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:betamount;"`
	Bet_amount 		  decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:bet_amount;"`
	Transactionamount decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:transactionamount;NOT NULL"`
	Amount decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:amount;"`
	Beforebalance     decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:beforebalance"`
	Proamount         decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:proamount"`
	AddOnamount       decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:addonamount"`
	Balance           decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:balance"`
	Turnover          decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:turnover"`
	Unix              int       `gorm:"column:unix"`
	CreatedAt         time.Time `gorm:"column:createdAt;NOT NULL"`
	UpdatedAt         time.Time `gorm:"column:updatedAt;NOT NULL"`
	DeletionAt    time.Time `gorm:"default:current_timestamp(3)";gorm:"column:deletionAt;NOT NULL"`
	Status            string    `gorm:"column:status;"`
	BankaccountID     int       `gorm:"column:bankaccount_id"`
	TransactionDate   string    `gorm:"type:varchar(20);column:transaction_date"`
	Channel           string    `gorm:"type:varchar(100);column:channel"`
	QrExpireDate      string    `gorm:"type:varchar(20);column:qr_expire_date"`
	StatementType     string    `gorm:"type:varchar(50);column:statement_type"`
	Detail            string    `gorm:"type:text;column:detail;"`
	Prefix            string    `gorm:"type:text;column:prefix;"`
	Ismanual          int       `gorm:"column:ismanual;default:0"`
}

func (m *BankStatement) TableName() string {
	return "BankStatement"
}

type SwaggerBankStatement struct {
 
	ID                int       `gorm:"column:id;NOT NULL"`
	Uid               string    `gorm:"type:varchar(255);column:uid;"` 
	Accountno         string    `gorm:"type:varchar(255)";gorm:"column:accountno"`
	Bankname          string    `gorm:"type:varchar(255)";gorm:"column:bankname"`
	Userid            int       `gorm:"column:userid;NOT NULL"`
	Walletid          int       `gorm:"column:walletid;NOT NULL"`
	Bankcode          int       `gorm:"column:bankcode"`
	BetAmount 		  decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:betamount;"`
	Bet_amount 		  decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:bet_amount;"`
	Transactionamount decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:transactionamount;NOT NULL"`
	Beforebalance     decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:beforebalance"`
	Balance           decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:balance"`
	Turnover          decimal.Decimal    `gorm:"type:numeric(10,2);gorm:"column:turnover"`
	Unix              int       `gorm:"column:unix"`
	CreatedAt         time.Time `gorm:"column:createdAt;NOT NULL"`
	UpdatedAt         time.Time `gorm:"column:updatedAt;NOT NULL"`
	DeletionAt    time.Time `gorm:"default:current_timestamp(3)";gorm:"column:deletionAt;NOT NULL"`
	Status            string    `gorm:"column:status;"`
	BankaccountID     int       `gorm:"column:bankaccount_id"`
	TransactionDate   string    `gorm:"type:varchar(20);column:transaction_date"`
	Channel           string    `gorm:"type:varchar(100);column:channel"`
	QrExpireDate      string    `gorm:"type:varchar(20);column:qr_expire_date"`
	StatementType     string    `gorm:"type:varchar(50);column:statement_type"`
	Detail            string    `gorm:"type:text;column:detail;"`
	Prefix            string    `gorm:"type:text;column:prefix;"`
	Ismanual          int       `gorm:"column:ismanual;default:0"`
}