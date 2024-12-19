package common

import (
	//"context"
   // "fmt"
   // "github.com/amalfra/etag"
   // "github.com/go-redis/redis/v8"
   //"github.com/gofiber/fiber/v2"
   //"github.com/valyala/fasthttp"
   // "strconv"
   "github.com/shopspring/decimal"
   //"github.com/Knetic/govaluate"
   // "github.com/streadway/amqp"
   // "github.com/tdewolff/minify/v2"
   // "github.com/tdewolff/minify/v2/js"
   // "github.com/valyala/fasthttp"
    _ "github.com/go-sql-driver/mysql"
   "gorm.io/driver/mysql"
   "hanoi/models"
   "gorm.io/gorm"
   //"hanoi/database"
   //"hanoi/handler/jwtn"
   //jtoken "github.com/golang-jwt/jwt/v4"
   //"github.com/golang-jwt/jwt"
   //jtoken "github.com/golang-jwt/jwt/v4"
   //"github.com/solrac97gr/basic-jwt-auth/config"
   //"github.com/solrac97gr/basic-jwt-auth/models"
   //"github.com/solrac97gr/basic-jwt-auth/repository"
   //"hanoi/repository"
   "encoding/json"
   //"log"
   // "net"
   // "net/http"
   "os"
   //"strconv"
   //"time"
   "strings"
   "fmt"
   "errors"
   //"github.com/go-resty/resty/v2"
)

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

type PartnerMember struct {
	ID              int             // ID ของสมาชิก
	WalletID        int             // ID ของกระเป๋าเงิน
	Username        string          // ชื่อผู้ใช้
	Password        string          // รหัสผ่าน
	ProviderPassword string          // รหัสผ่านจากผู้ให้บริการ
	Fullname        string          // ชื่อจริง
	Bankname        string          // ชื่อธนาคาร
	Banknumber      string          // หมายเลขบัญชีธนาคาร
	Balance         decimal.Decimal  // ยอดเงินในกระเป๋าเงิน
	Beforebalance   decimal.Decimal  // ยอดเงินก่อนหน้า
	Currency        string          // สกุลเงิน
	Token           string          // โทเค็น
	Role            string          // บทบาท
	Salt            string          // salt
	Status          int             // สถานะ
	Betamount       decimal.Decimal  // จำนวนเงินที่เดิมพัน
	Commission      decimal.Decimal  // ค่าคอมมิชชั่น
	Win             decimal.Decimal  // จำนวนเงินที่ชนะ
	Lose            decimal.Decimal  // จำนวนเงินที่แพ้
	Turnover        decimal.Decimal  // ยอดเดิมพันรวม
	TotalEarnings   decimal.Decimal  // ยอดรวมกำไร
	TotalTurnover   decimal.Decimal  // ยอดรวมการหมุนเวียน
	ProID           string          // ID ของโปรเจ็กต์
	PartnersKey     string          // คีย์ของพันธมิตร
	ProStatus       string          // สถานะของโปร
	ProActive       string          // สถานะการทำงานของโปร
	Prefix          string          // ค่าพรีฟิก
}

// type ParnterMember struct {
// 	//gorm.Model               // ใช้สำหรับให้ GORM จัดการฟิลด์ ID, CreatedAt, UpdatedAt, DeletedAt
// 	ID      int       `gorm:"column:id"`
// 	WalletID      int       `gorm:"column:walletid"`
// 	Username      string    `gorm:"column:username"`
// 	Password      string    `gorm:"column:password"`
// 	ProviderPassword string   `gorm:"column:provider_password"`
// 	Fullname      string    `gorm:"column:fullname"`
// 	Bankname      string    `gorm:"column:bankname"`
// 	Banknumber    string    `gorm:"column:banknumber"`
// 	Balance       decimal.Decimal   `gorm:"column:balance"`      // แนะนำให้ใช้ decimal.Decimal สำหรับค่าเงิน
// 	Beforebalance decimal.Decimal   `gorm:"column:beforebalance"` // แนะนำให้ใช้ decimal.Decimal สำหรับค่าเงิน
// 	Currency      string    `gorm:"column:currency"`
// 	Token         string    `gorm:"column:token"`
// 	Role          string    `gorm:"column:role"`
// 	Salt          string    `gorm:"column:salt"`
// 	Status        int       `gorm:"column:status"`
// 	Betamount     decimal.Decimal   `gorm:"column:betamount"`    // แนะนำให้ใช้ decimal.Decimal สำหรับค่าเงิน
// 	Commission    decimal.Decimal   `gorm:"column:commission"`   // แนะนำให้ใช้ decimal.Decimal สำหรับค่าเงิน
// 	Win           decimal.Decimal   `gorm:"column:win"`          // แนะนำให้ใช้ decimal.Decimal สำหรับค่าเงิน
// 	Lose          decimal.Decimal   `gorm:"column:lose"`         // แนะนำให้ใช้ decimal.Decimal สำหรับค่าเงิน
// 	Turnover      decimal.Decimal   `gorm:"column:turnover"`     // แนะนำให้ใช้ decimal.Decimal สำหรับค่าเงิน
// 	TotalEarnings decimal.Decimal   `gorm:"column:totalEarnings"`
// 	TotalTurnover decimal.Decimal   `gorm:"column:totalTurnover"`
// 	ProID         string    `gorm:"column:pro_id"`
// 	PartnersKey   string    `gorm:"column:partners_key"`
// 	ProStatus     string    `gorm:"column:pro_status"`
// 	ProActive     string    `gorm:"column:pro_active"`
// 	Prefix        string    `gorm:"column:prefix"`
//  }



var mysql_host = os.Getenv("MYSQL_HOST")
var mysql_user = os.Getenv("MYSQL_USER")
var mysql_pass = os.Getenv("MYSQL_ROOT_PASSWORD")
// ฟังก์ชั่นช่วยตรวจสอบ turnover
func CheckTurnover(db *gorm.DB, users *models.Users, pro_setting map[string]interface{}) (decimal.Decimal,error) {

	var promotionLog models.PromotionLog
	db.Where("userid = ? AND promotioncode = ? AND status = 1", users.ID, users.ProStatus).
		Order("created_at DESC").
		First(&promotionLog)


	var totalTurnover decimal.Decimal
	err := db.Model(&models.TransactionSub{}).
		Where("proid = ? AND membername = ? AND created_at >= ?", 
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
func ConnectMaster() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host, "master")
	
	fmt.Printf(" dsn: %s \n",dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}
	return db
}
func GetCommissionRate(prefix string) (decimal.Decimal,error) {
	var settings []models.Settings
	db := ConnectMaster()
	if db == nil {
		return decimal.NewFromFloat(0.0), errors.New("เกิดข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล")
	}
	db.Debug().Model(&settings).Where("`key` like ?", prefix+"%").Find(&settings)

	var commissionRate decimal.Decimal
	for _, setting := range settings {
		//fmt.Printf(" Setting: %+v",commissionValue)
		if setting.Key == strings.ToLower(prefix)+"_partner_commission" { // ตรวจสอบจาก prefix + "user_commission"
		   //fmt.Printf("setting.Value: %v \n",setting.Value)
			//commissionRate, _ = decimal.NewFromString(setting.Value) // แปลงค่าจาก string เป็น decimal
			if strings.HasSuffix(setting.Value, "%") {
				// ลบ '%' ออกก่อนแปลงเป็น decimal
				valueWithoutPercent := strings.TrimSuffix(setting.Value, "%")
				commissionRate, _ = decimal.NewFromString(valueWithoutPercent) // แปลงค่าจาก string เป็น decimal
			} else {
				commissionRate, _ = decimal.NewFromString(setting.Value) // แปลงค่าจาก string เป็น decimal
			}
			break
		}
	}

	return commissionRate,nil
}


func CalculatePartnerCommission(commissionRate decimal.Decimal, userTurnover decimal.Decimal) decimal.Decimal {
	//totalEarnings := decimal.NewFromFloat(0.0)
	//var partner models.Partner
	//db.First(&partner, partnerID)

	// คำนวณค่าคอมมิชชั่น
	return userTurnover.Mul(commissionRate.Div(decimal.NewFromFloat(100)))

	// บันทึกข้อมูลในตาราง AffiliateTracking
	// tracking := models.AffiliateTracking{
	// 	PartnerID:  partner.ID,
	// 	Turnover:   userTurnover,
	// 	Commission: commission,
	// }
	// db.Create(&tracking)

	// // อัปเดต TotalEarnings ใน Partners
	// totalEarnings = partner.TotalEarnings.Add(commission)
	// //db.Save(&partner)

	  
}

func GetProdetail(db *gorm.DB, procode string) (map[string]interface{}, error) {
	var promotion models.Promotion
	if err := db.Debug().Where("id = ?", procode).Find(&promotion).Error; err != nil {
		fmt.Printf("Error unmarshalling JSON: %v", err)
		return nil, err
	}
	if promotion.SpecificTime != "" {
			if err := json.Unmarshal([]byte(promotion.SpecificTime), &ProItem.ProType); err != nil {
			//log.Fatalf("Error unmarshalling JSON: %v", err)
			fmt.Printf("Error unmarshalling JSON: %v", err)
			return nil, err
			}
		} else {
			return nil, nil
		}
	
		 
		response := make(map[string]interface{}) 
		//fmt.Printf(" %s ",promotion)
			response["Type"] = ProItem.ProType.Type
			response["count"] = ProItem.UsageLimit
			response["MinTurnover"] = promotion.MinSpend
			response["Formular"] = promotion.Example
		    response["Name"] = promotion.Name
			response["TurnType"]=promotion.TurnType
		if ProItem.ProType.Type == "weekly" {
			response["Week"] = ProItem.ProType.DaysOfWeek
		}
	 
		return response, nil
	
}