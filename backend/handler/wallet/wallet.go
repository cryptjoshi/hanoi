package wallet

import (
	// "context"
	// "fmt"
	// "github.com/amalfra/etag"
	// "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	// "strconv"
	"github.com/shopspring/decimal"
	// "github.com/streadway/amqp"
	// "github.com/tdewolff/minify/v2"
	// "github.com/tdewolff/minify/v2/js"
	// "github.com/valyala/fasthttp"
	// _ "github.com/go-sql-driver/mysql"
	"hanoi/models"
	"gorm.io/gorm"
	//"hanoi/database"
	//"hanoi/handler/jwtn"
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
	// "os"
	// "strconv"
	//"time"
	//"strings"
	"fmt"
	//"errors"
)
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


 

// func GetStatement(c *fiber.Ctx) error {
// 	var statement []models.BankStatement
// 	db.Find(&statement)
// 	return c.Status(200).JSON(statement)
// }

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
			"status": false,
			"message": err_,
			"data": fiber.Map{ 
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
	// 		"status": false,
	// 		"message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
	// 		"data": fiber.Map{ 
	// 			"id": -1,
	// 		}})
	// } else {

		updates := map[string]interface{}{
			"status": BankStatement.Status,
				}
		if err := db.Model(&bankstatement).Updates(updates).Error; err != nil {
			return c.JSON(fiber.Map{
				"status": false,
				"message": err,
				"data": fiber.Map{ 
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
		"status": true,
		"data": fiber.Map{ 
			"id": bankstatement.ID,
			"beforebalance":bankstatement.Beforebalance,
			"balance": bankstatement.Balance,
		}})
	//}
 

}

func checkPro(db *gorm.DB, users *models.Users) (map[string]interface{}, error) {
	var Times struct {
		Type       string `json:"type"`
		Hours      string `json:"hours"`
		Minute     string `json:"minute"`
		DaysOfWeek string `json:"daysofweek"`
	}

	var promotion models.Promotion
	if err := db.Where("id = ?", users.ProStatus).Find(&promotion).Error; err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(promotion.SpecificTime), &Times); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	response := map[string]interface{}{}


	fmt.Println(Times.Type)
	if Times.Type == "first" {
		response["Formular"] = promotion.Example
	}
	return response, nil
}



	// fmt.Println(promotion)
	// fmt.Println(&users)

	// response := []interface{
	// 	MaxDiscount: promotion.MaxDiscount,
	// 	Unit: promotion.Unit,
	// 	Turnamount: promotion.Turnamount,
	// 	Widthdrawmax: promotion.Widthdrawmax
	// 	Includegames: promotion.Includegames,
	// 	Excludegames: promotino.Excludegames
	// }

	// return response


func AddStatement(c *fiber.Ctx) error {

	// user := c.Locals("user").(*jtoken.Token)
	// 	claims := user.Claims.(jtoken.MapClaims)
	var users models.Users
	//var promotion *models.Promotion // Change to pointer type
	// 	prefix := claims["Prefix"].(string)
	// 	if prefix == "" {
	// 		prefix,_ = jwt.GetPrefix(claims["Username"].(string))
	// 	}
	// db, _ := jwt.CheckDBConnection(c.Locals("db"),prefix)
	// _err := jwt.CheckedJWT(db,c);

	// if _err != nil {
	// 	log.Fatal(_err)
	// }

	db,_ := handler.GetDBFromContext(c)
	
	id := c.Locals("ID").(int)
	// if id == nil {
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"status": false,
	// 		"message": "WalletID is missing",
	// 	})
	// }
	//fmt.Printf("WalletID: %v ", id)

	// promotion, err := checkPro(db, &users) // Pass db directly
	// if err != nil {
	// 	fmt.Printf(" %s ",err)
	// }
	// fmt.Println(promotion)

	BankStatement := new(models.BankStatement)

	if err := c.BodyParser(BankStatement); err != nil {
		//fmt.Println(err)
		return c.Status(200).SendString(err.Error())
	}

	//fmt.Println(BankStatement)
	//db, _ := database.ConnectToDB(BankStatement.Prefix)
	 
	 
    if err_ := db.Where("walletid = ? ", id).First(&users).Error; err_ != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"message": err_,
			"data": fiber.Map{ 
				"id": -1,
			}})
    }

	promotion, err := checkPro(db, &users) // Pass db directly
	if err != nil {
		fmt.Printf(" %s ",err)
	}
	fmt.Printf(" %s ",promotion)

	// fmt.Println(BankStatement.Walletid)
	
	BankStatement.Userid = id
	BankStatement.Walletid = id
	BankStatement.BetAmount = BankStatement.BetAmount
	BankStatement.Beforebalance = users.Balance
	BankStatement.Balance = users.Balance.Add(BankStatement.Transactionamount)
	BankStatement.Bankname = users.Bankname
	BankStatement.Accountno = users.Banknumber
	//user.Username = user.Prefix + user.Username
	result := db.Create(&BankStatement); 
	
	if result.Error != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			"data": fiber.Map{ 
				"id": -1,
			}})
	} else {

		updates := map[string]interface{}{
			"Balance": BankStatement.Balance,
				}
	
		//db, _ = database.ConnectToDB(BankStatement.Prefix)
		_err := repository.UpdateUserFields(db, BankStatement.Userid, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		fmt.Println(_err)
		if _err != nil {
			fmt.Println("Error:", _err)
		} else {
			//fmt.Println("User fields updated successfully")
		}

 
 
	 
	 return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data": fiber.Map{ 
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
			   "status":           transaction.Status,
			   "channel": transaction.Channel,
			   "statement_type": transactionType,
			   "expire_date": transaction.CreatedAt,
			   "username": users.Username,
			   "membername": users.Fullname,
		   }
	   }
	
	   return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data": result,
	})
	 
}


func StartWallet(){

	

}
 

