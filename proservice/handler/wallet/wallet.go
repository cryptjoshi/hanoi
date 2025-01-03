package wallet

import (
	 "context"
	// "fmt"
	// "github.com/amalfra/etag"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	// "strconv"
	"github.com/shopspring/decimal"
	"github.com/Knetic/govaluate"
	// "github.com/streadway/amqp"
	// "github.com/tdewolff/minify/v2"
	// "github.com/tdewolff/minify/v2/js"
	// "github.com/valyala/fasthttp"
	// _ "github.com/go-sql-driver/mysql"
	"hanoi/models"
	"gorm.io/gorm"
	"hanoi/database"
	jwtn "hanoi/handler/jwtn"
	"hanoi/handler"
	//jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/golang-jwt/jwt"
	//jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	"hanoi/repository"
	"encoding/json"
    "log"
	// "net"
	// "net/http"
	"os"
	"strconv"
	"time"
	"strings"
	"fmt"
	"errors"
	"github.com/go-resty/resty/v2"
	"crypto/sha256"
)

type fistpayBody struct {
	Ref string `json:"ref`
	BankAccountName string `json:"bankAccountName"`
	Amount string `json:"amount"`
	BankCode string `json:"bankCode"`
	BankAccountNo string `json:"bankAccountNo"`
	MerchantURL string `json:"merchantURL"`
}

// var  key struct {
// 	Secret:"25320a8b-cb44-40a4-8456-f6dfff9b735c",
//     Access:"589235b3-0b5c-424b-84ad-2789e4132b33"
// }

type CallbackRequest struct {
	TransactionID   string  `json:"transactionID"`
	MerchantID      string  `json:"merchantID"`
	OwnerID         string  `json:"ownerID"`
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	Fee             float64 `json:"fee"`
	TransferAmount  float64 `json:"transferAmount"`
	BankAccountNo   string  `json:"bankAccountNo"`
	BankAccountName string  `json:"bankAccountName"`
	BankCode        string  `json:"bankCode"`
	Verify          int  `json:"verify"`
	CreateAt        string  `json:"createAt"`
	ExpiredAt       string  `json:"expiredAt"`
	IsExpired       int    `json:"isExpired"`
	Ref             string  `json:"ref"`
	Provider        string  `json:"provider"`
}


type CResponse struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`  
}
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BankBody struct {
	
	UserID           int             `json:"userid"`
	Username         string             `json:"username"`
    //TransactionAmount decimal.Decimal `json:"transactionamount"`
    Status           string             `json:"status"`
	Startdate        string 			`json:"startdate"`
	Stopdate        string 		  	`json:"stopdate"`
	Prefix           string           	`json:"prefix"`
	Channel        string 		  	`json:"channel"`

}
 
type PayInResponse struct {
	Link     string `json:"link"`
	TransactionID string `json:"transactionID"`
	MerchantID string `json:"merchantID"`
	Provider string `json:"provider"`
	Data     struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"data"`
  
}

type TokenRequest struct {
	SecretKey string `json:"secretKey"`
	AccessKey string `json:"accessKey"`
}

type TokenResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}


type PayInRequest struct {
	Ref            string `json:"ref"`
	BankAccountName string `json:"bankAccountName"`
	Amount         string `json:"amount"`
	BankCode       string `json:"bankCode"`
	BankAccountNo  string `json:"bankAccountNo"`
	MerchantURL    string `json:"merchantURL"`
}

var redis_master_host = os.Getenv("REDIS_HOST")
var redis_master_port = os.Getenv("REDIS_PORT")
var redis_master_password = os.Getenv("REDIS_PASSWORD")
var redis_slave_host = os.Getenv("REDIS_SLAVE_HOST")
var redis_slave_port = os.Getenv("REDIS_SLAVE_PORT")
var redis_slave_password = os.Getenv("REDIS_SLAVE_PASSWORD")
var redis_database = getEnv("REDIS_DATABASE", "0")
var secretKey = os.Getenv("PASSWORD_SECRET")
var rdb *redis.Client
var ctx = context.Background()
 

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", // redis_master_password,
		DB:       0,  // Use database 0
	})
}


// ฟังก์ชันที่ใช้บันทึกข้อมูลไปยัง Redis
func publishToRedis(client *redis.Client, channel string, data string) error {
	ctx := context.Background()
	return client.Publish(ctx, channel, data).Err()
}


func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}


func EvaluateExpression(expression string) (decimal.Decimal, error) {
    // Create a new evaluator
    expr, err := govaluate.NewEvaluableExpression(expression)
    if err != nil {
        return decimal.Zero, err
    }

    // Evaluate the expression
    result, err := expr.Evaluate(nil) // Pass any necessary parameters
    if err != nil {
        return decimal.Zero, err
    }

    // Convert the result to decimal.Decimal
    return decimal.NewFromFloat(result.(float64)), nil
}
 

func GetStatement(c *fiber.Ctx) error {
	BankStatement := new(models.BankStatement)

	if err := c.BodyParser(BankStatement); err != nil {
		fmt.Println(err)
		return c.Status(200).SendString(err.Error())
	}

	 
	db,_ := handler.GetDBFromContext(c)

 

	var bankstatement []models.BankStatement
    if err_ := db.Debug().Where("userid = ? ",c.Locals("ID")).Order("id desc").Find(&bankstatement).Error; err_ != nil {
		return c.Status(200).SendString(err_.Error())
    }

	return c.Status(200).JSON(fiber.Map{
		"Status": true,
		"Data": bankstatement,
	})
}

func UpdateStatement(c *fiber.Ctx) error {

	 
	

	BankStatement := new(models.BankStatement)

	if err := c.BodyParser(BankStatement); err != nil {
		fmt.Println(err)
		return c.Status(200).SendString(err.Error())
	}

	//fmt.Println(BankStatement)
	//db, _ := database.ConnectToDB(BankStatement.Prefix)
	 
	db,_ := handler.GetDBFromContext(c)

	var bankstatement models.BankStatement
    if err_ := db.Where("uid = ? ", BankStatement.Uid).First(&bankstatement).Error; err_ != nil {
		return c.JSON(fiber.Map{
			"Status": false,
			"Message": err_,
			"Data": fiber.Map{ 
				"id": -1,
			}})
    }
	 
	//BankStatement.Userid = users.Walletid
	// BankStatement.Beforebalance = users.Balance
	// BankStatement.Balance = users.Balance.Add(BankStatement.Transactionamount)
	// BankStatement.Bankname = users.Bankname
	// BankStatement.Accountno = users.Banknumber
	//user.Username = user.Prefix + user.Username
	//result := db.Create(&BankStatement); 
	
	// if result.Error != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
	// 		"Data": fiber.Map{ 
	// 			"id": -1,
	// 		}})
	// } else {

		updates := map[string]interface{}{
			"Status": BankStatement.Status,
				}
		if err := db.Model(&bankstatement).Updates(updates).Error; err != nil {
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": err,
				"Data": fiber.Map{ 
					"id": -1,
				}})
		}
		 
		// _err := repository.UpdateUserFields(db, BankStatement.Userid, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		// fmt.Println(_err)
		// if _err != nil {
		// 	fmt.Println("Error:", _err)
		// } else {
		// 	//fmt.Println("User fields updated successfully")
		// }

 
 
	 
	 return c.Status(200).JSON(fiber.Map{
		"Status": true,
		"Data": fiber.Map{ 
			"id": bankstatement.ID,
			"beforebalance":bankstatement.Beforebalance,
			"balance": bankstatement.Balance,
		}})
	//}
 

}
type Times struct {
		
	Type       string `json:"type"`
	Hours      string `json:"hours"`
	Minute     string `json:"minute"`
	DaysOfWeek []string `json:"daysofweek"`
}

var ProItem struct {
	UsageLimit int `json:"usagelimit"`
	ProType  Times `json:"protype"`
	Example string `json:"example"`
	Name string `json:"name"`
}
func GetProdetail(db *gorm.DB, procode string) (map[string]interface{}, error) {
	var promotion models.Promotion
	if err := db.Where("id = ?", procode).Find(&promotion).Error; err != nil {
		return nil, err
	}
	if promotion.SpecificTime != "" {
			if err := json.Unmarshal([]byte(promotion.SpecificTime), &ProItem.ProType); err != nil {
			//log.Fatalf("Error unmarshalling JSON: %v", err)
			return nil, err
			}
		} else {
			return nil, nil
		}
		response := make(map[string]interface{}) 
	 
		response["Type"] = ProItem.ProType.Type
		response["count"] = ProItem.UsageLimit
		response["Formular"] = promotion.Example
		response["Name"] = promotion.Name
		if ProItem.ProType.Type == "weekly" {
			response["Week"] = ProItem.ProType.DaysOfWeek
			//response["Week"] = strings.Join(ProItem.ProType.DaysOfWeek, ",") // แปลง array เป็น string ด้วย comma
			// หรือส่ง array ไปเลย
			//response["Week"] = ProItem.ProType.DaysOfWeek
		}
	 
		return response, nil
	
}

func CheckPro(db *gorm.DB, users *models.Users) (map[string]interface{}, error) {
	

	var promotion models.Promotion
	if err := db.Where("id = ?", users.ProStatus).Find(&promotion).Error; err != nil {
		return nil, err
	}
	
 
	if promotion.SpecificTime != "" {
	if err := json.Unmarshal([]byte(promotion.SpecificTime), &ProItem.ProType); err != nil {
		 
		return nil, err
	}
	} else {
		return nil, nil
	}
	// if err_ := json.Unmarshal([]byte(promotion.Example),&ProItem.Example); err_ != nil {
	// 	log.Fatal("Error_ unmarshalling JSON: %v",err_)
	// }

	//fmt.Printf("ProItem: %+v\n",ProItem)

	response := make(map[string]interface{}) // Changed to use make for clarity

	//fmt.Println(ProItem.ProType)
	var promotionlog models.PromotionLog
	// Use a single case statement to handle the different types
	
	var RowsAffected int64
	//db.Debug().Model(&settings).Select("id").Scan(&settings).Count(&RowsAffected)
    db.Debug().Model(&promotionlog).Where("promotioncode = ? and (userid=? or walletid=?) and status=1", users.ProStatus,users.ID,users.ID).Order("id desc").Scan(&promotionlog).Count(&RowsAffected)
	//fmt.Printf("244 line CheckPro\n")
    // fmt.Printf("RowsAffected: %d\n",RowsAffected)
    // fmt.Println(ProItem.UsageLimit)

	//fmt.Printf(" Promotion.UsageLimit: %+v\n",promotion.UsageLimit)
	// Check if promotionlog is not empty or has row affected = 1
	if int64(promotion.UsageLimit) > 0 && RowsAffected > int64(promotion.UsageLimit) { // Assuming ID is the primary key
		return nil, errors.New("คุณใช้งานโปรโมชั่นเกินจำนวนครั้งที่กำหนด")
	}  
		
	//fmt.Printf("254 line CheckPro\n")
	//fmt.Printf("ProItem.ProType: %v\n",ProItem.ProType)

	

	switch ProItem.ProType.Type {
	case "first", "once","weekly","daily","monthly":
		response["minDept"] = promotion.MinDept
		response["maxDept"] = promotion.MaxDiscount
		response["Widthdrawmax"] = promotion.MaxSpend
		response["Widthdrawmin"] = promotion.Widthdrawmin
		response["MinTurnover"] = promotion.MinSpend
		response["Example"] = promotion.Example
		//response["count"] = promotion.UsageLimit
		response["Formular"] = promotion.Example
		response["Name"] = promotion.Name
		response["MinSpendType"] = promotion.MinSpendType
		response["MinCreditType"] = promotion.MinCreditType
		response["MaxWithdrawType"] = promotion.MaxWithdrawType
		response["MinCredit"] = promotion.MinCredit
		response["TurnType"] = promotion.TurnType
		response["Zerobalance"] = promotion.Zerobalance
		response["CreatedAt"] = promotionlog.CreatedAt
		response["MaxUse"] = promotion.UsageLimit
		if ProItem.ProType.Type == "weekly" {
			response["Week"] = ProItem.ProType.DaysOfWeek
		}
	}

	//fmt.Printf(" %s ",response)

	return response, nil
}

func checkActived(db *gorm.DB,ID int) error { //user *models.Users) error {

	var bankstatement []models.BankStatement
	var RowsAffected int64
	 db.Debug().Model(&bankstatement).Where("userid=? or walletid=?",ID,ID).Order("id ASC").Scan(&bankstatement).Count(&RowsAffected)
	
	if RowsAffected >= 1 {
		updates := map[string]interface{}{
			"Actived": time.Now(),
			"Deposit": bankstatement[0].Transactionamount,
				}
	
		//db, _ = database.ConnectToDB(BankStatement.Prefix)
		_err := repository.UpdateUserFields(db, ID, updates) 
		if _err != nil {
			return _err
		}
	}


	return nil
}
 
func Deposit(db *gorm.DB,pro_setting map[string]string,BankStatement *models.BankStatement) (fiber.Map,error) {
	
	users  := models.Users{}
	
	//BankStatement := new(models.BankStatement)

	 Id,err := strconv.Atoi(pro_setting["Id"])
	 if err != nil {
		fmt.Println("Error:", err) 
	 }
	

	 if err := db.Debug().Model(models.Users{}).Where("id = ?",Id).First(&users).Error; err != nil {
		fmt.Println(err)
		return fiber.Map{},fmt.Errorf( "ระบบธนาคาร เกิดข้อผิดผลาด 446 %s",err)
	 }

	//balance,_ := decimal.NewFromString(pro_setting["balance"])
	before,_ := decimal.NewFromString(pro_setting["before_balance"])
    after,_ := decimal.NewFromString(pro_setting["after_balance"])
	bonus,_ := decimal.NewFromString(pro_setting["bonus_amount"])
	amount,_ := decimal.NewFromString(pro_setting["amount"])

	BankStatement.Userid = Id  
	BankStatement.Walletid  = Id
	BankStatement.Beforebalance = before
	BankStatement.Balance = after
	BankStatement.Proamount = bonus
	BankStatement.Transactionamount = amount
  
	
	if users.Bankname == "" && users.Banknumber == "" {
		return fiber.Map{},fmt.Errorf("หมายเลขบัญชีไม่ถูกต้อง")
	}
	
	BankStatement.Bankname = users.Bankname
	BankStatement.Accountno = users.Banknumber

	request := &PayInRequest{
		Ref:            users.Username,
		BankAccountName: users.Fullname,
		Amount:         pro_setting["amount"],
		BankCode:       users.Bankname,
		BankAccountNo:  users.Banknumber,
		MerchantURL:    "https://www.xn--9-twft5c6ayhzf2bxa.com/",
	}
	
	var result PayInResponse
	result, err = paying(request) // เรียกใช้ฟังก์ชัน paying พร้อมส่ง request
	if err != nil {
		return fiber.Map{},fmt.Errorf( "ระบบธนาคาร เกิดข้อผิดผลาด 486 %s",err.Error())
	}
	if result.TransactionID != "" {
		transactionID := result.TransactionID
		BankStatement.Uid = transactionID
		BankStatement.Prefix =  result.MerchantID
	}
	
	BankStatement.Status = "waiting"
	BankStatement.StatementType = "Deposit"
	BankStatement.Detail = result.Link
	// event statement log

	resultz := db.Debug().Create(&BankStatement); 
	if resultz.Error != nil {
		return fiber.Map{},fmt.Errorf( "ระบบธนาคาร เกิดข้อผิดผลาด 501 %s",resultz.Error)
	}



	promotionLog := models.PromotionLog{

		UserID: BankStatement.Userid,
		WalletID: BankStatement.Userid,
		StatementID: BankStatement.ID,
		Promotionname: pro_setting["Name"],
		Beforebalance: BankStatement.Beforebalance,
		//BetAmount: BankStatement.BetAmount,
		Promotioncode: pro_setting["proID"],
		Transactionamount: BankStatement.Transactionamount,
		Promoamount: BankStatement.Balance,
		Proamount: BankStatement.Proamount,
		Balance: BankStatement.Balance,
		Example: pro_setting["Example"],
		Uid: BankStatement.Uid,
		Status: 0,
		// Add other necessary fields for the promotion log
	}
	if err := db.Create(&promotionLog).Error; err != nil {
		return fiber.Map{},fmt.Errorf( "ระบบธนาคาร เกิดข้อผิดผลาด 525 %s",err)
	}
	 


	// if err := db.Create(&promotionLog).Error; err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message": "Failed to log promotion",
	// 		"Data": fiber.Map{
	// 			"id": -1,
	// 		}})
	// }
	return  fiber.Map{ 
			"id": BankStatement.Uid,
			"link": BankStatement.Detail,
			"beforebalance":BankStatement.Beforebalance,
			"balance": BankStatement.Balance,
		},nil
	 
}

func XDeposit(db *gorm.DB,userID string,BankStatement *models.BankStatement) (fiber.Map,error) {

	// user := c.Locals("user").(*jtoken.Token)
	// 	claims := user.Claims.(jtoken.MapClaims)
	var users models.Users
	 
	Id,err := strconv.Atoi(userID)
	if err != nil {
	   fmt.Println("Error:", err) 
	}
	 
	//BankStatement := new(models.BankStatement)
	 
    if err_ := db.Where("id = ? ", Id).First(&users).Error; err_ != nil {
		return  nil,err_ 
	
    }
	 




	BankStatement.Userid = Id
	BankStatement.Walletid = Id
	BankStatement.BetAmount = BankStatement.BetAmount
	BankStatement.Beforebalance = users.Balance
	if BankStatement.Amount.GreaterThan(decimal.Zero) {
	BankStatement.Transactionamount = BankStatement.Amount
	}
	deposit := BankStatement.Transactionamount
	 
	//BankStatement.ProID = users.ProStatus
	//turnoverdef = strings.Replace(users.MinTurnoverDef, "%", "", 1) 
	//var result decimal.Decimal
	//var percentValue decimal.Decimal
	//var percentStr = ""
	//var zeroBalance bool
	//var balanceIncrease decimal.Decimal
	//fmt.Printf("deposit: %v ",deposit)
	// fmt.Printf("TurnType: %v ",pro_setting["TurnType"])
	// fmt.Printf("MinCredit: %v ",pro_setting["MinCredit"])

	//fmt.Printf("pro_setting: %v \n",pro_setting["Zerobalance"])
	BankStatement.Balance = users.Balance.Add(deposit)
 	BankStatement.Bankname = users.Bankname
	BankStatement.Accountno = users.Banknumber
	BankStatement.Proamount = decimal.NewFromFloat(0.0)
	 
	//user.Username = user.Prefix + user.Username
	//fmt.Printf("692 line \n")
	//fmt.Println(BankStatement.Balance)
	//fmt.Println(deposit)

	// if BankStatement.Balance.IsZero() &&  deposit.LessThan(decimal.Zero) {
	// 	users.Turnover = decimal.Zero
	// 	users.ProStatus = ""
	// }

	// if users.Balance.LessThan(deposit.Abs()) && deposit.LessThan(decimal.Zero) {
	// 	fmt.Printf("724 line \n")
	// 	fmt.Printf("users.Balance: %v \n",users.Balance)
	// 	fmt.Printf("deposit.Abs(): %v \n",deposit.Abs())
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message": fmt.Sprintf("ยอดคงเหลือไม่พอ ไม่สามารถถอนเงินได้ %v %v !",users.Balance,users.Currency),
	// 		"Data": fiber.Map{
	// 			"id": -1,
	// 		}})
	// }
 
	// requestBody := map[string]interface{}{
	// 	"ref":            users.Username,
	// 	"bankAccountName": users.Fullname,
	// 	"amount":         deposit,
	// 	"bankCode":       users.Bankname,
	// 	"bankAccountNo":  users.Banknumber,
	// 	"merchantURL":    "https://www.xn--9-twft5c6ayhzf2bxa.com/",
	// }

	if users.Bankname == "" && users.Banknumber == "" {
		 
		return   nil, fmt.Errorf("หมายเลขบัญชีไม่ถูกต้อง") 
	}


	request := &PayInRequest{
		Ref:            users.Username,
		BankAccountName: users.Fullname,
		Amount:         deposit.String(),
		BankCode:       users.Bankname,
		BankAccountNo:  users.Banknumber,
		MerchantURL:    "https://www.xn--9-twft5c6ayhzf2bxa.com/",
	}
	
	var result PayInResponse
	result, err = paying(request) // เรียกใช้ฟังก์ชัน paying พร้อมส่ง request
	if err != nil {
		// จัดการข้อผิดพลาดที่เกิดขึ้น
		return   nil,err 
	}
	
	
	 

	if result.TransactionID != "" {
		transactionID := result.TransactionID
		BankStatement.Uid = transactionID
		BankStatement.Prefix =  result.MerchantID
	 

	} else {
		return   nil, fmt.Errorf("ไม่พบรหัสธุรกรรม")  
	}

	 
	
	BankStatement.Status = "waiting"
	BankStatement.StatementType = "Deposit"
	BankStatement.Detail = result.Link
	// event statement log

	resultz := db.Debug().Create(&BankStatement); 
	

	if resultz.Error != nil {
	 
			return   nil, resultz.Error

	 } else {

		  
	 
	 return  fiber.Map{ 
			"id": BankStatement.Uid,
			"link": BankStatement.Detail,
			"beforebalance":BankStatement.Beforebalance,
			"balance": BankStatement.Balance,
		},nil
	 
	 }

}

func Withdraw(c *fiber.Ctx) error {

    

    var users models.Users
    db, _ := handler.GetDBFromContext(c)
    id := c.Locals("ID").(int)
    BankStatement := new(models.BankStatement)

    // ตรวจสอบ request body
    if err := c.BodyParser(BankStatement); err != nil {
        return c.JSON(fiber.Map{
            "Status": false,
            "Message": "Invalid request body",
            "Data": fiber.Map{"id": -1},
        })
    }

	//fmt.Printf("Withdraw : %+v \n" ,BankStatement)
    //ตรวจสอบข้อมูลผู้ใช้
    if err := db.Where("walletid = ? or id = ?", id,id).First(&users).Error; err != nil {
        return c.JSON(fiber.Map{
            "Status": false,
            "Message": "ไม่พบข้อมูลผู้ใช้",
            "Data": fiber.Map{"id": -1},
        })
    }
	//BankStatement.Transactionamount = BankStatement.Amount
    withdraw := BankStatement.Transactionamount

	withdrawAbs := withdraw.Abs()
	// fmt.Printf("Users: %+v",users)
    // เพิ่มการตรวจสอบยอด ถอนกับยอดคงเหลือ
    if withdrawAbs.GreaterThan(users.Balance) {
		return c.JSON(fiber.Map{
			"Status": false,
			"Message": fmt.Sprintf("ยอดถอน  %v มากกว่ายอดคงเหลือในบัญชี (%v %v)",withdrawAbs, users.Balance, users.Currency),
			"Data": fiber.Map{"id": -1},
		})
    }

	
	
    // // เพิ่มการตรวจสอบ turnover ก่อนการถอน
    // if users.Turnover.IsZero() || users.Turnover.LessThan(decimal.Zero) {
	// 	var lastPromoLog models.PromotionLog
	// 	if err := db.Where("userid = ? AND promotioncode = ? AND status = 1", users.ID, users.ProStatus).
	// 		Order("created_at desc").
	// 		First(&lastPromoLog).Error; err != nil {
	// 		return c.JSON(fiber.Map{
	// 			"Status": false,
	// 			"Message": "ไม่พบข้อมูลโปรโมชั่น",
	// 			"Data": fiber.Map{"id": -1},
	// 		})
	// 	}
	
	// 	// คำนวณ turnover จาก TransactionSub
	// 	var totalTurnover decimal.Decimal
	// 	if err := db.Model(&models.TransactionSub{}).
	// 		Where("membername = ? AND proid = ? AND created_at >= ?", 
	// 			users.Username, 
	// 			users.ProStatus,
	// 			lastPromoLog.CreatedAt).
	// 		Select("COALESCE(SUM(turnover), 0)").
	// 		Scan(&totalTurnover).Error; err != nil {
	// 		return c.JSON(fiber.Map{
	// 			"Status": false,
	// 			"Message": "ไม่สามารถคำนวณยอดเทิร์น",
	// 			"Data": fiber.Map{"id": -1},
	// 		})
	// 	}
	
	// 	if totalTurnover.IsZero() || totalTurnover.LessThan(decimal.Zero) {
	// 		return c.JSON(fiber.Map{
	// 			"Status": false,
	// 			"Message": "ไม่สามารถถอนเงินได้ เนื่องจากยังไม่มียอดเทิร์น",
	// 			"Data": fiber.Map{"id": -1},
	// 		})
	// 	}
    // }
	fmt.Println("906 User ProStatus:",users.ProStatus)
    // ตรวจสอบโปรโมชั่น
    if users.ProStatus != "" {
        pro_setting, err := CheckPro(db, &users)
        if err != nil {
            return c.JSON(fiber.Map{
                "Status": false,
                "Message": "ไม่สามารถตรวจสอบโปรโมชั่นได้",
                "Data": fiber.Map{"id": -1},
            })
        }

        if pro_setting == nil {
            return c.JSON(fiber.Map{
                "Status": false,
                "Message": "ไม่พบข้อมูลโปรโมชั่น",
                "Data": fiber.Map{"id": -1},
            })
        }
	 
		if pro_setting["Widthdrawmin"].(decimal.Decimal).GreaterThan(users.Balance) {
            return c.JSON(fiber.Map{
                "Status": false,
                "Message": fmt.Sprintf("ยอดคงเหลือน้อยกว่ายอดถอนขั้นต่ำของโปรโมชั่น (%v %v)", pro_setting["Widthdrawmin"], users.Currency),
                "Data": fiber.Map{"id": -1},
            })
        }


        // ปรับยอดถอนตามเงื่อนไข
        //withdrawAbs := withdraw.Abs()
        if withdrawAbs.LessThan( pro_setting["Widthdrawmin"].(decimal.Decimal)) {
            withdraw = pro_setting["Widthdrawmin"].(decimal.Decimal).Neg()
        } else if withdrawAbs.GreaterThan(pro_setting["Widthdrawmax"].(decimal.Decimal)) {
            
			if pro_setting["MaxWithdrawType"] == "deposit" {
				withdraw = pro_setting["Widthdrawmax"].(decimal.Decimal).Mul(users.LastDeposit).Neg() 
			} else {
				withdraw = pro_setting["Widthdrawmax"].(decimal.Decimal).Mul(users.LastProamount).Neg() 
			}
        }
		BankStatement.Balance = decimal.Zero 
        // ตรวจสอบตาม turntype
        turnType, ok := pro_setting["TurnType"].(string)
        if !ok {
            return c.JSON(fiber.Map{
                "Status": false,
                "Message": "รูปแบบ TurnType ไม่ถูกต้อง",
                "Data": fiber.Map{"id": -1},
            })
        }

        switch turnType {
        case "turnover":
            totalTurnover,_ := CheckTurnover(db, &users, pro_setting);  
			minTurnover := pro_setting["MinTurnover"].(string)
			var baseAmount decimal.Decimal
            if pro_setting["MinSpendType"] == "deposit" {
                baseAmount = users.LastDeposit
            } else {
                baseAmount = users.LastDeposit.Add(users.LastProamount)
            }

            requiredTurnover, err := CalculateRequiredTurnover(minTurnover, baseAmount)
            
			fmt.Printf(" minTurnover: %v \n",minTurnover)
			fmt.Printf(" baseAmount: %v \n",baseAmount)
			fmt.Printf(" requiredTurnover: %v \n",requiredTurnover)
			fmt.Printf(" totalTurnover: %v \n",totalTurnover)
			fmt.Printf(" userTurnover: %v \n",users.Turnover)
			if err != nil {
				fmt.Printf("err  %s \n",err)
                return c.JSON(fiber.Map{
                    "Status": false,
                    "Message": "ไม่สามารถคำนวณยอดเทิร์นได้",
                    "Data": fiber.Map{"id": -1},
                })
            }

            if totalTurnover.LessThan(requiredTurnover.Mul(baseAmount)) {
                return c.JSON(fiber.Map{
                    "Status": false,
                    "Message": fmt.Sprintf("ยอดเทิร์นไม่เพียงพอ ต้องการ %v %v แต่มี %v %v", 
                        requiredTurnover.Mul(baseAmount),users.Currency, totalTurnover, users.Currency),
                    "Data": fiber.Map{"id": -1},
                })
            }
        case "turncredit":
			
            if err := checkTurnCredit(db, &users, pro_setting); err != nil {
                return c.JSON(fiber.Map{
                    "Status": false,
                    "Message": err.Error(),
                    "Data": fiber.Map{"id": -1},
                })
            }
			
        }
    } else {
		 
		 // เปลี่ยนจาก users.MinTurnover เป็น users.MinTurnoverDef
		 minTurnover, err := CalculateRequiredTurnover(users.MinTurnoverDef, users.LastDeposit)
		 if err != nil {
			 return c.JSON(fiber.Map{
				 "Status": false,
				 "Message": "ค่า MinTurnover ไม่ถูกต้อง",
				 "Data": fiber.Map{"id": -1},
			 })
		 }
		 type TurnoverResult struct {
			Turnover decimal.Decimal
		}
		
		 //var lastWithdrawTurnover decimal.Decimal
		 subQuery := db.Debug().
			 Table("BankStatement").
			 Select("TurnOver").
			 Where("(userid = ? OR walletid = ?) AND statement_type = ?", users.ID, users.ID, "Withdraw").
			 Order("created_at DESC").
			 Limit(1)
		 
		 // Query หลัก
		 var result TurnoverResult
		 err = db.Debug().
			 Table("TransactionSub").
			 Select("SUM(turnover) - COALESCE((?),0) as turnover", subQuery).
			 Where("MemberID = ?", users.ID).
			 Scan(&result).Error

		if err != nil {
			return c.JSON(fiber.Map{
				"Status": false,
				"Message":  err,
				"Data": fiber.Map{"id": -1},
			})
		}
		 
		 fmt.Printf(" 908 Users TurnOver: %v \n",result.Turnover)
		 fmt.Printf(" 909 minTurnover: %v \n",minTurnover)

		 
		 if result.Turnover.LessThan(minTurnover) {
			 return c.JSON(fiber.Map{
				 "Status": false,
				 "Message": fmt.Sprintf("ยอดเทิร์นไม่เพียงพอ ต้องการ %v %v แต่มี %v %v", 
					 minTurnover, users.Currency, users.Turnover, users.Currency),
				 "Data": fiber.Map{"id": -1},
			 })
		 }
		BankStatement.Balance = users.Balance.Sub(withdraw.Abs())
	}

    // บันทึกรายการและอัพเดทข้อมูล
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // บันทึกรายการถอน
    BankStatement.Userid = id
    BankStatement.Walletid = id
    BankStatement.Beforebalance = users.Balance
	BankStatement.ProStatus = users.ProStatus
    // ถ้ามีโปรโมชั่นให้ปรับเป็น 0
    BankStatement.Bankname = users.Bankname
    BankStatement.Accountno = users.Banknumber
    BankStatement.Transactionamount = withdraw

	request := &PayInRequest{
		Ref:            users.Username,
		BankAccountName: users.Fullname,
		Amount:         withdraw.Abs().String(),
		BankCode:       users.Bankname,
		BankAccountNo:  users.Banknumber,
		//MerchantURL:    "https://www.xn--9-twft5c6ayhzf2bxa.com/",
	}
	
	var result PayInResponse
	result, err := Payout(request) // เรียกใช้ฟังก์ชัน paying พร้อมส่ง request
	if err != nil {
		// จัดการข้อผิดพลาดที่เกิดขึ้น
		fmt.Println("Error in paying:", err)
	}
	
	
	fmt.Printf(" Result: %+v \n",result)

	if result.TransactionID != "" {
		transactionID := result.TransactionID
		BankStatement.Uid = transactionID
		BankStatement.Prefix =  result.MerchantID
		fmt.Println("Transaction ID:", transactionID)

	} else {
		fmt.Println("Transaction ID does not exist")
	}


	BankStatement.Status = "waiting"
	BankStatement.StatementType =  "Withdraw"
	
	 


    if err := tx.Create(&BankStatement).Error; err != nil {
        tx.Rollback()
        return c.JSON(fiber.Map{
            "Status": false,
            "Message": "ไม่สามารถบันทึกรายการได้",
            "Data": fiber.Map{"id": -1},
        })
    }

    // อัพเดทข้อมูลผู้ใช้
    updates := map[string]interface{}{
        "Balance": decimal.Zero,
        "LastWithdraw": withdraw,
    }
    fmt.Println("User ProStatus:",users.ProStatus)
    if users.ProStatus != "" {
        updates["ProStatus"] = ""
		if err := tx.Model(&users).Updates(updates).Error; err != nil {
			tx.Rollback()
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": "ไม่สามารถอัพเดทข้อมูลผู้ใช้ได้",
				"Data": fiber.Map{"id": -1},
			})
		}
    } else {
		updates["Balance"] = BankStatement.Balance
		if err := tx.Model(&users).Updates(updates).Error; err != nil {
			tx.Rollback()
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": "ไม่สามารถอัพเดทข้อมูลผู้ใช้ได้",
				"Data": fiber.Map{"id": -1},
			})
		}
	}

    

    if err := tx.Commit().Error; err != nil {
        return c.JSON(fiber.Map{
            "Status": false,
            "Message": "ไม่สามารถบันทึกข้อมูลได้",
            "Data": fiber.Map{"id": -1},
        })
    }

    return c.JSON(fiber.Map{
        "Status": true,
        "Message": "ถอนเงินสำเร็จ",
        "Data": fiber.Map{
            "id": BankStatement.ID,
            "beforebalance": BankStatement.Beforebalance,
            "balance": BankStatement.Balance,
        },
    })
}
func Webhook(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
	var requestBody CallbackRequest

	var deposit,withdraw decimal.Decimal

	message := "ถอนเงินสำเร็จ"		
	
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(200).SendString(err.Error())
	}

	//db,_ := handler.GetDBFromContext(c)

	//fmt.Printf("Body : %+v \n",&requestBody)

	db,_ := database.ConnectToDB(requestBody.MerchantID)

	var bankstatement models.BankStatement

	if err_ := db.Debug().Where("Uid = ? and status = ?",requestBody.TransactionID,"waiting").First(&bankstatement).Error; err_ != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status": false,
			"Message":  "ไม่พบรายการ หรือ มีการปรับปรุงรายการไปแล้ว!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
	}

 

	if  requestBody.Type == "payin" {

		if requestBody.Verify == 1 && requestBody.IsExpired == 0 {
			bankstatement.Status = "verified"
			}  else if requestBody.Verify == 0 && requestBody.IsExpired == 1 {
					bankstatement.Status = "expired"
					bankstatement.Transactionamount = decimal.NewFromFloat(0.0);
			} 


				updates := map[string]interface{}{
					"Balance": bankstatement.Balance,
					 
					//"Turnover": users.Turnover,
					//"ProStatus": users.ProStatus,
					}
				
				userupdate := map[string]interface{}{
					"Balance": bankstatement.Balance,
					//"ProID":"",
					//"Turnover": users.Turnover,
					//"ProStatus": users.ProStatus,
					}
				 
					_err := repository.UpdateUserFields(db, bankstatement.Userid, userupdate) // อัปเดตยูสเซอร์ที่มี ID = 1
				
				if _err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"Status": false,
					"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
					"Data": fiber.Map{ 
						"id": -1,
					}})
				}
				if err := db.Save(&bankstatement).Error; err != nil { // ใช้ db.Save เพื่ออัปเดต bankstatement
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"Status": false,
						"Message": "เกิดข้อผิดพลาดในการอัปเดตข้อมูล",
						"Data": fiber.Map{ 
							"id": -1,
						}})
				}
				//bankstatement.Beforebalance = bankstatement.Balance.Sub(bankstatement.Transactionamount)
			message = "ฝากเงินสำเร็จ"
			deposit = bankstatement.Transactionamount
			withdraw = decimal.Zero

			updates["status"] = 1
			
			if err := db.Debug().Model(&models.PromotionLog{}).Where("uid = ?",bankstatement.Uid).Updates(updates).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"Status": false,
					"Message": "เกิดข้อผิดพลาดในการอัปเดตข้อมูล",
					"Data": fiber.Map{ 
						"id": -1,
					}})
			}

	} else if requestBody.Type == "payout" {
		var balanceamount decimal.Decimal
		if requestBody.Verify == 1 && requestBody.IsExpired == 0 {
			bankstatement.Status = "verified"
			balanceamount = bankstatement.Balance
		}  else if requestBody.Verify == 0 && requestBody.IsExpired == 1 {
				bankstatement.Status = "expired"
				//bankstatement.Beforebalance = bankstatement.Balance.Sub(bankstatement.Transactionamount)
				balanceamount = bankstatement.Beforebalance
				bankstatement.Transactionamount = decimal.NewFromFloat(0.0);
				
				
		} 

		// updates := map[string]interface{}{
		// 	"Balance": decimal.Zero,
		// 	"LastWithdraw": withdraw,
		// }
	
		// if users.ProStatus != "" {
		// 	updates["ProStatus"] = ""
		// }
	
		// if err := tx.Model(&users).Updates(updates).Error; err != nil {
		// 	tx.Rollback()
		// 	return c.JSON(fiber.Map{
		// 		"Status": false,
		// 		"Message": "ไม่สามารถอัพเดทข้อมูลผู้ใช้ได้",
		// 		"Data": fiber.Map{"id": -1},
		// 	})
		// }
		updates := map[string]interface{}{
			"Balance": balanceamount,
			//"Turnover": users.Turnover,
			//"ProStatus": users.ProStatus,
			}
	
		
		_err := repository.UpdateUserFields(db, bankstatement.Userid, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		if _err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status": false,
			"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
		}
		if err := db.Save(&bankstatement).Error; err != nil { // ใช้ db.Save เพื่ออัปเดต bankstatement
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Status": false,
				"Message": "เกิดข้อผิดพลาดในการอัปเดตข้อมูล",
				"Data": fiber.Map{ 
					"id": -1,
				}})
		}
		withdraw = bankstatement.Transactionamount.Abs()
		deposit = decimal.Zero
	}


	updates := map[string]interface{}{
		
		//"Balance": bankstatement.Balance,
		//"LastWithdraw": BankStatement.Transactionamount,
		//"LastDeposit": BankStatement.Transactionamount,
		}
	
	if bankstatement.Status == "verified" {
		if bankstatement.Transactionamount.LessThan(decimal.Zero) {
		
		updates["LastWithdraw"] = bankstatement.Transactionamount
		
		} else {
		updates["LastDeposit"] = bankstatement.Transactionamount
		updates["LastProamount"] = bankstatement.Balance.Sub(bankstatement.Transactionamount)
		}
		updates["uid"] = bankstatement.Uid
		_err := repository.UpdateUserFields(db, bankstatement.Userid, updates)
		if _err != nil {
			return c.Status(200).JSON(fiber.Map{
				"Status": false,
				"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
				"Data": fiber.Map{ 
					"id": -1,
				}})
		} 

	    


		if err := checkActived(db,bankstatement.Userid); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message":  "actived deposit ข้อมูลไม่ได้!",
				"Data": fiber.Map{ 
					"id": -1,
				}})
		}

		data := map[string]interface{}{
			"Deposit": deposit,
			"Withdraw": withdraw,
			"Signup": 0,
			"Prefix": requestBody.MerchantID,
		}


		jsonData, _ := json.Marshal(data)

		// Compress Data
		// compressedData, err := jwtn.CompressData(jsonData)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// คีย์สำหรับการเข้ารหัส (ต้องมีความยาว 16, 24 หรือ 32 bytes)
		key := []byte(secretKey) // คีย์ที่ใช้งานเป็นต้องมีขนาด 16, 24 หรือ 32 bytes
		hashedKey := sha256.Sum256(key)

		encryptedData, err := jwtn.CompressAndEncrypt(jsonData, hashedKey[:])
		if err != nil {
			log.Fatal(err)
		}
		
		// Encrypt Data
		// encryptedData, err := jwtn.Encrypt(compressedData, hashedKey[:])
		// if err != nil {
		// 	log.Fatal(err)
		// }

		fmt.Printf("Data Before: %+v \n",data)
		fmt.Println("Encrypted and compressed data:", encryptedData)

		// decryptedData,err := jwtn.Decrypt(encryptedData,hashedKey[:])

		// decompressedData,err := jwtn.DecompressData(decryptedData)

		// fmt.Println("Decrypted and Decompress ",string(decompressedData))

		// Publish ข้อมูลไปยัง Redis
		channel := "balance_update_channel"
		// if rdb == nil { // เช็คว่า rdb ยังไม่ได้เชื่อมต่อ
		// 	rdb = redis.NewClient(&redis.Options{ // สร้างการเชื่อมต่อใหม่
		// 		Addr:     redis_master_host + ":" + redis_master_port,
		// 		Password: "", // redis_master_password,
		// 		DB:       0,  // ใช้ฐานข้อมูล 0
		// 	})
		// }
		err = publishToRedis(rdb, channel, string(jsonData))//encryptedData)
		if err != nil {
			log.Fatal("Failed to publish to Redis:", err)
		} else {
			fmt.Println("Data successfully published to Redis.")
		}

		if requestBody.Type == "payin" {
		if err := addTransaction(rdb, &bankstatement, requestBody.Type); err != nil {
			return c.Status(200).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),  // อ่านค่าข้อความของ error
				"Data": fiber.Map{
					"id": -1,
				},
			})
		}
		}  else if requestBody.Type == "payout" {
			userid := fmt.Sprintf("%s%d", bankstatement.Prefix, bankstatement.Userid)
			fmt.Printf("userid: %s \n",userid)
			if err := deleteBankStatementsByUser(rdb, userid); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Bank statements deleted successfully for user:", userid)
			}
			
			 
		}

		// now := time.Now()
		// nowRFC3339 := now.Format(time.RFC3339)

		// response := map[string]interface{}{
		//   //	"proID":         fmt.Sprintf("%d", pro_setting["Id"]),
		//     "status":        "verified",
		//     "timestamp":     nowRFC3339,
		// 	//"beforebalance": bankstatement.Balance.Sub(bankstatement.Transactionamount),
		// 	"withdraw": withdraw,
		// 	"deposit": deposit,
		// 	"balance": bankstatement.Balance,
		// 	"prefix": requestBody.MerchantID,
 
		//     }
	   
		// 	if requestBody.Type == "payout"  {
		// 		response["beforebalance"] = bankstatement.Balance.Add(bankstatement.Transactionamount.Abs())
		// 	} else {
		// 		response["beforebalance"] = bankstatement.Balance.Sub(bankstatement.Transactionamount.Abs())

		// 	}

		// currentStatementKey := fmt.Sprintf("%s:%s", bankstatement.Userid,bankstatement.Uid)	 
	   
		// // บันทึกโปรโมชั่นใหม่โดยเขียนทับโปรโมชั่นเก่า
		// if err := redisClient.HSet(ctx, currentStatementKey, response).Err(); err != nil {
		// 	return err
		// }

		return c.JSON(fiber.Map{
			"Status": true,
			"Message": message,
			"Data": fiber.Map{
				"id": bankstatement.Uid,
				"beforebalance": bankstatement.Beforebalance,
				"balance": bankstatement.Balance,
			},
		})
	} else if bankstatement.Status == "expired" {

		userkey := fmt.Sprintf("%s%d",requestBody.MerchantID, bankstatement.Userid)
		key := fmt.Sprintf("%s:current_promotion",  userkey)
		if err := redisClient.HSet(ctx, key, "status", "2", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
			fmt.Printf(" ไม่สามารถปรับสถานะได้ ")
		}
		
		updates["ProStatus"] = ""
		repository.UpdateUserFields(db, bankstatement.Userid, updates)
		return c.Status(200).JSON(fiber.Map{
			"Status": false,
			"Message":  "ธุรกรรม หมดอายุแล้ว!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
	}
	return nil
	}
}
func turn2Percentage(strvalue string) decimal.Decimal {
	var percentValue decimal.Decimal
	var percentStr = ""
	if strings.Contains(strvalue, "%") {
	
		percentStr = strings.TrimSuffix(strvalue, "%")
	    percentValue,_ = decimal.NewFromString(percentStr)
		percentValue = percentValue.Add(decimal.NewFromFloat(100)).Div(decimal.NewFromInt(100))
	} else {
		percentValue,_ = decimal.NewFromString(strvalue)
	}
	// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
	
	return percentValue
}
 
func deleteBankStatementsByUser(rdb *redis.Client, userid string) error {
    keyPattern := fmt.Sprintf("bank_statement:%s:*", userid)
    var cursor uint64

    for {
        // ใช้ SCAN เพื่อค้นหาคีย์ตาม pattern
        keys, nextCursor, err := rdb.Scan(ctx, cursor, keyPattern, 0).Result()
        if err != nil {
            return fmt.Errorf("could not scan keys: %v", err)
        }

        // ถ้ามีคีย์ที่ค้นพบ ให้กำจัดออก
        if len(keys) > 0 {
            if _, err := rdb.Del(ctx, keys...).Result(); err != nil {
                return fmt.Errorf("could not delete keys: %v", err)
            }
        }

        // อัปเดต cursor สำหรับการค้นหาครั้งถัดไป
        cursor = nextCursor
        // ถ้าค่า cursor กลับมาเป็น 0 แสดงว่าการค้นหาสิ้นสุดแล้ว
        if cursor == 0 {
            break
        }
    }

    return nil
}

// เพิ่มฟังก์ชั่นช่วยคำนวณยอดเทิร์นที่ต้องการ
func addTransaction(rdb *redis.Client, bankstatement *models.BankStatement, requestBodytype string) error {
    // ตรวจสอบว่า bankstatement เป็น nil หรือไม่
    if bankstatement == nil {
        return fmt.Errorf("bankstatement is nil")
    }

    now := time.Now()
    nowRFC3339 := now.Format(time.RFC3339)
    //fmt.Printf("BankStatement: %+v \n", bankstatement)

    // สร้างคีย์ใหม่ที่รวม uid และ userID
    userid := fmt.Sprintf("%s%d", bankstatement.Prefix, bankstatement.Userid)
    key := fmt.Sprintf("bank_statement:%s:%s", userid, bankstatement.Uid)
    //fmt.Printf("Key: %s \n", key)

    response := map[string]interface{}{
        "timestamp": nowRFC3339,
        "balance":   bankstatement.Balance.String(),
    }

    // ถ้า requestBodytype เป็น "payout"
    if requestBodytype == "payout" {
     
            response["before_balance"] = bankstatement.Balance.Add(bankstatement.Transactionamount.Abs()).String()
            response["deposit_amount"] = "0"
            response["withdraw_amount"] = bankstatement.Transactionamount.String()

    } else {
 
            response["before_balance"] = bankstatement.Balance.Sub(bankstatement.Transactionamount.Abs()).String()
            response["deposit_amount"] = bankstatement.Transactionamount.String()
			response["total_deposit_amount"] = bankstatement.Balance.String()
            response["withdraw_amount"] = "0"
			//response["deposit_count"] = "0"
			

    }

    // เก็บคีย์ธุรกรรมใน Hash
    err := rdb.HSet(ctx, key, response).Err()
    if err != nil {
        return fmt.Errorf("failed to set transaction in Redis: %v", err)
    }
	// transactionData, err := getTransaction(rdb, userid, bankstatement.Uid)
	// if err != nil {
    //     return fmt.Errorf("failed to set transaction in Redis: %v", err)
    // }
	
   

    return nil
}
func getTransaction(rdb *redis.Client, userID string, uid string) (map[string]string, error) {
    // สร้าง key ของธุรกรรม
    key := fmt.Sprintf("bank_statement:%s:%s", userID, uid)

    // ดึงข้อมูลธุรกรรมทั้งหมดจาก Redis
    transactionData, err := rdb.HGetAll(ctx, key).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to get transaction from Redis: %v", err)
    }

    return transactionData, nil
}
 
// เพิ่มฟังก์ชั่นช่วยคำนวณยอดเทิร์นที่ต้องการ
func CalculateRequiredTurnover(minTurnover string, lastDeposit decimal.Decimal) (decimal.Decimal, error) {
    if strings.Contains(minTurnover, "%") {
        percentStr := strings.TrimSuffix(minTurnover, "%")
        percentValue, err := decimal.NewFromString(percentStr)
        if err != nil {
            return decimal.Zero, err
        }
        return lastDeposit.Mul(percentValue.Div(decimal.NewFromInt(100))), nil
    } else {
		percentValue, err := decimal.NewFromString(minTurnover)
		if err != nil {
            return decimal.Zero, err
        }
		return lastDeposit.Mul(percentValue),nil
	}
    //return decimal.NewFromString(minTurnover)
}
func calculatePercentage(minTurnover string) (decimal.Decimal, error) {
    if strings.Contains(minTurnover, "%") {
        percentStr := strings.TrimSuffix(minTurnover, "%")
        percentValue, err := decimal.NewFromString(percentStr)
        if err != nil {
            return decimal.Zero, err
        }
        return percentValue.Add(decimal.NewFromInt(100)).Div(decimal.NewFromInt(100)), nil
    }
    return decimal.NewFromString(minTurnover)
}
// ฟังก์ชั่นช่วยตรวจสอบ turnover
func CheckTurnover(db *gorm.DB, users *models.Users, pro_setting map[string]interface{}) (decimal.Decimal,error) {

	var promotionLog models.PromotionLog
	db.Where("userid = ? AND promotioncode = ? AND status = 1", users.ID, users.ProStatus).
		Order("created_at DESC").
		First(&promotionLog)


	var totalTurnover decimal.Decimal
	err := db.Debug().Model(&models.TransactionSub{}).
		Where("proid = ? AND membername = ? AND date(created_at) >= date(?)", 
			users.ProStatus, 
			users.Username, 
			promotionLog.CreatedAt).
		Select("COALESCE(SUM(turnover), 0)").
		Scan(&totalTurnover).Error; 
		fmt.Printf( "1086 Err Check: %s",err)
	if err != nil {
		return decimal.Decimal(decimal.NewFromInt(0)),errors.New("ไม่สามารถคำนวณยอดเทิร์นได้")
	}
    // var totalTurnover decimal.Decimal
    // if err := db.Model(&models.TransactionSub{}).
    //     Where("proid = ? AND membername = ?", users.ProStatus, users.Username).
    //     Select("COALESCE(SUM(turnover), 0)").
    //     Scan(&totalTurnover).Error; err != nil {
    //     return errors.New("ไม่สามารถคำนวณยอดเทิร์นได้")
    // }

    minTurnoverStr, ok := pro_setting["MinTurnover"].(string)
    if !ok {
        return decimal.Decimal(decimal.NewFromInt(0)),errors.New("รูปแบบยอดเทิร์นขั้นต่ำไม่ถูกต้อง")
    }

    minTurnover, err := decimal.NewFromString(minTurnoverStr)
    if err != nil {
        return decimal.Decimal(decimal.NewFromInt(0)),errors.New("ไม่สามารถแปลงค่ายอดเทิร์นขั้นต่ำได้")
    }

    if totalTurnover.LessThan(minTurnover) {
        return decimal.Decimal(decimal.NewFromInt(0)),fmt.Errorf("ยอดเทิร์นโอเวอร์น้อยกว่ายอดเทิร์นโอเวอร์ขั้นต่ำ %v %v", minTurnover, users.Currency)
    }

    return totalTurnover,nil
}
func CheckTurnCredit(userID string,minCredit float64,baseAmount float64) error {
    //var lastCredit decimal.Decimal
	//var pro_balance decimal.Decimal
	//var createdAt time.Time
	//createdAt = time.Now() 
	// if pro_setting["CreatedAt"] != nil {
	// 	createdAt = pro_setting["CreatedAt"].(time.Time) 
	// }
  
	// fmt.Printf("ProSetting: %+v ",pro_setting)

	// if err := db.Model(&models.TransactionSub{}).
    //     Where("membername = ? AND deleted_at is null", users.Username).
    //     Order("id desc").
    //     Limit(1).
    //     Select("balance").
    //     Scan(&lastCredit).Error; err != nil {
    //     return errors.New("ไม่สามารถตรวจสอบยอดเครดิต")
    // }

	//db.Debug().Model(&models.TransactionSub{}).Select("balance").Where("membername = ? AND ProID = ? AND deleted_at is null and date(created_at) >= date(?)",users.Username,users.ProStatus,createdAt.Format("2006-01-02 15:04:05")).Limit(1).Order("id desc").Find(&pro_balance)

	 



	 
	

            // requiredTurnover, err := CalculateRequiredTurnover(minTurnover, baseAmount)
            
			// fmt.Printf(" minCredit: %v \n",minTurnover)
			// fmt.Printf(" baseAmount: %v \n",baseAmount)
			// fmt.Printf(" requiredCredit: %v \n",requiredTurnover)
			// //fmt.Printf(" totalTurnover: %v \n",totalTurnover)
			// //fmt.Printf(" userTurnover: %v \n",users.Turnover)
			// if err != nil {
			// 	fmt.Printf("err  %s \n",err)
            //     return c.JSON(fiber.Map{
            //         "Status": false,
            //         "Message": "ไม่สามารถคำนวณยอดเทิร์นได้",
            //         "Data": fiber.Map{"id": -1},
            //     })
            // }

            // if totalTurnover.LessThan(requiredTurnover.Mul(baseAmount)) {
            //     return c.JSON(fiber.Map{
            //         "Status": false,
            //         "Message": fmt.Sprintf("ยอดเทิร์นไม่เพียงพอ ต้องการ %v %v แต่มี %v %v", 
            //             requiredTurnover.Mul(baseAmount),users.Currency, totalTurnover, users.Currency),
            //         "Data": fiber.Map{"id": -1},
            //     })
            // }
	

    return  nil //fmt.Errorf("ยอดเครดิต %v น้อยกว่ายอดเครดิตขั้นต่ำ %v %v", pro_balance, requiredCredit, users.Currency) //nil
}
// ฟังก์ชั่นช่วยตรวจสอบ turncredit
func checkTurnCredit(db *gorm.DB, users *models.Users, pro_setting map[string]interface{}) error {
    //var lastCredit decimal.Decimal
	var pro_balance decimal.Decimal
	var createdAt time.Time
	createdAt = time.Now() 
	if pro_setting["CreatedAt"] != nil {
		createdAt = pro_setting["CreatedAt"].(time.Time) 
	}
  
	fmt.Printf("ProSetting: %+v ",pro_setting)

	// if err := db.Model(&models.TransactionSub{}).
    //     Where("membername = ? AND deleted_at is null", users.Username).
    //     Order("id desc").
    //     Limit(1).
    //     Select("balance").
    //     Scan(&lastCredit).Error; err != nil {
    //     return errors.New("ไม่สามารถตรวจสอบยอดเครดิต")
    // }

	db.Debug().Model(&models.TransactionSub{}).Select("balance").Where("membername = ? AND ProID = ? AND deleted_at is null and date(created_at) >= date(?)",users.Username,users.ProStatus,createdAt.Format("2006-01-02 15:04:05")).Limit(1).Order("id desc").Find(&pro_balance)


    minCreditStr, ok := pro_setting["MinCredit"].(string)
    if !ok {
        return errors.New("รูปแบบยอดเครดิตขั้นต่ำไม่ถูกต้อง")
    }

    minCredit, err := decimal.NewFromString(minCreditStr)
    if err != nil {
        return errors.New("ไม่สามารถแปลงค่ายอดเครดิตขั้นต่ำได้")
    }

	 
	var baseAmount decimal.Decimal
	
	if pro_setting["MinCreditType"] == "deposit" {
		baseAmount = users.LastDeposit
	} else {
		if pro_balance.IsZero() {
			pro_balance = users.LastProamount
		}
		fmt.Printf("ProBalance : %+v \n",pro_balance)
		baseAmount = users.LastDeposit.Add(users.LastProamount)	
	 
	}

            // requiredTurnover, err := CalculateRequiredTurnover(minTurnover, baseAmount)
            
			// fmt.Printf(" minCredit: %v \n",minTurnover)
			// fmt.Printf(" baseAmount: %v \n",baseAmount)
			// fmt.Printf(" requiredCredit: %v \n",requiredTurnover)
			// //fmt.Printf(" totalTurnover: %v \n",totalTurnover)
			// //fmt.Printf(" userTurnover: %v \n",users.Turnover)
			// if err != nil {
			// 	fmt.Printf("err  %s \n",err)
            //     return c.JSON(fiber.Map{
            //         "Status": false,
            //         "Message": "ไม่สามารถคำนวณยอดเทิร์นได้",
            //         "Data": fiber.Map{"id": -1},
            //     })
            // }

            // if totalTurnover.LessThan(requiredTurnover.Mul(baseAmount)) {
            //     return c.JSON(fiber.Map{
            //         "Status": false,
            //         "Message": fmt.Sprintf("ยอดเทิร์นไม่เพียงพอ ต้องการ %v %v แต่มี %v %v", 
            //             requiredTurnover.Mul(baseAmount),users.Currency, totalTurnover, users.Currency),
            //         "Data": fiber.Map{"id": -1},
            //     })
            // }
	requiredCredit := minCredit.Mul(baseAmount)
	fmt.Printf(" minCredit: %v \n",minCredit)
	fmt.Printf(" baseAmount: %v \n",baseAmount)
	fmt.Printf(" requiredCredit: %v \n",requiredCredit)

    if pro_balance.LessThan(requiredCredit) {
        return fmt.Errorf("ยอดเครดิต %v น้อยกว่ายอดเครดิตขั้นต่ำ %v %v", pro_balance, requiredCredit, users.Currency)
    }

    return  nil //fmt.Errorf("ยอดเครดิต %v น้อยกว่ายอดเครดิตขั้นต่ำ %v %v", pro_balance, requiredCredit, users.Currency) //nil
}
func XWithdraw(c *fiber.Ctx) error {

	// user := c.Locals("user").(*jtoken.Token)
	// 	claims := user.Claims.(jtoken.MapClaims)
	var users models.Users


	db,_ := handler.GetDBFromContext(c)
	
	id := c.Locals("ID").(int)

	// if err := db.AutoMigrate(&models.BankStatement{}); err != nil {
	// 	fmt.Println(err)
	// }
	BankStatement := new(models.BankStatement)

	if err := c.BodyParser(BankStatement); err != nil {
	 
		return c.Status(200).SendString(err.Error())
	}


	 
    if err_ := db.Where("walletid = ? ", id).First(&users).Error; err_ != nil {
		return c.JSON(fiber.Map{
			"Status": false,
			"Message": err_,
			"Data": fiber.Map{ 
				"id": -1,
			}})
    }

	pro_setting, err := CheckPro(db, &users) 
	if err != nil {
		
		return c.JSON(fiber.Map{
			"Status": false,
			"Message":  err.Error(),
			"Data": fiber.Map{
				"id": -1,
			}})
	}

	BankStatement.Userid = id
	BankStatement.Walletid = id
	BankStatement.BetAmount = BankStatement.BetAmount
	BankStatement.Beforebalance = users.Balance
	deposit := BankStatement.Transactionamount
	//BankStatement.ProID = users.ProStatus

 
	//turnoverdef = strings.Replace(users.MinTurnoverDef, "%", "", 1) 

	var result decimal.Decimal
	var percentValue decimal.Decimal
	var percentStr = ""
	//var zeroBalance bool
	//fmt.Printf("deposit: %v ",deposit)
	//fmt.Printf("TurnType: %v ",pro_setting["TurnType"])
	//fmt.Printf("MinCredit: %v ",pro_setting["MinCredit"])

	if pro_setting != nil {
		
			if pro_setting["TurnType"] == "turnover" {
					if strings.Contains(pro_setting["MinTurnover"].(string), "%") {
						percentStr = strings.TrimSuffix(pro_setting["MinTurnover"].(string), "%")
					
						//fmt.Printf(" MinturnoverDef : %s %",percentStr)
						// แปลงเป็น float64
						percentValue, _ = decimal.NewFromString(percentStr)
						
						// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
						percentValue = percentValue.Div(decimal.NewFromInt(100))
						
					} else {
						percentStr = pro_setting["MinTurnover"].(string)
						// fmt.Printf(" 492 MinturnoverDef : %s % \n",percentStr)
						percentValue, _ = decimal.NewFromString(percentStr)
					}
					
					//fmt.Printf(" Minturnover : %s ",percentStr)
					// แปลงเป็น float64
					//percentValue, _ := decimal.NewFromString(percentStr)
				
				
					// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
					

					fmt.Printf(" PercentValue: %v \n",percentValue)
					fmt.Printf(" Users.LastDeposit: %v \n",users.LastDeposit)
					fmt.Printf(" Users.LastProamount: %v \n",users.LastProamount)
					if pro_setting["MinSpendType"] == "deposit" {
						result = users.LastDeposit.Mul(percentValue)	
					} else {
						result = users.LastDeposit.Add(users.LastProamount).Mul(percentValue)	
					}

					fmt.Printf("pro_setting: %v \n",pro_setting)
					//minTurnover := users.MinTurnover
					fmt.Printf("bankstatement.Turnover: %v \n",BankStatement.Turnover)
					fmt.Printf("result: %v \n",result)
					//fmt.Printf("minTurnover: %v ",minTurnover)
				if (BankStatement.Turnover.GreaterThan(result) || BankStatement.Turnover.Equal(result)) && BankStatement.Turnover.GreaterThan(decimal.Zero) {
					if deposit.Abs().LessThanOrEqual(pro_setting["Widthdrawmax"].(decimal.Decimal)) {
					BankStatement.Balance = users.Balance.Add(deposit)
					} else {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"Status": false,
							"Message": fmt.Sprintf("ยอดเงินถอนสูงกว่ายอดถอนสูงสุดของโปรโมชั่น %v %v!",pro_setting["Widthdrawmax"],users.Currency),
							"Data": fiber.Map{
								"id": -1,
							}})
					}
					// clear turnover
					
				} else {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"Status": false,
						"Message": fmt.Sprintf("ยอดเทิร์นโอเวอร์น้อยกว่ายอดเทิร์นโอเวอร์ขั้นต่ำ %v %v !",result,users.Currency),
						"Data": fiber.Map{
							"id": -1,
						}})
				}
			} else {
			fmt.Printf("559 line \n")

			
			

			// var bankstatement models.BankStatement
			// db.Debug().Model(&models.BankStatement{}).Select("balance").Where("walletid = ? ",id).Limit(1).Order("id desc").Find(&bankstatement)
			//var transaction models.TransactionSub
			var tbalance decimal.Decimal
			db.Debug().Model(&models.TransactionSub{}).Select("balance").Where("membername = ? AND deleted_at is null and created_at > ?",users.Username,pro_setting["CreatedAt"].(time.Time).Format("2006-01-02 15:04:05")).Limit(1).Order("id desc").Find(&tbalance)
			if strings.Contains(pro_setting["MinCredit"].(string), "%") {
				percentStr = strings.TrimSuffix(pro_setting["MinCredit"].(string), "%")
			
				//fmt.Printf(" MinturnoverDef : %s %",percentStr)
				// แปลงเป็น float64
				percentValue, _ = decimal.NewFromString(percentStr)
				
				// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
				percentValue = percentValue.Add(decimal.NewFromInt(100)).Div(decimal.NewFromInt(100)).Mul(tbalance)
				
			} else {
				percentStr = pro_setting["MinCredit"].(string)
				// fmt.Printf(" 492 MinturnoverDef : %s % \n",percentStr)
				percentValue, _ = decimal.NewFromString(percentStr)
			}
			
			fmt.Printf("tbalance: %v \n",tbalance)
			fmt.Printf("percentValue: %v \n",percentValue)
			fmt.Printf("users.Balance: %v \n",users.Balance)
			if tbalance.LessThanOrEqual(percentValue) == true && users.Balance.GreaterThan(decimal.Zero) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": fmt.Sprintf("ยอดเครดิต %v น้อยกว่ายอดเครดิตขั้นต่ำ %v %v !",tbalance,pro_setting["MinCredit"],users.Currency),
				"Data": fiber.Map{
					"id": -1,
				}})
			} else {
				if deposit.Abs().GreaterThan(pro_setting["Widthdrawmax"].(decimal.Decimal)) {
					//BankStatement.Balance = users.Balance.Sub(pro_setting["Widthdrawmax"].(decimal.Decimal))
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"Status": false,
						"Message": fmt.Sprintf("ย���ดถอนมากว่ายอดถอนสูงสุด %v %v !",pro_setting["Widthdrawmax"],users.Currency),
						"Data": fiber.Map{
							"id": -1,
						}})
				} else if deposit.Abs().LessThan(pro_setting["Widthdrawmin"].(decimal.Decimal)) {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"Status": false,
						"Message": fmt.Sprintf("ยอดถอนน้อยกว่ายอดถอนต่ำสุด %v %v !",pro_setting["Widthdrawmin"],users.Currency),
						"Data": fiber.Map{
							"id": -1,
						}})
				
					} else {
					BankStatement.Balance = users.Balance.Add(deposit)
					fmt.Printf("607 line check user balance more than zero\n")
					if users.Balance.GreaterThan(deposit.Abs()) {
						users.Balance = decimal.Zero
						users.ProStatus = ""
						BankStatement.Balance = decimal.Zero
						
					}
				}
			}

			}
		// } else {
		// 	fmt.Printf("611 line %v \n",users.Balance.LessThanOrEqual(decimal.Zero))
		// 	fmt.Printf("612 line %v \n",pro_setting["ZeroBalance"])
		// 	if users.Balance.LessThanOrEqual(decimal.Zero) == false && pro_setting["ZeroBalance"] == 1 {
				
		// 		response := fiber.Map{
		// 			"Message": "ไม่สามารถ ฝากเงินเพิ่มได้ ขณะใช้งานโปรโมชั่น!",
		// 			"Status":  false,
		// 			"Data": fiber.Map{ 
		// 				"id": -1,
		// 			}}
		// 			return c.JSON(response)
		// 		} else {
		// 			fmt.Printf("622 line \n")
		// 			BankStatement.Balance = users.Balance.Add(deposit)
		// 		}
		// }
	 
		
	} else {
		
		//if deposit.LessThan(decimal.Zero) {

		// check deposit


		 if strings.Contains(users.MinTurnover, "%") {
				percentStr = strings.TrimSuffix(users.MinTurnover, "%")
				fmt.Printf(" MinturnoverDef : %s %",percentStr)
				// แปลงเป็น float64
				percentValue, _ = decimal.NewFromString(percentStr)
		 
		
			// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
				percentValue = percentValue.Div(decimal.NewFromInt(100))
				 
			} else {
				 percentStr = users.MinTurnover
				 fmt.Printf(" MinturnoverDef : %s ",percentStr)
			// แปลงเป็น float64
				percentValue, _ = decimal.NewFromString(percentStr)
		  
			}
			
			fmt.Printf(" Minturnover : %s ",percentStr)
			// แปลงเป็น float64
			//percentValue, _ := decimal.NewFromString(percentStr)
		 
		
			// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
			//percentValue = percentValue.Div(decimal.NewFromInt(100))
		
			// ใช้ในสูตรคำนวณ
			//baseValue := 500.0
			fmt.Printf(" PercentValue: %v ",percentValue)
			result := BankStatement.Transactionamount.Mul(percentValue)

			
			fmt.Printf(" Result : %v \n",result)
			fmt.Printf(" Turnover : %v \n",BankStatement.Turnover)
			
			if BankStatement.Turnover.LessThanOrEqual(result) || BankStatement.Turnover.IsZero() {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Status": false,
					"Message": fmt.Sprintf("ยอดเทิร์นโอเวอร์น้อยกว่ายอดเทิร์นโอเวอร์ขั้นต่ำ %v %v ของยอดฝาก !",users.MinTurnoverDef,users.Currency),
					"Data": fiber.Map{
						"id": -1,
					}})
			}  else {
				fmt.Printf("684 line \n")
				fmt.Printf("users.Balance: %v \n",users.Balance)
				fmt.Printf("deposit: %v \n",deposit)	
				fmt.Printf("users.Balance.LessThan(deposit.Abs()): %v \n",users.Balance.LessThan(deposit.Abs()))
				if users.Balance.LessThan(deposit.Abs()) && deposit.LessThan(decimal.Zero) {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"Status": false,
						"Message": fmt.Sprintf("ยอดคงเหลือไม่พอ ไม่สามารถถอนเงินได้ %v %v !",users.Balance,users.Currency),
						"Data": fiber.Map{
							"id": -1,
						}})
				} else if users.Balance.GreaterThan(decimal.Zero) {
					BankStatement.Balance = users.Balance.Add(deposit)
				}
				//BankStatement.Balance = users.Balance.Add(deposit)
			}
			
			//fmt.Printf("ผลลัพธ์: %.2f\n", result)


		// } 
		// else {
		// BankStatement.Balance = users.Balance.Add(deposit)
		// 	}	
		}

	BankStatement.Bankname = users.Bankname
	BankStatement.Accountno = users.Banknumber
	//user.Username = user.Prefix + user.Username
	
	fmt.Printf("692 line \n")
	fmt.Println(BankStatement.Balance)
	fmt.Println(deposit)
	

	if BankStatement.Balance.IsZero() &&  deposit.LessThan(decimal.Zero) {
		users.Turnover = decimal.Zero
		users.ProStatus = ""
	}

	// if users.Balance.LessThan(deposit.Abs()) && deposit.LessThan(decimal.Zero) {
	// 	fmt.Printf("724 line \n")
	// 	fmt.Printf("users.Balance: %v \n",users.Balance)
	// 	fmt.Printf("deposit.Abs(): %v \n",deposit.Abs())
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message": fmt.Sprintf("ยอดคงเหลือไม่พอ ไม่สามารถถอนเงินได้ %v %v !",users.Balance,users.Currency),
	// 		"Data": fiber.Map{
	// 			"id": -1,
	// 		}})
	// }

	resultz := db.Debug().Create(&BankStatement); 
	


	
	if resultz.Error != nil {
	 
			response := fiber.Map{
				"Message": "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
				"Status":  false,
				"Data":    "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			}
			return c.JSON(response)

	} else {

		updates := map[string]interface{}{
			"Balance": BankStatement.Balance,
			"Turnover": users.Turnover,
			"ProStatus": users.ProStatus,
			}
	
		//db, _ = database.ConnectToDB(BankStatement.Prefix)
	_err := repository.UpdateUserFields(db, BankStatement.Userid, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	if _err != nil {
		return c.Status(200).JSON(fiber.Map{
			"Status": false,
			"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
		}
     //if BankStatement.Transactionamount.LessThan(decimal.Zero) {
		updates["LastWithdraw"] = BankStatement.Transactionamount
	// } else {
	//	updates["LastDeposit"] = BankStatement.Transactionamount
	//	updates["LastProamount"] = BankStatement.Proamount
	// }
	 _err = repository.UpdateUserFields(db, BankStatement.Userid, updates)
		if _err != nil {
			return c.Status(200).JSON(fiber.Map{
				"Status": false,
				"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
				"Data": fiber.Map{ 
					"id": -1,
				}})
		} else {
			//fmt.Println("User fields updated successfully")
		}

 
	// if err := checkActived(db,&users); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message":  "actived deposit ข้อมูลไม่ได้!",
	// 		"Data": fiber.Map{ 
	// 			"id": -1,
	// 		}})
	// }
	 
	 return c.Status(200).JSON(fiber.Map{
		"Status": true,
		"Data": fiber.Map{ 
			"id": BankStatement.ID,
			"beforebalance":BankStatement.Beforebalance,
			"balance": BankStatement.Balance,
		}})
	}
 

}
func GetBankStatement(c *fiber.Ctx) error {

	body := new(BankBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	
	//prefix := body.Username[:3] 
	
	db,_ := handler.GetDBFromContext(c)
	//prefix := c.Locals("Prefix")
	//db, _ := database.ConnectToDB(prefix)
		channel := body.Channel

		if channel != "1stpay" {
			channel = "1stpay"
		}

		startDateStr := body.Startdate
		endDateStr := body.Stopdate
		// loc, _ := time.LoadLocation("Asia/Bangkok")
		 
		// startDate, _ := time.ParseInLocation("2006-01-02", startDateStr,loc)
		// endDate, _ := time.ParseInLocation("2006-01-02 15:04:05", endDateStr+" 23:59:59",loc)
		// currentDate := time.Now().Truncate(24 * time.Hour) // ใช้เวลาในวันนี้เพื่อเปรียบเทียบ

		// if startDate.After(currentDate) {
		// 	startDate = currentDate
		// }
		 
	
		 
	

		var statements []models.BankStatement
		 
		if body.Status == "all" {
			db.Debug().Select("uid,userid,createdAt,accountno,transactionamount,channel,walletid,status").Where(" channel= ? AND  DATE(createdat) BETWEEN ? AND ? ", channel, startDateStr, endDateStr).Order("id desc").Find(&statements)
		} else {
			db.Debug().Select("uid,userid,createdAt,accountno,transactionamount,channel,walletid,status").Where(" channel= ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?",channel, startDateStr, endDateStr,body.Status).Order("id desc").Find(&statements)
		}
	
	  // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
	  result := make([]fiber.Map, len(statements))

	  // วนลูปเพื่อประมวลผลแต่ละรายการ
	   for i, transaction := range statements {
		   // ตรวจสอบเงื่อนไขด้วย inline if-else
		   transactionType := func(amount decimal.Decimal,channel string) string {
			if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
				return "ถอน"
			}  
			return "ฝาก"
		}(transaction.Transactionamount,transaction.Channel)
		 //users models.Users
		 var users models.Users
		 db.Debug().Where("walletid = ?",transaction.Walletid).Find(&users)
		 amountFloat, _ := transaction.Transactionamount.Float64()
		 
		   // เก็บผลลัพธ์ใน slice
		   result[i] = fiber.Map{
				"uid": transaction.Uid,
			   "userid":           transaction.Userid,
			   "createdAt": transaction.CreatedAt,
			   "accountno": transaction.Accountno,
			   "bankname": users.Bankname,
			   "transactionamount": amountFloat,
			   "credit":  amountFloat,
			   "Status":           transaction.Status,
			   "channel": transaction.Channel,
			   "statement_type": transactionType,
			   "expire_date": transaction.CreatedAt,
			   "username": users.Username,
			   "membername": users.Fullname,
		   }
	   }
	
	   return c.Status(200).JSON(fiber.Map{
		"Status": true,
		"Data": result,
	})
	 
}
func makePostRequest(url string,token string, bodyData interface{}) (*fasthttp.Response, error) {
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
	//authHeader := createBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
	//if token != "" {
	req.Header.Add("Authorization","Bearer "+token)
	//}
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
	//authHeader := createBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
	//req.Header.Add("Authorization", authHeader)
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
func paying(request *PayInRequest) (PayInResponse, error) {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetDebug(true) // ตั้งค่า timeout เป็น 30 วินาที

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()


	url := "https://service.1stpay.co/api/payin"
	token := getToken()
	if token.Status {
		fmt.Println("Token received:", token.Data)
	} else {
		fmt.Println("Failed to get token")
	}
	var response PayInResponse

	authToken := fmt.Sprintf("%s %s", "Bearer", token.Data.(string))
 
	num, _ := strconv.Atoi(request.Amount)
	// สร้างคำร้อง
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authToken).
		//SetBody(string(body)). // ใช้ body ที่แปลงแล้ว
		SetBody(map[string]interface{}{
			"ref":            request.Ref,
			"bankAccountName": request.BankAccountName,
			"amount":         num,
			"bankCode":       request.BankCode,
			"bankAccountNo":  request.BankAccountNo,
			"merchantURL":    request.MerchantURL,
		}).
		SetResult(&response).
		Post(url)

	if err != nil {
		fmt.Println("Error making request:", err)
		return response, err
	}

	if resp.IsSuccess() && response.Link != "" {
		fmt.Println("Payment link:", response.Link)
		fmt.Println("Provider:", response.Provider)
		fmt.Println("Response data:", response.Data)
	} else {
		fmt.Println("Failed to get payment link")
	}
	return response, nil
}
func Payout(request *PayInRequest) (PayInResponse, error) {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetDebug(true) // ตั้งค่า timeout เป็น 30 วินาที

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()


	url := "https://service.1stpay.co/api/payout"
	token := getToken()
	if token.Status {
		fmt.Println("Token received:", token.Data)
	} else {
		fmt.Println("Failed to get token")
	}
	var response PayInResponse

	authToken := fmt.Sprintf("%s %s", "Bearer", token.Data.(string))
 
	num, _ := strconv.Atoi(request.Amount)
	// สร้างคำร้อง
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authToken).
		//SetBody(string(body)). // ใช้ body ที่แปลงแล้ว
		SetBody(map[string]interface{}{
			"ref":            request.Ref,
			"bankAccountName": request.BankAccountName,
			"amount":         num,
			"bankCode":       request.BankCode,
			"bankAccountNo":  request.BankAccountNo,
			"merchantURL":    request.MerchantURL,
		}).
		SetResult(&response).
		Post(url)

	if err != nil {
		fmt.Println("Error making request:", err)
		return response, err
	}

	if resp.IsSuccess() && response.Link != "" {
		fmt.Println("Payment link:", response.Link)
		fmt.Println("Provider:", response.Provider)
		fmt.Println("Response data:", response.Data)
	} else {
		fmt.Println("Failed to get payment link")
	}
	return response, nil
}
func NormalTurnover(prefix string,userid string) error {

	db,err := database.ConnectToDB(prefix)
	if err != nil {
		return  err
	}

	var users models.Users
	
	err = db.Debug().Where("id= ?", userid).Find(&users).Error
	if err != nil {
		return  err
	}

	minTurnover, err := CalculateRequiredTurnover(users.MinTurnoverDef, users.LastDeposit)
	if err != nil {
		return  err
	}
	type TurnoverResult struct {
	   Turnover decimal.Decimal
   }
   
	//var lastWithdrawTurnover decimal.Decimal
	subQuery := db.Debug().
		Table("BankStatement").
		Select("TurnOver").
		Where("(userid = ? OR walletid = ?) AND statement_type = ?", users.ID, users.ID, "Withdraw").
		Order("created_at DESC").
		Limit(1)
	
	// Query หลัก
	var result TurnoverResult
	err = db.Debug().
		Table("TransactionSub").
		Select("SUM(turnover) - COALESCE((?),0) as turnover", subQuery).
		Where("MemberID = ?", users.ID).
		Scan(&result).Error

   if err != nil {
		return  err
   }
	
	fmt.Printf(" Users TurnOver: %v \n",result.Turnover)
	fmt.Printf(" minTurnover: %v \n",minTurnover)

	
	if result.Turnover.LessThan(minTurnover) {
		return   fmt.Errorf("ยอดเทิร์นไม่เพียงพอ ต้องการ %v %v แต่มี %v %v", 
				minTurnover, users.Currency, users.Turnover, users.Currency)
	 
	}
	return nil
   //BankStatement.Balance = users.Balance.Sub(withdraw.Abs())
}
// const getTokenPrefix  = async (prefix:string) => {
// 	try {
  
// 	  let raw = ""
// 	  if(prefix.toLocaleLowerCase()=="ckd"){
// 		raw = JSON.stringify({
// 		  "secretKey":"25320a8b-cb44-40a4-8456-f6dfff9b735c",
// 		  "accessKey":"589235b3-0b5c-424b-84ad-2789e4132b33"
// 		})
// 	  } else {
// 	   raw =  JSON.stringify({
// 		  "secretKey":"970dca03-861b-41f5-bd51-024a6a7759a0",
// 		  "accessKey":"8c9ece18-a509-4e14-b1db-53856d56351b"
// 	  })
// 	  }
  
  
//    let res = await fetch("https://service.1stpay.co/api/auth",{
// 	  method: "POST",
// 	  headers: {
// 	  "Content-Type": "application/json",
// 	 // 'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZmx1a3h6eSIsImlhdCI6MTcyNDc2MzE3MSwiZXhwIjoxNzI0ODQ5NTcxfQ.77-N7Ex5mrjW6BAp9HBA2HbAtinELm0Zd7hVEhd0Ehw'
// 		  },
// 	  body: raw,
// 	  redirect: "follow"
//   })
//   const response = await res.json()
//   if(response.data)
//    return {status:true,data: response.data} 
//   else
// 	return {status:false,data: response.message}
//   }
//   catch(err){
// 	return {status:false}
//   }
//   }


// type PayInResponse struct {
// 	Link    string `json:"link"`
// 	Provider string `json:"provider"`
// 	Data     map[string]interface{} `json:"data"`
// }
func getToken() TokenResponse {
	url := "https://service.1stpay.co/api/auth"

	// reqData := TokenRequest{
	// 	SecretKey: "970dca03-861b-41f5-bd51-024a6a7759a0",
	// 	AccessKey: "8c9ece18-a509-4e14-b1db-53856d56351b",
	// }
	reqData := TokenRequest{
  		SecretKey:"25320a8b-cb44-40a4-8456-f6dfff9b735c",
        AccessKey:"589235b3-0b5c-424b-84ad-2789e4132b33",
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return TokenResponse{Status: false}
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(reqBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := fasthttp.Client{}

	if err := client.Do(req, resp); err != nil {
		fmt.Println("Error:", err)
		return TokenResponse{Status: false}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		fmt.Println("Error parsing response:", err)
		return TokenResponse{Status: false}
	}

	if data, ok := result["data"]; ok {
		return TokenResponse{Status: true, Data: data}
	}

	return TokenResponse{Status: false, Data: result["message"]}
}
 


 