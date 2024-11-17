package pg

import 
(
	"github.com/gofiber/fiber/v2"
	"pkd/models"
	"pkd/database"
	"pkd/handler"
	"pkd/repository"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"pkd/common"
	
	//jtoken "github.com/golang-jwt/jwt/v4"
	"fmt"
	"log"
	"encoding/json"
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

//http://ambsuperapi.com
//user : sunshinepgthb
//pass : Sunshine@688

func Index(c *fiber.Ctx) error {

	var user []models.Users
	database.Database.Find(&user)
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

	fmt.Printf("users: %s ",users.Token)
	fmt.Printf("request: %s ",request.SessionToken)

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
	prefix,perr := handler.GetPrefix(request.Username)
	if perr != nil {
			fmt.Println(perr)
			
	}
	var user models.Users
	db,cnn := database.ConnectToDB(prefix) 
	if cnn != nil {
		fmt.Printf("error: %s \n",cnn)
	}
		 dberr := db.Debug().Where("username = ?", request.Username).First(&user).Error
		 if dberr != nil {
			fmt.Println(dberr)
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
		 }
	
		
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
				rowsAffected := db.Debug().Model(&models.TransactionSub{}).Select("id").Where("GameRoundID = ? ",transaction.RoundId).Find(&c_transaction_found).RowsAffected
				fmt.Println(" GameRoundID RowAffected: ",rowsAffected)
				if rowsAffected == 0 {
							_err_  := db.Debug().Model(&models.TransactionSub{}).Create(xtransaction).Error
							if _err_ != nil {
								fmt.Println(_err_)
								//return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ไม่สามารถแทรกข้อมูลได้"})
							} 
							//_err_ := database.Database.Model(&models.TransactionSub{}).Create(xtransaction);
							//fmt.Printf("TransactionAmount %v \n",transactionAmount)
							updates := map[string]interface{}{
								"Balance": user.Balance.Add(transactionAmount),
								}

							repository.UpdateUserFields(db,user.ID, updates) 
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

var SECRET_KEY = os.Getenv("PASSWORD_SECRET")
var pg_prod_code = os.Getenv("PG_PRODUCT_ID")

var OPERATOR_CODE = "sunshinetest" //"sunshinepgthb"//"sunshinetest",
var SECRET_API_KEY = os.Getenv("PG_API_KEY") //"9dc857f4-2225-45ef-bf0f-665bcf7d4a1b" //os.Getenv("PG_API_KEY")
var PG_PROD_CODE= os.Getenv("PG_PRODUCT_ID")
var PG_API_URL = "https://test.ambsuperapi.com"//os.Getenv("PG_API_URL") //"https://prod_md.9977997.com"
var PG_PROD_URL = "https://api.hentory.io" 


func makePostRequest(url string, bodyData interface{}) (*fasthttp.Response, error) {
	// Marshal requestData struct เป็น JSON
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// สร้าง Request และ Response
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// ตั้งค่า URL, Method, และ Body
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	authHeader := common.CreateBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
	req.Header.Add("Authorization", authHeader)
	req.SetBody(jsonData)

	// ส่ง request
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)
	
	return resp, nil
}
func makeGetRequest(url string) (*fasthttp.Response, error) {
	// Marshal requestData struct เป็น JSON
	// jsonData, err := json.Marshal(bodyData)
	// if err != nil {
	// 	return nil, fmt.Errorf("error marshaling JSON: %v", err)
	// }

	// สร้าง Request และ Response
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// ตั้งค่า URL, Method, และ Body
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	authHeader := common.CreateBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
	req.Header.Add("Authorization", authHeader)
	//req.SetBody(jsonData)

	// ส่ง request
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)
	
	return resp, nil
}

func LaunchGame(c *fiber.Ctx) error {
	type BodyGame struct {
		ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
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
		//ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	//	Txns []TxnsRequest `json:"txns"`
	}
	type CResponse struct {
		Message string      `json:"message"`
		Status  bool        `json:"status"`
		Data    interface{} `json:"data"`  
	}

	var response CResponse
	// bodyRequest := new(BodyGame)

	// if err := c.BodyParser(&bodyRequest); err != nil {
	// 	fmt.Printf(" %s ", err.Error())
	// 	response := fiber.Map{
	// 		"Status":  false,
	// 		"Message": err.Error(),
	// 	}
	// 	return c.JSON(response)
	// }

	// fmt.Printf("Body: %s",bodyRequest.Body)
	//var tokenString := c.Get("Authorization")[7:]
	request := new(PgRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	var users models.Users
	users = handler.ValidateJWTReturn(request.SessionToken);

	//fmt.Printf("users: %v ",users)
	//fmt.Printf("request: %s ",request.SessionToken)
	// efargs = {
	// 	"OperatorCode": OPERATOR_CODE,
	// 	"MemberName": req.body.username,
	// 	"Password":   response.data.uid,
	// 	"ProductID": ProductID,
	// 	"GameType": GameType,
	// 	"GameID": GameID,
	// 	"LanguageCode": LanguageCode,
	// 	"Platform": Platform,
	// 	"Sign": hashSignature("LaunchGame",RequestTime),
	// 	"RequestTime": RequestTime
	// 	}
	var args = fiber.Map{
		"username": strings.ToLower(users.Username),//user.data.username,
		"productId":common.PG_PROD_CODE,
		"gameCode": request.ProductID,
		"isMobileLogin": true,
		"sessionToken": request.SessionToken,
		//"betLimit": [],
		"callbackUrl":"https://www.โชคดี789.com/lobby/slot/game?id=8888&type=1", //`${req.protocol}://${req.get('host')}${req.originalUrl}`
	}
	
	//fmt.Printf(" args : %s ",args)
	
	resp,err := makePostRequest(common.PG_API_URL+"/seamless/login",args)		
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	resultBytes := resp.Body()
	resultString := string(resultBytes)
	// แสดงผล string ที่ได้
	fmt.Println("Response body as string:", resultString)

	err = json.Unmarshal([]byte(resultString), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	respon := fiber.Map{
		"Status":  true,
		"Message": response.Message,
		"Data": response.Data,
	}
	return c.JSON(respon)
}

	//url := fmt.Sprintf(PG_PROD_URL,"/seamless/login")
