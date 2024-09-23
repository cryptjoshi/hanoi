package handler

import (
	// "context"
	// "fmt"
	// "github.com/amalfra/etag"
	// "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	// "github.com/streadway/amqp"
	// "github.com/tdewolff/minify/v2"
	// "github.com/tdewolff/minify/v2/js"
	// "github.com/valyala/fasthttp"
	// _ "github.com/go-sql-driver/mysql"
	"pkd/models"
	"pkd/database"
	jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/golang-jwt/jwt"
	//jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	"pkd/repository"
	// "log"
	// "net"
	// "net/http"
	"os"
	// "strconv"
	"time"
	//"strings"
	"fmt"
	"errors"
)

type Body struct {
	
	//UserID           int             `json:"userid"`
    //TransactionAmount decimal.Decimal `json:"transactionamount"`
	Username           string             `json:"username"`
	Status           string             `json:"status"`
	Startdate        string 			`json:"startdate"`
	Stopdate        string 		  	`json:"stopdate"`
	Prefix           string           	`json:"prefix`
	Channel        string 		  	`json:"channel"`

}
type UserTransactionSummary struct {
	Walletid           int             `json:"walletid"`
	UserName		string             `json:"username"`
	MemberName		string             `json:"membername"`
	BetAmount decimal.Decimal `json:"betamount"`
	WINLOSS decimal.Decimal `json:"winloss"`
	TURNOVER decimal.Decimal `json:"turnover"`
}
type FirstDeposit struct {
    WalletID          uint
    Username          string
    FirstDepositDate  time.Time
    Firstamount float64 `json:"firstamount"`
}

type Response struct {
    Message string      `json:"message"`
    Status  bool        `json:"status"`
    Data    interface{} `json:"data"` // ใช้ interface{} เพื่อรองรับข้อมูลหลายประเภทใน field data
}


func GetUsers(c *fiber.Ctx) error {
	var user []models.Users
	database.Database.Find(&user)
	response := Response{
		Message: "สำเร็จ!!",
		Status:  true,
		Data: fiber.Map{ 
			"users":user,
		}, 
	}
	return c.JSON(response)
 
}
func GetUser(c *fiber.Ctx) error {
	
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	var users []models.Users
	err := database.Database.Find(&users,claims["ID"]).Error
	
	if err != nil {
		response := Response{
			Message: "ไม่พบรหัสผู้ใช้งาน!!",
			Status:  false,
			
		}
		return c.JSON(response)
	}else {
		response := Response{
			Message: "สำเร็จ!!",
			Status:  true,
			Data: fiber.Map{ 
				"users":users,
			}, 
		}
		return c.JSON(response)
	}
	
}
func GetBalance(c *fiber.Ctx) error {
	
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	
	var users models.Users
	
	database.Database.Debug().Where("id= ?",claims["ID"]).Find(&users)
	fmt.Println(users)
	
	if users == (models.Users{}) {
	 
			response := Response{
				Message: "ไม่พบรหัสผู้ใช้งาน!!",
				Status:  false,
			}
			return c.JSON(response)
	}

	tokenString := c.Get("Authorization")[7:] 
	
	_err := ValidateJWT(tokenString);
	//fmt.Println(_err)
	 if _err != nil {
	   
		  response := Response{
			Message: "โทเคนไม่ถูกต้อง!!",
			Status:  false,
		}
		return c.JSON(response)
	}else {
	//fmt.Println('')
		response := Response{
			Status: true,
			Message: "สำเร็จ",
			Data: fiber.Map{ 
			"balance":users.Balance,
		}}
		return c.JSON(response)
	}
	
	
}
func GetBalanceFromID(c *fiber.Ctx) error {
	
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	var users models.Users
	var bankstatement models.BankStatement
	//fmt.Println(user.Username)
    // ดึงข้อมูลโดยใช้ Username และ Password
    if err := database.Database.Where("username = ?", user.Username).First(&users).Error; err != nil {
        return  errors.New("user not found")
    }
	result := database.Database.Debug().Select("Uid,transactionamount,status").Where("walletid = ? and channel=? ",users.ID,"1stpay").Order("id Desc").First(&bankstatement)
 
	if result.Error != nil {
		 
			fmt.Println("ไม่พบข้อมูล")
			return c.Status(200).JSON(fiber.Map{
				"status": true,
				"data": fiber.Map{ 
					"id": &users.ID,
					"token": &users.Token,
					"balance": &users.Balance,
					"withdraw": fiber.Map{
						"uid":"",
						"transactionamount":0,
						"status":"verified",
					},
				}})
		 
	} else {
		fmt.Printf("ข้อมูลที่พบ: %+v\n", bankstatement)

	 
		return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": fiber.Map{ 
				"id": &users.ID,
				"token": &users.Token,
				"balance": &users.Balance,
				"withdraw": fiber.Map{
					"uid":&bankstatement.Uid,
					"transactionamount": &bankstatement.Transactionamount,
					"status":&bankstatement.Status,
				},
			}})
	}
	//database.Database.Where("username = ?",user.Username).Find(&users)
	
	 
}
func AddUser(c *fiber.Ctx) error {
	
	var currency =  os.Getenv("CURRENCY")
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//user.Walletid = user.ID
	//user.Username = user.Prefix + user.Username + currency
	result := database.Database.Create(&user); 

	
	

	// ส่ง response เป็น JSON



	if result.Error != nil {
		response := Response{
			Message: "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
		} else {

			updates := map[string]interface{}{
				"Walletid":user.ID,
				"Preferredname": user.Username,
				"Username":user.Prefix + user.Username + currency,
			}
			if err := database.Database.Model(&user).Updates(updates).Error; err != nil {
				return errors.New("มีข้อผิดพลาด")
			}
		
		response := Response{
			Message: "เพิ่มยูสเซอร์สำเร็จ!",
			Status:  true,
			Data:    fiber.Map{ 
				"id": user.ID,
				"walletid":user.ID,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	}

	 
}
func GetUserStatement(c *fiber.Ctx) error {

	body := new(Body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	 
	//fmt.Println(body)
	

	user := c.Locals("user").(*jtoken.Token)

	claims := user.Claims.(jtoken.MapClaims)
	
	//fmt.Println(claims)
	//username := claims["username"].(string)
	id := claims["walletid"]

	tokenString := c.Get("Authorization")[7:] 
	
	  _err := ValidateJWT(tokenString);
	  //fmt.Println(_err)
	   if _err != nil {
		return c.JSON(fiber.Map{
					"status": false,
					"message": "โทเคนไม่ถูกต้อง!",
					})
			}
	 
		channel := body.Channel
		startDateStr := body.Startdate
		endDateStr := body.Stopdate
		 

		var statements []models.BankStatement
		 
		if body.Status == "all" {
			database.Database.Debug().Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? ",id, channel, startDateStr, endDateStr).Find(&statements)
		} else {
			database.Database.Debug().Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?", id,channel, startDateStr, endDateStr,body.Status).Find(&statements)
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
			amountFloat, _ := transaction.Transactionamount.Float64()
			   // เก็บผลลัพธ์ใน slice
			   result[i] = fiber.Map{
				   "userid":           transaction.Userid,
				   "createdAt": transaction.CreatedAt,
				   "transfer_amount": amountFloat,
				   "credit":  amountFloat,
				   "status":           transaction.Status,
				   "channel": transaction.Channel,
				   "statement_type": transactionType,
				   "expire_date": transaction.CreatedAt,
			   }
		   }
		
		   return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": result,
		})
	 
}
func GetIdStatement(c *fiber.Ctx) error {

	body := new(Body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	 
	 
	    username := body.Username
		channel := body.Channel
		startDateStr := body.Startdate
		endDateStr := body.Stopdate

	 

		var users models.Users
		database.Database.Debug().Where("username=?",username).Find(&users)

		 

		var statements []models.BankStatement
		 
		if body.Status == "all" {
			database.Database.Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount,COALESCE(bet_amount,0) as TURNOVER,createdAt,Beforebalance,Balance,Channel,Uid").Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? ",users.ID, channel, startDateStr, endDateStr).Find(&statements)
		} else {
			database.Database.Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount,COALESCE(bet_amount,0) as TURNOVER,createdAt,Beforebalance,Balance,Channel,Uid").Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?", users.ID,channel, startDateStr, endDateStr,body.Status).Find(&statements)
		}
		
		  // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
		  result := make([]fiber.Map, len(statements))
	
		  // วนลูปเพื่อประมวลผลแต่ละรายการ
		   for i, transaction := range statements {
			   // ตรวจสอบเงื่อนไขด้วย inline if-else
			   transactionType := func(amount decimal.Decimal,channel string) string {
				if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
					return "เสีย"
				}  
				return "ได้"
			}(transaction.Transactionamount,transaction.Channel)
			amountFloat, _ := transaction.Transactionamount.Float64()
			   // เก็บผลลัพธ์ใน slice
			   result[i] = fiber.Map{
				   "MemberID":           users.Username,
				   "MemberName":  users.Fullname,
				   "createdAt": transaction.CreatedAt,
				   "PayoutAmount": amountFloat,
				   "BetAmount": transaction.BetAmount,
				   "BeforeBalance": transaction.Beforebalance,
				   "AfterBalance": transaction.Balance,
				   "WINLOSS": transaction.Transactionamount,
				   "TURNOVER": transaction.BetAmount,
				   "credit":  amountFloat,
				   "status":           transaction.Status,
				   "GameRoundID": transaction.Uid,
				   "GameProvide": transaction.Channel,
				   "statement_type": transactionType,
				   "expire_date": transaction.CreatedAt,
			   }
		   }
		
		   return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": result,
		})
	 
}
func GetUserAll(c *fiber.Ctx) error {

	body := new(Body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	
	var users []models.Users
	
	database.Database.Debug().Select(" *,CASE WHEN Walletid in (Select distinct walletid from BankStatement Where channel='1stpay') THEN 1 ELSE 0 END As ProStatus ").Find(&users)

	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"data": fiber.Map{ 
			"data":users,
		}})
}
func GetUserAllStatement(c *fiber.Ctx) error {

	body := new(Body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	 
	 
		username := body.Username

		var users models.Users
		database.Database.Debug().Where("username=?",username).Find(&users)

		
		channel := body.Channel
		 
	
		startDateStr := body.Startdate
		endDateStr := body.Stopdate
		 
	
		// ตั้งค่าช่วงวันที่ในการค้นหา
		
	
		var statements []models.BankStatement
		 
		if channel == "game" {
			if body.Status == "all" {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",users.ID, startDateStr, endDateStr).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
			} else {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID , startDateStr, endDateStr,body.Status).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
			}
		} else {
			if body.Status == "all" {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? ",users.ID,channel, startDateStr, endDateStr).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
			} else {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID ,channel, startDateStr, endDateStr,body.Status).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
			}
	
		}
 
		 
		  // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
		  result := make([]fiber.Map, len(statements))
	
		  // วนลูปเพื่อประมวลผลแต่ละรายการ
		   for i, transaction := range statements {
			 
			
			   amountFloat, _ := transaction.Transactionamount.Float64()
			   balanceFloat, _ := transaction.Balance.Float64()
			   beforeFloat,_ := transaction.Beforebalance.Float64()
			   betFloat,_ := transaction.BetAmount.Float64()
			

			   // เก็บผลลัพธ์ใน slice
			   result[i] = fiber.Map{
		 		   "MemberID":           transaction.Walletid,
					"MemberName": users.Username,
					"UserName": users.Fullname,
		 		   "createdAt": transaction.CreatedAt,
				   "GameProvide": transaction.Channel,
				   "GameRoundID": transaction.Uid,
		 		   "BeforeBalance": beforeFloat,
		  		   "BetAmount": betFloat,
				   "WINLOSS": amountFloat,
				   "AfterBalance": balanceFloat,
				   "TURNOVER":betFloat,
		 		   "channel": transaction.Channel,
				   "status":transaction.Status,
	 
		 	   }
		    }
		
		   return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": result,
		})
	 
}
func GetUserDetailStatement(c *fiber.Ctx) error {

	body := new(Body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	 
	 
		username := body.Username

		var users models.Users
		database.Database.Debug().Where("username=?",username).Find(&users)

		
		channel := body.Channel
		 
		startDateStr := body.Startdate
		endDateStr := body.Stopdate
		 
		
	
		var statements []models.BankStatement
		 
		if channel == "game" {
			if body.Status == "all" {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",users.ID, startDateStr, endDateStr).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
			} else {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID , startDateStr, endDateStr,body.Status).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
			}
		} else {
			if body.Status == "all" {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? ",users.ID,channel, startDateStr, endDateStr).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
			} else {
				database.Database.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID ,channel, startDateStr, endDateStr,body.Status).Scan(&statements)
				//database.Database.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
			}
	
		}

		 
		
		   return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": statements,
		})
	 
}
func GetUserSumStatement(c *fiber.Ctx) error {

	body := new(Body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	  
	 
		channel := body.Channel
		startDateStr := body.Startdate
		endDateStr := body.Stopdate
 
		var summaries []UserTransactionSummary
		 
		//if body.Status == "all" {
		err := database.Database.Debug().Model(&models.BankStatement{}).Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,(select fullname FROM Users WHERE Users.id=BankStatement.walletid) as MemberName,COALESCE(sum(bet_amount),0) as BetAmount,sum(transactionamount) as WINLOSS,COALESCE(sum(bet_amount),0) as TURNOVER").Where("channel != ? AND  DATE(createdat) BETWEEN ? AND ? AND (status='101' OR status=0)", channel,startDateStr, endDateStr).Group("BankStatement.walletid").Scan(&summaries).Error
		
		if err != nil {
			fmt.Println(err)
			// ส่ง Error เป็น JSON
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// ส่งข้อมูล JSON
		return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data": summaries,
		})
					
		 
		
		   
	 
}
func UpdateToken(c *fiber.Ctx) error {

	

	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	tokenString := c.Get("Authorization")[7:]
	var users models.Users
	database.Database.Debug().Where("walletid = ?",claims["walletid"]).Find(&users)
	fmt.Println(users)
 
	updates := map[string]interface{}{
		"Token":tokenString,
	}

	 
	_err := repository.UpdateUserFields(database.Database, users.ID, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	if _err != nil {
		return c.Status(200).JSON(fiber.Map{
			"status": false,
			"message":  _err,
		})
	}  

	return c.Status(200).JSON(fiber.Map{
		"status": true,
		"message": "สำเร็จ!",
	 })
}
func GetFirstUsers(c *fiber.Ctx) error {

	type FirstGetResponse struct {
		//Counter           int             `json:"counter"`
		//username         string 	  `json:"username"`
		Walletid           int             `json:"walletid"`
		Firstamount       decimal.Decimal `json:"firstamount"`             
		Firstdate         string 	  `json:"firstdate"`
		
	}
	
	type FirstResponse struct {
		Counter           int             `json:"counter"`
		Active            int  `json:"active"`             
		Firstdate         string 	  `json:"firstdate"`
		
	}
	


	body := new(Body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	 
	 
		username := body.Username
		prefix := body.Prefix

		var users models.Users
		database.Database.Debug().Where("username=?",username).Find(&users)

		
		//channel := body.Channel
		 
	
		startDateStr := body.Startdate
		endDateStr := body.Stopdate
		 
		//var results []FirstGetResponse 
		// var result []struct {
		// 	Username          string
		// 	TransactionAmount float64
		// 	CreatedAt         time.Time
		// }
		
		// subQuery := database.Database.Model(&models.BankStatement{}).
		// Select("walletid, MIN(createdAt) AS firstdate").
		// Where("status = ? AND transactionamount > 0 AND deleted_at IS NULL", "verified").
		// Group("walletid").
		// Having("DATE(MIN(createdAt)) BETWEEN ? AND ?", startDateStr,endDateStr)
		
		// database.Database.Debug().Table("Users AS u").
		// Select("u.walletid,u.username, b.transactionamount AS firstamount, b.createdAt AS firstdate").
		// Joins("JOIN BankStatement AS b ON u.walletid = b.walletid").
		// Joins("JOIN (?) AS first_deposit ON b.walletid = first_deposit.walletid AND b.createdAt = first_deposit.firstdate", subQuery).
		// Where("u.prefix LIKE ?", prefix+"%").
		// Where("u.id IN (SELECT id FROM Users)").
		// Order("u.walletid").
		// Scan(&results)
		var firstDeposits []FirstDeposit
		// ตั้งค่าช่วงวันที่ในการค้นหา
		database.Database.Debug().Model(&models.Users{}).
        Select("Users.id, Users.username, MIN(BankStatement.createdAt) AS first_deposit_date, (SELECT c.transactionamount from BankStatement as c WHERE c.id=Users.id AND DATE(c.createdAt) BETWEEN ? AND ? and c.channel='1stpay' and status='verified' order by c.id ASC LIMIT 1) AS firstamount").
        Joins("JOIN BankStatement ON Users.id = BankStatement.walletid").
        Where("BankStatement.transactionamount > 0").
        Where("DATE(Users.createdAt) BETWEEN ? AND ? ", startDateStr, endDateStr,startDateStr, endDateStr).
        Group("Users.id, Users.username").
        Scan(&firstDeposits) 

	 

		// คำนวณยอดรวมของ first_deposit_amount
		var totalFirstDepositAmount float64
		for _, deposit := range firstDeposits {
			
			totalFirstDepositAmount += deposit.Firstamount
		}

		// แสดงยอดรวม
		fmt.Printf("Total First Deposit Amount: %.2f\n", totalFirstDepositAmount)

		// แสดงผลลัพธ์แต่ละรายการ
		for _, deposit := range firstDeposits {
			fmt.Printf("WalletID: %d, Username: %s, First Deposit Date: %s, First Deposit Amount: %.2f\n",
				deposit.WalletID, deposit.Username, deposit.FirstDepositDate.Format("2006-01-02 15:04:05"), deposit.Firstamount)
		}
	
		var statements FirstResponse
		//var firststate []FirstGetResponse 
		//var secondstate []FirstGetResponse 
		 
		database.Database.Model(&models.Users{}).Debug().Select("count(id) as counter").Where("DATE(createdat) BETWEEN ? AND ?  and username like ?",startDateStr, endDateStr,prefix+"%").Scan(&statements)
		//database.Database.Model(&models.BankStatement{}).Debug().Select("Walletid,min(createdAt) as Firstdate,transactionamount as Firstamount").Where("transactionamount>0 AND  DATE(createdat) BETWEEN ? AND ? and status=?",startDateStr, endDateStr,"verified").Group("walletid,transactionamount").Scan(&firststate)
		 
	
		 
		 
		//   // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
		//   result := make([]fiber.Map, len(firststate))
	
		// //   // วนลูปเพื่อประมวลผลแต่ละรายการ
		//    for i, transaction := range firststate {
		// // 	   // ตรวจสอบเงื่อนไขด้วย inline if-else
		// // 	   transactionType := func(amount decimal.Decimal,channel string) string {
		// // 		if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
		// // 			return "ถอน"
		// // 		}  
		// // 		return "ฝาก"
		// // 	}(transaction.Transactionamount,transaction.Channel)
		// // 	//fmt.Println(transaction.Bet_amount)
		// //fmt.Println(transaction)
		//  	amountFloat, _ := transaction.firstamount.Float64()
		// // 	balanceFloat, _ := transaction.Balance.Float64()
		// // 	beforeFloat,_ := transaction.Beforebalance.Float64()
		// // 	betFloat,_ := transaction.Bet_amount.Float64()
			

		// // 	   // เก็บผลลัพธ์ใน slice
		// 	   result[i] = fiber.Map{
		// 		     "walletid":           transaction.walletid,
		// 		//    "username": users.Username,
		// 		    "firstdate": transaction.firstdate,
		// 			"firstamount":amountFloat,
		// 		//    "beforebalnce": beforeFloat,
		// 		//    "betamount": betFloat,
		// 		//    "transfer_amount": amountFloat,
		// 		//    "balance": balanceFloat,
		// 		//    "credit":  amountFloat,
		// 		//    "status":           transaction.Status,
		// 		//    "channel": transaction.Channel,
		// 		//    "statement_type": transactionType,
		// 		//    "expire_date": transaction.CreatedAt,
		// 	   }
		//    }
			// if firststate==nil {
			// 	return c.Status(200).JSON(fiber.Map{
			// 		"status": true,
			// 		"message": "สำเร็จ!",
			// 		"data": fiber.Map{
			// 			"Counter": statements.Counter,
			// 			"Active": make([]FirstGetResponse, 0),
			// 		},
			// 	 })
				
			// } else {
	 
				return c.Status(200).JSON(fiber.Map{
					"status": true,
					"message": "สำเร็จ!",
					"data": fiber.Map{
						"Counter": statements.Counter,
						"Active": firstDeposits,
						"Firstamount":totalFirstDepositAmount,
					},
				})
		//}
	
}

