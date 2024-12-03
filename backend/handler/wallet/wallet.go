package wallet

import (
	 "context"
	// "fmt"
	// "github.com/amalfra/etag"
	// "github.com/go-redis/redis/v8"
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
	"strconv"
	"time"
	"strings"
	"fmt"
	"errors"
	"github.com/go-resty/resty/v2"
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

	fmt.Printf("ProItem: %+v\n",ProItem)

	response := make(map[string]interface{}) // Changed to use make for clarity

	//fmt.Println(ProItem.ProType)
	var promotionlog models.PromotionLog
	// Use a single case statement to handle the different types
	
	var RowsAffected int64
	//db.Debug().Model(&settings).Select("id").Scan(&settings).Count(&RowsAffected)
    db.Debug().Model(&promotionlog).Where("promotioncode = ? and (userid=? or walletid=?) and status=1", users.ProStatus,users.ID,users.ID).Order("id desc").Scan(&promotionlog).Count(&RowsAffected)
	fmt.Printf("244 line CheckPro\n")
     fmt.Printf("RowsAffected: %d\n",RowsAffected)
    // fmt.Println(ProItem.UsageLimit)

	fmt.Printf(" Promotion.UsageLimit: %+v\n",promotion.UsageLimit)
	// Check if promotionlog is not empty or has row affected = 1
	if int64(promotion.UsageLimit) > 0 && RowsAffected > int64(promotion.UsageLimit) { // Assuming ID is the primary key
		return nil, errors.New("คุณใช้งานโปรโมชั่นเกินจำนวนครั้งที่กำหนด")
	}  
		
	fmt.Printf("254 line CheckPro\n")
	fmt.Printf("ProItem.ProType: %v\n",ProItem.ProType)

	

	switch ProItem.ProType.Type {
	case "first", "once","weekly","daily","monthly":
		response["minDept"] = promotion.MinDept
		response["maxDept"] = promotion.MaxDiscount
		response["Widthdrawmax"] = promotion.MaxSpend
		response["Widthdrawmin"] = promotion.Widthdrawmin
		response["MinTurnover"] = promotion.MinSpend
		//response["count"] = promotion.UsageLimit
		response["Formular"] = promotion.Example
		response["Name"] = promotion.Name
		response["MinSpendType"] = promotion.MinSpendType
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

func checkActived(db *gorm.DB,user *models.Users) error {

	var bankstatement []models.BankStatement
	var RowsAffected int64
	 db.Debug().Model(&bankstatement).Where("userid=? or walletid=?",user.ID,user.ID).Order("id ASC").Scan(&bankstatement).Count(&RowsAffected)
	
	if RowsAffected >= 1 {
		updates := map[string]interface{}{
			"Actived": time.Now(),
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


func Deposit(c *fiber.Ctx) error {

	// user := c.Locals("user").(*jtoken.Token)
	// 	claims := user.Claims.(jtoken.MapClaims)
	var users models.Users
	db,_ := handler.GetDBFromContext(c)
	id := c.Locals("ID").(int)
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

	fmt.Printf("343 line Deposit\n")
	fmt.Printf("pro_setting: %v \n",pro_setting)




	if err != nil {
		fmt.Printf("err: %v \n",err)
		return c.JSON(fiber.Map{
			"Status": false,
			"Message":  err,
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
	//var result decimal.Decimal
	//var percentValue decimal.Decimal
	//var percentStr = ""
	var zeroBalance bool
	//fmt.Printf("deposit: %v ",deposit)
	// fmt.Printf("TurnType: %v ",pro_setting["TurnType"])
	// fmt.Printf("MinCredit: %v ",pro_setting["MinCredit"])

	fmt.Printf("pro_setting: %v \n",pro_setting["Zerobalance"])

	if pro_setting != nil {
		if pro_setting["Zerobalance"] == 1 {
			zeroBalance = users.Balance.IsZero() && deposit.GreaterThan(decimal.Zero)
		} else {
			zeroBalance =  users.Balance.LessThan(decimal.NewFromInt(1)) && deposit.GreaterThan(decimal.Zero)
		}
		fmt.Printf("386 line\n	")
		fmt.Printf("zeroBalance: %v \n",zeroBalance)
		fmt.Printf("pro_setting: %+v \n",pro_setting)
		if zeroBalance == true  {
		

		// New code to log to promotionlog
		//fmt.Printf("Prosetting: %v ",pro_setting)

		// Ensure pro_setting["Example"] is not nil before type assertion

		//fmt.Printf("deposit > 0 : %v ",deposit.GreaterThan(decimal.Zero))

	 
			

			Formular, ok := pro_setting["Formular"].(string)
			fmt.Printf("Formular: %v \n",Formular)
			if !ok {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Status": false,
					"Message": "Promotion example is not a valid string",
					"Data": fiber.Map{
						"id": -1,
					}})
			}
			//BankStatement.Proamount = BankStatement.Transactionamount
			// Calculate the new balance using the example from pro_setting

		

			
			minDept := pro_setting["minDept"].(decimal.Decimal)
			if deposit.LessThan(minDept) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Status": false,
					"Message": "ยอดเงินฝากน้อยกว่ายอดฝากขั้นต่ำของโปรโมชั่น",
					"Data": fiber.Map{
						"id": -1,
					}})
			}
			//fmt.Printf(" %v ", deposit)
			//fmt.Printf(" %v ", Formular)

			// Replace 'deposit' in the example with the actual value
			Formular = strings.Replace(Formular, "deposit", deposit.String(), 1) // Convert deposit to string if necessary
			fmt.Printf(" %v ",Formular)
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
			fmt.Printf("balanceIncrease: %v ",balanceIncrease)
			BankStatement.Proamount = balanceIncrease.Sub(deposit)



			// Update BankStatement.Balance

			fmt.Printf("wallet.go 447 line pro_setting[\"maxDept\"].(decimal.Decimal): %v \n",pro_setting["maxDept"].(decimal.Decimal))
			if balanceIncrease.Sub(deposit).GreaterThan(pro_setting["maxDept"].(decimal.Decimal)) {
				BankStatement.Balance = users.Balance.Add(deposit.Add(pro_setting["maxDept"].(decimal.Decimal)))
			} else {
				fmt.Printf("wallet.go 453 line balanceIncrease: %v \n",balanceIncrease)
				BankStatement.Balance = users.Balance.Add(balanceIncrease)
			}

			fmt.Printf("wallet.go 451 line BankStatement.Balance: %v \n",BankStatement.Balance)

			promotionLog := models.PromotionLog{

				UserID: BankStatement.Userid,
				WalletID: BankStatement.Userid,
				StatementID: BankStatement.ID,
				Promotionname: pro_setting["Name"].(string),
				Beforebalance: BankStatement.Beforebalance,
				//BetAmount: BankStatement.BetAmount,
				Promotioncode: users.ProStatus,
				Transactionamount: deposit,
				Promoamount: balanceIncrease,
				Proamount: balanceIncrease.Sub(deposit),
				Balance: BankStatement.Balance,
				Example: Formular,
				Status: 1,
				// Add other necessary fields for the promotion log
			}

			
			if BankStatement.Balance.LessThan(balanceIncrease) {
				promotionLog.Proamount = BankStatement.Balance.Sub(users.Balance).Sub(deposit)
			}
			if err := db.Create(&promotionLog).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Status": false,
					"Message": "Failed to log promotion",
					"Data": fiber.Map{
						"id": -1,
					}})
			}
	 
		 
	
			if users.Balance.LessThan(decimal.NewFromInt(1)) == false && pro_setting["Zerobalance"] == 0  {
				
				response := fiber.Map{
					"Message": "ไม่สามารถ ฝากเงินเพิ่มได้ ยังมียอดเงินฝากคงเหลือ!",
					"Status":  false,
					"Data": fiber.Map{ 
						"id": -1,
					}}
					return c.JSON(response)
				}
			
			 
		} else if deposit.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Status": false,
			"Message": "ยอดเงินฝากไม่ถูกต้อง!",
			"Data": fiber.Map{
				"id": -1,
			}})
		} else
		// else if deposit.LessThan(decimal.Zero) {
		// 	//fmt.Printf("MinTurnover: %v \n",pro_setting["MinTurnover"])
		// 	if pro_setting["TurnType"] == "turnover" {
		// 			if strings.Contains(pro_setting["MinTurnover"].(string), "%") {
		// 				percentStr = strings.TrimSuffix(pro_setting["MinTurnover"].(string), "%")
					
		// 				//fmt.Printf(" MinturnoverDef : %s %",percentStr)
		// 				// แปลงเป็น float64
		// 				percentValue, _ = decimal.NewFromString(percentStr)
						
		// 				// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
		// 				percentValue = percentValue.Div(decimal.NewFromInt(100))
						
		// 			} else {
		// 				percentStr = pro_setting["MinTurnover"].(string)
		// 				// fmt.Printf(" 492 MinturnoverDef : %s % \n",percentStr)
		// 				percentValue, _ = decimal.NewFromString(percentStr)
		// 			}
					
		// 			//fmt.Printf(" Minturnover : %s ",percentStr)
		// 			// แปลงเป็น float64
		// 			//percentValue, _ := decimal.NewFromString(percentStr)
				
				
		// 			// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
					

		// 			fmt.Printf(" PercentValue: %v \n",percentValue)
		// 			fmt.Printf(" Users.LastDeposit: %v \n",users.LastDeposit)
		// 			fmt.Printf(" Users.LastProamount: %v \n",users.LastProamount)
		// 			if pro_setting["MinSpendType"] == "deposit" {
		// 				result = users.LastDeposit.Mul(percentValue)	
		// 			} else {
		// 				result = users.LastDeposit.Add(users.LastProamount).Mul(percentValue)	
		// 			}

		// 			fmt.Printf("pro_setting: %v \n",pro_setting)
		// 			//minTurnover := users.MinTurnover
		// 			fmt.Printf("bankstatement.Turnover: %v \n",BankStatement.Turnover)
		// 			fmt.Printf("result: %v \n",result)
		// 			//fmt.Printf("minTurnover: %v ",minTurnover)
		// 		if (BankStatement.Turnover.GreaterThan(result) || BankStatement.Turnover.Equal(result)) && BankStatement.Turnover.GreaterThan(decimal.Zero) {
		// 			if deposit.Abs().LessThanOrEqual(pro_setting["Widthdrawmax"].(decimal.Decimal)) {
		// 			BankStatement.Balance = users.Balance.Add(deposit)
		// 			} else {
		// 				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 					"Status": false,
		// 					"Message": fmt.Sprintf("ยอดเงินถอนสูงกว่ายอดถอนสูงสุดของโปรโมชั่น %v %v!",pro_setting["Widthdrawmax"],users.Currency),
		// 					"Data": fiber.Map{
		// 						"id": -1,
		// 					}})
		// 			}
		// 			// clear turnover
					
		// 		} else {
		// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 				"Status": false,
		// 				"Message": fmt.Sprintf("ยอดเทิร์นโอเวอร์น้อยกว่ายอดเทิร์นโอเวอร์ขั้นต่ำ %v %v !",result,users.Currency),
		// 				"Data": fiber.Map{
		// 					"id": -1,
		// 				}})
		// 		}
		// 	} else {
		// 	fmt.Printf("559 line \n")

			
			

		// 	// var bankstatement models.BankStatement
		// 	// db.Debug().Model(&models.BankStatement{}).Select("balance").Where("walletid = ? ",id).Limit(1).Order("id desc").Find(&bankstatement)
		// 	//var transaction models.TransactionSub
		// 	var tbalance decimal.Decimal
		// 	db.Debug().Model(&models.TransactionSub{}).Select("balance").Where("membername = ? AND deleted_at is null and created_at > ?",users.Username,pro_setting["CreatedAt"].(time.Time).Format("2006-01-02 15:04:05")).Limit(1).Order("id desc").Find(&tbalance)
		// 	if strings.Contains(pro_setting["MinCredit"].(string), "%") {
		// 		percentStr = strings.TrimSuffix(pro_setting["MinCredit"].(string), "%")
			
		// 		//fmt.Printf(" MinturnoverDef : %s %",percentStr)
		// 		// แปลงเป็น float64
		// 		percentValue, _ = decimal.NewFromString(percentStr)
				
		// 		// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
		// 		percentValue = percentValue.Add(decimal.NewFromInt(100)).Div(decimal.NewFromInt(100)).Mul(tbalance)
				
		// 	} else {
		// 		percentStr = pro_setting["MinCredit"].(string)
		// 		// fmt.Printf(" 492 MinturnoverDef : %s % \n",percentStr)
		// 		percentValue, _ = decimal.NewFromString(percentStr)
		// 	}
			
		// 	fmt.Printf("tbalance: %v \n",tbalance)
		// 	fmt.Printf("percentValue: %v \n",percentValue)
		// 	fmt.Printf("users.Balance: %v \n",users.Balance)
		// 	if tbalance.LessThanOrEqual(percentValue) == true && users.Balance.GreaterThan(decimal.Zero) {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"Status": false,
		// 		"Message": fmt.Sprintf("ยอดเครดิต %v น้อยกว่ายอดเครดิตขั้นต่ำ %v %v !",tbalance,pro_setting["MinCredit"],users.Currency),
		// 		"Data": fiber.Map{
		// 			"id": -1,
		// 		}})
		// 	} else {
		// 		if deposit.Abs().GreaterThan(pro_setting["Widthdrawmax"].(decimal.Decimal)) {
		// 			//BankStatement.Balance = users.Balance.Sub(pro_setting["Widthdrawmax"].(decimal.Decimal))
		// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 				"Status": false,
		// 				"Message": fmt.Sprintf("ยอดถอนมากว่ายอดถอนสูงสุด %v %v !",pro_setting["Widthdrawmax"],users.Currency),
		// 				"Data": fiber.Map{
		// 					"id": -1,
		// 				}})
		// 		} else {
		// 			BankStatement.Balance = users.Balance.Add(deposit)
		// 			fmt.Printf("607 line check user balance more than zero\n")
		// 			if users.Balance.GreaterThan(deposit.Abs()) {
		// 				users.Balance = decimal.Zero
		// 				users.ProStatus = ""
		// 				BankStatement.Balance = decimal.Zero
						
		// 			}
		// 		}
		// 	}

		// 	}
		// } 
		  {
			fmt.Printf("611 line %v \n",users.Balance.LessThanOrEqual(decimal.Zero))
			fmt.Printf("612 line %v \n",pro_setting["Zerobalance"])
			if users.Balance.LessThanOrEqual(decimal.Zero) == false && pro_setting["Zerobalance"] == 1 {
				
				response := fiber.Map{
					"Message": "ไม่สามารถ ฝากเงินเพิ่มได้ ขณะใช้งานโปรโมชั่น!",
					"Status":  false,
					"Data": fiber.Map{ 
						"id": -1,
					}}
					return c.JSON(response)
				} else if users.Balance.GreaterThan(decimal.NewFromFloat(0.1)) == true && pro_setting["Zerobalance"] == 0 {
				
					response := fiber.Map{
						"Message": "ไม่สามารถ ฝากเงินเพิ่มได้ ขณะใช้งานโปรโมชั่น!",
						"Status":  false,
						"Data": fiber.Map{ 
							"id": -1,
						}}
						return c.JSON(response)
				
				 } else {
					fmt.Printf(" wallet.go 637 line \n")
					BankStatement.Balance = users.Balance.Add(deposit)
				}
		}
	 
		
	} else {
		BankStatement.Balance = users.Balance.Add(deposit)
	}	
	

	BankStatement.Bankname = users.Bankname
	BankStatement.Accountno = users.Banknumber
	if pro_setting != nil {
		BankStatement.Proamount = BankStatement.Balance.Sub(BankStatement.Transactionamount) 
	} else {
		BankStatement.Proamount = decimal.NewFromFloat(0.0)
	}
	//user.Username = user.Prefix + user.Username
	fmt.Printf("692 line \n")
	fmt.Println(BankStatement.Balance)
	fmt.Println(deposit)

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
	BankStatement.StatementType = "Deposit"
	// event statement log

	resultz := db.Debug().Create(&BankStatement); 
	

	if resultz.Error != nil {
	 
			response := fiber.Map{
				"Message": "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
				"Status":  false,
				"Data":    "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			}
			return c.JSON(response)

	 } else {

// 		updates := map[string]interface{}{
// 			"Balance": BankStatement.Balance,
// 			"Turnover": users.Turnover,
// 			"ProStatus": users.ProStatus,
// 			}
	
 
// 	 _err := repository.UpdateUserFields(db, BankStatement.Userid, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
// 	if _err != nil {
// 		return c.Status(200).JSON(fiber.Map{
// 			"Status": false,
// 			"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
// 			"Data": fiber.Map{ 
// 				"id": -1,
// 			}})
// 		}

//     //  if BankStatement.Transactionamount.LessThan(decimal.Zero) {
// 	// 	updates["LastWithdraw"] = BankStatement.Transactionamount
// 	//  } else {
// 		updates["LastDeposit"] = BankStatement.Transactionamount
// 		updates["LastProamount"] = BankStatement.Balance.Sub(BankStatement.Transactionamount)
// 	 //}
// 	 _err = repository.UpdateUserFields(db, BankStatement.Userid, updates)
// 		if _err != nil {
// 			return c.Status(200).JSON(fiber.Map{
// 				"Status": false,
// 				"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
// 				"Data": fiber.Map{ 
// 					"id": -1,
// 				}})
// 		} else {
// 			//fmt.Println("User fields updated successfully")
// 		}

 
// 	if err := checkActived(db,&users); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"Status": false,
// 			"Message":  "actived deposit ข้อมูลไม่ได้!",
// 			"Data": fiber.Map{ 
// 				"id": -1,
// 			}})
// 	}
	 
	 return c.Status(200).JSON(fiber.Map{
		"Status": true,
		"Data": fiber.Map{ 
			"id": BankStatement.Uid,
			"beforebalance":BankStatement.Beforebalance,
			"balance": BankStatement.Balance,
		},
	})
	
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

    //ตรวจสอบข้อมูลผู้ใช้
    if err := db.Where("walletid = ? or id = ?", id,id).First(&users).Error; err != nil {
        return c.JSON(fiber.Map{
            "Status": false,
            "Message": "ไม่พบข้อมูลผู้ใช้",
            "Data": fiber.Map{"id": -1},
        })
    }
	
    withdraw := BankStatement.Transactionamount

	withdrawAbs := withdraw.Abs()

    // เพิ่มการตรวจสอบยอดถอนกับยอดคงเหลือ
    if withdrawAbs.GreaterThan(users.Balance) {
        return c.JSON(fiber.Map{
            "Status": false,
            "Message": fmt.Sprintf("ยอดถอนมากกว่ายอดคงเหลือในบัญชี (%v %v)", users.Balance, users.Currency),
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
            withdraw = pro_setting["Widthdrawmax"].(decimal.Decimal).Neg()
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
            totalTurnover,_ := checkTurnover(db, &users, pro_setting);  
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
		 if users.Turnover.LessThan(minTurnover) {
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
	
    // ถ้ามีโปรโมชั่นให้ปรับเป็น 0
    BankStatement.Bankname = users.Bankname
    BankStatement.Accountno = users.Banknumber
    BankStatement.Transactionamount = withdraw
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

    if users.ProStatus != "" {
        updates["ProStatus"] = ""
    }

    if err := tx.Model(&users).Updates(updates).Error; err != nil {
        tx.Rollback()
        return c.JSON(fiber.Map{
            "Status": false,
            "Message": "ไม่สามารถอัพเดทข้อมูลผู้ใช้ได้",
            "Data": fiber.Map{"id": -1},
        })
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


func Webhook(c *fiber.Ctx) error {
		
  		var requestBody CallbackRequest

		
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(200).SendString(err.Error())
		}

  		//db,_ := handler.GetDBFromContext(c)
  		db,_ := database.ConnectToDB(requestBody.MerchantID)

  		var bankstatement models.BankStatement

    	if err_ := db.Debug().Where("Uid = ?",requestBody.TransactionID).First(&bankstatement).Error; err_ != nil {
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": err_,
				"Data": fiber.Map{ 
					"id": -1,
				}})
		}

// 		bankstatement.Uid = requestBody.TransactionID

		if requestBody.Verify == 1 && requestBody.IsExpired == 0 {
				bankstatement.Status = "verified"
		}  else if requestBody.Verify == 0 && requestBody.IsExpired == 1 {
				bankstatement.Status = "expired"
				//bankstatement.Beforebalance = bankstatement.Balance.Sub(bankstatement.Transactionamount)
				bankstatement.Transactionamount = decimal.NewFromFloat(0.0);
		} 

		if  requestBody.Type == "payin" {
			 		updates := map[string]interface{}{
						"Balance": bankstatement.Balance,
						//"Turnover": users.Turnover,
						//"ProStatus": users.ProStatus,
						}
				
					
					 _err := repository.UpdateUserFields(db, bankstatement.Userid, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
					if _err != nil {
					return c.Status(200).JSON(fiber.Map{
						"Status": false,
						"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
						"Data": fiber.Map{ 
							"id": -1,
						}})
					}
					if err := db.Save(&bankstatement).Error; err != nil { // ใช้ db.Save เพื่ออัปเดต bankstatement
						return c.JSON(fiber.Map{
							"Status": false,
							"Message": "เกิดข้อผิดพลาดในการอัปเดตข้อมูล",
							"Data": fiber.Map{ 
								"id": -1,
							}})
					}
				 	//bankstatement.Beforebalance = bankstatement.Balance.Sub(bankstatement.Transactionamount)
				
		}


		

//     //  if BankStatement.Transactionamount.LessThan(decimal.Zero) {
// 	// 	updates["LastWithdraw"] = BankStatement.Transactionamount
// 	//  } else {
// 		updates["LastDeposit"] = BankStatement.Transactionamount
// 		updates["LastProamount"] = BankStatement.Balance.Sub(BankStatement.Transactionamount)
// 	 //}
// 	 	_err = repository.UpdateUserFields(db, BankStatement.Userid, updates)
// 		if _err != nil {
// 			return c.Status(200).JSON(fiber.Map{
// 				"Status": false,
// 				"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
// 				"Data": fiber.Map{ 
// 					"id": -1,
// 				}})
// 		} 

 
// 		if err := checkActived(db,&users); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"Status": false,
// 				"Message":  "actived deposit ข้อมูลไม่ได้!",
// 				"Data": fiber.Map{ 
// 					"id": -1,
// 				}})
// 		}

		return c.JSON(fiber.Map{
			"Status": true,
			"Message": "ถอนเงินสำเร็จ",
			"Data": fiber.Map{
				"id": bankstatement.Uid,
				"beforebalance": bankstatement.Beforebalance,
				"balance": bankstatement.Balance,
			},
		})
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

// เพิ่มฟังก์ชั่นช่วยคำนวณยอดเทิร์นที่ต้องการ
func CalculateRequiredTurnover(minTurnover string, lastDeposit decimal.Decimal) (decimal.Decimal, error) {
    if strings.Contains(minTurnover, "%") {
        percentStr := strings.TrimSuffix(minTurnover, "%")
        percentValue, err := decimal.NewFromString(percentStr)
        if err != nil {
            return decimal.Zero, err
        }
        return lastDeposit.Mul(percentValue.Div(decimal.NewFromInt(100))), nil
    }
    return decimal.NewFromString(minTurnover)
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
func checkTurnover(db *gorm.DB, users *models.Users, pro_setting map[string]interface{}) (decimal.Decimal,error) {

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

// ฟังก์ชั่นช่วยตรวจสอบ turncredit
func checkTurnCredit(db *gorm.DB, users *models.Users, pro_setting map[string]interface{}) error {
    var lastCredit decimal.Decimal
    if err := db.Model(&models.TransactionSub{}).
        Where("membername = ? AND deleted_at is null", users.Username).
        Order("id desc").
        Limit(1).
        Select("balance").
        Scan(&lastCredit).Error; err != nil {
        return errors.New("ไม่สามารถตรวจสอบยอดเครดิต")
    }

    minCreditStr, ok := pro_setting["MinCredit"].(string)
    if !ok {
        return errors.New("รูปแบบยอดเครดิตขั้นต่ำไม่ถูกต้อง")
    }

    minCredit, err := decimal.NewFromString(minCreditStr)
    if err != nil {
        return errors.New("ไม่สามารถแปลงค่ายอดเครดิตขั้นต่ำได้")
    }

    if lastCredit.LessThan(minCredit) {
        return fmt.Errorf("ยอดเครดิต %v น้อยกว่ายอดเครดิตขั้นต่ำ %v %v", lastCredit, minCredit, users.Currency)
    }

    return nil
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
					"Message": fmt.Sprintf("ยอดเทิร์นโอเวอร์น้อยกว่ายอดเทิร์นโอเวอร์ขั้นต่ำ %v %v ของยอดฝากล่าสุด !",users.MinTurnoverDef,users.Currency),
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
func payout(amount string) (error) {

	return nil
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
 


 