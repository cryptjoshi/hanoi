package users

import (
	// 	// "context"
	// 	// "fmt"
	// 	// "github.com/amalfra/etag"
	// 	// "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	// 	// "github.com/streadway/amqp"
	// 	// "github.com/tdewolff/minify/v2"
	// 	// "github.com/tdewolff/minify/v2/js"
	// 	// "github.com/valyala/fasthttp"
	// 	// _ "github.com/go-sql-driver/mysql"
	jtoken "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"hanoi/database"
	"hanoi/handler"
	"hanoi/handler/jwtn"
	"hanoi/models"
	wallet "hanoi/handler/wallet"
	// 	//"github.com/golang-jwt/jwt"
	// 	//jtoken "github.com/golang-jwt/jwt/v4"
	// 	//"github.com/solrac97gr/basic-jwt-auth/config"
	// 	//"github.com/solrac97gr/basic-jwt-auth/models"
	// 	//"github.com/solrac97gr/basic-jwt-auth/repository"
	"github.com/go-redis/redis/v8" 
	"hanoi/repository"
	"hanoi/encrypt"
	"context" 
	//"log"
	// 	// "net"
	// 	// "net/http"
	"os"
	// 	// "strconv"
	"time"
	"fmt"
	"strings"
	"encoding/json"
	//"errors"
)
var redis_master_host = "redis" //os.Getenv("REDIS_HOST")
var redis_master_port = "6379"  //os.Getenv("REDIS_PORT")
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

// @Summary Login user
// @Description Get a list of all users in the database.
// @Tags users
// @Produce json
// @Success 200 {array} models.SwaggerUser
// @Router /users/login [post]
// @Param user body models.Users true "User registration info"
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

	// fmt.Printf("%s",loginRequest)

	db, err := database.ConnectToDB(loginRequest.Prefix)
	//db.AutoMigrate(&models.BankStatement{},&models.PromotionLog{})
	db.Migrator().CreateTable(&models.Promotion{},&models.PromotionLog{})
	//db.AutoMigrate(&models.Promotion{})
	//database.MigrationPromotion(db)
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

	//day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID":          user.ID,
		"Walletid":    user.Walletid,
		"Username":    user.Username,
		"Fullname":    user.Fullname,
		"Role":        user.Role,
		"PartnersKey": user.PartnersKey,
		"Prefix":      user.Prefix,
		//"exp":   time.Now().Add(day * 1).Unix(),
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", //redis_master_password,
		DB:       0,  // ใช้ database 0
	})

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
	// return c.JSON(models.LoginResponse{
	// 	Token: t,
	// })

}

// @Summary Get all users
// @Description Get a list of all users in the database.
// @Tags users
// @Produce json
// @Success 200 {array} models.SwaggerUser
// @Router /users/all [post]
// @Param user body models.SwaggerBody true "User registration info"
// @param Authorization header string true "Bearer token"
func GetUsers(c *fiber.Ctx) error {

	user := new(models.SwaggerBody)

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

// @Summary Register new user
// @Description Register user in the database.
// @Tags users
// @Produce json
// @Accept json
// @Success 200 {object} models.SwaggerUser
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /users/register [post]
// @Param user body models.Users true "User registration info"
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

// @Summary Get userinfo
// @Description Get userinfo in the database.
// @Tags users
// @Produce json
// @Accept json
// @Security BearerAuth
// @Success 200 {object} models.SwaggerUser
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /users/info [post]
// @Param user body models.Users true "User user info"
// @param Authorization header string true "Bearer token"
func GetUser(c *fiber.Ctx) error {

	var users models.Users

	db, _err := handler.GetDBFromContext(c)
	prefix := c.Locals("Prefix")
	fmt.Println("prefix:", prefix)
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
	// if err := db.AutoMigrate(&models.BankStatement{});err != nil {
	// 	fmt.Errorf("Tables schema migration not successfully\n")
	// }
	// type Summary struct {
	// 	Turnover decimal.Decimal `json:"turnover"`
	// 	createdAt time.Time `json:"createdat"`
	// }
	//var summary models.BankStatement
	
	//db.Debug().Model(&models.BankStatement{}).Select("turnover,createdat").Where("userid= ?", users.ID).Last(&summary)
	//db.Debug().Model(&models.BankStatement{}).Select("turnover, createdAt").Where("userid= ? and transactionamount<0", users.ID).Order("createdAt DESC").Limit(1).Scan(&summary)
	
		 
	 //fmt.Println(summary.CreatedAt.Format("2006-01-02"))

	//createdate := summary.createdAt.Format("2006-01-02")
	// fmt.Printf("turnover: %v \n",summary.Turnover)
	// fmt.Printf("createdAt: %v \n",summary.CreatedAt)
	// fmt.Printf("342 line GetUser \n")

	// if summary.Turnover.LessThanOrEqual(decimal.Zero) {
	//  	db.Debug().Model(&models.TransactionSub{}).Select("COALESCE(sum(BetAmount),0) as turnover").Where("membername= ? and deleted_at is null", users.Username).Scan(&summary)
	// } else {
	// 	db.Debug().Model(&models.TransactionSub{}).Select("COALESCE(sum(BetAmount),0) as turnover").Where("membername= ? and created_at > ? and deleted_at is null", users.Username,summary.CreatedAt.Format("2006-01-02 15:04:05")).Scan(&summary)
	// }
	
	

	var promotionLog models.PromotionLog
	promotionLogs,err := getPromotionsWithUserID(users.ID,db)
	if err != nil {
		return c.JSON(fiber.Map{
			"Status": false,
			"Message":  err,
			"Data": fiber.Map{"id": -1},
		})
	}

	//fmt.Printf("PromotionLogs : %+v \n",promotionLogs)

	db.Debug().Model(&models.PromotionLog{}).Where("userid = ? AND promotioncode = ? AND status = 1", users.ID, users.ProStatus).
		Order("created_at DESC").
		First(&promotionLog)
 


	var totalTurnover decimal.Decimal
 
	if err := db.Debug().Model(&models.TransactionSub{}).
		Where("proid = ? AND membername = ? AND Date(created_at) >= Date(?)", 
			users.ProStatus, 
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
 

	// type TurnoverResult struct {
	// 	Turnover decimal.Decimal
	// }
	
	//  //var lastWithdrawTurnover decimal.Decimal
	//  subQuery := db.Debug().
	// 	 Table("BankStatement").
	// 	 Select("TurnOver").
	// 	 Where("(userid = ? OR walletid = ?) AND statement_type = ?", users.ID, users.ID, "Withdraw").
	// 	 Order("created_at DESC").
	// 	 Limit(1)
	 
	//  // Query หลัก
	//  var result TurnoverResult
	//  err = db.Debug().
	// 	 Table("TransactionSub").
	// 	 Select("SUM(turnover) - COALESCE((?),0) as turnover", subQuery).
	// 	 Where("MemberID = ?", users.ID).
	// 	 Scan(&result).Error

	// if err != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message":  err,
	// 		"Data": fiber.Map{"id": -1},
	// 	})
	// }
	 
	//  fmt.Printf(" Users TurnOver: %v \n",result.Turnover)
	//  fmt.Printf(" minTurnover: %v \n",minTurnover)



	var promotion models.Promotion
	//fmt.Println(summary.Turnover)
	if users.ProStatus != "" {
		db.Debug().Model(&models.Promotion{}).Select("Includegames,Excludegames").Where("ID = ?",users.ProStatus).Scan(&promotion)
	}
	pro_setting, err := wallet.CheckPro(db, &users) 
	if err != nil {
		//fmt.Printf("388 error: %s \n",err)
		
	updates := map[string]interface{}{
		"pro_status": "",
	}

	// อัปเดตข้อมูลยูสเซอร์
	repository.UpdateFieldsUserString(db, users.Username, updates)


		return c.JSON(fiber.Map{
			"Status": false,
			"Message":  err.Error(),
			"Data": fiber.Map{
				"id":         users.ID,
				"fullname":   users.Fullname,
				"banknumber": users.Banknumber,
				"bankname":   users.Bankname,
				"username":   strings.ToUpper(users.Username),
				"balance":    users.Balance,
				"prefix":     users.Prefix,
				"promotionlog":  promotionLogs,
				"includegames": "13,8,2,5,7,1,3,2",
				"excludegames": "",
			}})
	}

	var minTurnover string
	if pro_setting["MinTurnover"] != nil {
		minTurnover = fmt.Sprintf("%v", pro_setting["MinTurnover"])
	} else {
		minTurnover = "0" // ค่าเริ่มต้นเมื่อเป็น nil
	}

	var baseAmount decimal.Decimal
	if pro_setting["MinSpendType"] == "deposit" {
		baseAmount = users.LastDeposit
	} else {
		baseAmount = users.LastDeposit.Add(users.LastProamount)
	}

	if minTurnover == "" {
		minTurnover = "0"
	}
	//fmt.Printf("508 minTurnover: %+v \n",minTurnover)
	//fmt.Printf("509 baseAmount: %+v \n",baseAmount)


	requiredTurnover, err := wallet.CalculateRequiredTurnover(minTurnover, baseAmount)
 
	if err != nil {
		return c.JSON(fiber.Map{
			"Status": false,
			"Message": "ไม่สามารถคำนวณยอดเทิร์นได้",
			"Data": fiber.Map{"id": -1},
		})
	}

	// var transaction models.TransactionSub
	// db.Debug().Model(&models.TransactionSub{}).Select("COALESCE(balance,0) as balance").Where("membername= ? and deleted_at is null", users.Username).Scan(&transaction)
	// var pro_balance decimal.Decimal
	// db.Debug().Model(&models.TransactionSub{}).
	// 	Select("COALESCE(balance, 0) as balance").
	// 	Where("membername = ? AND deleted_at is null AND ProID=? and created_at > ?", users.Username,users.ProStatus,time.Now().Format("2006-01-02 15:04:05")).Order("id DESC").Limit(1).Scan(&pro_balance)
	// createdAt := time.Now()
	// if pro_setting["CreatedAt"] != nil {
	// 	if t, ok := pro_setting["CreatedAt"].(time.Time); ok {
	// 		createdAt = t
	// 	}
	// }
	// pro_setting["CreatedAt"] = createdAt.Format("2006-01-02 15:04:05")
	
	var pro_balance decimal.Decimal
	var createdAt time.Time
	createdAt = time.Now() 
	if pro_setting["CreatedAt"] != nil {
		createdAt = pro_setting["CreatedAt"].(time.Time) 
	}
	db.Debug().Model(&models.TransactionSub{}).Select("balance").Where("membername = ? AND deleted_at is null and created_at > ?",users.Username,createdAt.Format("2006-01-02 15:04:05")).Limit(1).Order("id desc").Find(&pro_balance)
	

	updates := map[string]interface{}{
		"ProBalance": pro_balance,
	}
	if err := db.Debug().Model(&users).Where("id=? and pro_status=?",users.ID,users.ProStatus).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
		}

	
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
			"prefix":     users.Prefix,
			"turnover":   totalTurnover,//summary.Turnover,
			"minturnover": requiredTurnover,
			"lastdeposit": users.LastDeposit,
			"lastproamount": users.LastProamount,
			"referredby": users.ReferredBy,
			"lastwithdraw": users.LastWithdraw,
			"promotionlog": promotionLogs,
			"pro_status": users.ProStatus,
			"pro_balance": pro_balance,
			"includegames": promotion.Includegames,
			"excludegames": promotion.Excludegames,
		}}
	return c.JSON(response)
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
	if err = db.Debug().Model(models.Promotion{}).Where("status=1 and date(end_date) > date(?)",createdAt.Format("2006-01-02 15:04:05")).Scan(&promotions).Error; err != nil {
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

	var promotion models.Promotion

	err := db.Debug().Model(&models.Promotion{}).Where("id = ?",user.ProStatus).First(&promotion).Error

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
			"lastwithdraw": users.LastWithdraw,
			"pro_status": users.ProStatus,
			"includegames": promotion.Includegames,
			"excludegames": promotion.Excludegames,
		}}
	return c.JSON(response)
}
// @Summary Get user balance
// @Description Get user balance in the database.
// @Tags users
// @Produce json
// @Accept json
// @Security BearerAuth
// @Success 200 {object} models.SwaggerUser
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /users/balance [post]
// @Param user body models.Users true "User balance info"
// @param Authorization header string true "Bearer token"
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
		}}
	return c.JSON(response)
}

// @Summary User Logout
// @Description Get user balance in the database.
// @Tags users
// @Produce json
// @Accept json
// @Security BearerAuth
// @Success 200 {object} models.SwaggerUser
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /users/logout [post]
// @Param user body models.Users true "User Logout"
// @param Authorization header string true "Bearer token"
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

// @Summary Get user Transaction
// @Description Get user Transaction in the database.
// @Tags users
// @Produce json
// @Accept json
// @Security BearerAuth
// @Success 200 {object} models.SwaggerTransactionSub
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /users/transactions [post]
// @Param user body models.TransactionSub true "User Transaction info"
// @param Authorization header string true "Bearer token"
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

// @Summary Get user Statement
// @Description Get user Statement in the database.
// @Tags users
// @Produce json
// @Accept json
// @Security BearerAuth
// @Success 200 {object} models.SwaggerBankStatement
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /users/statement [post]
// @Param user body models.BankStatement true "User Bank Statement info"
// @param Authorization header string true "Bearer token"
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


	proStatusValue, exists := body["pro_status"]
	if !exists {
		return c.JSON(fiber.Map{
			"Status":  false,
			"Message": "pro_status not found",
		})
	}

	userPro := user.ProStatus

	fmt.Printf(" userPrp: %v \n",userPro)
	fmt.Printf(" user.Balance.IsZero(): %v \n",user.Balance.IsZero())
	fmt.Printf(" user.Balance.LessThan: %v \n",user.Balance.LessThan(decimal.NewFromFloat(0.9)))


	if user.Balance.IsZero() || user.Balance.LessThan(decimal.NewFromFloat(0.9)) {
			userPro = ""
	}
	
	pro_setting, err := handler.GetProdetail(db,userPro)
	fmt.Printf("ProSetting: %+v \n",pro_setting)

	if err != nil {
		return c.JSON(fiber.Map{
			"Status": false,
			"Message": "ไม่สามารถตรวจสอบโปรโมชั่นได้",
			"Data": fiber.Map{"id": -1},
		})
	}

	if pro_setting == nil  {


		proStatus := fmt.Sprintf("%v", proStatusValue)
		fmt.Printf("proStatus: %v \n",proStatus)
	   //check maxuse
	   var usageCount int64
	   if err := db.Debug().Model(&models.PromotionLog{}).
		   Where("userid = ? AND promotioncode = ? AND status = 1", user.ID, proStatus).
		   Count(&usageCount).Error; err != nil {
		   return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			   "Status": false,
			   "Message": "ไม่สามารถตรวจสอบประวัติการใช้โปรโมชั่นได้",
			   "Data": fiber.Map{"id": -1},
		   })
	   }


	  
	//    if prodetail["Type"] == "first" {
	// 	// ตรวจสอบว่ามีการใช้โปรโมชั่นแล้วหรือไม่
	// 		if usageCount > 1 {
	// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 				"Status": false,
	// 				"Message": "คุณใช้โปรโมชั่นนี้แล้ว",
	// 				"Data": fiber.Map{"id": -1},
	// 			})
	// 		}
	// 	}


	  

	   prodetail,_ := handler.GetProdetail(db,proStatus)
		// ตรวจสอบกับค่า MaxUse จาก pro_setting
		maxUse, ok := prodetail["MaxUse"].(int64)

		fmt.Printf("maxUse: %v \n",maxUse)
		fmt.Printf("usageCount: %v \n",usageCount)

		if ok && usageCount >= maxUse {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": "คุณใช้โปรโมชั่นนี้เกินจำนวนครั้งที่กำหนดแล้ว \n โปรโมชั่นจะปรับอัตโนมัติ",
				"Data": fiber.Map{"id": -1},
			})
		}

		fmt.Printf("Balance: %v \n",user.Balance)
		fmt.Printf("Round: %v \n",user.Balance.Floor())

		one,_ := decimal.NewFromString("0")

		if user.Balance.Floor().GreaterThan(one) {
			return c.Status(400).JSON(fiber.Map{
				"Message": "คุณมียอดคงเหลือมากกว่าศูนย์",
				"Status":  false,
				"Data":    "คุณมียอดคงเหลือมากกว่าศูนย์",
			})
		}

	   updates := map[string]interface{}{
		   "MinTurnover": prodetail["MinTurnover"],
		   "ProStatus": proStatus,
	   }
	   repository.UpdateFieldsUserString(db, user.Username, updates)
	   return c.Status(200).JSON(fiber.Map{
		   "Status":  true,
		   "Message": "อัปเดตข้อมูลสำเร็จ!",
	   })

		// updates = map[string]interface{}{
		// 	"MinTurnover": pro_setting["MinTurnover"],
		// 	"ProStatus": proStatus,
		// }
		// repository.UpdateFieldsUserString(db, user.Username, updates)
		// return c.Status(200).JSON(fiber.Map{
		// 	"Status":  true,
		// 	"Message": "อัปเดตข้อมูลสำเร็จ!",
		// })
		
	} else {

		

		return c.Status(400).JSON(fiber.Map{
			"Message": "คุณใช้งานโปรโมชั่นอยู่",
			"Status":  false,
			"Data":    "คุณใช้งานโปรโมชั่นอยู่",
		})
		 
	}

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


// ... rest of the code ...

// type Body struct {

// 	//UserID           int             `json:"userid"`
//     //TransactionAmount decimal.Decimal `json:"transactionamount"`
// 	Username           string             `json:"username"`
// 	Status           string             `json:"Status"`
// 	Startdate        string 			`json:"startdate"`
// 	Stopdate        string 		  	`json:"stopdate"`
// 	Prefix           string           	`json:"prefix`
// 	Channel        string 		  	`json:"channel"`

// }
// type UserTransactionSummary struct {
// 	Walletid           int             `json:"walletid"`
// 	UserName		string             `json:"username"`
// 	MemberName		string             `json:"membername"`
// 	BetAmount decimal.Decimal `json:"betamount"`
// 	WINLOSS decimal.Decimal `json:"winloss"`
// 	TURNOVER decimal.Decimal `json:"turnover"`
// }
// type FirstDeposit struct {
//     WalletID          uint
//     Username          string
//     FirstDepositDate  time.Time
//     Firstamount float64 `json:"firstamount"`
// }

// type Response struct {
//     Message string      `json:"message"`
//     Status  bool        `json:"Status"`
//     Data    interface{} `json:"data"` // ใช้ interface{} เพื่อรองรับข้อมูลหลายประเภทใน field data
// }

// func GetUsers(c *fiber.Ctx) error {
// 	var user []models.Users
//  db, _ := database.ConnectToDB(users.Prefix)
// 	db.Find(&user)
// 	response := fiber.Map{
// 		"Message": "สำเร็จ!!",
// 		"Status":  true,
// 		"Data": fiber.Map{
// 			"users":user,
// 		},
// 	}
// 	return c.JSON(response)

// }
// func GetUser(c *fiber.Ctx) error {

// 	user := c.Locals("user").(*jtoken.Token)
// 	claims := user.Claims.(jtoken.MapClaims)

// 	var users []models.Users
//   db, _ := database.ConnectToDB(users.Prefix)
// 	err := db.Find(&users,claims["ID"]).Error

// 	if err != nil {
// 		response := fiber.Map{
// 			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
// 			"Status":  false,

// 		}
// 		return c.JSON(response)
// 	}else {
// 		response := fiber.Map{
// 			"Message": "สำเร็จ!!",
// 			"Status":  true,
// 			"Data": fiber.Map{
// 				"users":users,
// 			},
// 		}
// 		return c.JSON(response)
// 	}

// }

// func GetUserByID(c *fiber.Ctx) error {

// 	user := c.Locals("user").(*jtoken.Token)
// 	claims := user.Claims.(jtoken.MapClaims)
// 	response := fiber.Map{
// 		"Message": "ไม่พบรหัสผู้ใช้งาน!!",
// 		"Status":  false,
// 	}

// 	var users models.Users

// 	db.Debug().Where("id= ?",claims["ID"]).Find(&users)

// 	if users == (models.Users{}) {

// 		response = fiber.Map{
// 			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
// 			"Status":  false,
// 		}

// 	}

// 	tokenString := c.Get("Authorization")[7:]

// 	_err := handler.validateJWT(tokenString);
// 	//fmt.Println(_err)
// 	 if _err != nil {

// 		  response = fiber.Map{
// 			"Message": "โทเคนไม่ถูกต้อง!!",
// 			"Status":  false,
// 		}

// 	}else {

// 		response = fiber.Map{
// 			"Status": true,
// 			"Message": "สำเร็จ",
// 			"Data": fiber.Map{
// 			   "userid":   users.ID,
// 			   "username": users.Username,
// 			   "fullname": users.Fullname,
// 			   "createdAt": users.CreatedAt,
// 			   "backnumber": users.Banknumber,
// 			   "bankname": users.Bankname,
// 			   "balance": users.Balance,
// 			   "Status": users.Status,

// 			},
// 		}

// 	}

// 	 return c.JSON(response)
// }
// func GetBalanceFromID(c *fiber.Ctx) error {

// 	user := new(models.Users)

// 	if err := c.BodyParser(user); err != nil {
// 		return c.Status(200).SendString(err.Error())
// 	}
// 	var users models.Users
// 	var bankstatement models.BankStatement
// 	//fmt.Println(user.Username)
//     // ดึงข้อมูลโดยใช้ Username และ Password
//     if err := db.Where("username = ?", user.Username).First(&users).Error; err != nil {
//         return  errors.New("user not found")
//     }
// 	result := db.Debug().Select("Uid,transactionamount,status").Where("walletid = ? and channel=? ",users.ID,"1stpay").Order("id Desc").First(&bankstatement)

// 	if result.Error != nil {

// 			fmt.Println("ไม่พบข้อมูล")
// 			return c.Status(200).JSON(fiber.Map{
// 				"Status": true,
// 				"data": fiber.Map{
// 					"id": &users.ID,
// 					"token": &users.Token,
// 					"balance": &users.Balance,
// 					"withdraw": fiber.Map{
// 						"uid":"",
// 						"transactionamount":0,
// 						"Status":"verified",
// 					},
// 				}})

// 	} else {
// 		fmt.Printf("ข้อมูลที่พบ: %+v\n", bankstatement)

// 		return c.Status(200).JSON(fiber.Map{
// 			"Status": true,
// 			"data": fiber.Map{
// 				"id": &users.ID,
// 				"token": &users.Token,
// 				"balance": &users.Balance,
// 				"withdraw": fiber.Map{
// 					"uid":&bankstatement.Uid,
// 					"transactionamount": &bankstatement.Transactionamount,
// 					"Status":&bankstatement.Status,
// 				},
// 			}})
// 	}
// 	//db.Where("username = ?",user.Username).Find(&users)

// }
// func AddUser(c *fiber.Ctx) error {

// 	var currency =  os.Getenv("CURRENCY")
// 	user := new(models.Users)

// 	if err := c.BodyParser(user); err != nil {
// 		return c.Status(200).SendString(err.Error())
// 	}
// 	//user.Walletid = user.ID
// 	//user.Username = user.Prefix + user.Username + currency
// 	result := db.Create(&user);

// 	// ส่ง response เป็น JSON

// 	if result.Error != nil {
// 		response := fiber.Map{
// 			"Message": "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
// 			"Status":  false,
// 			"Data":    fiber.Map{
// 				"id": -1,
// 			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
// 		}
// 		return c.JSON(response)
// 		} else {

// 			updates := map[string]interface{}{
// 				"Walletid":user.ID,
// 				"Preferredname": user.Username,
// 				"Username":user.Prefix + user.Username + currency,
// 			}
// 			if err := db.Model(&user).Updates(updates).Error; err != nil {
// 				return errors.New("มีข้อผิดพลาด")
// 			}

// 		response := fiber.Map{
// 			"Message": "เพิ่มยูสเซอร์สำเร็จ!",
// 			"Status":  true,
// 			"Data":    fiber.Map{
// 				"id": user.ID,
// 				"walletid":user.ID,
// 			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
// 		}
// 		return c.JSON(response)
// 	}

// }

// func GetUserStatement(c *fiber.Ctx) error {

// 	body := new(Body)
// 	if err := c.BodyParser(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":err.Error(),
// 		})
// 	}
// 	response := fiber.Map{
// 		"Status": false,
// 		"Message": "สำเร็จ",
// 		"Data": map[string]interface{}{},
// 	}

// 	user := c.Locals("user").(*jtoken.Token)

// 	claims := user.Claims.(jtoken.MapClaims)

// 	//fmt.Println(claims)
// 	//username := claims["username"].(string)
// 	id := claims["walletid"]

// 	  tokenString := c.Get("Authorization")[7:]

// 	  _err := handler.validateJWT(tokenString);
// 	  //fmt.Println(_err)
// 	   if _err != nil {
// 			response = fiber.Map{
// 				"Status": false,
// 				"Message": "ไม่สำเร็จ",
// 				//"Data": map[string]interface{}
// 				}
// 			} else {

// 		channel := body.Channel
// 		startDateStr := body.Startdate
// 		endDateStr := body.Stopdate

// 		var statements []models.BankStatement

// 		if body.Status == "all" {
// 			db.Debug().Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? ",id, channel, startDateStr, endDateStr).Find(&statements)
// 		} else {
// 			db.Debug().Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?", id,channel, startDateStr, endDateStr,body.Status).Find(&statements)
// 		}

// 		  // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
// 		  result := make([]fiber.Map, len(statements))

// 		  // วนลูปเพื่อประมวลผลแต่ละรายการ
// 		   for i, transaction := range statements {
// 			   // ตรวจสอบเงื่อนไขด้วย inline if-else
// 			   transactionType := func(amount decimal.Decimal,channel string) string {
// 				if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
// 					return "ถอน"
// 				}
// 				return "ฝาก"
// 			}(transaction.Transactionamount,transaction.Channel)
// 			amountFloat, _ := transaction.Transactionamount.Float64()
// 			   // เก็บผลลัพธ์ใน slice
// 			   result[i] = fiber.Map{
// 				   "userid":           transaction.Userid,
// 				   "createdAt": transaction.CreatedAt,
// 				   "transfer_amount": amountFloat,
// 				   "credit":  amountFloat,
// 				   "Status":           transaction.Status,
// 				   "channel": transaction.Channel,
// 				   "statement_type": transactionType,
// 				   "expire_date": transaction.CreatedAt,
// 			   }
// 		   }

// 		   response = fiber.Map{
// 			"Status": true,
// 			"Message": "สำเร็จ",
// 			"Data": result,
// 			}
// 		}
// 	return c.JSON(response)

// }
// func GetIdStatement(c *fiber.Ctx) error {

// 	body := new(Body)
// 	if err := c.BodyParser(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":err.Error(),
// 		})
// 	}

// 	    username := body.Username
// 		channel := body.Channel
// 		startDateStr := body.Startdate
// 		endDateStr := body.Stopdate

// 		var users models.Users
// 		db.Debug().Where("username=?",username).Find(&users)

// 		var statements []models.BankStatement

// 		if body.Status == "all" {
// 			db.Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount,COALESCE(bet_amount,0) as TURNOVER,createdAt,Beforebalance,Balance,Channel,Uid").Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? ",users.ID, channel, startDateStr, endDateStr).Find(&statements)
// 		} else {
// 			db.Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount,COALESCE(bet_amount,0) as TURNOVER,createdAt,Beforebalance,Balance,Channel,Uid").Where("id=? AND channel <> ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?", users.ID,channel, startDateStr, endDateStr,body.Status).Find(&statements)
// 		}

// 		  // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
// 		  result := make([]fiber.Map, len(statements))

// 		  // วนลูปเพื่อประมวลผลแต่ละรายการ
// 		   for i, transaction := range statements {
// 			   // ตรวจสอบเงื่อนไขด้วย inline if-else
// 			   transactionType := func(amount decimal.Decimal,channel string) string {
// 				if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
// 					return "เสีย"
// 				}
// 				return "ได้"
// 			}(transaction.Transactionamount,transaction.Channel)
// 			amountFloat, _ := transaction.Transactionamount.Float64()
// 			   // เก็บผลลัพธ์ใน slice
// 			   result[i] = fiber.Map{
// 				   "MemberID":           users.Username,
// 				   "MemberName":  users.Fullname,
// 				   "createdAt": transaction.CreatedAt,
// 				   "PayoutAmount": amountFloat,
// 				   "BetAmount": transaction.BetAmount,
// 				   "BeforeBalance": transaction.Beforebalance,
// 				   "AfterBalance": transaction.Balance,
// 				   "WINLOSS": transaction.Transactionamount,
// 				   "TURNOVER": transaction.BetAmount,
// 				   "credit":  amountFloat,
// 				   "Status":           transaction.Status,
// 				   "GameRoundID": transaction.Uid,
// 				   "GameProvide": transaction.Channel,
// 				   "statement_type": transactionType,
// 				   "expire_date": transaction.CreatedAt,
// 			   }
// 		   }

// 		   return c.Status(200).JSON(fiber.Map{
// 			"Status": true,
// 			"data": result,
// 		})

// }
// func GetUserAll(c *fiber.Ctx) error {

// 	body := new(Body)
// 	if err := c.BodyParser(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":err.Error(),
// 		})
// 	}

// 	var users []models.Users

// 	db.Debug().Select(" *,CASE WHEN Walletid in (Select distinct walletid from BankStatement Where channel='1stpay') THEN 1 ELSE 0 END As ProStatus ").Find(&users)

// 	return c.Status(200).JSON(fiber.Map{
// 		"Status": true,
// 		"data": fiber.Map{
// 			"data":users,
// 		}})
// }
// func GetUserAllStatement(c *fiber.Ctx) error {

// 	body := new(Body)
// 	if err := c.BodyParser(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":err.Error(),
// 		})
// 	}

// 		username := body.Username

// 		var users models.Users
// 		db.Debug().Where("username=?",username).Find(&users)

// 		channel := body.Channel

// 		startDateStr := body.Startdate
// 		endDateStr := body.Stopdate

// 		// ตั้งค่าช่วงวันที่ในการค้นหา

// 		var statements []models.BankStatement

// 		if channel == "game" {
// 			if body.Status == "all" {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",users.ID, startDateStr, endDateStr).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
// 			} else {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID , startDateStr, endDateStr,body.Status).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
// 			}
// 		} else {
// 			if body.Status == "all" {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? ",users.ID,channel, startDateStr, endDateStr).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
// 			} else {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,createdAt(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,Beforebalance,balance,transactionamount,COALESCE(bet_amount,0) as TURNOVER,channel,status").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID ,channel, startDateStr, endDateStr,body.Status).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
// 			}

// 		}

// 		  // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
// 		  result := make([]fiber.Map, len(statements))

// 		  // วนลูปเพื่อประมวลผลแต่ละรายการ
// 		   for i, transaction := range statements {

// 			   amountFloat, _ := transaction.Transactionamount.Float64()
// 			   balanceFloat, _ := transaction.Balance.Float64()
// 			   beforeFloat,_ := transaction.Beforebalance.Float64()
// 			   betFloat,_ := transaction.BetAmount.Float64()

// 			   // เก็บผลลัพธ์ใน slice
// 			   result[i] = fiber.Map{
// 		 		   "MemberID":           transaction.Walletid,
// 					"MemberName": users.Username,
// 					"UserName": users.Fullname,
// 		 		   "createdAt": transaction.CreatedAt,
// 				   "GameProvide": transaction.Channel,
// 				   "GameRoundID": transaction.Uid,
// 		 		   "BeforeBalance": beforeFloat,
// 		  		   "BetAmount": betFloat,
// 				   "WINLOSS": amountFloat,
// 				   "AfterBalance": balanceFloat,
// 				   "TURNOVER":betFloat,
// 		 		   "channel": transaction.Channel,
// 				   "Status":transaction.Status,

// 		 	   }
// 		    }

// 		   return c.Status(200).JSON(fiber.Map{
// 			"Status": true,
// 			"data": result,
// 		})

// }
// func GetUserDetailStatement(c *fiber.Ctx) error {

// 	body := new(Body)
// 	if err := c.BodyParser(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":err.Error(),
// 		})
// 	}

// 		username := body.Username

// 		var users models.Users
// 		db.Debug().Where("username=?",username).Find(&users)

// 		channel := body.Channel

// 		startDateStr := body.Startdate
// 		endDateStr := body.Stopdate

// 		var statements []models.BankStatement

// 		if channel == "game" {
// 			if body.Status == "all" {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",users.ID, startDateStr, endDateStr).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
// 			} else {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID , startDateStr, endDateStr,body.Status).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
// 			}
// 		} else {
// 			if body.Status == "all" {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? ",users.ID,channel, startDateStr, endDateStr).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? ",Users.id, startDate, endDate).Find(&statements)
// 			} else {
// 				db.Model(&models.BankStatement{}).Debug().Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,COALESCE(bet_amount,0) as BetAmount,transactionamount as WINLOSS,COALESCE(bet_amount,0) as TURNOVER").Where("walletid = ? AND channel = ? AND  DATE(createdat) BETWEEN ? AND ? and status = ?",users.ID ,channel, startDateStr, endDateStr,body.Status).Scan(&statements)
// 				//db.Debug().Where("walletid = ? AND channel<>'1stpay' AND  DATE(createdat) BETWEEN ? AND ? and status = ?",Users.id , startDate, endDate,body.Status).Find(&statements)
// 			}

// 		}

// 		   return c.Status(200).JSON(fiber.Map{
// 			"Status": true,
// 			"data": statements,
// 		})

// }
// func GetUserSumStatement(c *fiber.Ctx) error {

// 	body := new(Body)
// 	if err := c.BodyParser(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":err.Error(),
// 		})
// 	}

// 		channel := body.Channel
// 		startDateStr := body.Startdate
// 		endDateStr := body.Stopdate

// 		var summaries []UserTransactionSummary

// 		//if body.Status == "all" {
// 		err := db.Debug().Model(&models.BankStatement{}).Select("walletid,(select username FROM Users WHERE Users.id=BankStatement.walletid) as UserName,(select fullname FROM Users WHERE Users.id=BankStatement.walletid) as MemberName,COALESCE(sum(bet_amount),0) as BetAmount,sum(transactionamount) as WINLOSS,COALESCE(sum(bet_amount),0) as TURNOVER").Where("channel != ? AND  DATE(createdat) BETWEEN ? AND ? AND (status='101' OR status=0)", channel,startDateStr, endDateStr).Group("BankStatement.walletid").Scan(&summaries).Error

// 		if err != nil {
// 			fmt.Println(err)
// 			// ส่ง Error เป็น JSON
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}

// 		// ส่งข้อมูล JSON
// 		return c.Status(200).JSON(fiber.Map{
// 			"Status": true,
// 			"data": summaries,
// 		})

// }
// func UpdateToken(c *fiber.Ctx) error {

// 	user := c.Locals("user").(*jtoken.Token)
// 	claims := user.Claims.(jtoken.MapClaims)

// 	tokenString := c.Get("Authorization")[7:]
// 	var users models.Users
// 	db.Debug().Where("walletid = ?",claims["walletid"]).Find(&users)
// 	fmt.Println(users)

// 	updates := map[string]interface{}{
// 		"Token":tokenString,
// 	}

// 	_err := repository.UpdateUserFields(db, users.ID, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
// 	if _err != nil {
// 		return c.Status(200).JSON(fiber.Map{
// 			"Status": false,
// 			"Message":  _err,
// 		})
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"Status": true,
// 		"Message": "สำเร็จ!",
// 	 })
// }
// func GetFirstUsers(c *fiber.Ctx) error {

// 	type FirstGetResponse struct {
// 		//Counter           int             `json:"counter"`
// 		//username         string 	  `json:"username"`
// 		Walletid           int             `json:"walletid"`
// 		Firstamount       decimal.Decimal `json:"firstamount"`
// 		Firstdate         string 	  `json:"firstdate"`

// 	}

// 	type FirstResponse struct {
// 		Counter           int             `json:"counter"`
// 		Active            int  `json:"active"`
// 		Firstdate         string 	  `json:"firstdate"`

// 	}

// 	body := new(Body)
// 	if err := c.BodyParser(body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":err.Error(),
// 		})
// 	}

// 		username := body.Username
// 		prefix := body.Prefix

// 		var users models.Users
// 		db.Debug().Where("username=?",username).Find(&users)

// 		//channel := body.Channel

// 		startDateStr := body.Startdate
// 		endDateStr := body.Stopdate

// 		//var results []FirstGetResponse
// 		// var result []struct {
// 		// 	Username          string
// 		// 	TransactionAmount float64
// 		// 	CreatedAt         time.Time
// 		// }

// 		// subQuery := db.Model(&models.BankStatement{}).
// 		// Select("walletid, MIN(createdAt) AS firstdate").
// 		// Where("status = ? AND transactionamount > 0 AND deleted_at IS NULL", "verified").
// 		// Group("walletid").
// 		// Having("DATE(MIN(createdAt)) BETWEEN ? AND ?", startDateStr,endDateStr)

// 		// db.Debug().Table("Users AS u").
// 		// Select("u.walletid,u.username, b.transactionamount AS firstamount, b.createdAt AS firstdate").
// 		// Joins("JOIN BankStatement AS b ON u.walletid = b.walletid").
// 		// Joins("JOIN (?) AS first_deposit ON b.walletid = first_deposit.walletid AND b.createdAt = first_deposit.firstdate", subQuery).
// 		// Where("u.prefix LIKE ?", prefix+"%").
// 		// Where("u.id IN (SELECT id FROM Users)").
// 		// Order("u.walletid").
// 		// Scan(&results)
// 		var firstDeposits []FirstDeposit
// 		// ตั้งค่าช่วงวันที่ในการค้นหา
// 		db.Debug().Model(&models.Users{}).
//         Select("Users.id, Users.username, MIN(BankStatement.createdAt) AS first_deposit_date, (SELECT c.transactionamount from BankStatement as c WHERE c.id=Users.id AND DATE(c.createdAt) BETWEEN ? AND ? and c.channel='1stpay' and status='verified' order by c.id ASC LIMIT 1) AS firstamount").
//         Joins("JOIN BankStatement ON Users.id = BankStatement.walletid").
//         Where("BankStatement.transactionamount > 0").
//         Where("DATE(Users.createdAt) BETWEEN ? AND ? ", startDateStr, endDateStr,startDateStr, endDateStr).
//         Group("Users.id, Users.username").
//         Scan(&firstDeposits)

// 		// คำนวณยอดรวมของ first_deposit_amount
// 		var totalFirstDepositAmount float64
// 		for _, deposit := range firstDeposits {

// 			totalFirstDepositAmount += deposit.Firstamount
// 		}

// 		// แสดงยอดรวม
// 		fmt.Printf("Total First Deposit Amount: %.2f\n", totalFirstDepositAmount)

// 		// แสดงผลลัพธ์แต่ละรายการ
// 		for _, deposit := range firstDeposits {
// 			fmt.Printf("WalletID: %d, Username: %s, First Deposit Date: %s, First Deposit Amount: %.2f\n",
// 				deposit.WalletID, deposit.Username, deposit.FirstDepositDate.Format("2006-01-02 15:04:05"), deposit.Firstamount)
// 		}

// 		var statements FirstResponse
// 		//var firststate []FirstGetResponse
// 		//var secondstate []FirstGetResponse

// 		db.Model(&models.Users{}).Debug().Select("count(id) as counter").Where("DATE(createdat) BETWEEN ? AND ?  and username like ?",startDateStr, endDateStr,prefix+"%").Scan(&statements)
// 		//db.Model(&models.BankStatement{}).Debug().Select("Walletid,min(createdAt) as Firstdate,transactionamount as Firstamount").Where("transactionamount>0 AND  DATE(createdat) BETWEEN ? AND ? and status=?",startDateStr, endDateStr,"verified").Group("walletid,transactionamount").Scan(&firststate)

// 		//   // สร้าง slice เพื่อเก็บผลลัพธ์หลังจากตรวจสอบเงื่อนไข
// 		//   result := make([]fiber.Map, len(firststate))

// 		// //   // วนลูปเพื่อประมวลผลแต่ละรายการ
// 		//    for i, transaction := range firststate {
// 		// // 	   // ตรวจสอบเงื่อนไขด้วย inline if-else
// 		// // 	   transactionType := func(amount decimal.Decimal,channel string) string {
// 		// // 		if amount.LessThan(decimal.NewFromInt(0)) { // ใช้ LessThan สำหรับการเปรียบเทียบ
// 		// // 			return "ถอน"
// 		// // 		}
// 		// // 		return "ฝาก"
// 		// // 	}(transaction.Transactionamount,transaction.Channel)
// 		// // 	//fmt.Println(transaction.Bet_amount)
// 		// //fmt.Println(transaction)
// 		//  	amountFloat, _ := transaction.firstamount.Float64()
// 		// // 	balanceFloat, _ := transaction.Balance.Float64()
// 		// // 	beforeFloat,_ := transaction.Beforebalance.Float64()
// 		// // 	betFloat,_ := transaction.Bet_amount.Float64()

// 		// // 	   // เก็บผลลัพธ์ใน slice
// 		// 	   result[i] = fiber.Map{
// 		// 		     "walletid":           transaction.walletid,
// 		// 		//    "username": users.Username,
// 		// 		    "firstdate": transaction.firstdate,
// 		// 			"firstamount":amountFloat,
// 		// 		//    "beforebalnce": beforeFloat,
// 		// 		//    "betamount": betFloat,
// 		// 		//    "transfer_amount": amountFloat,
// 		// 		//    "balance": balanceFloat,
// 		// 		//    "credit":  amountFloat,
// 		// 		//    "Status":           transaction.Status,
// 		// 		//    "channel": transaction.Channel,
// 		// 		//    "statement_type": transactionType,
// 		// 		//    "expire_date": transaction.CreatedAt,
// 		// 	   }
// 		//    }
// 			// if firststate==nil {
// 			// 	return c.Status(200).JSON(fiber.Map{
// 			// 		"Status": true,
// 			// 		"Message": "สำเร็จ!",
// 			// 		"data": fiber.Map{
// 			// 			"Counter": statements.Counter,
// 			// 			"Active": make([]FirstGetResponse, 0),
// 			// 		},
// 			// 	 })

// 			// } else {

// 				return c.Status(200).JSON(fiber.Map{
// 					"Status": true,
// 					"Message": "สำเร็จ!",
// 					"data": fiber.Map{
// 						"Counter": statements.Counter,
// 						"Active": firstDeposits,
// 						"Firstamount":totalFirstDepositAmount,
// 					},
// 				})
// 		//}

// }
