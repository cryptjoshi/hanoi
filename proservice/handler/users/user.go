package users

import (

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	jtoken "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"hanoi/database"
	"hanoi/handler"
	"hanoi/handler/jwtn"
	"hanoi/models"
	"github.com/go-redis/redis/v8" 
	"hanoi/repository"
	"hanoi/encrypt"
	"context" 
	"os"
	"time"
	"fmt"
	"strings"
	"encoding/json"

)

type ErrorResponse struct {
	Status  bool   `json:"Status"`
	Message string `json:"message"`
}
var ctx = context.Background() 
type Body struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	//TransactionAmount decimal.Decimal `json:"transactionamount"`
	Password  string `json:"password"`
	Status    string `json:"Status"`
	Startdate string `json:"startdate"`
	Stopdate  string `json:"stopdate"`
	Prefix    string `json:"prefix`
	Channel   string `json:"channel"`
	Provide   string `json:"provide"`
}

var jwtSecret = os.Getenv("PASSWORD_SECRET")

 



 
func Login(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	loginRequest := new(Body)

	if err := c.BodyParser(loginRequest); err != nil {
		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!" + err.Error(),
			"Status":  false,
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}
	var user models.Users

 

	db, err := database.ConnectToDB(loginRequest.Prefix)
 
//	db.Migrator().CreateTable(&models.Promotion{},&models.PromotionLog{},&models.Users{})
	 
	
	err = db.Where("preferredname = ? AND password = ?", loginRequest.Username, loginRequest.Password).First(&user).Error;
	 
	if err != nil {
		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
			"Status":  false,
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	if err != nil {
		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
			"Status":  false,
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	 

	claims := jtoken.MapClaims{
		"ID":          user.ID,
		"Walletid":    user.Walletid,
		"Username":    user.Username,
		"Fullname":    user.Fullname,
		"Role":        user.Role,
		"PartnersKey": user.PartnersKey,
		"Prefix":      user.Prefix,
	 
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecret))

	updates := map[string]interface{}{
		"Token": t,
	}

	// อัปเดตข้อมูลยูสเซอร์
	_err := repository.UpdateUserFields(db, user.ID, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	if _err != nil {
		if err != nil {
			response := fiber.Map{
				"Message": "ดึงข้อมูลผิดพลาด",
				"Status":  false,
				"Data":    err.Error(),
			}
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
	} else {
		fmt.Println("User fields updated successfully")
	}

	if err != nil {
		response := fiber.Map{
			"Message": "ดึงข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}
	// เชื่อมต่อกับ Redis
	rdb :=  handler.CreateRedisClient()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return  err
		// อาจจะ return error หรือจัดการข้อผิดพลาดตามที่คุณต้องการ
	} else {
		fmt.Println("เชื่อมต่อ Redis สำเร็จ:", pong)
	}

	dbname,db_err := database.GetDBName(user.Prefix)
	if db_err != nil {
		fmt.Printf("Error %s",db_err.Error())
	}

	if err := db.AutoMigrate(&models.TransactionSub{});err != nil {
		fmt.Errorf("Tables schema migration not successfully\n")
	}

	// บันทึกข้อมูลลง Redis
	loginDate := time.Now().Format("2006-01-02 15:04:05") // กำหนดวันที่ login
	loginData := map[string]interface{}{
		"login_date": loginDate,
		"username":   user.Username,
		"prefix":     user.Prefix,
		"database":   dbname, // หรือใช้ชื่อฐานข้อมูลที่ต้องการ
	}

	// อัปเดตข้อมูลใน Redis
	err = rdb.HMSet(ctx,user.Username, loginData).Err()
	if err != nil {
		fmt.Println("Error saving to Redis:", err)
	}


	response := fiber.Map{
		"Token":  t,
		"Data": user,
		"Status": true,
	}
	return c.JSON(response)
 

}

 
func GetUsers(c *fiber.Ctx) error {

	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		//fmt.Println(err)
		return c.Status(200).SendString(err.Error())
	}

	db, _err := handler.GetDBFromContext(c)
	prefix := c.Locals("Prefix")
	if _err != nil {

		response := fiber.Map{
			"Message": "ไม่พบการเชื่อมต่อดาต้าเบส!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}
		return c.JSON(response)
	}

	var users []models.Users

	result := db.Find(&users)

	if result.Error != nil {

		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่่ พบข้อมูล Prefix!",
			Status:  false,
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := fiber.Map{
		"Message": "สำเร็จ!",
		"Status":  true,
		"Data":    &users,
	}

	return c.JSON(response)
}

 
func Register(c *fiber.Ctx) error {

	var currency = os.Getenv("CURRENCY")
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(200).SendString(err.Error())
	}

	// ตรวจสอบว่า user.ReferredBy ไม่พบใน Partner

	// fmt.Printf(" %s ",user)
	db, conn := database.ConnectToDB(user.Prefix)
	if conn != nil {
		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่่ พบข้อมูล Prefix!",
			Status:  false,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	
	var partner models.Partner // สมมุติว่า models.Partner คือโมเดลที่ใช้สำหรับ Partner
	if err := db.Debug().Where("affiliatekey = ?", user.ReferredBy).First(&partner).Error; err != nil {
		// ถ้าไม่พบให้ตั้งค่า ReferredBy เป็นค่าว่าง
		user.ReferredBy = ""
	}
	//fmt.Printf(" partner: %+v \n",partner)

	seedPhrase  := CheckSeed(db)
	//seedPhrase,_ := encrypt.GenerateAffiliateCode(5) //handler.GenerateReferralCode(user.Username,1)

	fmt.Printf("SeedPhase  %s\n", strings.ToUpper(user.Prefix) + user.Username + currency) 
   
	user.ReferralCode = seedPhrase

	fmt.Printf("user: %s \n",user)
	
	result := db.Debug().Create(&user)

	if result.Error != nil {
		fmt.Printf(" Error: %+v \n",result.Error)
		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Status:  false,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	} else {

		updates := map[string]interface{}{
			"Walletid":      user.ID,
			"Preferredname": user.Username,
			"Username":      strings.ToUpper(user.Prefix) + user.Username + currency,
			"Currency":      currency,
			"MinTurnoverDef": "5%",
			"Actived": nil,
			"ReferredBy": user.ReferredBy,
		}
		if err := db.Debug().Model(&user).Where("id = ?",user.ID).Updates(updates).Error; err != nil {
			response := ErrorResponse{
				Message: "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
				Status:  false,
			}

			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := fiber.Map{
			"Message": "เพิ่มยูสเซอร์สำเร็จ!",
			"Status":  true,
			"Data": fiber.Map{
				"id":       user.ID,
				"walletid": user.ID,
				"Username": user.Username,
			},
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

 
func GetUser(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {

	var users models.Users

	db, _err := handler.GetDBFromContext(c)
	prefix := c.Locals("Prefix")
	//fmt.Println("prefix:", prefix)
	if _err != nil {

		// response := fiber.Map{
		// "Message": "โทเคนไม่ถูกต้อง!!",
		// "Status":  false,
		// "Data": fiber.Map{
		// "prefix": prefix,
		// 	},
		// }
		// return c.JSON(response)
		db, _ = database.ConnectToDB(prefix.(string))
	}
	id := c.Locals("Walletid")



	u_err := db.Debug().Where("id= ?", id).Find(&users).Error

	if u_err != nil {

		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}

		return c.JSON(response)
	}
	 
	
	

	var promotionLog models.PromotionLog
	promotionLogs,err := getPromotionsWithUserID(users.ID,db)
	if err != nil {
		return c.JSON(fiber.Map{
			"Status": false,
			"Message":  err,
			"Data": fiber.Map{"id": -1},
		})
	}

	fmt.Printf("PromotionLogs : %+v \n",promotionLogs)

	db.Debug().Model(&models.PromotionLog{}).Where("userid = ? AND promotioncode = ? AND status = 1", users.ID, users.ProID).
		Order("created_at DESC").
		First(&promotionLog)
 

	var count_trans float64
	var totalTurnover decimal.Decimal
 
	fmt.Printf("Users.ProID: %s \n",users.ProID)
	key := fmt.Sprintf("%s:current_promotion", fmt.Sprintf("%s%d",c.Locals("prefix"), users.ID))
	
	pro_uid,_ := getRedisKV(redisClient,key,"uid")
	status,err := getRedisKV(redisClient,key,"status")
    //proID,_ := getRedisKV(redisClient,key,"proID")
	if err != nil {
		fmt.Printf("ไม่พบโปรโมชั่น %s \n",err.Error())
		// if err := redisClient.HSet(ctx, , "status", "2", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
		// 	fmt.Printf(" ไม่สามารถปรับสถานะได้ ")
		// }
		
	} else {
		fmt.Printf(" สถานะ:  %s \n",status)
		fmt.Printf(" key:  %s \n",key)
	// 	if err := redisClient.HSet(ctx, key, "status", "0", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
	// 		fmt.Printf(" ไม่สามารถปรับสถานะได้ ")
	 	}
	if status != "" {
	if err := db.Debug().Model(&models.TransactionSub{}).
		Where("membername = ? AND created_at >= ?", 
			users.Username, 
			promotionLog.CreatedAt).
		Select("COALESCE(SUM(turnover), 0)").
		Scan(&totalTurnover).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Message": "ไม่สามารถคำนวณยอดเทิร์นได้ !",
				"Status": false,
				"Data": "เกิดข้อผิดพลาด!",
			})
			
		}
	} else {
		// var lastdate models.BankStatement
		// if err := db.Debug().Model(&models.BankStatement{}).
		// Where("userid = ? and status='verified'",id).
		// Select("updatedAt").
		// Order("id desc").
		// First(&lastdate).Error; err != nil {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"Message": "ไม่สามารถคำนวณยอดเทิร์นได้ !",
		// 		"Status": false,
		// 		"Data": "เกิดข้อผิดพลาด!",
		// 	})
		// }

		if err := db.Debug().Model(&models.TransactionSub{}).
		Where("membername = ? AND  created_at >= (select MAX(updatedat) from BankStatement Where userid = ? and Status='Verified' )", 
			users.Username, 
			users.ID).
		Select("COALESCE(SUM(turnover), 0)").
		Scan(&totalTurnover).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Message": "ไม่สามารถคำนวณยอดเทิร์นได้ !",
				"Status": false,
				"Data": "เกิดข้อผิดพลาด!",
			})
			
		}

		

			if err := db.Debug().Model(&models.BankStatement{}).
			Where("userid = ? AND  createdAt  >= (select MAX(updatedat) from BankStatement Where userid = ? and Status='Verified' and statement_type='Withdraw' )", 
				users.ID, 
				users.ID).
			Select("COALESCE(count(id), 0)").
			Scan(&count_trans).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Message": "ไม่สามารถคำนวณยอดเทิร์นได้ !",
					"Status": false,
					"Data": "เกิดข้อผิดพลาด!",
				})
				
			}
	}
 
		

	var promotion models.Promotion
	//fmt.Println(summary.Turnover)
	if users.ProStatus != "" {
		db.Debug().Model(&models.Promotion{}).Select("Includegames,Excludegames").Where("ID = ?",users.ProID).Scan(&promotion)
	}
	
 
	balance64,_ := users.Balance.Float64()

	if balance64 < 1  {


		fmt.Printf("User.Status %v \n",users.ProStatus)
		//currentPromotionKey := fmt.Sprintf("%s%s",c.Locals("prefix"),userID)
		// check redis status
		r_status,_ := getRedisKV(redisClient,key,"status")
		fmt.Printf("r.Status %v \n",r_status)
		if r_status == "1" {
			// if err := redisClient.HSet(ctx, key, "status","2").Err(); err != nil {
			// 	return fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
			// }
		
		// updates := map[string]interface{}{
		// 	"ProStatus": "2",
		// 	"Balance": 0,
		// }
 
		// repository.UpdateUserFields(db, users.ID, updates)
		updates := map[string]interface{}{
			"status": "2",
		}
       if err = db.Debug().Model(&models.PromotionLog{}).Where("uid=?",pro_uid).Updates(updates).Error; err != nil {
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": err,
				"Data": fiber.Map{ 
					"id": -1,
				}})
		}
		key := fmt.Sprintf("%s%d", users.Prefix, users.ID)
		currentPromotionKey := fmt.Sprintf("%s:current_promotion", key)
		fmt.Printf("line Key %s \n", currentPromotionKey)

		// ลบคีย์จาก Redis โดยตรง
		err := redisClient.Del(ctx, currentPromotionKey).Err()
		if err != nil {
			fmt.Println("Error clearing Redis keys:", err)
		} else {
			fmt.Println("Successfully cleared Redis keys for user:", users.ID)
		}

		fmt.Printf("line 598 \n")
		}
		return c.JSON(fiber.Map{
			"Status": true,
			"Message":  fmt.Errorf("ข้อมูลยูสเซอร์ เกิดข้อผิดพลาด"),
			"Data": fiber.Map{
				"id":         users.ID,
				"fullname":   users.Fullname,
				"banknumber": users.Banknumber,
				"bankname":   users.Bankname,
				"username":   strings.ToUpper(users.Username),
				"balance":    users.Balance,
				"lastdeposit": users.LastDeposit,
				"lastproamount":users.LastProamount,
				"pro_status": "0",
				"pro_id": "",
				"prefix":     users.Prefix,
				"minturnoverdef": users.MinTurnoverDef,
				"deposit_count": count_trans,
				"turnover": "0",
				"currency":   users.Currency,
				"promotionlog":  promotionLogs,
				"includegames": "13,8,2,5,7,1,3,2",
				"excludegames": "",
			}})


	} else if users.ProStatus != "2" {

	fmt.Printf("line 624")
	
	//fmt.Printf("data: %+v\n", users)
	response := fiber.Map{
		"Status":  true,
		"Message": "สำเร็จ",
		"Data": fiber.Map{
			"id":         users.ID,
			"fullname":   users.Fullname,
			"banknumber": users.Banknumber,
			"bankname":   users.Bankname,
			"username":   strings.ToUpper(users.Username),
			"balance":    users.Balance,
			"currency":   users.Currency,
			"prefix":     users.Prefix,
			"turnover":   totalTurnover.Round(2),//summary.Turnover,
			"minturnoverdef": users.MinTurnoverDef,
			"deposit_count": count_trans,
			"lastdeposit": users.LastDeposit,
			"lastproamount": users.LastProamount,
			"referredby": users.ReferredBy,
			"lastwithdraw": users.LastWithdraw,
			"promotionlog": promotionLogs,
			"pro_status": users.ProStatus,
			"pro_id": pro_uid,
			//"pro_balance": pro_balance,
			"includegames": promotion.Includegames,
			"excludegames": promotion.Excludegames,
		}}
	return c.JSON(response)
	} else {
		fmt.Printf("line 654")
		timestamp,_ := getRedisKV(redisClient,key,"timestamp")

		if err := db.Debug().Model(&models.TransactionSub{}).
		Where(" membername = ? AND Date(created_at) >= Date(?)", 
			users.Username, 
			timestamp).
		Select("COALESCE(SUM(turnover), 0)").
		Scan(&totalTurnover).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Message": "ไม่สามารถคำนวณยอดเทิร์นได้ !",
				"Status": false,
				"Data": "เกิดข้อผิดพลาด!",
			})
			
		
		}


		return c.JSON(fiber.Map{
			"Status": true,
			"Message":   "ดึงข้อมูลสำเร็จ",
			"Data": fiber.Map{
				"id":         users.ID,
				"fullname":   users.Fullname,
				"banknumber": users.Banknumber,
				"bankname":   users.Bankname,
				"username":   strings.ToUpper(users.Username),
				"balance":    users.Balance,
				"pro_status": "0",
				"pro_id": pro_uid,
				"prefix":     users.Prefix,
				"currency":   users.Currency,
				"deposit_count": count_trans,
				"turnover": totalTurnover.Round(2),
				"minturnoverdef": users.MinTurnoverDef,
				"promotionlog":  promotionLogs,
				"includegames": "13,8,2,5,7,1,3,2",
				"excludegames": "",
			}})

	}
}
}


func getRedisKV(redisClient *redis.Client,currentPromotionKey string,key string) (string, error) {
    //currentPromotionKey := fmt.Sprintf("%s:%s", userID,key)

    // ดึงข้อมูล "Example" จาก Redis
    value, err := redisClient.HGet(ctx, currentPromotionKey, key).Result()
    if err != nil && err != redis.Nil {
        return "", fmt.Errorf("error fetching value: %v", err)
    }

    // ถ้าไม่พบค่า ให้คืนค่าเป็นสตริงว่าง
    if err == redis.Nil {
        return "", nil
    }

    return value, nil
}

func getPromotionsWithUserID(userID int, db *gorm.DB) ([]models.Promotion, error) {
	var goodpromotions []models.Promotion
	var promotionLogs []models.PromotionLog
	var createdAt time.Time
	createdAt = time.Now() 
	type Promotion struct {
		SpecificTime json.RawMessage `json:"SpecificTime"`
		// เพิ่มฟิลด์อื่น ๆ ที่คุณต้องการที่นี่
	}
	// ดึง userID จาก c.Locals
	//userID := c.Locals("userid").(int)

	// ดึงข้อมูลจาก promotionlog ตาม userID
	err := db.Where("userid = ?", userID).Find(&promotionLogs).Error
	if err != nil {
		return nil, err
	}

	// ตรวจสอบสถานะของ user
	var user models.Users
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	var promotions []models.Promotion
	if err = db.Debug().Model(models.Promotion{}).Where("status=1 and date(end_date) >= date(?)",createdAt.Format("2006-01-02 15:04:05")).Scan(&promotions).Error; err != nil {
		return nil, err
	}
	
	// ตรวจสอบประเภทของ promotion
	for _, log := range promotions {
		var promotionlog []models.PromotionLog
		err := db.Debug().Model(models.PromotionLog{}).Where("promotioncode = ? and userid = ?",log.ID,userID).Scan(&promotionlog).Error
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
		}
		//if len(promotionlog) == 0 {
		//	continue // ข้ามไปถ้าไม่พบ promotion
		//}
		// if err := json.Unmarshal([]byte(promotionData), &promotion); err != nil {
		// 	fmt.Println("Error unmarshalling JSON:", err)
		// 	return
		// }

		var specificTime map[string]interface{}
		if err := json.Unmarshal([]byte(log.SpecificTime), &specificTime); err != nil {
			fmt.Println("Error unmarshalling SpecificTime:", err)
			
		}

		//fmt.Printf("Promotion: %+v \n",promotion)
		//fmt.Printf("Actived: %v \n",user.Actived)
		//fmt.Printf("Type: %v \n",specificTime["type"])

		// ถ้า user มีการ actived แล้ว ให้ตัด promotion ชนิด first ออก
		// if user.Actived != nil && specificTime["type"] == "first" {
			
		// 	continue // ข้าม promotion ชนิด first
		// }

		// เปรียบเทียบข้อมูลระหว่าง promotionlog และ promotion
		if specificTime["type"] == "first" {
			if user.Actived != nil {
				continue
			}
			// ถ้าเป็น type first ให้ดึงข้อมูล promotion ที่ id ไม่ตรงกับ promotioncode ใน promotionlog
			//if string(promotion.ID) != log.Promotioncode {
				goodpromotions = append(goodpromotions, log)
				
			//}
		} else {

			if specificTime["type"] == "month" {

				var currentMonthPromotionLogs []models.PromotionLog
				currentMonth := time.Now().Month()
				for _, log := range promotionlog {
					if log.CreatedAt.Month() == currentMonth {
						currentMonthPromotionLogs = append(currentMonthPromotionLogs, log)
					}
				}

				if len(currentMonthPromotionLogs) < int(log.UsageLimit) {
					goodpromotions = append(goodpromotions, log)
				}
			
			}  else if specificTime["type"] == "weekly" {
				// กรอง promotionlog ที่มี createdat ตรงกับสัปดาห์ปัจจุบัน
				var currentDayPromotionLogs []models.PromotionLog
				currentWeekday := time.Now().Weekday()
				for _, log := range promotionlog {
					if log.CreatedAt.Weekday() == currentWeekday {
						currentDayPromotionLogs = append(currentDayPromotionLogs, log)
					}
				}

				// ตรวจสอบจำนวนครั้งที่ใช้ในวันปัจจุบัน
				if len(currentDayPromotionLogs) < int(log.UsageLimit) {
					goodpromotions = append(goodpromotions, log)
				}

			} else if specificTime["type"] == "once" {
				// สำหรับ once ให้ตรวจสอบว่า promotionlog มีการใช้งานแล้วหรือไม่
				var usageCount int
				for _, log := range promotionlog {
					if log.CreatedAt.Year() == time.Now().Year() && log.CreatedAt.YearDay() == time.Now().YearDay() {
						usageCount++
					}
				}

				if usageCount < int(log.UsageLimit) {
					goodpromotions = append(goodpromotions, log)
				}
			} else {

			// ถ้าเป็น type อื่น ให้ตรวจสอบจำนวน rows ใน promotionlog
			// var count int64
			// db.Model(&models.PromotionLog{}).Where("promotioncode = ? AND userid = ?", promotion.Promotioncode, userID).Count(&count)
			//fmt.Printf("Count: %v \n",len(promotionlog))
			//fmt.Printf("UsageLimit: %v \n",log.UsageLimit)
			if len(promotionlog)< int(log.UsageLimit) {
				goodpromotions = append(goodpromotions, log)
			}
			}
			
		}
		//fmt.Printf("Promotions: %+v \n",goodpromotions)
		// เปรียบเทียบข้อมูลเพิ่มเติมที่คุณต้องการ
		// เช่น เปรียบเทียบ Promoamount ใน PromotionLog กับ MaxDiscount ใน Promotion
		//if log.Promoamount.GreaterThan(promotion.MaxDiscount) {
			// ทำอะไรบางอย่างถ้า Promoamount มากกว่า MaxDiscount
			// เช่น อาจจะเพิ่มโปรโมชั่นนี้ในรายการ
		//	promotions = append(promotions, promotion)
		//}
	}

	return goodpromotions, nil
}



func  GetPromotionByUser(c *fiber.Ctx) (error) {
	

	 

	// if err := c.BodyParser(user); err != nil {
	// 	return c.Status(200).SendString(err.Error())
	// }
	// //fmt.Printf(" %s ",user.Username)
	// db, _ := database.ConnectToDB(user.Prefix)


	// //db, _err := handler.GetDBFromContext(c)
	// prefix := c.Locals("Prefix")

	db, _err := handler.GetDBFromContext(c)
	 
	if db == nil {
		fmt.Println(_err)
		response := fiber.Map{
			"Status":  false,
			"Message": "โทเคนไม่ถูกต้อง!!",
		}
		return c.JSON(response)
	}
 
	userID := c.Locals("ID")
	//var promotionlog = []models.PromotionLog{}
	 	// ตรวจสอบสถานะของ user
	var user models.Users
	if err := db.Debug().First(&user, userID).Error; err != nil {
		response := fiber.Map{
			"Message": err,
			"Status":  false,
			"Data": fiber.Map{
				"id": -1,
			},
		}

		return c.JSON(response)
	}

	updates := map[string]interface{}{
		//"Balance": user.Balance.Add(transactionAmount),
		//"ProID":user.ProStatus,
		}
	one,_ := decimal.NewFromString("1")
	if user.Balance.LessThan(one) {//&& user.ProID == "1" {
		if user.ProID == "1"{
		updates["pro_status"] = ""
		} else {
			updates["pro_status"] = user.ProStatus
		}
	} else {
		updates["pro_status"] = user.ProStatus
	}
	repository.UpdateUserFields(db,userID.(int), updates) 

	var promotion models.Promotion

	err := db.Debug().Model(&models.Promotion{}).Where("id = ?",updates["pro_status"].(string)).First(&promotion).Error

	//var promotionLogs []models.PromotionLog

	// ดึง userID จาก c.Locals
	//userID := c.Locals("userid").(int)

	// ดึงข้อมูลจาก promotionlog ตาม userID
	//err := db.Where("userid = ?", userID).Find(&promotionLogs).Error
	if err != nil {
		response := fiber.Map{
			"Message": "ไม่พบข้อมูล",
			"Status":  false,
			"Data": fiber.Map{ 
			},
		}

		return c.JSON(response)
	}




	// // ตรวจสอบประเภทของ promotion
	// for _, log := range promotionLogs {
	// 	var promotion models.Promotion
	// 	err := db.First(&promotion, log.Promotioncode).Error
	// 	if err != nil {
	// 		continue // ข้ามไปถ้าไม่พบ promotion
	// 	}
	// 	// ถ้า user มีการ actived แล้ว ให้ตัด promotion ชนิด first ออก
	// 	if user.Actived != nil && promotion.Typetime == "first" {
	// 		continue // ข้าม promotion ชนิด first
	// 	}

	// 	// เปรียบเทียบข้อมูลระหว่าง promotionlog และ promotion
	// 	if promotion.Typetime == "first" {
	// 		// ถ้าเป็น type first ให้ดึงข้อมูล promotion ที่ id ไม่ตรงกับ promotioncode ใน promotionlog
	// 		if string(promotion.ID) != log.Promotioncode {
	// 			promotions = append(promotions, promotion)
	// 		}
	// 	} else {
	// 		// ถ้าเป็น type อื่น ให้ตรวจสอบจำนวน rows ใน promotionlog
	// 		var count int64
	// 		db.Model(&models.PromotionLog{}).Where("promotioncode = ? AND userid = ?", promotion.Promotioncode, userID).Count(&count)

	// 		if count < int64(promotion.UsageLimit) {
	// 			promotions = append(promotions, promotion)
	// 		}
	// 	}

	// 	// เปรียบเทียบข้อมูลเพิ่มเติมที่คุณต้องการ
	// 	// เช่น เปรียบเทียบ Promoamount ใน PromotionLog กับ MaxDiscount ใน Promotion
	// 	if log.Promoamount.GreaterThan(promotion.MaxDiscount) {
	// 		// ทำอะไรบางอย่างถ้า Promoamount มากกว่า MaxDiscount
	// 		// เช่น อาจจะเพิ่มโปรโมชั่นนี้ในรายการ
	// 		promotions = append(promotions, promotion)
	// 	}
	// }
	response := fiber.Map{
		"Status":  true,
		"Message": "สำเร็จ",
		"Data": promotion,
	}
	return c.JSON(response)
}

func GetUserByUsername(c *fiber.Ctx) error {
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//fmt.Printf(" %s ",user.Username)
	db, conn := database.ConnectToDB(user.Prefix)


	//db, _err := handler.GetDBFromContext(c)
	prefix := c.Locals("Prefix")
	//fmt.Println("prefix:", prefix)
	if conn != nil {

		// response := fiber.Map{
		// "Message": "โทเคนไม่ถูกต้อง!!",
		// "Status":  false,
		// "Data": fiber.Map{
		// "prefix": prefix,
		// 	},
		// }
		// return c.JSON(response)
		db, _ = database.ConnectToDB(prefix.(string))
	}
	//id := c.Locals("Walletid")
	var users models.Users
	u_err := db.Debug().Where("username= ?", user.Username).Find(&users).Error

	if u_err != nil {

		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}

		return c.JSON(response)
	}

	// type Summary struct {
	// 	Turnover decimal.Decimal `json:"turnover"`
	// 	createdAt time.Time `json:"createdat"`
	// }
	var summary models.BankStatement
	
	//db.Debug().Model(&models.BankStatement{}).Select("turnover,createdat").Where("userid= ?", users.ID).Last(&summary)
	db.Debug().Model(&models.BankStatement{}).Select("turnover, createdAt").Where("userid= ? and transactionamount<0", users.ID).Order("createdAt DESC").Limit(1).Scan(&summary)
	
		 
	 //fmt.Println(summary.CreatedAt.Format("2006-01-02"))

	//createdate := summary.createdAt.Format("2006-01-02")
	fmt.Println(summary.Turnover)
	fmt.Println(summary.CreatedAt)
	if summary.Turnover.LessThanOrEqual(decimal.Zero) {
	 	db.Debug().Model(&models.TransactionSub{}).Select("COALESCE(sum(BetAmount),0) as turnover").Where("membername= ? and deleted_at is null", users.Username).Scan(&summary)
	} else {
		db.Debug().Model(&models.TransactionSub{}).Select("COALESCE(sum(BetAmount),0) as turnover").Where("membername= ? and created_at > ? and deleted_at is null", users.Username,summary.CreatedAt.Format("2006-01-02 15:04:05")).Scan(&summary)
	}
	
	var promotion models.Promotion
	//fmt.Println(summary.Turnover)
	if users.ProStatus != "" {
		db.Debug().Model(&models.Promotion{}).Select("Includegames,Excludegames").Where("ID = ?",users.ProStatus).Scan(&promotion)
	}




	response := fiber.Map{
		"Status":  true,
		"Message": "สำเร็จ",
		"Data": fiber.Map{
			"id":         users.ID,
			"fullname":   users.Fullname,
			"banknumber": users.Banknumber,
			"bankname":   users.Bankname,
			"username":   strings.ToUpper(users.Username),
			"balance":    users.Balance,
			"prefix":     users.Prefix,
			"turnover":   summary.Turnover,
			"minturnover": users.MinTurnover,
			"lastdeposit": users.LastDeposit,
			"lastproamount": users.LastProamount,
			"minturnoverdef": users.MinTurnoverDef,
			"lastwithdraw": users.LastWithdraw,
			"pro_status": users.ProStatus,
			"includegames": promotion.Includegames,
			"excludegames": promotion.Excludegames,
		}}
	return c.JSON(response)
}
 
func GetBalance(c *fiber.Ctx) error {

	var users models.Users

	db, _err := handler.GetDBFromContext(c)
	prefix := c.Locals("Prefix")
	if _err != nil {

		response := fiber.Map{
			"Message": "โทเคนไม่ถูกต้อง!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}
		return c.JSON(response)
	}
	id := c.Locals("Walletid")
	u_err := db.Debug().Where("id= ?", id).Find(&users).Error

	if u_err != nil {

		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}

		return c.JSON(response)
	}

	response := fiber.Map{
		"Status":  true,
		"Message": "สำเร็จ",
		"Data": fiber.Map{
			"id":         users.ID,
			"fullname":   users.Fullname,
			"banknumber": users.Banknumber,
			"bankname":   users.Bankname,
			"username":   strings.ToUpper(users.Username),
			"balance":    users.Balance,
			"prefix":     users.Prefix,
			"minturnoverdef": users.MinTurnoverDef,
		}}
	return c.JSON(response)
}

 
func Logout(c *fiber.Ctx) error {

	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	fmt.Println("Claims : ", claims)

	username := claims["Username"].(string)

	prefix, _ := jwt.GetPrefix(username)
	db, _ := jwt.CheckDBConnection(c.Locals("db"), prefix)

	updates := map[string]interface{}{
		"Token": "",
	}

	// อัปเดตข้อมูลยูสเซอร์
	repository.UpdateFieldsUserString(db, username, updates)

	response := fiber.Map{
		"Message": "ออกจากระบบสำเร็จ!",
		"Status":  true,
		"Data": fiber.Map{
			"id":     -1,
			"prefix": prefix,
		},
	}
	return c.JSON(response)

}

 
func GetUserTransaction(c *fiber.Ctx) error {

	body := new(Body)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	response := fiber.Map{
		"Status":  false,
		"Message": "สำเร็จ",
		"Data":    map[string]interface{}{},
	}

	db, _err := handler.GetDBFromContext(c)
	prefix := c.Locals("Prefix")
	if _err != nil {

		response := fiber.Map{
			"Message": "โทเคนไม่ถูกต้อง!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}
		return c.JSON(response)
	}
	id := c.Locals("Walletid")
	provide := body.Provide
	startDateStr := body.Startdate
	endDateStr := body.Stopdate

	var statements []models.TransactionSub

	if body.Status == "all" {
		db.Debug().Where("id=? AND GameProvide = ? AND  DATE(createdat) BETWEEN ? AND ? ", id, provide, startDateStr, endDateStr).Find(&statements)
	} else {
		db.Debug().Where("id=? AND GameProvide = ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?", id, provide, startDateStr, endDateStr, body.Status).Find(&statements)
	}

	// สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
	result := make([]fiber.Map, len(statements))

	// วนลูปเพื่อประมวลผลแต่ละรายการ
	for i, transaction := range statements {
		// ตรวจสอบเงื่อนไขด้วย inline if-else
		transactionType := func(amount decimal.Decimal) string {
			if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
				return "เสีย"
			}
			return "ได้"
		}(transaction.TransactionAmount)
		amountFloat, _ := transaction.TransactionAmount.Float64()
		// เก็บผลลัพธ์ใน slice
		result[i] = fiber.Map{
			"userid":          transaction.MemberID,
			"createdAt":       transaction.CreatedAt,
			"transfer_amount": amountFloat,
			"credit":          amountFloat,
			"Status":          transaction.Status,
			//	"channel": transaction.Channel,
			"statement_type": transactionType,
			"expire_date":    transaction.CreatedAt,
		}
	}

	//    response = fiber.Map{
	// 	"Status": true,
	// 	"Message": "ไม่สำเร็จ",

	// 	}

	return c.JSON(response)

}

 
func GetUserStatement(c *fiber.Ctx) error {

	body := new(Body)
	//prefix := ""
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	response := fiber.Map{
		"Status":  false,
		"Message": "สำเร็จ",
		"Data":    map[string]interface{}{},
	}

	//user := c.Locals("username")//.(*jtoken.Token)
	//claims := user.Claims.(jtoken.MapClaims)

	//fmt.Println(&claims)

	// if claims["Prefix"] != nil {
	// 	prefix = claims["Prefix"].(string)
	// } else {
	// 	prefix,_ = jwt.GetPrefix(claims["Username"].(string))
	// }

	// if prefix == "" {
	// 	prefix,_ = jwt.GetPrefix(claims["Username"].(string))
	// }

	//db, _ := jwt.CheckDBConnection(c.Locals("db"),prefix)
	//_err := jwt.CheckedJWT(db,c);
	dbInterface := c.Locals("db")
	if dbInterface == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No database connection found",
		})
	}

	// แปลง interface{} ให้เป็น *gorm.DB
	db, ok := dbInterface.(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid database connection",
		})
	}

	// prefix := c.Locals("prefix")
	// if _err != nil {
	// 	response := fiber.Map{
	// 	"Message": "โทเคนไม่ถูกต้อง!!",
	// 	"Status":  false,
	// 	"Data": fiber.Map{
	// 	"prefix": prefix,
	// 		},
	// 	}
	// 	return c.JSON(response)
	// }

	id := c.Locals("Walletid") //claims["walletid"]
	channel := body.Channel
	startDateStr := body.Startdate
	endDateStr := body.Stopdate

	var statements []models.BankStatement

	if body.Status == "all" {
		db.Debug().Where("userid=? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? ", id, channel, startDateStr, endDateStr).Find(&statements)
	} else {
		db.Debug().Where("userid=? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?", id, channel, startDateStr, endDateStr, body.Status).Find(&statements)
	}

	// สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
	result := make([]fiber.Map, len(statements))

	// วนลูปเพื่อประมวลผลแต่ละรายการ
	for i, transaction := range statements {
		// ตรวจสอบเงื่อนไขด้วย inline if-else
		transactionType := func(amount decimal.Decimal, channel string) string {
			if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
				return "ถอน"
			}
			return "ฝาก"
		}(transaction.Transactionamount, transaction.Channel)
		amountFloat, _ := transaction.Transactionamount.Float64()
		// เก็บผลลัพธ์ใน slice
		result[i] = fiber.Map{
			"userid":          transaction.Userid,
			"createdAt":       transaction.CreatedAt,
			"transfer_amount": amountFloat,
			"credit":          amountFloat,
			"Status":          transaction.Status,
			"channel":         transaction.Channel,
			"statement_type":  transactionType,
			"expire_date":     transaction.CreatedAt,
		}
	}

	return c.JSON(response)

}

func GetBalanceSum(c *fiber.Ctx) error {

	//var users models.Users

	db, _err := handler.GetDBFromContext(c)
	prefix := c.Locals("Prefix")
	if _err != nil {

		response := fiber.Map{
			"Message": "โทเคนไม่ถูกต้อง!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}
		return c.JSON(response)
	}
	//id := c.Locals("Walletid")
	var sum decimal.Decimal
	u_err := db.Debug().Table("Users").Select("sum(balance)").Where("deposit is not NULL").Row().Scan(&sum)
	//u_err := db.Debug().Table("Users").Select("sum(balance)").Where("deposit is not NULL").Find(&users).Error

	if u_err != nil {

		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
			"Status":  false,
			"Data": fiber.Map{
				"prefix": prefix,
			},
		}

		return c.JSON(response)
	}

	response := fiber.Map{
		"Status":  true,
		"Message": "สำเร็จ",
		"Data": fiber.Map{
			"balance": sum,
		}}
	return c.JSON(response)
}
type UpdateBody struct {
	ID   string         `json:"id"`
	Body map[string]interface{} `json:"body"`
}
// ... existing code ...

func UpdateUser(c *fiber.Ctx) error {
	// Parse the request body into a map
	body := make(map[string]interface{})
	if err := c.BodyParser(&body); err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": err.Error(),
		}
		return c.JSON(response)
	}

	// Get the username from the context
	username := c.Locals("username").(string)
	
	db, _err := handler.GetDBFromContext(c)
	if _err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": "โทเคนไม่ถูกต้อง!!",
		}
		return c.JSON(response)
	}

	var user models.Users
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
		}
		return c.JSON(response)
	}
 

	// Update the user with the provided fields
	fmt.Printf("Body: %s",body)
	if err := db.Debug().Model(&user).Updates(body).Error; err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": "ไม่สามารถอัปเดตข้อมูลได้: " + err.Error(),
		}
		return c.JSON(response)
	}

	response := fiber.Map{
		"Status":  true,
		"Message": "อัปเดตข้อมูลสำเร็จ!",
	}
	return c.JSON(response)
}


func checkProlog(pro_id string,userid string) {
	
}

func UpdateUserPro(c *fiber.Ctx) error {

	type PBody struct {
		Prefix    string `json:"prefix"`
		Prostatus string `json:"pro_status"`
	}

	rdb := handler.CreateRedisClient()
	defer rdb.Close()

	//	body :=   make(map[string]interface{})
	body := PBody{}

	if err := c.BodyParser(&body); err != nil {
		fmt.Println("Error parsing body:", err.Error())
		response := fiber.Map{
			"Status":  false,
			"Message": err.Error(),
		}
		return c.JSON(response)
	}
	
	// Get the username from the context
	username := c.Locals("username").(string)
	
	fmt.Printf("Body: %+v \n",body)
	fmt.Printf("username: %v \n",username)

	db, _err := handler.GetDBFromContext(c)
	if _err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": "โทเคนไม่ถูกต้อง!!",
		}
		return c.JSON(response)
	}

	var user models.Users
	err := db.Debug().Where("username = ?", username).First(&user).Error
	if err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
		}
		return c.JSON(response)
	}

	if user.Balance.Floor().GreaterThan(decimal.NewFromFloat(0)) {
		return c.Status(400).JSON(fiber.Map{
			"Message": "คุณมียอดคงเหลือมากกว่าศูนย์",
			"Status":  false,
			"Data":    "คุณมียอดคงเหลือมากกว่าศูนย์",
		})
	}
	currentTime := time.Now()
	// uid := ""
	
	// if err := db.Debug().Model(&models.BankStatement{}).
	// 	Select("uid").
	// 	Where("userid = ? and status='verified'",user.ID).
	// 	Order("id desc").
	// 	First(&uid).Error; err != nil {
	// 		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		// 	"Status": false,
	// 		// 	"Message": "ไม่พบรายการฝากที่ตรงกับโปรโมชั่น",
	// 		// 	"Data": fiber.Map{"id": -1},
	// 		// })
	// 		fmt.Printf("Error %s \n",err.Error())
	// 	}
	
	/// check promax
	var usageCount int64
	if err := db.Debug().Model(&models.PromotionLog{}).
		Where("userid = ? AND promotioncode = ? AND status = 2 and date(created_at)=date(?)", user.ID, body.Prostatus,currentTime.Format("2006-01-02")).
		Count(&usageCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status": false,
			"Message": "ไม่สามารถตรวจสอบประวัติการใช้โปรโมชั่นได้",
			"Data": fiber.Map{"id": -1},
		})
	}

	// response["Id"] = ProItem.ID
	// 		response["Type"] = ProItem.ProType.Type
	// 		response["count"] = ProItem.UsageLimit
	// 		response["MinTurnover"] = promotion.MinSpend
	// 		response["Formular"] = promotion.Example
	// 	    response["Name"] = promotion.Name
	// 		response["TurnType"]=promotion.TurnType
	// 	if ProItem.ProType.Type == "weekly" {
	// 		response["Week"] = ProItem.ProType.DaysOfWeek
	// 	}


	prodetail,_ := handler.GetProdetail(db,body.Prostatus)
	// ตรวจสอบกับค่า MaxUse จาก pro_setting
	maxUse, ok := prodetail["MaxUse"].(int64)

	fmt.Println("วันที่:", currentTime.Format("2006-01-02"))
	fmt.Printf("ไอดี: %v \n",prodetail["Id"])
	fmt.Printf("ชื่อ: %v \n",prodetail["Name"])
	fmt.Printf("ประเภท: %v \n",prodetail["Type"])
	fmt.Printf("ประเภท Turn: %v \n",prodetail["TurnType"])
	fmt.Printf("ใช้งานได้สูงสุด: %v \n",maxUse)
	fmt.Printf("ใช้แล้วกี่ครั้ง: %v \n",usageCount)

	if ok && usageCount >= maxUse {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Status": false,
			"Message": "คุณใช้โปรโมชั่นนี้เกินจำนวนครั้งที่กำหนดแล้ว \n โปรโมชั่นจะปรับอัตโนมัติ",
			"Data": fiber.Map{"id": -1},
		})
	}


	if err := selectPromotion(rdb, string(user.ID), body.Prostatus); err != nil {
			return c.JSON(fiber.Map{
			"Status": false,
			"Message": fmt.Sprintf("ไม่สามารถตรวจสอบโปรโมชั่นได้ %s ",err.Error()),
			"Data": fiber.Map{"id": -1},
		})
	}


	// proStatusValue := body.Prostatus
	// userPro := user.ProStatus

	// fmt.Printf(" userPrp: %v \n",userPro)
	// fmt.Printf(" user.Balance.IsZero(): %v \n",user.Balance.IsZero())
	// fmt.Printf(" user.Balance.LessThan: %v \n",user.Balance.LessThan(decimal.NewFromFloat(1)))


	// if user.Balance.IsZero() || user.Balance.LessThan(decimal.NewFromFloat(1)) {
	// 		userPro = ""
	// }
	
	// pro_setting, err := handler.GetProdetail(db,userPro)
	// fmt.Printf("ProSetting: %+v \n",pro_setting)

	// if err != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message": "ไม่สามารถตรวจสอบโปรโมชั่นได้",
	// 		"Data": fiber.Map{"id": -1},
	// 	})
	// }

	// if pro_setting == nil  {


	// 	proStatus := fmt.Sprintf("%v", proStatusValue)
	// 	fmt.Printf("proStatus: %v \n",proStatus)
	//    //check maxuse
	//    var usageCount int64
	//    if err := db.Debug().Model(&models.PromotionLog{}).
	// 	   Where("userid = ? AND promotioncode = ? AND status = 1", user.ID, proStatus).
	// 	   Count(&usageCount).Error; err != nil {
	// 	   return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		   "Status": false,
	// 		   "Message": "ไม่สามารถตรวจสอบประวัติการใช้โปรโมชั่นได้",
	// 		   "Data": fiber.Map{"id": -1},
	// 	   })
	//    }

	//    prodetail,_ := handler.GetProdetail(db,proStatus)
	// 	// ตรวจสอบกับค่า MaxUse จาก pro_setting
	// 	maxUse, ok := prodetail["MaxUse"].(int64)

	// 	fmt.Printf("maxUse: %v \n",maxUse)
	// 	fmt.Printf("usageCount: %v \n",usageCount)

	// 	if ok && usageCount >= maxUse {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"Status": false,
	// 			"Message": "คุณใช้โปรโมชั่นนี้เกินจำนวนครั้งที่กำหนดแล้ว \n โปรโมชั่นจะปรับอัตโนมัติ",
	// 			"Data": fiber.Map{"id": -1},
	// 		})
	// 	}

	// 	fmt.Printf("Balance: %v \n",user.Balance)
	// 	fmt.Printf("Round: %v \n",user.Balance.Floor())

	// 	one,_ := decimal.NewFromString("0")

	// 	if user.Balance.Floor().GreaterThan(one) {
	// 		return c.Status(400).JSON(fiber.Map{
	// 			"Message": "คุณมียอดคงเหลือมากกว่าศูนย์",
	// 			"Status":  false,
	// 			"Data":    "คุณมียอดคงเหลือมากกว่าศูนย์",
	// 		})
	// 	}

	   updates := map[string]interface{}{
		   "MinTurnover": prodetail["MinTurnover"],
		   "ProStatus": body.Prostatus,
		   "ProID": "0",
	   }
	   repository.UpdateFieldsUserString(db, user.Username, updates)
	   return c.Status(200).JSON(fiber.Map{
		   "Status":  true,
		   "Message": "อัปเดตข้อมูลสำเร็จ!",
	   })

	   
		
	// } else {

		

	// 	return c.Status(400).JSON(fiber.Map{
	// 		"Message": "คุณใช้งานโปรโมชั่นอยู่",
	// 		"Status":  false,
	// 		"Data":    "คุณใช้งานโปรโมชั่นอยู่",
	// 	})
		 
	// }

	 
}
 

func CheckSeed(db *gorm.DB) string {

	var seedPhrase string
	for {
		seedPhrase, _ = encrypt.GenerateAffiliateCode(5) // สร้าง affiliate key ใหม่
		rowsAffected := db.Debug().Model(&models.Users{}).Where("referral_code = ?", seedPhrase).RowsAffected
		if rowsAffected == 0 { // ถ้าไม่ซ้ำ
			break // ออกจากลูป
		}
	}
	return seedPhrase
}

func selectPromotion(redisClient *redis.Client, userID string, proID string) error {
	key := fmt.Sprintf("%s:%s", userID, proID)

	// ตรวจสอบสถานะโปรโมชั่นก่อนหน้านี้
	oldPromotion, err := redisClient.Keys(ctx, fmt.Sprintf("%s:*", userID)).Result()
	if err != nil {
		return err
	}

	fmt.Printf("oldPromotion: %s\n", oldPromotion)

	var hasActivePromotion bool
	for _, oldKey := range oldPromotion {
		fmt.Printf("oldKey: %s\n", oldKey)
		oldStatus, err := redisClient.Get(ctx, oldKey).Result()
		fmt.Printf("oldStatus: %s\n", oldStatus)
		if err != nil {
			return err
		}
		if oldStatus == "1" { // เช็คว่าสถานะไม่เป็น 2
			hasActivePromotion = true
			break
		}
	}

	// หากพบโปรโมชั่นที่ยังไม่มีสถานะเป็น 2 ให้ไม่อนุญาตให้เปลี่ยน
	if hasActivePromotion {
		return fmt.Errorf("cannot select new promotion while the previous promotion is still active")
	}
	
	// กำหนดสถานะโปรโมชั่นใหม่
	// if err := redisClient.Set(ctx, key, "0", 0).Err(); err != nil { // ตั้งสถานะใหม่เป็น 0
	// 	return err
	// }

	// ใช้ EXPIRE เพื่อกำหนด TTL 24 ชั่วโมง (86400 วินาที)
	if err := redisClient.SetEX(ctx, key, "0", 24*time.Hour).Err(); err != nil {
	return err
}

	fmt.Printf("New promotion selected with ID: %s\n", proID)
	return nil
}

// ฟังก์ชันอัปเดตสถานะเป็น 1 (ฝากเงินสำเร็จ)
func depositSuccess(redisClient *redis.Client, userID string, proID string) error {
	key := fmt.Sprintf("%s:%s", userID, proID)
	if err := redisClient.Set(ctx, key, "1", 0).Err(); err != nil {
		return err
	}
	fmt.Printf("Deposit successful. Updated promotion %s with status 1\n", proID)
	return nil
}

// ฟังก์ชันอัปเดตสถานะเป็น 2 (ถอนเงินหรือยอดเงินคงเหลือเป็นศูนย์)
func withdrawOrZeroBalance(redisClient *redis.Client, userID string, proID string) error {
	key := fmt.Sprintf("%s:%s", userID, proID)
	if err := redisClient.Set(ctx, key, "2", 0).Err(); err != nil {
		return err
	}
	fmt.Printf("Withdraw or zero balance. Updated promotion %s with status 2\n", proID)
	return nil
}
