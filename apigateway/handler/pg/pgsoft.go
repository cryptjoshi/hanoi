package pg

import 
(
	"github.com/gofiber/fiber/v2"
	"apigateway/models"
	"apigateway/database"
	"apigateway/handler"
	"apigateway/repository"
	"github.com/shopspring/decimal"
	//jtoken "github.com/golang-jwt/jwt/v4"
	"fmt"
	
	"os"
	"strings"
)
type Response struct {
    Message string      `json:"message"`
    Status  bool        `json:"status"`
    Data    interface{} `json:"data"` // ใช้ interface{} เพื่อรองรับข้อมูลหลายประเภทใน field data
}
 
type TxnsRequest struct {

		Status string `json:"status"`
		RoundId string `json:"roundid"`
		BetAmount decimal.Decimal `json:"betamount"`
		PayoutAmount decimal.Decimal `json:"payoutamount"`
		GameCode string `json:"gamecode"`
		PlayInfo string `json:"playinfo"`
		TxnId string `json:"txnid"`
		TurnOver decimal.Decimal `json:"turnover"`
		IsEndRound bool `json:"isendround"`
        IsFeatureBuy bool `json:"isfeaturebuy"` 
        IsFeature bool `json:"isfeature"`
}

type PgRequest struct {
	Id string `json:"id"`
	TimestampMillis int `json:"timestampmillis"`
	ProductID string `json:"productid`
	Currency string `json:"currency"`
	Username string `json:"username"`
	SessionToken string `json:"sessiontoken"`
	StatusCode  int  `json:"statuscode"`
	Balance  decimal.Decimal `json:"balance"`
	Txns []TxnsRequest `json:"txns"`
}

type PGResponse struct {
	statusCode  int  `json:"statuscode"`
    ErrorMessage string `json:"errormessage"`
    Balance  decimal.Decimal `json:"balance"`
    BeforeBalance decimal.Decimal `json:"beforebalance"`
}
type ResponseBalance struct {
	BetAmount decimal.Decimal `json:"betamount"`
	BeforeBalance decimal.Decimal `json:"beforebalance"`
	Balance decimal.Decimal `json:"balance"`
}

// var PG_API_KEY = "9dc857f4-2225-45ef-bf0f-665bcf7d4a1b"  
// var PG_API_KEY= "31d3cc58-4e34-4dc4-9c45-b8abe6a1b0d2"
var SECRET_KEY = os.Getenv("PASSWORD_SECRET")
var pg_prod_code = os.Getenv("PG_PRODUCT_ID")



func Index(c *fiber.Ctx) error {

	//var user []models.Users
	
	//database.Database.Find(&user)
	response := Response{
		Message: "Welcome to PGSoft!!",
		Status:  true,
		Data: fiber.Map{}, 
	}
	 
	return c.JSON(response)
   
}
 

func GetBalance(c *fiber.Ctx) error {

 

	request := new(PgRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	var users models.Users
	users = handler.ValidateJWTReturn(request.SessionToken);

 

	balanceFloat, _ := users.Balance.Float64()
	if users.Token == request.SessionToken {
		
		response := fiber.Map{
			"statusCode": 0,
			"id": request.Id,
			"timestampMillis": request.TimestampMillis+100,
			"productId": request.ProductID,
			"currency": request.Currency,
			"username": strings.ToUpper(request.Username),
			"balance": balanceFloat,
		}
		return c.JSON(response)
	}else {
		response := fiber.Map{
			"statusCode": 30001,
			"id": request.Id,
			"timestampMillis": request.TimestampMillis +100,
			"productId": request.ProductID,
			"currency": request.Currency,
			"username": strings.ToUpper(request.Username),
			"balance": decimal.NewFromFloat(0),
		}
		return c.JSON(response)
	}

 

	
}

func PlaceBet(c *fiber.Ctx) error {
	
	request := new(PgRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	response := fiber.Map{
		"statusCode": 0,
		"id": request.Id,
		"timestampMillis": request.TimestampMillis +100,
		"productId": request.ProductID,
		"currency": request.Currency,
		"username": strings.ToUpper(request.Username),
		"balance": decimal.NewFromFloat(0),
	}
		var user models.Users
		db, _ := database.ConnectToDB(request.Username)
		db.Where("username = ?", request.Username).First(&user)
		
		 for _, transaction := range request.Txns {
			
			transactionAmount := func(betamount decimal.Decimal,payoutamount decimal.Decimal,status string,feature bool) decimal.Decimal {
				 if status == "OPEN" {
					return betamount.Neg()
				 } else if feature == true {
					return payoutamount.Sub(betamount)
				 } else {
					return payoutamount.Sub(betamount)
				 }
			}(transaction.BetAmount,transaction.PayoutAmount,transaction.Status,transaction.IsFeatureBuy)

			// fmt.Printf(" IsFeatureBuy: %s ",transaction.IsFeatureBuy)
			
			
			xtransaction := map[string]interface{}{
				"MemberID" : user.ID,
				"MemberName":strings.ToUpper(request.Username),
				"ProductID":1,//productId,
				"ProviderID":1,
				"WagerID":0,
				"CurrencyID":0,//currency=0THB,
				"GameCode":transaction.GameCode,
				"PlayInfo":transaction.PlayInfo,
				"GameID":transaction.GameCode,
				"GameRoundID":transaction.RoundId,
				"BetAmount":transaction.BetAmount,
				//"TxnsID":transaction.TxnId,
				"TransactionID":0,
				"PayoutAmount":transaction.PayoutAmount,
				"PayoutDetail":transaction.PlayInfo,
				"SettlementDate":request.TimestampMillis,
				"Status":0,//status-0=SETTLED,
				//BeforeBalance:beforeBalance,
			   // Balance:beforeBalance-betAmount,
				"OperatorCode":pg_prod_code,
				"OperatorID":1,//1=PGGAME
				"ProviderLineID":1,//1-PGAME
				"GameType":1,//1=PGGAME
				"ValidBetAmount":transaction.BetAmount,
				"TransactionAmount":transactionAmount,
				"TurnOver":transaction.TurnOver,
				"CommissionAmount":0,
				"JackpotAmount":0,
				"JPBet":0,
				"MessageID":"",
				"Sign":"",
				"RequestTime":request.TimestampMillis,
				"IsFeature":transaction.IsFeature,
				"IsEndRound":transaction.IsEndRound, 
				"IsFeatureBuy":transaction.IsFeatureBuy, 
				"GameProvide": "PGSOFT",
				"BeforeBalance":user.Balance,
				"Balance":user.Balance.Add(transactionAmount),
			  } 

			
	  		
			if user.Balance.LessThan(transaction.BetAmount) {

				response := fiber.Map{
					"statusCode": 10002,
					"id": request.Id,
					"timestampMillis": request.TimestampMillis +100,
					"productId": request.ProductID,
					"currency": request.Currency,
					"balanceBefore": 0,
                	"balanceAfter": 0,
					"username": strings.ToUpper(request.Username),
					"message": "Balance incorrect",
				}	
				return c.JSON(response)
			} else 
			{
				var c_transaction_found models.TransactionSub
				db, _ := database.ConnectToDB(request.Username)

				rowsAffected := db.Debug().Model(&models.TransactionSub{}).Select("id").Where("GameRoundID = ? ",transaction.RoundId).Find(&c_transaction_found).RowsAffected
				fmt.Println(" GameRoundID RowAffected: ",rowsAffected)
				if rowsAffected == 0 {
							_err_  := db.Table("TransactionSub").Create(xtransaction).Error
							if _err_ != nil {
								fmt.Println(_err_)
								//return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ไม่สามารถแทรกข้อมูลได้"})
							} 
							//_err_ := database.Database.Model(&models.TransactionSub{}).Create(xtransaction);
							fmt.Println(transactionAmount)
							updates := map[string]interface{}{
								"Balance": user.Balance.Add(transactionAmount),
								}
							repository.UpdateFieldsUserString(request.Username, updates) 
							balanceBeforeFloat, _ := user.Balance.Float64()
							balanceAfterFloat, _ := user.Balance.Add(transactionAmount).Float64()
							response := fiber.Map{
								"statusCode": 0,
								"id": request.Id,
								"timestampMillis": request.TimestampMillis +100,
								"productId": request.ProductID,
								"currency": request.Currency,
								"balanceBefore": balanceBeforeFloat,
								"balanceAfter": balanceAfterFloat,
								"username": strings.ToUpper(request.Username),
							}
						
						return c.JSON(response)
					} else {
						// balanceBeforeFloat, _ := c_transaction_found.BeforeBalance.Float64()
						balanceAfterFloat, _ := user.Balance.Float64()
						// fmt.Println("---------------------------------------------")
						// fmt.Println("GameRoundID:",c_transaction_found.GameRoundID)
						// fmt.Println("---------------------------------------------")
						// fmt.Println("user Balance:",balanceBeforeFloat)
						// fmt.Println("user Balance:",balanceAfterFloat)
						// fmt.Println("user Balance:",user.Balance)
						// fmt.Println("---------------------------------------------")
						
						response := fiber.Map{
							"statusCode": 0,
							"id": request.Id,
							"timestampMillis": request.TimestampMillis +100,
							"productId": request.ProductID,
							"currency": request.Currency,
							"balanceBefore": balanceAfterFloat,
							"balanceAfter": balanceAfterFloat,
							"username": strings.ToUpper(request.Username),
							"message": "Balance incorrect",
						}
						return c.JSON(response)	
					}
			}
		 }
		 return c.JSON(response)		 
}

 