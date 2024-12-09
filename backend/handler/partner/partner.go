package partner

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
	//wallet "hanoi/handler/wallet"
	// 	//"github.com/golang-jwt/jwt"
	// 	//jtoken "github.com/golang-jwt/jwt/v4"
	// 	//"github.com/solrac97gr/basic-jwt-auth/config"
	// 	//"github.com/solrac97gr/basic-jwt-auth/models"
	// 	//"github.com/solrac97gr/basic-jwt-auth/repository"
	"hanoi/repository"
	"hanoi/encrypt"
	//"log"
	// 	// "net"
	// 	// "net/http"
	"os"
	// 	// "strconv"
	//"time"
	"fmt"
	"strings"
	//"errors"
)

type ErrorResponse struct {
	Status  bool   `json:"Status"`
	Message string `json:"message"`
}

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

type Dbstruct struct {
	DBName   string   `json:"dbname"`
	Prefix   string   `json:"prefix"`
	Username string   `json:"username"`
	Dbnames  []string `json:"dbnames"`
}
type PartnerBody struct {
	Prefix string `json:"prefix"`
	ID     int    `json:"id"`
	Body   struct {
		Fullname   string `json:"fullname"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Status     int    `json:"status"`
		Bankname   string `json:"bankname"`
		Banknumber string `json:"banknumber"`
		ProStatus  string `json:"prostatus"`
		MinTurnoverDef string `json:"minturnoverdef"`
	}
}

type RequestBody struct {
	Prefix string      `json:"prefix"`
	Body   models.Partner    `json:"body"` // หรือใช้โครงสร้างที่เหมาะสมกับข้อมูลใน body
}

var jwtSecret = os.Getenv("PASSWORD_SECRET")

// @Summary Login user
// @Description Get a list of all users in the database.
// @Tags users
// @Produce json
// @Success 200 {array} models.SwaggerUser
// @Router /users/login [post]
// @Param user body models.Partner true "User registration info"
func Signing(c *fiber.Ctx) error {
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
	var user models.Partner

	// fmt.Printf("%s",loginRequest)

	db, err := database.ConnectToDB(loginRequest.Prefix)
	//db.AutoMigrate(&models.BankStatement{},&models.PromotionLog{})
	//db.Migrator().CreateTable(&models.PromotionLog{})
	//db.AutoMigrate(&models.Partner{})
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
		//"Walletid":    user.Walletid,
		"Username":    user.Name,
		//"Role":        user.Role,
		"PartnersKey": user.AffiliateKey,
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
	response := fiber.Map{
		"Token":  t,
		"Status": true,
	}
	return c.JSON(response)
	// return c.JSON(models.LoginResponse{
	// 	Token: t,
	// })

}

func Login(c *fiber.Ctx) error {


	loginRequest := new(Body)

	if err := c.BodyParser(loginRequest); err != nil {
		response := fiber.Map{
			"Message": "ไม่พบรหัสผู้ใช้งาน!!" + err.Error(),
			"Status":  false,
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}





	//db, err := database.ConnectToDB(loginRequest.Prefix)

// ค้นหา key ที่อยู่ใน database master จากตาราง settings
	//var settings []models.Settings

	// ใช้ username เพื่อค้นหาค่าที่มีชื่อ username เป็นเริ่มต้นตรงกับ key
	// usernamePrefix := loginRequest.Username[:3] // ใช้ 3 ตัวแรกของ username
	// if len(loginRequest.Username) < 3 {
	// 	usernamePrefix = loginRequest.Username // หาก username สั้นกว่า 3 ตัว ให้ใช้ username ทั้งหมด
	// }

	// fmt.Printf(" %s ",user)

	// settings,err := database.GetMaster(usernamePrefix)
	// if err != nil || len(settings) == 0 {
	// 	response := fiber.Map{
	// 		"Message": "ไม่พบการตั้งค่าที่ตรงกับชื่อผู้ใช้งาน!",
	// 		"Status":  false,
	// 	}
	// 	return c.Status(fiber.StatusUnauthorized).JSON(response)
	// }

	// fmt.Printf(" %s ",settings)
	var databaseList,err := getDataList()
	var t string
	if err != nil || len(databaseList) == 0 {
		response := fiber.Map{
			"Message": "ไม่พบการตั้งค่าที่ตรงกับชื่อผู้ใช้งาน!",
			"Status":  false,
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	for _, dbInfo := range databaseList {
		dbConnection, connErr := database.ConnectToDB(dbInfo.Prefix) // เรียกใช้ connectDBP
		if connErr == nil {
			// เก็บ db ใน c.Locals
			c.Locals("db", dbConnection)

			// ค้นหาผู้ใช้ใน partners
			var partner models.Partner
			err = dbConnection.Where("preferredname = ? AND password = ?", loginRequest.Username, loginRequest.Password).First(&partner).Error
			if err == nil {
				// ผู้ใช้พบแล้ว ทำการดำเนินการต่อ

				claims := jtoken.MapClaims{
					"ID":          partner.ID,
					//"Walletid":    user.Walletid,
					"Username":    partner.Name,
					//"Role":        user.Role,
					"PartnersKey": partner.AffiliateKey,
					"Prefix":      partner.Prefix,
					//"exp":   time.Now().Add(day * 1).Unix(),
				}
			 
				token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
			
				t, _ = token.SignedString([]byte(jwtSecret))
				// response := fiber.Map{
				// 	"Token":  t,
				// 	"Status": true,
				// }
				// return c.JSON(response)
				break
			}
		}
	}

	// err = db.Raw("SELECT * FROM settings WHERE `key` LIKE ?", usernamePrefix+"%").Scan(&settings).Error
	// if err != nil || len(settings) == 0 {
	// 	response := fiber.Map{
	// 		"Message": "ไม่พบการตั้งค่าที่ตรงกับชื่อผู้ใช้งาน!",
	// 		"Status":  false,
	// 	}
	// 	return c.Status(fiber.StatusUnauthorized).JSON(response)
	// }

	// // เชื่อมต่อกับ database ตาม key ที่พบ
	// for _, setting := range settings {
	// 	// เชื่อมต่อกับฐานข้อมูลที่เกี่ยวข้อง
	// 	dbConnection, connErr := database.connectDBP(setting.Key) // เรียกใช้ connectDBP
	// 	if connErr == nil {
	// 		// เก็บ db ใน c.Locals
	// 		c.Locals("db", dbConnection)
	// 		break
	// 	}
	// }

	// ตรวจสอบว่ามีการเชื่อมต่อหรือไม่
	// if c.Locals("db") == nil {
	// 	response := fiber.Map{
	// 		"Message": "ไม่สามารถเชื่อมต่อกับฐานข้อมูลที่เกี่ยวข้อง!",
	// 		"Status":  false,
	// 	}
	// 	return c.Status(fiber.StatusInternalServerError).JSON(response)
	// }
	// var partner models.Partner
	// err = dbConnection.Where("preferredname = ? AND password = ?", loginRequest.Username, loginRequest.Password).First(&partner).Error;
	 
	// if err != nil {
	// 	response := fiber.Map{
	// 		"Message": "ไม่พบรหัสผู้ใช้งาน!!",
	// 		"Status":  false,
	// 	}
	// 	return c.Status(fiber.StatusUnauthorized).JSON(response)
	// }

	// if err != nil {
	// 	response := fiber.Map{
	// 		"Message": "ไม่พบรหัสผู้ใช้งาน!!",
	// 		"Status":  false,
	// 	}
	// 	return c.Status(fiber.StatusUnauthorized).JSON(response)
	// }

	// //day := time.Hour * 24

	// claims := jtoken.MapClaims{
	// 	"ID":          partner.ID,
	// 	//"Walletid":    user.Walletid,
	// 	"Username":    partner.Name,
	// 	//"Role":        user.Role,
	// 	"PartnersKey": partner.AffiliateKey,
	// 	"Prefix":      partner.Prefix,
	// 	//"exp":   time.Now().Add(day * 1).Unix(),
	// }
 
	// token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	// t, err := token.SignedString([]byte(jwtSecret))

	// updates := map[string]interface{}{
	// 	"Token": t,
	// }

	// อัปเดตข้อมูลยูสเซอร์
	// _err := repository.UpdateUserFields(dbConnection, user.ID, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	// if _err != nil {
	// 	if err != nil {
	// 		response := fiber.Map{
	// 			"Message": "ดึงข้อมูลผิดพลาด",
	// 			"Status":  false,
	// 			"Data":    err.Error(),
	// 		}
	// 		return c.Status(fiber.StatusUnauthorized).JSON(response)
	// 	}
	// } else {
	// 	fmt.Println("User fields updated successfully")
	// }


	response := fiber.Map{
		"Token":  t,
		"Status": true,
	}
	return c.JSON(response)
}

// @Summary Get all users
// @Description Get a list of all users in the database.
// @Tags users
// @Produce json
// @Success 200 {array} models.SwaggerUser
// @Router /users/all [post]
// @Param user body models.SwaggerBody true "User registration info"
// @param Authorization header string true "Bearer token"
func GetPartners(c *fiber.Ctx) error {

	body := new(Dbstruct)
	if err := c.BodyParser(body); err != nil {
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

 
	db, err := database.ConnectToDB(body.Prefix)
	if err != nil {
		response := fiber.Map{
			"Message": "ติดต่อฐานข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

	var users []models.Partner

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
// @Param user body models.Partner true "Partner registration info"
func Register(c *fiber.Ctx) error {

	var currency = os.Getenv("CURRENCY")
	// type RequestBody struct {
	// 	Prefix string      `json:"prefix"`
	// 	Body   models.Partner    `json:"body"` // หรือใช้โครงสร้างที่เหมาะสมกับข้อมูลใน body
	// }

	var partner RequestBody

	//user := new(models.Partner)

	if err := c.BodyParser(&partner); err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": "ไม่สามารถแปลงข้อมูลได้: " + err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// fmt.Printf(" %s ",user)
	db, conn := database.ConnectToDB(partner.Prefix)
	if conn != nil {
		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่่ พบข้อมูล Prefix!",
			Status:  false,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	fmt.Printf("partner: %+v \n",partner.Body)
	//seedPhrase,_ := encrypt.GenerateAffiliateCode(5) //handler.GenerateReferralCode(user.Username,1)

	//fmt.Printf("SeedPhase  %s\n", seedPhrase) 

	//user.AffiliateKey = seedPhrase

	
	if partner.Body.Name == "" {
		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Status:  false,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	result := db.Debug().Create(&partner.Body)
	
	if result.Error != nil {
		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Status:  false,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	} else {
		fmt.Printf("Result : %+v \n",result)
		updates := map[string]interface{}{
			//"Partnerid":      partner.Body.ID,
			//"Preferredname": partner.Body.Name,
			"Username":      strings.ToUpper(partner.Body.Username) + currency,
			"Currency":      currency,
			//"Actived": nil,
			//"AffiliateKey": partner.Body.AffiliateKey,
		}
		if err := db.Debug().Model(&models.Partner{}).Where("id = ?",partner.Body.ID).Updates(updates).Error; err != nil {
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
				"id":       partner.Body.ID,
				"partnerid": partner.Body.ID,
				"Partnername": partner.Body.Name,
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
// @Param user body models.Partner true "User user info"
// @param Authorization header string true "Bearer token"
func GetPartner(c *fiber.Ctx) error {

	var users models.Partner

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
	
 
	// var promotionLog models.PromotionLog
	// db.Where("userid = ? AND promotioncode = ? AND status = 1", users.ID, users.ProStatus).
	// 	Order("created_at DESC").
	// 	First(&promotionLog)

	// fmt.Printf("PromotionLog: %+v \n",promotionLog)

	// var totalTurnover decimal.Decimal
	// if err := db.Model(&models.TransactionSub{}).
	// 	Where("proid = ? AND membername = ? AND created_at >= ?", 
	// 		users.ProStatus, 
	// 		users.Username, 
	// 		promotionLog.CreatedAt).
	// 	Select("COALESCE(SUM(turnover), 0)").
	// 	Scan(&totalTurnover).Error; err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"Message": "ไม่สามารถคำนวณยอดเทิร์นได้ !",
	// 			"Status": false,
	// 			"Data": "เกิดข้อผิดพลาด!",
	// 		})
			
		
	// }


	// var promotion models.Promotion
	// //fmt.Println(summary.Turnover)
	// if users.ProStatus != "" {
	// 	db.Debug().Model(&models.Promotion{}).Select("Includegames,Excludegames").Where("ID = ?",users.ProStatus).Scan(&promotion)
	// }
	// pro_setting, err := wallet.CheckPro(db, &users) 
	// if err != nil {
	// 	fmt.Printf("388 error: %s \n",err)
	// 	return c.JSON(fiber.Map{
	// 		"status": false,
	// 		"message":  err.Error(),
	// 		"data": fiber.Map{
	// 			"id": -1,
	// 		}})
	// }

	// var minTurnover string
	// if pro_setting["MinTurnover"] != nil {
	// 	minTurnover = fmt.Sprintf("%v", pro_setting["MinTurnover"])
	// } else {
	// 	minTurnover = "0" // ค่าเริ่มต้นเมื่อเป็น nil
	// }

	// var baseAmount decimal.Decimal
	// if pro_setting["MinSpendType"] == "deposit" {
	// 	baseAmount = users.LastDeposit
	// } else {
	// 	baseAmount = users.LastDeposit.Add(users.LastProamount)
	// }

	// if minTurnover == "" {
	// 	minTurnover = "0"
	// }
	// fmt.Printf("minTurnover: %+v \n",minTurnover)
	// fmt.Printf("baseAmount: %+v \n",baseAmount)


	// requiredTurnover, err := wallet.CalculateRequiredTurnover(minTurnover, baseAmount)
 
	// if err != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message": "ไม่สามารถคำนวณยอดเทิร์นได้",
	// 		"Data": fiber.Map{"id": -1},
	// 	})
	// }

	// // var transaction models.TransactionSub
	// // db.Debug().Model(&models.TransactionSub{}).Select("COALESCE(balance,0) as balance").Where("membername= ? and deleted_at is null", users.Username).Scan(&transaction)
	// // var pro_balance decimal.Decimal
	// // db.Debug().Model(&models.TransactionSub{}).
	// // 	Select("COALESCE(balance, 0) as balance").
	// // 	Where("membername = ? AND deleted_at is null AND ProID=? and created_at > ?", users.Username,users.ProStatus,time.Now().Format("2006-01-02 15:04:05")).Order("id DESC").Limit(1).Scan(&pro_balance)
	// // createdAt := time.Now()
	// // if pro_setting["CreatedAt"] != nil {
	// // 	if t, ok := pro_setting["CreatedAt"].(time.Time); ok {
	// // 		createdAt = t
	// // 	}
	// // }
	// // pro_setting["CreatedAt"] = createdAt.Format("2006-01-02 15:04:05")
	
	// var pro_balance decimal.Decimal
	// var createdAt time.Time
	// createdAt = time.Now() 
	// if pro_setting["CreatedAt"] != nil {
	// 	createdAt = pro_setting["CreatedAt"].(time.Time) 
	// }
	// db.Debug().Model(&models.TransactionSub{}).Select("balance").Where("membername = ? AND deleted_at is null and created_at > ?",users.Username,createdAt.Format("2006-01-02 15:04:05")).Limit(1).Order("id desc").Find(&pro_balance)
	

	// updates := map[string]interface{}{
	// 	"ProBalance": pro_balance,
	// }
	// if err := db.Debug().Model(&users).Where("id=? and pro_status=?",users.ID,users.ProStatus).Updates(updates).Error; err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"status": false,
	// 			"message": err.Error(),
	// 		})
	// 	}

	
	//fmt.Printf("data: %+v\n", users)
	response := fiber.Map{
		"Status":  true,
		"Message": "สำเร็จ",
		"Data": fiber.Map{
			"id":         users.ID,
			"fullname":   users.Name,
			"banknumber": users.Banknumber,
			"bankname":   users.Bankname,
			"username":   strings.ToUpper(users.Username),
			"balance":    users.Balance,
			"prefix":     users.Prefix,
			"affiliatekey": users.AffiliateKey,
		}}
	return c.JSON(response)
}
func GetPartnerById(c *fiber.Ctx) error {
    body := new(PartnerBody)
    if err := c.BodyParser(body); err != nil {
        response := fiber.Map{
            "Message": "รับข้อมูลผิดพลาด",
            "Status":  false,
            "Data":    err.Error(),
        }
        return c.JSON(response)
    }

    db, err := database.ConnectToDB(body.Prefix)
    if err != nil {
        response := fiber.Map{
            "Message": "ติดต่อฐานข้อมูลผิดพลาด",
            "Status":  false,
            "Data":    err.Error(),
        }
        return c.JSON(response)
    }

    user := models.Partner{}
    err = db.Debug().First(&user, body.ID).Error
    if err != nil {
        response := fiber.Map{
            "Message": "ดึงข้อมูลผิดพลาด",
            "Status":  false,
            "Data":    err.Error(),
        }
        return c.JSON(response)
    }

    response := fiber.Map{
        "Message": "ดึงข้อมูลสำเร็จ",
        "Status":  true,
        "Data":    user,
    }
    return c.JSON(response)
}
func GetPartnerByUsername(c *fiber.Ctx) error {
	user := new(models.Partner)

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
	var users models.Partner
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
	
	//var promotion models.Promotion
	//fmt.Println(summary.Turnover)
	// if users.ProStatus != "" {
	// 	db.Debug().Model(&models.Promotion{}).Select("Includegames,Excludegames").Where("ID = ?",users.ProStatus).Scan(&promotion)
	// }




	response := fiber.Map{
		"Status":  true,
		"Message": "สำเร็จ",
		"Data": fiber.Map{
			"id":         users.ID,
			"fullname":   users.Name,
			"banknumber": users.Banknumber,
			"bankname":   users.Bankname,
			"username":   strings.ToUpper(users.Username),
			"balance":    users.Balance,
			"prefix":     users.Prefix,
			// "turnover":   summary.Turnover,
			// "minturnover": users.MinTurnover,
			// "lastdeposit": users.LastDeposit,
			// "lastproamount": users.LastProamount,
			// "lastwithdraw": users.LastWithdraw,
			// "pro_status": users.ProStatus,
			// "includegames": promotion.Includegames,
			// "excludegames": promotion.Excludegames,
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
// @Param user body models.Partner true "User balance info"
// @param Authorization header string true "Bearer token"
func GetBalance(c *fiber.Ctx) error {

	var users models.Partner

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
			"fullname":   users.Name,
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
// @Param user body models.Partner true "User Logout"
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
func GetPartnerTransaction(c *fiber.Ctx) error {

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
func GetPartnerStatement(c *fiber.Ctx) error {

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

	//var users models.Partner

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

func UpdatePartner(c *fiber.Ctx) error {
	// Parse the request body into a map

	var data RequestBody

	//var partner RequestBody


	//body = make(map[string]interface{})
	if err := c.BodyParser(&data); err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": err.Error(),
		}
		return c.JSON(response)
	}
	db, conn := database.ConnectToDB(data.Prefix)
	if conn != nil {
		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่่ พบข้อมูล Prefix!",
			Status:  false,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	fmt.Printf("data: %+v \n",data.Body)
 

	// // Get the username from the context
	// username := c.Locals("username").(string)
	
	// db, _err := handler.GetDBFromContext(c)
	// if _err != nil {
	// 	response := fiber.Map{
	// 		"Status":  false,
	// 		"Message": "โทเคนไม่ถูกต้อง!!",
	// 	}
	// 	return c.JSON(response)
	// }
	//db, _ := database.ConnectToDB(data.Prefix)
	
	var user models.Partner

	err := db.Where("username = ? ", data.Body.Username).First(&user).Error
	if err != nil {
		response := fiber.Map{
			"Status":  false,
			"Message": "ไม่พบรหัสผู้ใช้งาน!!",
		}
		return c.JSON(response)
	}
 

	// Update the user with the provided fields
	// fmt.Printf("Body: %s",data.Body)
	if err := db.Debug().Model(&user).Updates(data.Body).Error; err != nil {
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

func GetSeed(c *fiber.Ctx) error {
	var seedPhrase string
	//db, _ := handler.GetDBFromContext(c)
	user := new(models.Partner)

	if err := c.BodyParser(user); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//fmt.Printf(" %s ",user.Username)
	db, conn := database.ConnectToDB(user.Prefix)
	if conn != nil {
		response := ErrorResponse{
			Message: "เกิดข้อผิดพลาดไม่่ พบข้อมูล Prefix!",
			Status:  false,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	//db, _err := handler.GetDBFromContext(c)
	//prefix := c.Locals("Prefix")
	
 
	seedPhrase = CheckSeed(db)

	return c.JSON(fiber.Map{
		"Status":  true,
		"Message": "สร้าง affiliate key สำเร็จ!",
		"Data": fiber.Map{
			"affiliatekey": seedPhrase,
		},
	})
}

func CheckSeed(db *gorm.DB) string {

	var seedPhrase string
	for {
		seedPhrase, _ = encrypt.GenerateAffiliateCode(5) // สร้าง affiliate key ใหม่
		rowsAffected := db.Debug().Model(&models.Partner{}).Where("affiliatekey = ?", seedPhrase).RowsAffected
		if rowsAffected == 0 { // ถ้าไม่ซ้ำ
			break // ออกจากลูป
		}
	}
	return seedPhrase
}

type DatabaseInfo struct {
	Prefix string   `json:"prefix"`
	Names  []string `json:"names"`
}

func getDataList() ([]DatabaseInfo,error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host)
	fmt.Printf(" DSN: %s \n",dsn)
	// Connect to MySQL without a specific database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil,err
	}
	
	// Query to get all databases
	groupedDatabases := make(map[string][]string)

	rows, err := db.Raw("SHOW DATABASES").Rows()
	fmt.Printf("Rows: %v \n",rows)
	fmt.Printf("Err: %v \n",err)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	systemDatabases := map[string]bool{
		"information_schema": true,
		"mysql":              true,
		"performance_schema": true,
		"sys":                true,
	}

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			continue
		}
		fmt.Printf("DBName: %v",dbName)
		fmt.Printf("SysDBName: %v",systemDatabases[dbName])
		// Include databases with underscore in their names and exclude system databases
		//
		if strings.Contains(dbName, "_") && !systemDatabases[dbName] {
			parts := strings.SplitN(dbName, "_", 2)
			if len(parts) == 2 {
				prefix := parts[0]
				if _, exists := groupedDatabases[prefix]; !exists {
					groupedDatabases[prefix] = []string{}
				}
				groupedDatabases[prefix] = append(groupedDatabases[prefix], dbName)
			}
		}
	}
	var databaseList []DatabaseInfo
	for prefix, names := range groupedDatabases {
		databaseList = append(databaseList, DatabaseInfo{
			Prefix: prefix,
			Names:  names,
		})
	}
	return databaseList,nil

}