package ef

import 
(
	"github.com/gofiber/fiber/v2"
	"pkd/models"
	"pkd/database"
	"pkd/repository"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"crypto/md5"
	"encoding/hex"
	//"os"
	"encoding/json"
	"time"
	"log"
	"pkd/common"
	"pkd/handler"
	//"strconv"
	//"repository"
	"strings"
	"fmt"
	
)


type Balance struct {
    BetAmount decimal.Decimal
}

type User struct {
    Balance decimal.Decimal
}

type Response struct {
    Message string      `json:"message"`
    Status  bool        `json:"status"`
    Data    interface{} `json:"data"` // ใช้ interface{} เพื่อรองรับข้อมูลหลายประเภทใน field data
}

type EFResponse struct {
	ErrorCode  int  `json:"errorcode"`
    ErrorMessage string `json:"errormessage"`
    Balance  decimal.Decimal `json:"balance"`
    BeforeBalance decimal.Decimal `json:"beforebalance"`
}

type EFBody struct {
	MemberName string `json:"membername"`
	OperatorCode string `json:"operatorcode"`
	ProductID int `json:"productid"`
	MessageID string `json:"messageid"`
	Sign string `json:"sign"`
	RequestTime  string `json:"requesttime"`
	
} 

type EFBodyTransaction struct {
	MemberName string `json:"membername"`
	OperatorCode string `json:"operatorcode"`
	ProductID string `json:"productid"`
	MessageID string `json:"messageid"`
	Sign string `json:"sign"`
	RequestTime  string `json:"requesttime"`
	Transactions []models.TransactionSub `json:"transactions"`
	
} 
type EFTransaction struct {
	Transactions []models.TransactionSub `json:"transactions"`
	
}

type ResponseBalance struct {
	BetAmount decimal.Decimal `json:"betamount"`
	BeforeBalance decimal.Decimal `json:"beforebalance"`
	Balance decimal.Decimal `json:"balance"`
}
// ฟังก์ชันตัวอย่างใน efinity.go
//const EF_SECRET_KEY="456Ayb" //product
//var EF_SECRET_KEY="1g1bb3" //stagging
//var OPERATOR_CODE = os.Getenv("EF_OPERATOR")

//var EF_API_URL = os.Getenv("INFINITY_STAG_URL") //"https://swmd.6633663.com/"
//var OPERATOR_CODE = os.Getenv("INFINITY_OPERATOR_CODE") || "E293"
//var INFINITY_PROD_URL  = os.Getenv("INFINITY_PROD_URL") // "https://prod_md.9977997.com"
//var INFINITY_STAG_URL = "https://swmd.6633663.com"  //os.Getenv("INFINITY_STAG_URL") // "https://stag_md.9977997.com"
//var DEVELOPMENT_SECRET_KEY = os.Getenv("DEVELOPMENT_SECRET_KEY")    || "1g1bb3"  //staging
//var PRODUCTION_SECRET_KEY= os.Getenv("PRODUCTION_SECRET_KEY")  || "456Ayb" //product

//var USER_FIX = os.Getenv("USER_FIX") 
//var PASS_FIX = os.Getenv("PASS_FIX") 

func parseTime(layout, value string) (time.Time, error) {
    return time.Parse(layout, value)
}

func Index(c *fiber.Ctx) error {

	var user []models.Users
	// db, _ := database.GetDatabaseConnection(user[0].Username)
	// db.Find(&user)
	response := Response{
		Message: "Welcome to Efinity!!",
		Status:  true,
		Data: fiber.Map{ 
			"users":user,
		}, 
	}
	 
	return c.JSON(response)
   
}

func CheckSign(Signature string,methodName string,requestTime string) bool {

	
	//requestTime := "2024-09-15T12:00:00Z"
	//methodName := "MethodName"
	
	secretKey := common.EF_SECRET_KEY

	// สร้างข้อมูลที่ต้องใช้ hash
	data := common.EF_OPERATOR_CODE + requestTime + strings.ToLower(methodName) + secretKey

	 
	// สร้าง MD5 hash
	hash := md5.New()
	hash.Write([]byte(data))

	// เปลี่ยน hash เป็นรูปแบบ hexadecimal string
	hashInHex := hex.EncodeToString(hash.Sum(nil))

	// fmt.Println("data:",data)
	// fmt.Println("RequestTime",requestTime)
	// fmt.Println("Operator Hash:", OPERATOR_CODE)
	// fmt.Println("SecretKey Hash:", secretKey)
	// fmt.Println("Sign Hash:", Signature)
	// fmt.Println("MD5 Hash:", hashInHex)
	
	
	return Signature == hashInHex	
}

func GetBalance(c *fiber.Ctx) error {
	


	
	body := new(EFBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	db, err := database.GetDatabaseConnection(body.MemberName)
	if err != nil {
		response := EFResponse {
			ErrorCode: 16,
			ErrorMessage: "Faild",
			Balance: decimal.NewFromFloat(0),
			BeforeBalance: decimal.NewFromFloat(0),
		}
		return c.JSON(response)
	}
	
	//}
	if CheckSign(body.Sign,"getbalance",body.RequestTime) == true {
			var users models.Users
		 
			 
			if err := db.Where("username = ?", body.MemberName).First(&users).Error; err != nil {
				
				response := EFResponse {
					ErrorCode: 16,
					ErrorMessage: "Faild",
					Balance: decimal.NewFromFloat(0),
					BeforeBalance: decimal.NewFromFloat(0),
				}
				return c.JSON(response)
				// return  errors.New("user not found")
			} else {
					response := EFResponse{
						ErrorCode:0,
						ErrorMessage:"Success",
						Balance: users.Balance,
						BeforeBalance: decimal.NewFromFloat(0),
					}
				return c.JSON(response)
			}
	} else {
		response := EFResponse{
			ErrorCode:1004,
			ErrorMessage:"API Invalid Sign",
			Balance: decimal.NewFromFloat(0),
			BeforeBalance: decimal.NewFromFloat(0),
		}
		return c.JSON(response)
	 
	}
}
func AddBuyOut(transactionsub models.BuyInOut,membername string) Response {


	response := Response{
		Status: false,
		Message:"Success",
		Data: ResponseBalance{
			BetAmount: decimal.NewFromFloat(0),
			BeforeBalance: decimal.NewFromFloat(0),
			Balance: decimal.NewFromFloat(0),
		},
	}
	 
	var users models.Users
	db, _ := database.GetDatabaseConnection(membername)
    if err_ := db.Where("username = ? ", membername).First(&users).Error; err_ != nil {
		response = Response{
			Status: false,
			Message: "ไม่พบข้อมูล",
			Data:map[string]interface{}{
				"id": -1,
			},
    	}
	}

    transactionsub.GameProvide = "EFINITY"
    transactionsub.MemberName = membername
	transactionsub.ProductID = transactionsub.ProductID
	transactionsub.BetAmount = transactionsub.BetAmount
	transactionsub.BeforeBalance = users.Balance
	transactionsub.Balance = users.Balance.Add(transactionsub.TransactionAmount)
	transactionsub.ProID = users.ProStatus
	
	result := db.Create(&transactionsub); 
	//fmt.Println(result)
	if result.Error != nil {
		response = Response{
			Status: false,
			Message:  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Data: map[string]interface{}{ 
				"id": -1,
			}}
	} else {

		updates := map[string]interface{}{
			"Balance": transactionsub.Balance,
				}
	
		 
		_err :=  repository.UpdateFieldsUserString(membername, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		//fmt.Println(_err)
		if _err != nil {
			fmt.Println("Error:", _err)
		} else {
			fmt.Println("User fields updated successfully")
		}

 
 
	 
	  response = Response{
		Status: true,
		Message: "สำเร็จ",
		Data: ResponseBalance{
			BeforeBalance: transactionsub.BeforeBalance,
			Balance:       transactionsub.Balance,
		},
		}
	}
	return response

}
func AddBuyInOut(transaction models.BuyInOut,membername string) Response {


	response := Response{
		Status: false,
		Message:"Success",
		Data: ResponseBalance{
			BetAmount: decimal.NewFromFloat(0),
			BeforeBalance: decimal.NewFromFloat(0),
			Balance: decimal.NewFromFloat(0),
		},
	}
	 
	var users models.Users
	db, _ := database.GetDatabaseConnection(membername)
    if err_ :=db.Where("username = ? ", membername).First(&users).Error; err_ != nil {
		response = Response{
			Status: false,
			Message: "ไม่พบข้อมูล",
			Data:map[string]interface{}{
				"id": -1,
			},
    	}
	}
	// fmt.Println("----------------")
	// fmt.Println(transaction.TransactionAmount)
	// fmt.Println("----------------")
    transaction.GameProvide = "EFINITY"
    transaction.MemberName = membername
	transaction.ProductID = transaction.ProductID
	//transactionsub.BetAmount = transactionsub.BetAmount
	transaction.BeforeBalance = users.Balance
	transaction.Balance = users.Balance.Add(transaction.TransactionAmount)
	transaction.ProID = users.ProStatus
	result := db.Create(&transaction); 
	
	
	fmt.Print(result)
	
	if result.Error != nil {
		response = Response{
			Status: false,
			Message:  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Data: map[string]interface{}{ 
				"id": -1,
			}}
	} else {

		updates := map[string]interface{}{
			"Balance": transaction.Balance,
				}
	
		 
		_err :=  repository.UpdateFieldsUserString(membername, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		// fmt.Println(_err)
		if _err != nil {
			fmt.Println("Error:", _err)
		} else {
			//fmt.Println("User fields updated successfully")
		}

 
 
	 
	  response = Response{
		Status: true,
		Message: "สำเร็จ",
		Data: ResponseBalance{
			BeforeBalance: transaction.BeforeBalance,
			Balance:       transaction.Balance,
		},
		}
	}
	return response

}
func AddTransactions(transactionsub models.TransactionSub,membername string) Response {

	//fmt.Printf("Add transactionsub: %+v \n",transactionsub)
	response := Response{
		Status: false,
		Message:"Success",
		Data: ResponseBalance{
			BetAmount: decimal.NewFromFloat(0),
			BeforeBalance: decimal.NewFromFloat(0),
			Balance: decimal.NewFromFloat(0),
		},
	}
	 
	var users models.Users
	db,_ := database.GetDatabaseConnection(membername)
    if err_ := db.Where("username = ? ", membername).First(&users).Error; err_ != nil {
		response = Response{
			Status: false,
			Message: "ไม่พบข้อมูล",
			Data:map[string]interface{}{
				"id": -1,
			},
    	}
	}
	//fmt.Printf("ProID : %v \n",users)
    transactionsub.GameProvide = "EFINITY"
    transactionsub.MemberName = membername
	transactionsub.ProductID = transactionsub.ProductID
	transactionsub.BetAmount = transactionsub.BetAmount
	// if transactionsub.TurnOver.IsZero() && transactionsub.Status == 100 {
	// 	transactionsub.TurnOver = transactionsub.BetAmount
	// } else if transactionsub.GameType == 1 && transactionsub.TurnOver.IsZero() && transactionsub.Status == 101 {
	// 	transactionsub.TurnOver = transactionsub.BetAmount
	// }
	transactionsub.BeforeBalance = users.Balance
	transactionsub.Balance = users.Balance.Add(transactionsub.TransactionAmount)
	transactionsub.ProID = users.ProStatus
	result := db.Create(&transactionsub); 
	//fmt.Println(result)
	if result.Error != nil {
		response = Response{
			Status: false,
			Message:  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Data: map[string]interface{}{ 
				"id": -1,
			}}
	} else {

		updates := map[string]interface{}{
			"Balance": transactionsub.Balance,
		}
		
		// one,_ := decimal.NewFromString("1")
		// if transactionsub.Balance.LessThan(one) {
		// 	updates["pro_status"] = ""
		// }
		 
		_err :=  repository.UpdateFieldsUserString(membername, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
 
		if _err != nil {
			fmt.Println("Error:", _err)
		} else {
			fmt.Println("User fields updated successfully")
		}
		
 
 
	 
	  response = Response{
		Status: true,
		Message: "สำเร็จ",
		Data: ResponseBalance{
			BeforeBalance: transactionsub.BeforeBalance,
			Balance:       transactionsub.Balance,
		},
		}
	}
	return response

}
func PlaceBet(c *fiber.Ctx) error {

	request := new(models.TransactionsRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"placebet",request.RequestTime) == true { 

		
		var user models.Users
		
		db,_ := database.GetDatabaseConnection(request.MemberName)
		 for _, transaction := range request.Transactions {
			transaction.IsAction = "PlaceBet"
			db.Where("username = ?", request.MemberName).First(&user)
	  		
			if user.Balance.LessThan(transaction.BetAmount) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			{
				var c_transaction_found models.TransactionSub
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
				if rowsAffected == 0 {
							if transaction.TurnOver.IsZero() {
								transaction.TurnOver = transaction.BetAmount
							}
							result := AddTransactions(transaction,request.MemberName)
							responseBalance, _ := result.Data.(ResponseBalance)
						
							response = EFResponse{
								ErrorCode:    0,
								ErrorMessage: "สำเร็จ",
								Balance:      responseBalance.Balance,
								BeforeBalance: responseBalance.BeforeBalance,
						}
					} else {
						response = EFResponse{
							ErrorCode:16,
							ErrorMessage:"รายการซ้ำ" ,
							Balance:  decimal.NewFromFloat(0),
							BeforeBalance: decimal.NewFromFloat(0),
						}
					}
			}
		 }
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func GameResult(c *fiber.Ctx) error {
	
	// body := new(EFBodyTransaction)
	// if err := c.BodyParser(body); err != nil {
	// 	return c.Status(200).SendString(err.Error())
	// }

	response := EFResponse{
		ErrorCode:0,
		ErrorMessage:"ไม่พบรายการ",
		Balance:  decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}

	var request models.TransactionsRequest
	
	body := c.Body()

	// แปลง JSON body เป็น struct
	if err := json.Unmarshal(body, &request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON format")
	}



	for _, transaction := range request.Transactions { 

 	transaction.IsAction = "GameResult"

	// ตรวจสอบ ว่า มี transactions เดิมอยู่มั้ย
    var _transaction_found models.TransactionSub
	db,_ := database.GetDatabaseConnection(request.MemberName)
	_terr := db.Model(&models.TransactionSub{}).Where("WagerID = ?",transaction.WagerID).Scan(&_transaction_found).RowsAffected
	
	if _terr == 0 {
		//ตรวจสอบว่า เป็น transactions buyin buyout หรือไม่
		var buyinout models.BuyInOut;
		_berr := db.Model(&models.BuyInOut{}).Where("WagerID = ?",transaction.WagerID).Scan(&buyinout).RowsAffected

		// ถ้าเป็น buyin buyout
		if _berr >  0 {
			 result := AddTransactions(transaction,request.MemberName)
			 responseBalance, _ := result.Data.(ResponseBalance)
			 
			 response = EFResponse{
				ErrorCode:    0,
				ErrorMessage: "สำเร็จ" ,
				Balance:      responseBalance.Balance,
				BeforeBalance: responseBalance.BeforeBalance,
			}
		} else {
			response = EFResponse{
				ErrorCode:16,
				ErrorMessage:"ไม่พบรายการ",
				Balance:  transaction.TransactionAmount,
				BeforeBalance: transaction.TransactionAmount,
			}
		}

	} else {
		var c_transaction_found models.TransactionSub
		rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
		 if rowsAffected == 0 {
			result := AddTransactions(transaction,request.MemberName)
			responseBalance, _ := result.Data.(ResponseBalance)
			 
			 response = EFResponse{
				ErrorCode:    0,
				ErrorMessage: "สำเร็จ",
				Balance:      responseBalance.Balance,
				BeforeBalance: responseBalance.BeforeBalance,
			}
		 } else {
			response = EFResponse{
				ErrorCode:16,
				ErrorMessage:"รายการซ้ำ" ,
				Balance:  decimal.NewFromFloat(0),
				BeforeBalance: decimal.NewFromFloat(0),
			}
		 }
	
		}
			 
 
	}
	 
	return c.JSON(response)
}
func RollBack(c *fiber.Ctx) error {

	request := new(models.TransactionsRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"rollback",request.RequestTime) == true { 

		
		var user models.Users
		
		db,_ := database.GetDatabaseConnection(request.MemberName)
		 for _, transaction := range request.Transactions {
			transaction.IsAction = "Rollback"
			db.Where("username = ?", request.MemberName).First(&user)
	  		
			if user.Balance.LessThan(transaction.BetAmount) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			 {
				var c_transaction_found models.TransactionSub
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
				
				 if rowsAffected == 0 {
					result := AddTransactions(transaction,request.MemberName)
					responseBalance, _ := result.Data.(ResponseBalance)
					 
					 response = EFResponse{
						ErrorCode:    0,
						ErrorMessage: "สำเร็จ",
						Balance:      responseBalance.Balance,
						BeforeBalance: responseBalance.BeforeBalance,
					}
				 } else {
					response = EFResponse{
						ErrorCode:16,
						ErrorMessage:"รายการซ้ำ" ,
						Balance:  decimal.NewFromFloat(0),
						BeforeBalance: decimal.NewFromFloat(0),
					}
				 }
			
				}
			
			// {
					
			// 		result := AddTransactions(transaction,request.MemberName)
			// 		responseBalance, _ := result.Data.(ResponseBalance)
				
			// 		response = EFResponse{
			// 			ErrorCode:    0,
			// 			ErrorMessage: "สำเร็จ",
			// 			Balance:      responseBalance.Balance,
			// 			BeforeBalance: responseBalance.BeforeBalance,
			// 	}
			
			// }
		 }
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func CancelBet(c *fiber.Ctx) error {

	request := new(models.TransactionsRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"cancelbet",request.RequestTime) == true { 

		
		var user models.Users
		
		db,_ := database.GetDatabaseConnection(request.MemberName)
		 for _, transaction := range request.Transactions {
			transaction.IsAction = "CancelBet"
			db.Where("username = ?", request.MemberName).First(&user)
	  		
			if user.Balance.LessThan(transaction.BetAmount) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			 {
				var c_transaction_found models.TransactionSub
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
				
				 if rowsAffected == 0 {
					result := AddTransactions(transaction,request.MemberName)
					responseBalance, _ := result.Data.(ResponseBalance)
					 
					 response = EFResponse{
						ErrorCode:    0,
						ErrorMessage: "สำเร็จ",
						Balance:      responseBalance.Balance,
						BeforeBalance: responseBalance.BeforeBalance,
					}
				 } else {
					response = EFResponse{
						ErrorCode:16,
						ErrorMessage:"รายการซ้ำ" ,
						Balance:  decimal.NewFromFloat(0),
						BeforeBalance: decimal.NewFromFloat(0),
					}
				 }
			
				}
			
			// {
					
			// 		result := AddTransactions(transaction,request.MemberName)
			// 		responseBalance, _ := result.Data.(ResponseBalance)
				
			// 		response = EFResponse{
			// 			ErrorCode:    0,
			// 			ErrorMessage: "สำเร็จ",
			// 			Balance:      responseBalance.Balance,
			// 			BeforeBalance: responseBalance.BeforeBalance,
			// 	}
			
			// }
		 }
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func Bonus(c *fiber.Ctx) error {
	request := new(models.TransactionsRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"bonus",request.RequestTime) == true { 

		
		var user models.Users
		
		db,_ := database.GetDatabaseConnection(request.MemberName)
		 for _, transaction := range request.Transactions {
			transaction.IsAction = "Bonus"
			db.Where("username = ?", request.MemberName).First(&user)
	  		
			if user.Balance.LessThan(transaction.BetAmount) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			{
				var c_transaction_found models.TransactionSub
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
				if rowsAffected == 0 {

							result := AddTransactions(transaction,request.MemberName)
							responseBalance, _ := result.Data.(ResponseBalance)
						
							response = EFResponse{
								ErrorCode:    0,
								ErrorMessage: "สำเร็จ",
								Balance:      responseBalance.Balance,
								BeforeBalance: responseBalance.BeforeBalance,
						}
					} else {
						response = EFResponse{
							ErrorCode:16,
							ErrorMessage:"รายการซ้ำ" ,
							Balance:  decimal.NewFromFloat(0),
							BeforeBalance: decimal.NewFromFloat(0),
						}
					}
			}
		 }
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func Jackpot(c *fiber.Ctx) error {
	request := new(models.TransactionsRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"jackpot",request.RequestTime) == true { 

		
		var user models.Users
		db,_ := database.GetDatabaseConnection(request.MemberName)
		
		 for _, transaction := range request.Transactions {
			transaction.IsAction = "JackPot"
			db.Where("username = ?", request.MemberName).First(&user)
	  		
			if user.Balance.LessThan(transaction.BetAmount) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			{
				var c_transaction_found models.TransactionSub
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
				if rowsAffected == 0 {

							result := AddTransactions(transaction,request.MemberName)
							responseBalance, _ := result.Data.(ResponseBalance)
						
							response = EFResponse{
								ErrorCode:    0,
								ErrorMessage: "สำเร็จ",
								Balance:      responseBalance.Balance,
								BeforeBalance: responseBalance.BeforeBalance,
						}
					} else {
						response = EFResponse{
							ErrorCode:16,
							ErrorMessage:"รายการซ้ำ" ,
							Balance:  decimal.NewFromFloat(0),
							BeforeBalance: decimal.NewFromFloat(0),
						}
					}
			}
		 }
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func PushBet(c *fiber.Ctx) error {
	request := new(models.TransactionsRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"pushbet",request.RequestTime) == true { 

		
		var user models.Users
		
		db,_ := database.GetDatabaseConnection(request.MemberName)
		 for _, transaction := range request.Transactions {
			transaction.IsAction = "PushBet"
			db.Where("username = ?", request.MemberName).First(&user)
	  		
			if user.Balance.LessThan(transaction.BetAmount) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			{
				var c_transaction_found models.TransactionSub
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("TransactionID = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
				if rowsAffected == 0 {
							if transaction.TurnOver.IsZero() {
								transaction.TurnOver = transaction.BetAmount
							}
							result := AddTransactions(transaction,request.MemberName)
							responseBalance, _ := result.Data.(ResponseBalance)
						
							response = EFResponse{
								ErrorCode:    0,
								ErrorMessage: "สำเร็จ",
								Balance:      responseBalance.Balance,
								BeforeBalance: responseBalance.BeforeBalance,
						}
					} else {
						if transaction.TurnOver.IsZero() {
							transaction.TurnOver = transaction.BetAmount
						}
						    multi_result := AddTransactions(transaction,request.MemberName)
							multi_responseBalance, _ := multi_result.Data.(ResponseBalance)
						
							response = EFResponse{
								ErrorCode:    0,
								ErrorMessage: "สำเร็จ",
								Balance:      multi_responseBalance.Balance,
								BeforeBalance: multi_responseBalance.BeforeBalance,
						}
					}
			}
		 }
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func MobileLogin(c *fiber.Ctx) error {

	type Authorized struct {
		MemberName string `json:"membername"`
		Password string `json:"password"`
	}
	response := EFResponse{
		ErrorCode:16,
		ErrorMessage:"กรุณาตรวจสอบ ชื่อผู้ใช้ และ รหัสผ่าน อีกครั้ง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}

	request := new(Authorized)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	var user models.Users
	db,_ := database.GetDatabaseConnection(request.MemberName)
	rowsAffected := db.Where("username = ? and password = ?", request.MemberName,request.Password).First(&user).RowsAffected
	  	
	if rowsAffected == 0 {
		response = EFResponse{
			ErrorCode:16,
			ErrorMessage:"กรุณาตรวจสอบ ชื่อผู้ใช้ และ รหัสผ่าน อีกครั้ง",
			Balance: decimal.NewFromFloat(0),
			BeforeBalance: decimal.NewFromFloat(0),
		}
	} else {

		response = EFResponse{
			ErrorCode:0,
			ErrorMessage:"สำเร็จ",
			Balance:      user.Balance,
			BeforeBalance: decimal.NewFromFloat(0),
		}
	}

	return c.JSON(response)
	
}

func BuyIn(c *fiber.Ctx) error {

	request := new(models.BuyInOutRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//fmt.Println(request.Transaction)

	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ main",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"buyin",request.RequestTime) == true { 

		
		var user models.Users
		
		
		db,_ := database.GetDatabaseConnection(request.MemberName)

			db.Where("username = ?", request.MemberName).First(&user)
	  		//fmt.Println(&request)
			if user.Balance.LessThan(request.Transaction.TransactionAmount.Abs()) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			{
				var c_transaction_found models.BuyInOut
				rowsAffected := db.Model(&models.BuyInOut{}).Select("id").Where("transaction_id = ? ",request.Transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
				if rowsAffected == 0 {

							result := AddBuyInOut(request.Transaction,request.MemberName)
							responseBalance, _ := result.Data.(ResponseBalance)
						
							response = EFResponse{
								ErrorCode:    0,
								ErrorMessage: "สำเร็จ",
								Balance:      responseBalance.Balance,
								BeforeBalance: responseBalance.BeforeBalance,
							}
					} else {
						response = EFResponse{
							ErrorCode:16,
							ErrorMessage:"รายการซ้ำ" ,
							Balance:  decimal.NewFromFloat(0),
							BeforeBalance: decimal.NewFromFloat(0),
						}
					}
			}
		 
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func BuyOut(c *fiber.Ctx) error {

	request := new(models.BuyInOutRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	
	response := EFResponse{
		ErrorCode: 0,
		ErrorMessage:"สำเร็จ",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	
	
	
	if CheckSign(request.Sign,"buyout",request.RequestTime) == true { 

		
		var user models.Users
		
		
		db,_ := database.GetDatabaseConnection(request.MemberName)
			
			db.Where("username = ?", request.MemberName).First(&user)
	  		
			if user.Balance.LessThan(request.Transaction.BetAmount) {

				response = EFResponse{
					ErrorCode:    1001,
					ErrorMessage: "Insufficient Balance",
					Balance:      user.Balance,
					BeforeBalance: decimal.NewFromFloat(0),
				}		
			} else 
			{
				var c_transaction_found models.BuyInOut
				rowsAffected := db.Model(&models.BuyInOut{}).Where("transaction_id = ? ",request.Transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
				if rowsAffected == 0 {

							result := AddBuyInOut(request.Transaction,request.MemberName)
							responseBalance, _ := result.Data.(ResponseBalance)
						
							response = EFResponse{
								ErrorCode:    0,
								ErrorMessage: "สำเร็จ",
								Balance:      responseBalance.Balance,
								BeforeBalance: responseBalance.BeforeBalance,
						}
					} else {
						response = EFResponse{
							ErrorCode:16,
							ErrorMessage:"รายการซ้ำ" ,
							Balance:  decimal.NewFromFloat(0),
							BeforeBalance: decimal.NewFromFloat(0),
						}
					}
			}
	 
		} else {
	response = EFResponse{
		ErrorCode:1004,
		ErrorMessage:"โทเคนไม่ถูกต้อง",
		Balance: decimal.NewFromFloat(0),
		BeforeBalance: decimal.NewFromFloat(0),
	}
	
	}
	return c.JSON(response)
}
func makePostRequest(url string, bodyData interface{}) (*fasthttp.Response, error) {
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
	//authHeader := createBasicAuthHeader(OPERATOR_CODE, SECRET_API_KEY)
	//req.Header.Add("Authorization", authHeader)
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
	//authHeader := createBasicAuthHeader(OPERATOR_CODE, SECRET_API_KEY)
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
func hashSignature(MethodName string,requestTime string) string {

	hash := md5.New()
    hash.Write([]byte(common.EF_OPERATOR_CODE + requestTime + strings.ToLower(MethodName) + common.EF_SECRET_KEY))
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	return md5Hash
}
func GetGameList(c *fiber.Ctx) error {
	type BodyGame struct {
		ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
			
	}

	type PgRequest struct {
		Id string `json:"id"`
		TimestampMillis int `json:"timestampmillis"`
		ProductID string `json:"productid`
		Currency string `json:"currency"`
		Username string `json:"username"`
		SessionToken string `json:"sessiontoken"`
		StatusCode  int  `json:"statuscode"`
		Balance  decimal.Decimal `json:"balance"`
		//ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	//	Txns []TxnsRequest `json:"txns"`
	}
	type CResponse struct {
		Message string      `json:"message"`
		Status  bool        `json:"status"`
		Data    interface{} `json:"data"`  
	}

	type CEResponse struct {
		Message interface{}  `json:"message"`
		Status  bool        `json:"status"`
		Data    interface{} `json:"data"` 
		ProviderGames    interface{} `json:"providergames"`  
	}

	var response CEResponse

	// bodyRequest := new(BodyGame)
	// if err := c.BodyParser(&bodyRequest); err != nil {
	// 	fmt.Printf(" %s ", err.Error())
	// 	response := fiber.Map{
	// 		"Status":  false,
	// 		"Message": err.Error(),
	// 	}
	// 	return c.JSON(response)
	// }

	// fmt.Printf("Body: %s",bodyRequest.Body)
	//var tokenString := c.Get("Authorization")[7:]
	request := new(BodyGame)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//var users models.Users
	//users = ValidateJWTReturn(request.SessionToken);

	//fmt.Printf("users: %v ",users)
	//fmt.Printf("request: %s ",request.SessionToken)

	// var args = fiber.Map{
	// 	"username": strings.ToLower(users.Username),//user.data.username,
	// 	"productId":pg_prod_code,
	// 	"gameCode": request.ProductID,
	// 	"isMobileLogin": true,
	// 	"sessionToken": request.SessionToken,
	// 	//"betLimit": [],
	// 	"callbackUrl":"https://www.โชคดี789.com/lobby/slot/game?id=8888&type=1", //`${req.protocol}://${req.get('host')}${req.originalUrl}`
	// }
	var RequestTime = time.Now().Format("20060102150405")
	var args = fiber.Map{
		"OperatorCode": common.EF_OPERATOR_CODE,
		"MemberName": common.USER_FIX,//req.body.username,
		"Password":   common.PASS_FIX,//user.data.uid,
		"ProductID": request.ProductID,
		"GameType": request.GameType,
		"LanguageCode": request.LanguageCode,
		"Platform": request.Platform,
		"Sign": hashSignature("LaunchGame",RequestTime),
		"RequestTime": RequestTime,
		}
	
	resp,err := makePostRequest(common.INFINITY_PROD_URL+"/Seamless/GetGameList",args)		
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	resultBytes := resp.Body()
	resultString := string(resultBytes)
	// แสดงผล string ที่ได้
	//fmt.Println("Response body as string:", resultString)

	

	err = json.Unmarshal([]byte(resultString), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	//url := fmt.Sprintf(PG_PROD_URL,"/seamless/login")

	respon := fiber.Map{
		"Status":  true,
		"Message": "ดึงข้อมูลสำเร็จ",
		"Data": fiber.Map{
			"games":response.ProviderGames,
		},
	}
	return c.JSON(respon)
}


// var SECRET_KEY = os.Getenv("PASSWORD_SECRET")
// var pg_prod_code = os.Getenv("PG_PRODUCT_ID")

// var OPERATOR_CODE = "sunshinetest" //"sunshinepgthb"//"sunshinetest",
// var SECRET_API_KEY = os.Getenv("PG_API_KEY") //"9dc857f4-2225-45ef-bf0f-665bcf7d4a1b" //os.Getenv("PG_API_KEY")
// var PG_PROD_CODE= os.Getenv("PG_PRODUCT_ID")
// var PG_API_URL = "https://test.ambsuperapi.com"//os.Getenv("PG_API_URL") //"https://prod_md.9977997.com"
// var PG_PROD_URL = "https://api.hentory.io" 


// func makePostRequest(url string, bodyData interface{}) (*fasthttp.Response, error) {
// 	// Marshal requestData struct เป็น JSON
// 	jsonData, err := json.Marshal(bodyData)
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshaling JSON: %v", err)
// 	}

// 	// สร้าง Request และ Response
// 	req := fasthttp.AcquireRequest()
// 	resp := fasthttp.AcquireResponse()

// 	// ตั้งค่า URL, Method, และ Body
// 	req.SetRequestURI(url)
// 	req.Header.SetMethod("POST")
// 	req.Header.SetContentType("application/json")
// 	authHeader := common.CreateBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
// 	req.Header.Add("Authorization", authHeader)
// 	req.SetBody(jsonData)

// 	// ส่ง request
// 	client := &fasthttp.Client{}
// 	if err := client.Do(req, resp); err != nil {
// 		return nil, fmt.Errorf("error making POST request: %v", err)
// 	}

// 	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
// 	fasthttp.ReleaseRequest(req)
	
// 	return resp, nil
// }
// func makeGetRequest(url string) (*fasthttp.Response, error) {
// 	// Marshal requestData struct เป็น JSON
// 	// jsonData, err := json.Marshal(bodyData)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("error marshaling JSON: %v", err)
// 	// }

// 	// สร้าง Request และ Response
// 	req := fasthttp.AcquireRequest()
// 	resp := fasthttp.AcquireResponse()

// 	// ตั้งค่า URL, Method, และ Body
// 	req.SetRequestURI(url)
// 	req.Header.SetMethod("GET")
// 	req.Header.SetContentType("application/json")
// 	authHeader := common.CreateBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
// 	req.Header.Add("Authorization", authHeader)
// 	//req.SetBody(jsonData)

// 	// ส่ง request
// 	client := &fasthttp.Client{}
// 	if err := client.Do(req, resp); err != nil {
// 		return nil, fmt.Errorf("error making POST request: %v", err)
// 	}

// 	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
// 	fasthttp.ReleaseRequest(req)
	
// 	return resp, nil
// }

func LaunchGame(c *fiber.Ctx) error {
	type BodyGame struct {
		ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	}

	type EfRequest struct {
		Id string `json:"id"`
		TimestampMillis int `json:"timestampmillis"`
		ProductID string `json:"productid`
		Currency string `json:"currency"`
		Username string `json:"username"`
		SessionToken string `json:"sessiontoken"`
		StatusCode  int  `json:"statuscode"`
		Balance  decimal.Decimal `json:"balance"`
		//ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	//	Txns []TxnsRequest `json:"txns"`
	}
	type CResponse struct {
		// ErrorMessage string `json:"message"`
		// Status  bool        `json:"status"`
		Url    string `json:"url"`
		ErrorCode int `json:"errorcode"`
		ErrorMessage interface{} `json:"errormessage"`
	}

	var response CResponse
	// bodyRequest := new(BodyGame)

	// if err := c.BodyParser(&bodyRequest); err != nil {
	// 	fmt.Printf(" %s ", err.Error())
	// 	response := fiber.Map{
	// 		"Status":  false,
	// 		"Message": err.Error(),
	// 	}
	// 	return c.JSON(response)
	// }

	// fmt.Printf("Body: %s",bodyRequest.Body)
	//var tokenString := c.Get("Authorization")[7:]
	request := new(EfRequest)
	if err := c.BodyParser(request); err != nil {
		respon := fiber.Map{
			"Status":  false,
			"Message": "BodyParser Error",
			"Data": response,
		}
		return c.JSON(respon)
	}
	var users models.Users
	users = handler.ValidateJWTReturn(request.SessionToken);

	//fmt.Printf("users: %s ",users)
	//fmt.Printf("request: %s ",request.SessionToken)
	var RequestTime = time.Now().Format("20060102150405")
	// var args = fiber.Map{
	// 	"OperatorCode": common.EF_OPERATOR_CODE,
	// 	"MemberName": common.USER_FIX,//req.body.username,
	// 	"Password":   common.PASS_FIX,//user.data.uid,
	// 	"ProductID": request.ProductID,
	// 	"GameType": request.GameType,
	// 	"LanguageCode": request.LanguageCode,
	// 	"Platform": request.Platform,
	// 	"Sign": hashSignature("LaunchGame",RequestTime),
	// 	"RequestTime": RequestTime,
	// 	}
	var  args = fiber.Map{
		"OperatorCode": common.EF_OPERATOR_CODE,
		"MemberName": users.Username,
		"Password":   users.Password,
		"ProductID": request.ProductID,
		"GameType": request.GameType,
		"GameID": request.GameID,
		"LanguageCode": request.LanguageCode,
		"Platform": request.Platform,
		"Sign": hashSignature("LaunchGame",RequestTime),
		"RequestTime": RequestTime,
		}
	// var args = fiber.Map{
	// 	"username": strings.ToLower(users.Username),//user.data.username,
	// 	"productId":common.PG_PROD_CODE,
	// 	"gameCode": request.ProductID,
	// 	"isMobileLogin": true,
	// 	"sessionToken": request.SessionToken,
	// 	//"betLimit": [],
	// 	"callbackUrl":"https://www.โชคดี789.com/lobby/slot/game?id=8888&type=1", //`${req.protocol}://${req.get('host')}${req.originalUrl}`
	// }
	
	 
	
	resp,err := makePostRequest(common.INFINITY_PROD_URL+"/Seamless/LaunchGame",args)		
	if err != nil {
		respon := fiber.Map{
			"Status":  false,
			"Message": "RequestError",
			"Data": response,
		}
		return c.JSON(respon)
	}
	resultBytes := resp.Body()
	resultString := string(resultBytes)
	// แสดงผล string ที่ได้
	//fmt.Println("Response body as string:", resultString)

	err = json.Unmarshal([]byte(resultString), &response)
	if err != nil {
		respon := fiber.Map{
			"Status":  false,
			"Message": "Unmarshalling Error",
			"Data": response,
		}
		return c.JSON(respon)
	}

	respon := fiber.Map{
		"Status":  true,
		"Message": "Success",
		"Data": response,
	}
	return c.JSON(respon)
}