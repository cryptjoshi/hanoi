package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)


type BuyInOut struct {
	gorm.Model
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
}

func (m *BuyInOut) TableName() string {
	return "BuyInOut"
}


type Buyinout struct {
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

type BuyInOutRequest struct {
    Transaction  BuyInOut `json:"Transaction"`
    MemberName   string         `json:"MemberName"`
    OperatorCode string         `json:"OperatorCode"`
    ProductID    int            `json:"ProductID"`
    MessageID    string         `json:"MessageID"`
    Sign         string         `json:"Sign"`
    RequestTime  string         `json:"RequestTime"`
}

// Struct หลักที่ประกอบด้วยรายการ Transactions


