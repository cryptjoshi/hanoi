package ef

import 
(
	"github.com/gofiber/fiber/v2"
	"hanoi/models"
	"hanoi/database"
	"hanoi/repository"
	"github.com/shopspring/decimal"
	"crypto/md5"
	"encoding/hex"
	"os"
	"encoding/json"
	"time"
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
const EF_SECRET_KEY="1g1bb3" //stagging
var OPERATOR_CODE = os.Getenv("EF_OPERATOR")

func parseTime(layout, value string) (time.Time, error) {
    return time.Parse(layout, value)
}

func Index(c *fiber.Ctx) error {

	//var user []models.Users
	//db.Find(&user)
	response := Response{
		Message: "Welcome to Efinity!!",
		Status:  true,
		Data: fiber.Map{ 
			//"users":user,
		}, 
	}
	 
	return c.JSON(response)
   
}

func CheckSign(Signature string,methodName string,requestTime string) bool {

	
	//requestTime := "2024-09-15T12:00:00Z"
	//methodName := "MethodName"
	
	secretKey := EF_SECRET_KEY

	// สร้างข้อมูลที่ต้องใช้ hash
	data := OPERATOR_CODE + requestTime + strings.ToLower(methodName) + secretKey

	 
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
	
	 
	
	if CheckSign(body.Sign,"getbalance",body.RequestTime) == true {
			var users models.Users
			db, _ := database.ConnectToDB(body.MemberName)
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
	db, _ := database.ConnectToDB(membername)
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
	
		 
		  repository.UpdateFieldsUserString(db,membername, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		//fmt.Println(_err)
		// if _err != nil {
		// 	fmt.Println("Error:", _err)
		// } else {
		// 	//fmt.Println("User fields updated successfully")
		// }

 
 
	 
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
	db, _ := database.ConnectToDB(membername)
    if err_ := db.Where("username = ? ", membername).First(&users).Error; err_ != nil {
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
	
		 
		  repository.UpdateFieldsUserString(db,membername, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		//fmt.Println(_err)
		// if _err != nil {
		// 	fmt.Println("Error:", _err)
		// } else {
		// 	//fmt.Println("User fields updated successfully")
		// }

 
 
	 
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
	db, _ := database.ConnectToDB(membername)
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
	
		 
		  repository.UpdateFieldsUserString(db,membername, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		//fmt.Println(_err)
		// if _err != nil {
		// 	fmt.Println("Error:", _err)
		// } else {
		// 	//fmt.Println("User fields updated successfully")
		// }

 
 
	 
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
		
		db, _ := database.ConnectToDB(request.MemberName)
		 for _, transaction := range request.Transactions {
			
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
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("transaction_id = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
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

 

	// ตรวจสอบ ว่า มี transactions เดิมอยู่มั้ย
    var _transaction_found models.TransactionSub
	
	db, _ := database.ConnectToDB(request.MemberName)

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
		rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("transaction_id = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
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
		
		db, _ := database.ConnectToDB(request.MemberName)
		 for _, transaction := range request.Transactions {
			
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
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("transaction_id = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
				
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
		
		db, _ := database.ConnectToDB(request.MemberName)
		 for _, transaction := range request.Transactions {
			
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
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("transaction_id = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
				
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
		
		db, _ := database.ConnectToDB(request.MemberName)
		 for _, transaction := range request.Transactions {
			
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
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("transaction_id = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
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
		
		db, _ := database.ConnectToDB(request.MemberName)
		 for _, transaction := range request.Transactions {
			
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
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("transaction_id = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
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
		
		db, _ := database.ConnectToDB(request.MemberName)
		 for _, transaction := range request.Transactions {
			
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
				rowsAffected := db.Model(&models.TransactionSub{}).Select("id").Where("transaction_id = ? ",transaction.TransactionID).Find(&c_transaction_found).RowsAffected
		
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
	db, _ := database.ConnectToDB(request.MemberName)
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
		
		
		db, _ := database.ConnectToDB(request.MemberName)

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
		
		
		db, _ := database.ConnectToDB(request.MemberName)
			
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