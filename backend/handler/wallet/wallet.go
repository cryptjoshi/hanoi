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
	//"strconv"
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
    db.Debug().Model(&promotionlog).Where("promotioncode = ? and (userid=? or walletid=?) and status=1", users.ProStatus,users.ID,users.ID).Order("id desc").Scan(&promotionlog).Count(&RowsAffected)
	fmt.Println("244 line CheckPro")
     fmt.Println(RowsAffected)
    // fmt.Println(ProItem.UsageLimit)


	// Check if promotionlog is not empty or has row affected = 1
	if int64(ProItem.UsageLimit) > 0 && RowsAffected >= int64(ProItem.UsageLimit) { // Assuming ID is the primary key
		return nil, nil
	}  
		
	fmt.Println("254 line CheckPro")
	fmt.Println(ProItem.ProType)

	

	switch ProItem.ProType.Type {
	case "first", "once","week":
		response["minDept"] = promotion.MinDept
		response["maxDept"] = promotion.MaxDiscount
		response["Widthdrawmax"] = promotion.MaxSpend
		response["MinTurnover"] = promotion.MinSpend
		response["count"] = ProItem.UsageLimit
		response["Formular"] = promotion.Example
		response["Name"] = promotion.Name
		response["MinSpendType"] = promotion.MinSpendType
		response["MinCredit"] = promotion.MinCredit
		response["TurnType"] = promotion.TurnType
		response["ZeroBalance"] = promotion.ZeroBalance
		response["CreatedAt"] = promotionlog.CreatedAt
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
			"status": false,
			"message": err_,
			"data": fiber.Map{ 
				"id": -1,
			}})
    }
	pro_setting, err := CheckPro(db, &users) 

	fmt.Printf("343 line Deposit\n")
	fmt.Printf("pro_setting: %v \n",pro_setting)

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
	deposit := BankStatement.Transactionamount
	//BankStatement.ProID = users.ProStatus
	//turnoverdef = strings.Replace(users.MinTurnoverDef, "%", "", 1) 
	var result decimal.Decimal
	var percentValue decimal.Decimal
	var percentStr = ""
	var zeroBalance bool
	//fmt.Printf("deposit: %v ",deposit)
	fmt.Printf("TurnType: %v ",pro_setting["TurnType"])
	fmt.Printf("MinCredit: %v ",pro_setting["MinCredit"])

	if pro_setting != nil {
		if pro_setting["ZeroBalance"] == 1 {
			zeroBalance = users.Balance.IsZero() && deposit.GreaterThan(decimal.Zero)
		} else {
			zeroBalance = users.Balance.LessThan(decimal.NewFromInt(1)) && deposit.GreaterThan(decimal.Zero)
		}
		fmt.Println("386 line")
		fmt.Println(zeroBalance)
		if zeroBalance == true  {
		

		// New code to log to promotionlog
		//fmt.Printf("Prosetting: %v ",pro_setting)

		// Ensure pro_setting["Example"] is not nil before type assertion

		//fmt.Printf("deposit > 0 : %v ",deposit.GreaterThan(decimal.Zero))

	 

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
			if balanceIncrease.GreaterThan(pro_setting["maxDept"].(decimal.Decimal)) {
				BankStatement.Balance = users.Balance.Add(deposit.Add(pro_setting["maxDept"].(decimal.Decimal)))
			} else {
				BankStatement.Balance = users.Balance.Add(balanceIncrease)
			}

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
				Status: 1,
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
	 
		 
	
			if users.Balance.IsZero() == false {
				
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
		} else if deposit.LessThan(decimal.Zero) {
			//fmt.Printf("MinTurnover: %v \n",pro_setting["MinTurnover"])
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
						"Message": fmt.Sprintf("ยอดถอนมากว่ายอดถอนสูงสุด %v %v !",pro_setting["Widthdrawmax"],users.Currency),
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
		} else {
			fmt.Printf("611 line %v \n",users.Balance.LessThanOrEqual(decimal.Zero))
			fmt.Printf("612 line %v \n",pro_setting["ZeroBalance"])
			if users.Balance.LessThanOrEqual(decimal.Zero) == false && pro_setting["ZeroBalance"] == 1 {
				
				response := fiber.Map{
					"Message": "ไม่สามารถ ฝากเงินเพิ่มได้ ขณะใช้งานโปรโมชั่น!",
					"Status":  false,
					"Data": fiber.Map{ 
						"id": -1,
					}}
					return c.JSON(response)
				} else {
					fmt.Printf("622 line \n")
					BankStatement.Balance = users.Balance.Add(deposit)
				}
		}
	 
		
	} else {
		BankStatement.Balance = users.Balance.Add(deposit)
	}	
	

	BankStatement.Bankname = users.Bankname
	BankStatement.Accountno = users.Banknumber
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
	
	// 	//db, _ = database.ConnectToDB(BankStatement.Prefix)
	 _err := repository.UpdateUserFields(db, BankStatement.Userid, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	if _err != nil {
		return c.Status(200).JSON(fiber.Map{
			"Status": false,
			"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
		}
    //  if BankStatement.Transactionamount.LessThan(decimal.Zero) {
	// 	updates["LastWithdraw"] = BankStatement.Transactionamount
	//  } else {
	// 	updates["LastDeposit"] = BankStatement.Transactionamount
	// 	updates["LastProamount"] = BankStatement.Proamount
	//  }
	//  _err = repository.UpdateUserFields(db, BankStatement.Userid, updates)
	// 	if _err != nil {
	// 		return c.Status(200).JSON(fiber.Map{
	// 			"Status": false,
	// 			"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
	// 			"Data": fiber.Map{ 
	// 				"id": -1,
	// 			}})
	// 	} else {
	// 		//fmt.Println("User fields updated successfully")
	// 	}

 
	// if err := checkActived(db,&users); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message":  "actived deposit ข้อมูลไม่ได้!",
	// 		"Data": fiber.Map{ 
	// 			"id": -1,
	// 		}})
	}
	 
	 return c.Status(200).JSON(fiber.Map{
		"Status": true,
		"Data": fiber.Map{ 
			"id": BankStatement.ID,
			"beforebalance":BankStatement.Beforebalance,
			"balance": BankStatement.Balance,
		},
	})
	
 

}
func Withdraw(c *fiber.Ctx) error {

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
	deposit := BankStatement.Transactionamount
	//BankStatement.ProID = users.ProStatus

 
	//turnoverdef = strings.Replace(users.MinTurnoverDef, "%", "", 1) 

	var result decimal.Decimal
	var percentValue decimal.Decimal
	var percentStr = ""
	//var zeroBalance bool
	//fmt.Printf("deposit: %v ",deposit)
	fmt.Printf("TurnType: %v ",pro_setting["TurnType"])
	fmt.Printf("MinCredit: %v ",pro_setting["MinCredit"])
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
						"Message": fmt.Sprintf("ยอดถอนมากว่ายอดถอนสูงสุด %v %v !",pro_setting["Widthdrawmax"],users.Currency),
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



 

