package gc

import 
(
	"github.com/gofiber/fiber/v2"
	"hanoi/models"
	"hanoi/database"
	//"pkd/handler"
	"hanoi/repository"
	"hanoi/encrypt"
	"crypto/md5"
	//"crypto/des"
	//"crypto/cipher"
	//"bytes"
	"net/url" // เพิ่มการนำเข้าแพ็คเกจ net/url
	"encoding/base64"
	"encoding/json"
	//"encoding/hex"
	//"encoding/json"
	//"pkd/repository"
	"github.com/shopspring/decimal"
	//jtoken "github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4"
	"fmt"
	"time"
	"log"
	"math/big"
	//"strconv"
	//"os"
	"strings"
	"github.com/valyala/fasthttp"
	"github.com/go-resty/resty/v2" // เพิ่มการนำเข้าแพ็คเกจ resty
)
var (
	CLIENT_ID     = "6342e1be-fa03-456f-8d2d-8e1c9513c351" 
	CLIENT_SECRET = "6d83ac42"
	SYSTEMCODE    = "tsxthb"
	WEBID         = "tsxthb"
	DESKEY 		  = "9c62a148"
	DESIV 		  =	"8e014099"
)
type Response struct {
    Message string      `json:"message"`
    Status  bool        `json:"status"`
    Data    interface{} `json:"data"` // ใช้ interface{} เพื่อรองรับข้อมูลหลายประเภทใน field data
}


type GResponse struct {
	MsgID    int    `json:"msgId"`
	Message  string `json:"message"`
	Data     GResponseData   `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

// Struct สำหรับข้อมูลใน field "data"
type GResponseData struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}
type GRequestData struct  {
	SystemCode string `json:"systemcode"`
	WebId      string  `json:"webid"`
	DataList   []interface{} `json:"datalist"`
}

type RequestEncrypt struct {
	Data  string `json:"data"`
	Key string `json:"key"`
	Iv string `json:"iv"`
}
 
type GResult struct {
	MsgId int `json:"msgid"`
	Message string `json:"message"`
	Data GResponseData `json:"data"`
	Timestamp int `json:"timestamp"`
}
type GClaims struct {
	SystemCode    string `json:"systemcode"`
	WebID         string `json:"webid"`
	MemberAccount string `json:"memberaccount"`
	TokenType     string `json:"tokentype"`
	jwt.RegisteredClaims
}
type GTransaction struct {
	Id string `json:"id"`
	Balance decimal.Decimal `json:"balance"`
	Amount decimal.Decimal `json:"amount"`
	ReferenceId string `json:"referenceId"`
	TargetId string `json:"targetId"`
}
type GCGame struct {
	DeskId string `json:"deskId"`
	GameName string `json:"gameName"`
	Shoe string `json:"shoe"`
	Run string `json:"run"`
}
type GcRequest struct {
	Id string `json:"id"`
	Systemcode string `json:"systemcode"`
	Webid string `json:"webid"`
	Account string `json:"account"`
	Requestid string `json:"requestid"`
	Token string `json:"token"`
	Username string `json:"username"`
	Transaction  GTransaction `json:"transaction"`
	Game GCGame `json:"game"`
	
	// Id string `json:"id"`
	// TimestampMillis int `json:"timestampmillis"`
	// ProductID string `json:"productid`
	// Currency string `json:"currency"`
	// Username string `json:"username"`
	// SessionToken string `json:"sessiontoken"`
	// StatusCode  int  `json:"statuscode"`
	// Balance  decimal.Decimal `json:"balance"`
	//Txns []TxnsRequest `json:"txns"`
}

// ฟังก์ชันตัวอย่างใน gclub.go
type LoginRquest struct {
	BackUrl        string `json:"BackUrl"`
	GroupLimitID   string `json:"GroupLimitID"`
	ItemNo         string `json:"ItemNo"`
	Lang           string `json:"Lang"`
	MemberAccount  string `json:"MemberAccount"`
	SystemCode     string `json:"SystemCode"`
	WebId          string `json:"WebId"`
}

type CResponse struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	Data    GResponseData `json:"data"`  
}

type ResponseBalance struct {
	BetAmount decimal.Decimal `json:"betamount"`
	BeforeBalance decimal.Decimal `json:"beforebalance"`
	Balance decimal.Decimal `json:"balance"`
}

var API_URL_G = "http://rcgapiv2.rcg666.com/"
var API_URL_PROXY = "http://api.tsxbet.info:8001"




func Index(c *fiber.Ctx) error {

	//var user []models.Users
	
	response := Response{
		Status: true,
		Message: "GCLUB OK",
		Data:   []interface{}{},
		} 
	 
	// database.Database.Find(&user)
	// response = GResponse{
	// 	MsgId: 0,
	// 	Message: "OK",
	// 	Data: GData{
	// 	  SystemCode: "DocDemoSystem",
	// 	  WebId: "DocDemoWeb",
	// 	  DataList: []interface{}{},
	// 	},
	// 	TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
	// }
	 
	return c.JSON(response)
   
}

func GFetch(apiurl,account string) (CResponse,error) {
	
	apiURL := API_URL_PROXY + apiurl

	data := fiber.Map{
		"SystemCode": SYSTEMCODE,
		"WebId": WEBID,
		"MemberAccount": account,
		"ItemNo": "1",
		"BackUrl": "https://tsx.bet/",
		"GroupLimitID": "1,4,12",
		"Lang": "th-TH",
	}

	unx := time.Now().UnixNano() / int64(time.Millisecond)
	
	// แปลง data เป็น JSON string ก่อนเข้ารหัส
	jsonData, err := json.Marshal(data)
	if err != nil {
		//return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return CResponse{}, fmt.Errorf("error making POST request: %w", err)
	}

	// เข้ารหัส DES
	des, err := encrypt.EncryptDES(string(jsonData), DESKEY, DESIV)
	if err != nil {
		//return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return CResponse{}, fmt.Errorf("error making POST request: %w", err)
	}

	// สร้าง MD5 hash
	md5Hash := fmt.Sprintf("%s%s%d%s", CLIENT_ID, CLIENT_SECRET, unx, des)
	hash := md5.Sum([]byte(md5Hash))
	enc := base64.StdEncoding.EncodeToString(hash[:])

	// สร้าง data ที่เข้ารหัส
	// encryptedData := fiber.Map{
	// 	"enc": enc,
	// 	"des": des,
	// 	"unx": unx,
	// }

	// fmt.Printf("encryptedData: %+v \n",encryptedData)

	// ใช้ข้อมูลที่เข้ารหัสในการทำ POST request
	resp, err := resty.New().R().
		SetHeader("X-API-ClientID", CLIENT_ID).
		SetHeader("X-API-Signature", enc).
		SetHeader("X-API-Timestamp", fmt.Sprintf("%d", unx)).
		SetHeader("Content-Type", "application/json").
		SetBody(url.QueryEscape(des)). // ใช้ url.QueryEscape แทน encodeURIComponent
		Post(apiURL)

	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	// }

	if err != nil {
		return CResponse{}, fmt.Errorf("error making POST request: %w", err) // แก้ไขที่นี่
	}

	// ทำการ decrypt ผลลัพธ์ที่ได้จาก response
	decryptedResult, _ := encrypt.DecryptDES(string(resp.Body()), DESKEY, DESIV)
	
	// แปลงผลลัพธ์ที่ decrypt เป็น JSON
	var response CResponse
	if err := json.Unmarshal([]byte(decryptedResult), &response); err != nil {
		if err != nil {
			return CResponse{}, fmt.Errorf("error making POST request: %w", err) // แก้ไขที่นี่
		}
	}

	return response,nil
}
func CheckUser(c *fiber.Ctx) error {
	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}

	var authToken,_ = Gsign(request.Account,24,"AuthToken")
	var sessionToken,_ = Gsign(request.Account,1,"SessionToken")
	
	var users models.Users
	db,err := database.GetDatabaseConnection(strings.ToUpper(request.Account))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	fmt.Printf(" Account : %+v \n",strings.ToUpper(request.Account))
	//users = handler.ValidateJWTReturn(request.SessionToken);
	//db, _ := database.ConnectToDB(request.Account)
	var rowsAffected = db.Debug().Where("username = ? AND g_token = ?", strings.ToUpper(request.Account),request.Token).First(&users).RowsAffected
  
	if rowsAffected == 0 {
		var response = fiber.Map{
			"msgId": 4,
			"message": "Invalid GameToken",
			"data": fiber.Map{
				"requestId": request.Requestid,
				"account": request.Account,
				"token": request.Token,
			},
			"timestamp": time.Now().Unix(),
		  }
		  return c.Status(400).JSON(response)
	} else {
		var response = fiber.Map{
			"msgId": 0,
			"message": "OK",
			"data": fiber.Map{
				"requestId": request.Requestid,
				"account": request.Account,
				"authToken": authToken,
				"sessionToken": sessionToken,
			},
			"timestamp": time.Now().Unix(),
			}
		return c.Status(200).JSON(response)

	}	
}
func GetBalance(c *fiber.Ctx) error {

	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//var users models.Users

	tokenString := c.Get("Authorization")[7:] 
	claims, _ :=  Gverify(tokenString)	
	//fmt.Printf("Claims: %+v \n",claims)
	//users,_err = handler.Gverify(tokenString);
	var user models.Users
	db,err := database.GetDatabaseConnection(strings.ToUpper(claims.MemberAccount))
	if err != nil {
		fmt.Printf("err: %+v \n",err)
	}
	//db, _ := database.ConnectToDB(claims.MemberAccount)
	db.Debug().Where("username = ?", claims.MemberAccount).First(&user)
    balanceFloat, _ := user.Balance.Float64()
	var response = fiber.Map{
		"msgId": 0,
		"message": "OK",
		"data": fiber.Map{
			"status": 0,
			"requestId": request.Requestid,
			"account": claims.MemberAccount,
			"balance": balanceFloat,
		},
		"timestamp": time.Now().Unix(),
		}

	return c.Status(200).JSON(response)
	
}
func AddTransactions(transactionsub models.TransactionSub,membername string) Response {

	//fmt.Printf("Add transactionsub: %+v \n",transactionsub)
	response := Response{
		Status: false,
		Message:"Success",
		Data: ResponseBalance{
			BetAmount: decimal.NewFromFloat(0),
			BeforeBalance: decimal.NewFromFloat(0),
			Balance: decimal.NewFromFloat(0),
		},
	}
	 
	var users models.Users
	db,_ := database.GetDatabaseConnection(membername)
    if err_ := db.Where("username = ? ", membername).First(&users).Error; err_ != nil {
		response = Response{
			Status: false,
			Message: "ไม่พบข้อมูล",
			Data:map[string]interface{}{
				"id": -1,
			},
    	}
	}
	//fmt.Printf("ProID : %v \n",users)
    transactionsub.GameProvide = "GCLUB"
    transactionsub.MemberName = membername
	transactionsub.ProductID = transactionsub.ProductID
	transactionsub.BetAmount = transactionsub.BetAmount
	// if transactionsub.TurnOver.IsZero() && transactionsub.Status == 100 {
	// 	transactionsub.TurnOver = transactionsub.BetAmount
	// } else if transactionsub.GameType == 1 && transactionsub.TurnOver.IsZero() && transactionsub.Status == 101 {
	// 	transactionsub.TurnOver = transactionsub.BetAmount
	// }
	if transactionsub.Status == 101 {
		transactionsub.TurnOver =  decimal.Zero
		transactionsub.BetAmount = decimal.Zero
	}
	transactionsub.BeforeBalance = users.Balance
	transactionsub.Balance = users.Balance.Add(transactionsub.TransactionAmount)
	transactionsub.ProID = users.ProStatus
	result := db.Create(&transactionsub); 
	//fmt.Println(result)
	if result.Error != nil {
		response = Response{
			Status: false,
			Message:  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Data: map[string]interface{}{ 
				"id": -1,
			}}
	} else {

		updates := map[string]interface{}{
			"Balance": transactionsub.Balance,
				}
		one,_ := decimal.NewFromString("1") 
		if transactionsub.Balance.LessThan(one) {
			updates["pro_status"] = ""
		}
		_err :=  repository.UpdateFieldsUserString(db,membername, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
 
		if _err != nil {
			fmt.Println("Error:", _err)
		} else {
			fmt.Println("User fields updated successfully")
		}

 
 
	 
	  response = Response{
		Status: true,
		Message: "สำเร็จ",
		Data: ResponseBalance{
			BeforeBalance: transactionsub.BeforeBalance,
			Balance:       transactionsub.Balance,
		},
		}
	}
	return response

}
func Debit(c  *fiber.Ctx) error {
	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//var users models.Users

	tokenString := c.Get("Authorization")[7:] 
	claims, _ :=  Gverify(tokenString)	
	
	//users,_err = handler.Gverify(tokenString);
	var user models.Users
	db,err := database.GetDatabaseConnection(strings.ToUpper(claims.MemberAccount))
	if err != nil {
		fmt.Printf(" 422 err: %+v \n",err)
	}
	//db, _ := database.ConnectToDB(claims.MemberAccount)
	db.Where("username = ?", claims.MemberAccount).First(&user)


    balanceFloat, _ := user.Balance.Float64()

	//fmt.Printf("Request: %+v \n",request)

	var c_transaction_found models.TransactionSub
	rowsAffected := db.Debug().Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",request.Transaction.Id).Find(&c_transaction_found).RowsAffected



	if rowsAffected == 0 {
		c_transaction_found.MemberID = user.ID
        // ปรับใช้ข้อมูลที่ให้มา
        c_transaction_found.MemberName = user.Username
        c_transaction_found.OperatorCode = SYSTEMCODE
        c_transaction_found.OperatorID = 3
        c_transaction_found.ProductID = 2
        c_transaction_found.ProviderID = 0
        c_transaction_found.ProviderLineID = 0
        c_transaction_found.WagerID = 0
        c_transaction_found.CurrencyID = 7
        c_transaction_found.GameType = 2
        c_transaction_found.GameID = request.Game.GameName // สมมุติว่า request มี GameName
        c_transaction_found.GameNumber = request.Game.Run // สมมุติว่า request มี Run
        c_transaction_found.GameRoundID = request.Game.Shoe // สมมุติว่า request มี Shoe
        c_transaction_found.ValidBetAmount = decimal.NewFromFloat(0.0)
        c_transaction_found.BetAmount = request.Transaction.Amount // สมมุติว่า request มี Amount
        c_transaction_found.TurnOver = request.Transaction.Amount
        c_transaction_found.TransactionAmount = decimal.Zero.Sub(request.Transaction.Amount) // สมมุติว่า request มี Amount
        c_transaction_found.TransactionID = request.Transaction.Id
        c_transaction_found.PayoutAmount = decimal.NewFromFloat(0.0)
        c_transaction_found.PayoutDetail = ""
		c_transaction_found.CommissionAmount = decimal.NewFromBigInt(big.NewInt(0), 0) // ใช้ big.NewInt(0) เพื่อสร้างค่า
		c_transaction_found.JackpotAmount = decimal.NewFromBigInt(big.NewInt(0), 0) // ใช้ big.NewInt(0) เพื่อสร้างค่า
        c_transaction_found.SettlementDate = time.Now().Format("20060102150405") // ใช้เวลาปัจจุบันในรูปแบบที่ต้องการ
        c_transaction_found.JPBet = decimal.NewFromFloat(0.0)
        c_transaction_found.Status = 100
        c_transaction_found.BeforeBalance = user.Balance
        c_transaction_found.Balance = user.Balance.Sub(request.Transaction.Amount) // สมมุติว่า request มี Amount
        c_transaction_found.MessageID = request.Transaction.ReferenceId // แปลง decimal.Decimal เป็น string // สมมุติว่า request มี ReferenceId
        c_transaction_found.Sign = ""
        c_transaction_found.RequestTime = time.Now().Format("20060102150405") // ใช้เวลาปัจจุบันในรูปแบบที่ต้องการ
		c_transaction_found.IsAction = "Debit"

		fmt.Printf("467 request.Transaction.Amount %v \n",user.Balance)
		fmt.Printf("468 request.Transaction.Amount %v \n",request.Transaction.Amount)

		if user.Balance.GreaterThanOrEqual(request.Transaction.Amount) {

			result := AddTransactions(c_transaction_found,user.Username)
			fmt.Printf("Add Debit Result: %+v \n",result)
			return c.Status(200).JSON(fiber.Map{
				"msgId": 0,
				"message": "OK",
				"data": fiber.Map{
					//"status": 0,
					"requestId": request.Requestid,
					"account": claims.MemberAccount,
					"transaction": fiber.Map{
						"id": request.Transaction.Id,
						"balance": balanceFloat,
					},
				},
				"timestamp": time.Now().Unix(),
				})
		} else {
			
				return c.Status(400).JSON(fiber.Map{
					"msgId": 1,
					"message": "Amount_over_balance",
					"data": fiber.Map{
						"requestId": request.Requestid,
						"account": claims.MemberAccount,
					},
					"timestamp": time.Now().Unix(),
					})
		}
	}  else {
		return c.Status(400).JSON(fiber.Map{
			"msgId": 2,
			"message": "TransactionId_Duplicate",
			"data": nil,
			"timestamp": time.Now().Unix(),
			})
	}

	// var response = fiber.Map{
	// 	"msgId": 0,
	// 	"message": "OK",
	// 	"data": fiber.Map{
	// 		//"status": 0,
	// 		"requestId": request.Requestid,
	// 		"account": claims.MemberAccount,
	// 		"transaction": fiber.Map{
	// 			"id": request.Transaction.Id,
	// 			"balance": balanceFloat,
	// 		},
	// 	},
	// 	"timestamp": time.Now().Unix(),
	// 	}

	// return c.Status(200).JSON(response)
}


func Credit(c  *fiber.Ctx) error {
	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//var users models.Users

	tokenString := c.Get("Authorization")[7:] 
	claims, _ :=  Gverify(tokenString)	
	//fmt.Printf("Claims: %+v \n",claims)
	//users,_err = handler.Gverify(tokenString);
	var user models.Users
	db,err := database.GetDatabaseConnection(strings.ToUpper(claims.MemberAccount))
	if err != nil {
		fmt.Printf("err: %+v \n",err)
	}
	//db, _ := database.ConnectToDB(claims.MemberAccount)
	db.Where("username = ?", claims.MemberAccount).First(&user)


   // balanceFloat, _ := user.Balance.Float64()
	var c_transaction_found models.TransactionSub
	rowsAffected := db.Debug().Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",request.Transaction.Id).Find(&c_transaction_found).RowsAffected



	if rowsAffected == 0 {
		c_transaction_found.MemberID = user.ID
        // ปรับใช้ข้อมูลที่ให้มา
        c_transaction_found.MemberName = user.Username
        c_transaction_found.OperatorCode = SYSTEMCODE
        c_transaction_found.OperatorID = 3
        c_transaction_found.ProductID = 2
        c_transaction_found.ProviderID = 0
        c_transaction_found.ProviderLineID = 0
        c_transaction_found.WagerID = 0
        c_transaction_found.CurrencyID = 7
        c_transaction_found.GameType = 2
        c_transaction_found.GameID = request.Game.GameName // สมมุติว่า request มี GameName
        c_transaction_found.GameNumber = request.Game.Run // สมมุติว่า request มี Run
        c_transaction_found.GameRoundID = request.Game.Shoe // สมมุติว่า request มี Shoe
        c_transaction_found.ValidBetAmount = decimal.NewFromFloat(0.0)
        c_transaction_found.BetAmount = request.Transaction.Amount // สมมุติว่า request มี Amount
        c_transaction_found.TurnOver = request.Transaction.Amount
        c_transaction_found.TransactionAmount = request.Transaction.Amount // สมมุติว่า request มี Amount
        c_transaction_found.TransactionID = request.Transaction.Id
        c_transaction_found.PayoutAmount = decimal.NewFromFloat(0.0)
        c_transaction_found.PayoutDetail = ""
		c_transaction_found.CommissionAmount = decimal.NewFromBigInt(big.NewInt(0), 0) // ใช้ big.NewInt(0) เพื่อสร้างค่า
		c_transaction_found.JackpotAmount = decimal.NewFromBigInt(big.NewInt(0), 0) // ใช้ big.NewInt(0) เพื่อสร้างค่า
        c_transaction_found.SettlementDate = time.Now().Format("20060102150405") // ใช้เวลาปัจจุบันในรูปแบบที่ต้องการ
        c_transaction_found.JPBet = decimal.NewFromFloat(0.0)
        c_transaction_found.Status = 101
        c_transaction_found.BeforeBalance = user.Balance
        c_transaction_found.Balance = user.Balance.Add(request.Transaction.Amount) // สมมุติว่า request มี Amount
        c_transaction_found.MessageID = request.Transaction.ReferenceId// แปลง decimal.Decimal เป็น string // สมมุติว่า request มี ReferenceId
        c_transaction_found.Sign = ""
        c_transaction_found.RequestTime = time.Now().Format("20060102150405") // ใช้เวลาปัจจุบันในรูปแบบที่ต้องการ
		c_transaction_found.IsAction = "Credit"
	 

		result := AddTransactions(c_transaction_found,user.Username)
		fmt.Printf("Credit Add Result: %+v \n",result)
		//c_transaction_found.Balance 
		balanceFloat, _ := c_transaction_found.Balance.Float64()
			return c.Status(200).JSON(fiber.Map{
				"msgId": 0,
				"message": "OK",
				"data": fiber.Map{
					//"status": 0,
					"requestId": request.Requestid,
					"account": claims.MemberAccount,
					"transaction": fiber.Map{
						"id": request.Transaction.Id,
						"balance": balanceFloat,
					},
				},
				"timestamp": time.Now().Unix(),
				})
	 
	}  else {
		return c.Status(400).JSON(fiber.Map{
			"msgId": 2,
			"message": "TransactionId_Duplicate",
			"data": nil,
			"timestamp": time.Now().Unix(),
			})
	}

	// var response = fiber.Map{
	// 	"msgId": 0,
	// 	"message": "OK",
	// 	"data": fiber.Map{
	// 		//"status": 0,
	// 		"requestId": request.Requestid,
	// 		"account": claims.MemberAccount,
	// 		"transaction": fiber.Map{
	// 			"id": request.Transaction.Id,
	// 			"balance": balanceFloat,
	// 		},
	// 	},
	// 	"timestamp": time.Now().Unix(),
	// 	}

	// return c.Status(200).JSON(response)
}


func CancelBet(c  *fiber.Ctx) error {
	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//var users models.Users

	tokenString := c.Get("Authorization")[7:] 
	claims, _ :=  Gverify(tokenString)	
	fmt.Printf("Claims: %+v \n",claims)
	//users,_err = handler.Gverify(tokenString);
	var user models.Users
	db,err := database.GetDatabaseConnection(strings.ToUpper(claims.MemberAccount))
	if err != nil {
		fmt.Printf("err: %+v \n",err)
	}
	//db, _ := database.ConnectToDB(claims.MemberAccount)
	db.Where("username = ?", claims.MemberAccount).First(&user)


    balanceFloat, _ := user.Balance.Float64()
	 
	var c_transaction_found models.TransactionSub
	rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",request.Transaction.Id).Find(&c_transaction_found).RowsAffected
	rowsAffected_b := db.Model(&models.TransactionSub{}).Select("id").Where("TargetID = ? ",request.Transaction.TargetId).Find(&c_transaction_found).RowsAffected



	if rowsAffected == 0 && rowsAffected_b == 0 {

		rowsAffected_c := db.Model(&models.TransactionSub{}).Select("id").Where("TargetID = ? and GameNumber = ? ",request.Transaction.TargetId,request.Game.Run).Find(&c_transaction_found).RowsAffected

		if rowsAffected_c == 0 {
			return c.Status(400).JSON(fiber.Map{
				"msgId": 4,
				"message": "targetId not Found",
				"data": nil,
				"timestamp": time.Now().Unix(),
				})
		} else {

		c_transaction_found.MemberID = user.ID
        // ปรับใช้ข้อมูลที่ให้มา
        c_transaction_found.MemberName = user.Username
        c_transaction_found.OperatorCode = SYSTEMCODE
        c_transaction_found.OperatorID = 3
        c_transaction_found.ProductID = 2
        c_transaction_found.ProviderID = 0
        c_transaction_found.ProviderLineID = 0
        c_transaction_found.WagerID = 0
        c_transaction_found.CurrencyID = 7
        c_transaction_found.GameType = 2
        c_transaction_found.GameID = request.Game.GameName // สมมุติว่า request มี GameName
        c_transaction_found.GameNumber = request.Game.Run // สมมุติว่า request มี Run
        c_transaction_found.GameRoundID = request.Game.Shoe // สมมุติว่า request มี Shoe
        c_transaction_found.ValidBetAmount = decimal.NewFromFloat(0.0)
        c_transaction_found.BetAmount = decimal.NewFromFloat(0.0)// สมมุติว่า request มี Amount
        c_transaction_found.TurnOver = decimal.NewFromFloat(0.0)
        c_transaction_found.TransactionAmount = c_transaction_found.BetAmount // สมมุติว่า request มี Amount
        c_transaction_found.TransactionID = request.Transaction.Id
        c_transaction_found.PayoutAmount = decimal.NewFromFloat(0.0)
        c_transaction_found.PayoutDetail = ""
		c_transaction_found.CommissionAmount = decimal.NewFromBigInt(big.NewInt(0), 0) // ใช้ big.NewInt(0) เพื่อสร้างค่า
		c_transaction_found.JackpotAmount = decimal.NewFromBigInt(big.NewInt(0), 0) // ใช้ big.NewInt(0) เพื่อสร้างค่า
        c_transaction_found.SettlementDate = time.Now().Format("20060102150405") // ใช้เวลาปัจจุบันในรูปแบบที่ต้องการ
        c_transaction_found.JPBet = decimal.NewFromFloat(0.0)
        c_transaction_found.Status = 102
        c_transaction_found.BeforeBalance = user.Balance
        c_transaction_found.Balance = user.Balance // สมมุติว่า request มี Amount
        c_transaction_found.MessageID = request.Transaction.TargetId// แปลง decimal.Decimal เป็น string // สมมุติว่า request มี ReferenceId
        c_transaction_found.Sign = ""
        c_transaction_found.RequestTime = time.Now().Format("20060102150405") // ใช้เวลาปัจจุบันในรูปแบบที่ต้องการ
		c_transaction_found.IsAction = "CancelBet"
	 

		result := AddTransactions(c_transaction_found,user.Username)
		fmt.Printf("Result: %+v \n",result)
			return c.Status(200).JSON(fiber.Map{
				"msgId": 0,
				"message": "OK",
				"data": fiber.Map{
					//"status": 0,
					"requestId": request.Requestid,
					"account": claims.MemberAccount,
					"transaction": fiber.Map{
						"id": request.Transaction.Id,
						"balance": balanceFloat,
					},
				},
				"timestamp": time.Now().Unix(),
				})
			}
	}  else {
		if rowsAffected == 0 {
		return c.Status(400).JSON(fiber.Map{
			"msgId": 1,
			"message": "TransactionId_Duplicate",
			"data": nil,
			"timestamp": time.Now().Unix(),
			})
		} else {
			return c.Status(400).JSON(fiber.Map{
				"msgId": 4,
				"message": "TargetId_Duplicate",
				"data": nil,
				"timestamp": time.Now().Unix(),
				})
		}
	}

	// var response = fiber.Map{
	// 	"msgId": 0,
	// 	"message": "OK",
	// 	"data": fiber.Map{
	// 		//"status": 0,
	// 		"requestId": request.Requestid,
	// 		"account": claims.MemberAccount,
	// 		"transaction": fiber.Map{
	// 			"id": request.Transaction.Id,
	// 			"balance": balanceFloat,
	// 		},
	// 	},
	// 	"timestamp": time.Now().Unix(),
	// 	}

	// return c.Status(200).JSON(response)
}
func Test(c* fiber.Ctx) error {

	 // "ak+xb8pip08kqqijH/vcAYZ56//9nZWqm/Tu7E2ZpjL4zaHQo91QP+F6wbsZfEhgAH02smpi470="

	 
	data := encrypt.Data{
		Token:  "VBBF",
		Amount: 100.0,
		TranID: "T1_1700154",
	}

	// แปลง struct เป็น JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}



	encode,_err_ := encrypt.EncryptDESAndMD5(string(jsonData),"12345678","98765432","CF7861C7-556F-499A-890C-F9C7C4190266","p@ssw0rd")
	if _err_ != nil {
		fmt.Println(_err_)
	}
	
	//encode := &encrypt.ECResult{} // สมมุติว่า encode เป็น pointer

	result := *encode 
	
	fmt.Printf("Encrypted Data: %+v\n", result)
	//eresult := encrypt.ECResult(encode)
	
	decode,_err := encrypt.DecryptDES(result.Des,"12345678","98765432")
	if _err != nil {
		fmt.Printf("Error Data: %+v\n", _err)
	} else {
		fmt.Printf("Decrypted Data: %+v\n", decode)
	}
	return c.JSON("Test")
}
func Login(c *fiber.Ctx) (error) {
	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "โทเคนหมดอายุ หรือไม่ถูกต้อง",
		})
	}
	var users models.Users

	tokenString := c.Get("Authorization")[7:] 
	claims, _ :=  Gverify(tokenString)
	
	if claims != nil {
	fmt.Println(claims)
	}
	db,err := database.GetDatabaseConnection(strings.ToUpper(request.Account))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	fmt.Printf(" Account : %+v \n",strings.ToUpper(request.Account))

	//db, _ := database.ConnectToDB(request.Account)
	db.Where("username = ?", strings.ToUpper(request.Account)).First(&users)
	//fmt.Println(users)
	loginResponse,err := loging(request.Account)
	updates := map[string]interface{}{
		"g_token": loginResponse.Data.Token,
		}

		repository.UpdateFieldsUserString(db,strings.ToUpper(request.Account), updates) 
	//loginResponse, err := parseLoginResponse(responseString)
	if err != nil {
		
		log.Fatal("Post Error",err)
	 	 
	}

	// แสดงผล
	// fmt.Printf("MsgID: %s\n", loginResponse.MsgID)
	// fmt.Printf("Message: %s\n", loginResponse.Message)
	// fmt.Printf("Data: %s\n", loginResponse.Data)
	// fmt.Printf("TimeStamp: %s\n", loginResponse.Timestamp)
	// fmt.Printf("GroupLimitID: %s\n", loginResponse.GroupLimitID)
	// fmt.Printf("ItemNo: %s\n", loginResponse.ItemNo)
	// fmt.Printf("Lang: %s\n", loginResponse.Lang)
	// fmt.Printf("MemberAccount: %s\n", loginResponse.MemberAccount)
	// fmt.Printf("SystemCode: %s\n", loginResponse.SystemCode)
	// fmt.Printf("WebId: %s\n", loginResponse.WebId)
 
	//  if err != nil {
	// 	log.Fatal("Post Error",err)
	 	 
	//  	} else {
	// 	str_resp := string(resp.Body())
	// 	desenc_str,_xerr := encrypt.DecryptDES(str_resp,handler.DESKEY,handler.DESIV)
		
	// 	if _xerr != nil {
	// 		log.Fatal("Post Error",_xerr)
	// 	}

	// 	fmt.Println(desenc_str)
	// 	//return  str_resp
	//  }

	//var user models.Users
	//database.Database.Where("username = ?", strings.ToUpper(request.Account)).First(&user)
	return c.JSON(loginResponse)
}
func LaunchGame(c *fiber.Ctx) error {

	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	// //var users models.Users
	
	// response := CResponse{
	// 	Status:false,
	// 	Message:"",
	// 	Data: GResponseData{},
	// }
	
	db,err := database.GetDatabaseConnection(strings.ToUpper(request.Username))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	var response,_err = loging(request.Username)
	
	if _err != nil {
	return c.Status(fiber.StatusBadRequest).SendString(_err.Error())
	}

	var user models.Users
	if response.Message == "OK" {
		db.Where("username = ?", strings.ToUpper(request.Username)).First(&user)

		fmt.Printf("User: %+v \n",user)

		// response := fiber.Map{
		// Status:true,
		// Message:"",
		// Data: data_response.data,
		// }
		//token := response.Data
		updates := map[string]interface{}{
		"g_token": response.Data.Token,
		}

		repository.UpdateFieldsUserString(db,strings.ToUpper(request.Username), updates) 
		response.Status =  true
		return c.JSON(response)
	} else if response.Message == "MEMBER_NOT_EXISTS" {
		var xresponse,_ = CreateOrUpdate(request.Username)
		if xresponse.Message == "OK" {
			var login,err = loging(request.Username)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString(err.Error())
			}
			updates := map[string]interface{}{
				"g_token": login.Data.Token,
				}
		
				repository.UpdateFieldsUserString(db,strings.ToUpper(request.Username), updates) 
				login.Status = true
			return c.JSON(login)
		}
	}
	// var loginResponse GResponse

	// jsonData, err := json.Marshal(data)

	// strdata := RequestEncrypt{
	// 	Data: string(jsonData),
	// 	Key: DESKEY,
	// 	Iv: DESIV,
	// }

	
	

	//fmt.Printf(" DB : %+v \n",db)
	
	//db, _ := database.ConnectToDB(request.Username)

	// strdata := fiber.Map{
	// 	"account": request.Username,
	// }
	// resp,err := makePostRequest("http://gservice:9003/LaunchGame",strdata)
	// if err != nil {
	// 	log.Fatalf("Error making POST request: %v", err)
	// }
	// bodyBytes := resp.Body()
	// bodyString := string(bodyBytes)
	

	 

 	// err = json.Unmarshal([]byte(bodyString), &response)
	// if err != nil {
	// 	fmt.Println("Error unmarshalling JSON:", err)
	// 	return err
	// }

	// //var users models.Users
	// updates := map[string]interface{}{
	// 	"g_token": response.Data.Token,
	// 	}

	// repository.UpdateFieldsUserString(db,request.Account, updates) 

 
	
	 
	 return c.JSON(response)
}
func loging(account string) (CResponse,error) {
	 
	
	apiURL := API_URL_PROXY + "/api/Player/Login"

	data := fiber.Map{
		"SystemCode": SYSTEMCODE,
		"WebId": WEBID,
		"MemberAccount": account,
		"ItemNo": "1",
		"BackUrl": "https://tsx.bet/",
		"GroupLimitID": "1,4,12",
		"Lang": "th-TH",
	}

	unx := time.Now().UnixNano() / int64(time.Millisecond)
	
	// แปลง data เป็น JSON string ก่อนเข้ารหัส
	jsonData, err := json.Marshal(data)
	if err != nil {
		//return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return CResponse{}, fmt.Errorf("error making POST request: %w", err)
	}

	// เข้ารหัส DES
	des, err := encrypt.EncryptDES(string(jsonData), DESKEY, DESIV)
	if err != nil {
		//return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return CResponse{}, fmt.Errorf("error making POST request: %w", err)
	}

	// สร้าง MD5 hash
	md5Hash := fmt.Sprintf("%s%s%d%s", CLIENT_ID, CLIENT_SECRET, unx, des)
	hash := md5.Sum([]byte(md5Hash))
	enc := base64.StdEncoding.EncodeToString(hash[:])

	// สร้าง data ที่เข้ารหัส
	// encryptedData := fiber.Map{
	// 	"enc": enc,
	// 	"des": des,
	// 	"unx": unx,
	// }

	//fmt.Printf("encryptedData: %+v \n",encryptedData)

	// ใช้ข้อมูลที่เข้ารหัสในการทำ POST request
	resp, err := resty.New().R().
		SetHeader("X-API-ClientID", CLIENT_ID).
		SetHeader("X-API-Signature", enc).
		SetHeader("X-API-Timestamp", fmt.Sprintf("%d", unx)).
		SetHeader("Content-Type", "application/json").
		SetBody(url.QueryEscape(des)). // ใช้ url.QueryEscape แทน encodeURIComponent
		Post(apiURL)

	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	// }

	if err != nil {
		return CResponse{}, fmt.Errorf("error making POST request: %w", err) // แก้ไขที่นี่
	}

	// ทำการ decrypt ผลลัพธ์ที่ได้จาก response
	decryptedResult, _ := encrypt.DecryptDES(string(resp.Body()), DESKEY, DESIV)
	
	// แปลงผลลัพธ์ที่ decrypt เป็น JSON
	var response CResponse
	if err := json.Unmarshal([]byte(decryptedResult), &response); err != nil {
		if err != nil {
			return CResponse{}, fmt.Errorf("error making POST request: %w", err) // แก้ไขที่นี่
		}
	}

	return response,nil

}
func GetUserOnline(c *fiber.Ctx) (error){
	// var data = fiber.Map{
	// 	"SystemCode": handler.SYSTEMCODE,
	// 	"WebId": handler.WEBID,
	// }
	response := Response{
		Status:false,
		Data: encrypt.GData{},
	}
	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	log.Fatal("Error marshalling JSON:", err)
	// }

	// encode,_enerr := encrypt.EncryptDES([]byte(jsonData),[]byte(handler.DESKEY),[]byte(handler.DESIV))
	// //EncryptDESAndMD5(string(jsonData),handler.DESKEY,handler.DESIV,handler.CLIENT_ID,handler.CLIENT_SECRET)
	
	// if _enerr != nil {
	// 	return _enerr
	// }
	// resultx,err_ := encrypt.DecryptDES([]byte(encode.Des),[]byte(handler.DESKEY),[]byte(handler.DESIV))
	// if err_ != nil {
	// 	return err_
	// }
	// //result := *encode 
	//  fmt.Println(resultx)
	//  resp,err := MakePostRequest(API_URL_PROXY+"/api/Player/GetPlayerOnlineList",encode)	
	//  str_resp := string(resp.Body())

	// //  decode,err := encrypt.DecryptDES(str_resp,[]byte(handler.DESKEY),[]byte(handler.DESIV))
	// //  if err !=nil {
	// // 	 response = Response{
	// // 		 Status:false,
	// // 		 Data: err,
	// // 	 }
	// //  }
	  
 

	return c.JSON(response)
}

func CreateOrUpdate(account string) (*CResponse,error){
	apiURL := API_URL_PROXY + "/api/Player/CreateOrSetUser"

	data := fiber.Map{
		"SystemCode": SYSTEMCODE,
		"WebId": WEBID,
		"MemberAccount": account,
		"MemberName": account,
		"StopBalance": -1,
		"BetLimitGroup": "1,4,12",
		"Currency": "THB",
		"Language": "th-TH",
		"OpenGameList": "ALL",
	}

	unx := time.Now().UnixNano() / int64(time.Millisecond)
	
	// แปลง data เป็น JSON string ก่อนเข้ารหัส
	jsonData, err := json.Marshal(data)
	if err != nil {
		//return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return &CResponse{}, fmt.Errorf("error making POST request: %w", err)
	}

	// เข้ารหัส DES
	des, err := encrypt.EncryptDES(string(jsonData), DESKEY, DESIV)
	if err != nil {
		//return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return &CResponse{}, fmt.Errorf("error making POST request: %w", err)
	}

	// สร้าง MD5 hash
	md5Hash := fmt.Sprintf("%s%s%d%s", CLIENT_ID, CLIENT_SECRET, unx, des)
	hash := md5.Sum([]byte(md5Hash))
	enc := base64.StdEncoding.EncodeToString(hash[:])

	// สร้าง data ที่เข้ารหัส
	// encryptedData := fiber.Map{
	// 	"enc": enc,
	// 	"des": des,
	// 	"unx": unx,
	// }

	//fmt.Printf("encryptedData: %+v \n",encryptedData)

	// ใช้ข้อมูลที่เข้ารหัสในการทำ POST request
	resp, err := resty.New().R().
		SetHeader("X-API-ClientID", CLIENT_ID).
		SetHeader("X-API-Signature", enc).
		SetHeader("X-API-Timestamp", fmt.Sprintf("%d", unx)).
		SetHeader("Content-Type", "application/json").
		SetBody(url.QueryEscape(des)). // ใช้ url.QueryEscape แทน encodeURIComponent
		Post(apiURL)

	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	// }

	if err != nil {
		return &CResponse{}, fmt.Errorf("error making POST request: %w", err) // แก้ไขที่นี่
	}

	// ทำการ decrypt ผลลัพธ์ที่ได้จาก response
	decryptedResult, _ := encrypt.DecryptDES(string(resp.Body()), DESKEY, DESIV)
	
	// แปลงผลลัพธ์ที่ decrypt เป็น JSON
	var response CResponse
	if err := json.Unmarshal([]byte(decryptedResult), &response); err != nil {
		if err != nil {
			return &CResponse{}, fmt.Errorf("error making POST request: %w", err) // แก้ไขที่นี่
		}
	}

	return &response,nil
}


func fastPost(url,clienid string,encoded *encrypt.ECResult) (*fasthttp.Response,error) {

	//url = "http://api.tsxbet.info:8001/api/Player/Login"
	method := "POST"

	// สร้าง request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	req.Header.Set("X-API-ClientID",clienid )
	req.Header.Set("X-API-Signature", encoded.Enc)
	req.Header.Set("X-API-Timestamp",string(encoded.Unx))

	// สร้าง response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	req.SetBody([]byte(encoded.Des))
	// ส่ง request
	 
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)

	return resp, nil
}
func Gsign(account string, expire int64, tokenType string) (string, error) {
	// ถ้า expire คือ ชั่วโมง (h) ให้แปลงเป็นวินาที
	expire = expire * 3600

	// กำหนดเวลาสำหรับ nbf, iat, exp
	now := time.Now().Unix()
	expirationTime := now + expire

	claims := &GClaims{
		SystemCode:   SYSTEMCODE,
		WebID:          WEBID,
		MemberAccount: account,
		TokenType:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Unix(now, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(now, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := fmt.Sprintf("%s%s", CLIENT_ID, CLIENT_SECRET)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func Gverify(tokenString string) (*GClaims, error) {
	secret := fmt.Sprintf("%s%s", CLIENT_ID, CLIENT_SECRET)

	token, err := jwt.ParseWithClaims(tokenString, &GClaims{}, func(token *jwt.Token) (interface{}, error) {
		 
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า Token เป็นของ Claims ที่เรากำหนดไว้
	if claims, ok := token.Claims.(*GClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
func parseLoginResponse(responseString string) (GResponse, error) {
	var loginResponse GResponse

	// แปลง JSON string เป็น struct
	err := json.Unmarshal([]byte(responseString), &loginResponse)
	if err != nil {
		return GResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	return loginResponse, nil
}
func responseToString(resp *fasthttp.Response) string {
	// อ่าน body ของ response
	body := resp.Body()

	// แปลงเป็น string และคืนค่า
	return string(body)
}
func parseResponseToECResult(resp *fasthttp.Response) (*encrypt.ECResult, error) {
	// อ่าน response body
	body := resp.Body()

	// สร้าง struct ECResult เปล่า
	var result encrypt.ECResult

	// Unmarshal JSON response ลงใน struct
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	
	return &result, nil
}
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
func GPostRequest(url ,clienid string,ecresult *encrypt.ECResult) (*fasthttp.Response, error) {
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
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set("X-API-ClientID",clienid )
	req.Header.Set("X-API-Signature", ecresult.Enc)
	req.Header.Set("X-API-Timestamp",string(ecresult.Unx))
	req.SetBody([]byte(ecresult.Des))

	// ส่ง request
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)
	
	return resp, nil
}