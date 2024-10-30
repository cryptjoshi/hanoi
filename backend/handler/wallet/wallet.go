package wallet

import (
	// "context"
	// "fmt"
	// "github.com/amalfra/etag"
	// "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
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
    //"log"
	// "net"
	// "net/http"
	// "os"
	// "strconv"
	"time"
	"strings"
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

func evaluateExpression(expression string) (decimal.Decimal, error) {
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
    if err_ := db.Debug().Where("userid = ? ",c.Locals("ID")).Find(&bankstatement).Error; err_ != nil {
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
type Times struct {
		
	Type       string `json:"type"`
	Hours      string `json:"hours"`
	Minute     string `json:"minute"`
	DaysOfWeek string `json:"daysofweek"`
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
		if ProItem.ProType.Type == "week" {
			response["Week"] = ProItem.ProType.DaysOfWeek
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
		//log.Fatalf("Error unmarshalling JSON: %v", err)
		return nil, err
	}
	} else {
		return nil, nil
	}
	// if err_ := json.Unmarshal([]byte(promotion.Example),&ProItem.Example); err_ != nil {
	// 	log.Fatal("Error_ unmarshalling JSON: %v",err_)
	// }



	response := make(map[string]interface{}) // Changed to use make for clarity

	//fmt.Println(ProItem.ProType)
	var promotionlog models.PromotionLog
	// Use a single case statement to handle the different types
	
	var RowsAffected int64
	//db.Debug().Model(&settings).Select("id").Scan(&settings).Count(&RowsAffected)
    db.Debug().Model(&promotionlog).Where("promotioncode = ? and (userid=? or walletid=?)", users.ProStatus,users.ID,users.ID).Scan(&promotionlog).Count(&RowsAffected)

    // fmt.Println(RowsAffected)
    // fmt.Println(ProItem.UsageLimit)


	// Check if promotionlog is not empty or has row affected = 1
	if int64(ProItem.UsageLimit) > 0 && RowsAffected == int64(ProItem.UsageLimit) { // Assuming ID is the primary key
		return nil, nil
	}  
		
	

	

	switch ProItem.ProType.Type {
	case "first", "once","week":
		response["minDept"] = promotion.MinDept
		response["count"] = ProItem.UsageLimit
		response["Formular"] = promotion.Example
		response["Name"] = promotion.Name
		if ProItem.ProType.Type == "week" {
			response["Week"] = ProItem.ProType.DaysOfWeek
		}
	}

	//fmt.Printf(" %s ",response)

	return response, nil
}

func checkActived(db *gorm.DB,user *models.Users) error {

	var bankstatement []models.BankStatement
	var RowsAffected int64
	 db.Debug().Where("userid=? or walletid=?").Order("id ASC").Scan(&bankstatement).Count(&RowsAffected)

	if RowsAffected >= 1 {
		updates := map[string]interface{}{
			"Activated": time.Now(),
			"Deposit": bankstatement[0].Transactionamount,
				}
	
		//db, _ = database.ConnectToDB(BankStatement.Prefix)
		_err := repository.UpdateUserFields(db, user.ID, updates) 
		if _err != nil {
			return _err
		}
	}


	return nil
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


	db,_ := handler.GetDBFromContext(c)
	
	id := c.Locals("ID").(int)


	BankStatement := new(models.BankStatement)

	if err := c.BodyParser(BankStatement); err != nil {
		//fmt.Println(err)
		return c.Status(200).SendString(err.Error())
	}


	 
    if err_ := db.Where("walletid = ? ", id).First(&users).Error; err_ != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"message": err_,
			"data": fiber.Map{ 
				"id": -1,
			}})
    }

	pro_setting, err := CheckPro(db, &users) 
	if err != nil {
		
		return c.JSON(fiber.Map{
			"status": false,
			"message":  err.Error(),
			"data": fiber.Map{
				"id": -1,
			}})
	}

	BankStatement.Userid = id
	BankStatement.Walletid = id
	BankStatement.BetAmount = BankStatement.BetAmount
	BankStatement.Beforebalance = users.Balance
	
	if pro_setting != nil &&  users.Balance.IsZero() {
		

		// New code to log to promotionlog
		//fmt.Printf("Prosetting: %v ",pro_setting)

		// Ensure pro_setting["Example"] is not nil before type assertion
		Formular, ok := pro_setting["Formular"].(string)
		
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": false,
				"message": "Promotion example is not a valid string",
				"data": fiber.Map{
					"id": -1,
				}})
		}
		//BankStatement.Proamount = BankStatement.Transactionamount
		// Calculate the new balance using the example from pro_setting

	 

		deposit := BankStatement.Transactionamount
		minDept := pro_setting["minDept"].(decimal.Decimal)
		if deposit.LessThan(minDept) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": false,
				"message": "ยอดเงินฝากน้อยกว่ายอดฝากขั้นต่ำของโปรโมชั่น",
				"data": fiber.Map{
					"id": -1,
				}})
		}
		//fmt.Printf(" %v ", deposit)
		//fmt.Printf(" %v ", Formular)

		// Replace 'deposit' in the example with the actual value
		Formular = strings.Replace(Formular, "deposit", deposit.String(), 1) // Convert deposit to string if necessary
		//fmt.Printf(" %v ",Formular)
		// Evaluate the expression (you may need to implement a function to evaluate the string expression)
		balanceIncrease, err := evaluateExpression(Formular) // Implement this function to evaluate the expression
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": "การตั้งค่าสูตรไม่ถูกต้อง!",
				"Data": fiber.Map{
					"id": -1,
				}})
		}
		//fmt.Printf("balanceIncrease: %v ",balanceIncrease)
		BankStatement.Proamount = balanceIncrease.Sub(deposit)
		// Update BankStatement.Balance
		BankStatement.Balance = users.Balance.Add(balanceIncrease)
		promotionLog := models.PromotionLog{

			UserID: BankStatement.Userid,
			WalletID: BankStatement.Userid,
			Promotionname: pro_setting["Name"].(string),
			Beforebalance: BankStatement.Beforebalance,
			//BetAmount: BankStatement.BetAmount,
			Promotioncode: users.ProStatus,
			Transactionamount: deposit,
			Promoamount: balanceIncrease,
			Proamount: balanceIncrease.Sub(deposit),
			Balance: BankStatement.Balance,
			Example: Formular,
			// Add other necessary fields for the promotion log
		}
		if err := db.Create(&promotionLog).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": "Failed to log promotion",
				"Data": fiber.Map{
					"id": -1,
				}})
		}
	} else {
	
		if users.Balance.IsZero() {
		BankStatement.Balance = users.Balance.Add(BankStatement.Transactionamount)
		} else {
			fmt.Printf(" Promotion used: %v", pro_setting)
			fmt.Printf(" Balance is Zero:" ,users.Balance.IsZero())
			response := fiber.Map{
				"Message": "ไม่สามารถ ฝากเงินเพิ่มได้ ขณะใช้งานโปรโมชั่น!",
				"Status":  false,
				"Data": fiber.Map{ 
					"id": -1,
				}}
				return c.JSON(response)
			}
			
			 
		}
	
	BankStatement.Bankname = users.Bankname
	BankStatement.Accountno = users.Banknumber
	//user.Username = user.Prefix + user.Username
	result := db.Create(&BankStatement); 
	
	if result.Error != nil {
	 
			response := fiber.Map{
				"Message": "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
				"Status":  false,
				"Data":    "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			}
			return c.JSON(response)

	} else {

		updates := map[string]interface{}{
			"Balance": BankStatement.Balance,
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
		} else {
			//fmt.Println("User fields updated successfully")
		}

 
	if err := checkActived(db,&users); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Status": false,
			"Message":  "actived deposit ข้อมูลไม่ได้!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
	}
	 
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
			   "status":           transaction.Status,
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



 

