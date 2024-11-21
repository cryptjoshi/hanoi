package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)


type TransactionSub struct {
	//gorm.Model
	ID                int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL"`
	MemberID          int       `gorm:"type:varchar(100);column:MemberID"`
	MemberName        string    `gorm:"type:varchar(255);column:MemberName"`
	CurrencyID        int       `gorm:"column:CurrencyID"`
	TransactionAmount decimal.Decimal    `gorm:"column:TransactionAmount"`
	Status            int       `gorm:"column:Status"`
	BeforeBalance     decimal.Decimal    `gorm:"column:BeforeBalance"`
	Balance           decimal.Decimal    `gorm:"column:Balance"`
	Sign              string    `gorm:"type:text;column:Sign"`
	RequestTime       string    `gorm:"type:varchar(50);column:RequestTime"`
	OperatorCode      string    `gorm:"type:varchar(50);column:OperatorCode"`
	OperatorID        int       `gorm:"column:OperatorID"`
	ProductID         int64     `gorm:"column:ProductID"`
	ProviderID        int       `gorm:"column:ProviderID"`
	ProviderLineID    int       `gorm:"column:ProviderLineID"`
	WagerID           int64     `gorm:"column:WagerID"`
	GameType          int       `gorm:"column:GameType"`
	GameID            string    `gorm:"type:varchar(255);column:GameID"`
	GameRoundID       string    `gorm:"type:varchar(255);column:GameRoundID"`
	ValidBetAmount    decimal.Decimal    `gorm:"column:ValidBetAmount"`
	BetAmount         decimal.Decimal    `gorm:"column:BetAmount"`
	PayoutAmount      decimal.Decimal    `gorm:"column:PayoutAmount"`
	PayoutDetail      string    `gorm:"type:text;column:PayoutDetail"`
	CommissionAmount  decimal.Decimal    `gorm:"column:CommissionAmount"`
	JackpotAmount     decimal.Decimal    `gorm:"column:JackpotAmount"`
	SettlementDate    string `gorm:"column:SettlementDate"`
	JPBet             decimal.Decimal       `gorm:"column:JPBet"`
	AfterBalance      decimal.Decimal    `gorm:"column:AfterBalance"`
	MessageID         string    `gorm:"type:varchar(255);column:MessageID"`
	CreatedAt         time.Time       `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt         time.Time       `gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt         gorm.DeletedAt  `gorm:"column:deleted_at;type:timestamp null"`
	TransactionID     string    `gorm:"type:varchar(255);column:TransactionID"`
	IsEndRound        int       `gorm:"column:IsEndRound"`
	IsFeatureBuy      int       `gorm:"column:IsFeatureBuy"`
	IsFeature         int       `gorm:"column:IsFeature"`
	IsAction          string    `gorm:"type:varchar(50);column:IsAction"`
	GameProvide       string    `gorm:"type:varchar(255);column:GameProvide"`
	GameNumber        string    `gorm:"type:varchar(20);column:GameNumber"`
	Prefix			  string    `gorm:"column:Prefix"`
	TurnOver		  decimal.Decimal    `gorm:"column:TurnOver"`
 	ProID			  string    `gorm:"column:ProID"`
	AffiliateID      string          `gorm:"type:varchar(50);column:affiliate_id"`         // รหัส Affiliated ที่เชื่อมโยง
	AffiliateTurnover decimal.Decimal `gorm:"type:decimal(15,2);column:affiliate_turnover;default:0"` // Turnover ที่เกิดจาก Affiliate
	PartnerID int `gorm:"column:partner_id"` // บันทึก partner ที่เกี่ยวข้องกับ transaction
}

func (m *TransactionSub) TableName() string {
	return "TransactionSub"
}


type Transaction struct {
	gorm.Model
    MemberID          int             `json:"MemberID"`
    OperatorID        int             `json:"OperatorID"`
    ProductID         int             `json:"ProductID"`
    ProviderID        int             `json:"ProviderID"`
    ProviderLineID    int             `json:"ProviderLineID"`
    WagerID           int64           `json:"WagerID"`
    CurrencyID        int             `json:"CurrencyID"`
    GameType          int             `json:"GameType"`
    GameID            string          `json:"GameID"`
    GameRoundID       string          `json:"GameRoundID"`
    BetAmount         decimal.Decimal `json:"BetAmount"`
	BeforeBalance     decimal.Decimal `json:"BeforeBalance"`
	Balance     	  decimal.Decimal `json:"Balance"`
    ValidBetAmount    decimal.Decimal `json:"ValidBetAmount"`
    Fee               decimal.Decimal `json:"Fee"`
    JPBet             decimal.Decimal `json:"JPBet"`
    PayoutAmount      decimal.Decimal `json:"PayoutAmount"`
    CommissionAmount  decimal.Decimal `json:"CommissionAmount"`
    JackpotAmount     decimal.Decimal `json:"JackpotAmount"`
    PayoutDetail      interface{}     `json:"PayoutDetail"` // ใช้ interface{} เพราะข้อมูลอาจเป็น null
    Data              interface{}     `json:"Data"`         // ใช้ interface{} เพราะข้อมูลอาจเป็น null
    Status            int             `json:"Status"`
    CreatedOn         string       `json:"CreatedOn"`
    ModifiedOn        string       `json:"ModifiedOn"`
    SettlementDate    *time.Time      `json:"SettlementDate"` // ใช้ pointer เพราะข้อมูลอาจเป็น null
    TransactionID     string          `json:"TransactionID"`
    TransactionAmount decimal.Decimal `json:"TransactionAmount"`
}


func (m *Transaction) TableName() string {
	return "Transaction"
}

// Struct หลักที่ประกอบด้วยรายการ Transactions
type TransactionsRequest struct {
    Transactions []TransactionSub `json:"Transactions"`
    MemberName   string         `json:"MemberName"`
    OperatorCode string         `json:"OperatorCode"`
    ProductID    int            `json:"ProductID"`
    MessageID    string         `json:"MessageID"`
    Sign         string         `json:"Sign"`
    RequestTime  string         `json:"RequestTime"`
}

type SwaggerTransactionSub struct {
 
	ID                int       `gorm:"column:id;NOT NULL"`
	MemberID          int       `gorm:"type:varchar(100)";gorm:"column:MemberID"`
	MemberName        string    `gorm:"type:varchar(255)";gorm:"column:MemberName"`
	CurrencyID        int       `gorm:"column:CurrencyID"`
	TransactionAmount decimal.Decimal    `gorm:"column:TransactionAmount"`
	Status            int       `gorm:"column:Status"`
	BeforeBalance     decimal.Decimal    `gorm:"column:BeforeBalance"`
	Balance           decimal.Decimal    `gorm:"column:Balance"`
	Sign              string    `gorm:"type:text";gorm:"column:Sign"`
	RequestTime       string    `gorm:"type:varchar(50)";gorm:"column:RequestTime"`
	OperatorCode      string    `gorm:"type:varchar(50)";gorm:"column:OperatorCode"`
	OperatorID        int       `gorm:"column:OperatorID"`
	ProductID         int64     `gorm:"column:ProductID"`
	ProviderID        int       `gorm:"column:ProviderID"`
	ProviderLineID    int       `gorm:"column:ProviderLineID"`
	WagerID           int64     `gorm:"column:WagerID"`
	GameType          int       `gorm:"column:GameType"`
	GameID            string    `gorm:"type:varchar(255)";gorm:"column:GameID"`
	GameRoundID       string    `gorm:"type:varchar(255)";gorm:"column:GameRoundID"`
	ValidBetAmount    decimal.Decimal    `gorm:"column:ValidBetAmount"`
	BetAmount         decimal.Decimal    `gorm:"column:BetAmount"`
	PayoutAmount      decimal.Decimal    `gorm:"column:PayoutAmount"`
	PayoutDetail      string    `gorm:"type:text";gorm:"column:PayoutDetail"`
	CommissionAmount  decimal.Decimal    `gorm:"column:CommissionAmount"`
	JackpotAmount     decimal.Decimal    `gorm:"column:JackpotAmount"`
	SettlementDate    string `gorm:"column:SettlementDate"`
	JPBet             decimal.Decimal       `gorm:"column:JPBet"`
	AfterBalance      decimal.Decimal    `gorm:"column:AfterBalance"`
	MessageID         string    `gorm:"type:varchar(255)";gorm:"column:MessageID"`
	CreatedAt         string `gorm:"column:createdAt;NOT NULL"`
	UpdatedAt         string `gorm:"column:updatedAt;NOT NULL"`
	DeletedAt         string `gorm:"column:deletedAt"`
	TransactionID     string    `gorm:"type:varchar(255)";gorm:"column:TransactionID"`
	IsEndRound        int       `gorm:"column:IsEndRound"`
	IsFeatureBuy      int       `gorm:"column:IsFeatureBuy"`
	IsFeature         int       `gorm:"column:IsFeature"`
	IsAction          string    `gorm:"type:varchar(50)";gorm:"column:IsAction"`
	GameProvide       string    `gorm:"type:varchar(255)";gorm:"column:GameProvide"`
	GameNumber        string    `gorm:"type:varchar(20)";gorm:"column:GameNumber"`
	Prefix			  string    `gorm:"column:Prefix"`
}

func (t *TransactionSub) BeforeCreate(tx *gorm.DB) error {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	t.CreatedAt = time.Now().In(loc)
	t.UpdatedAt = time.Now().In(loc)
	return nil
}

// BeforeUpdate จะทำงานก่อนที่จะอัพเดทข้อมูล
func (t *TransactionSub) BeforeUpdate(tx *gorm.DB) error {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	t.UpdatedAt = time.Now().In(loc)
	return nil
}