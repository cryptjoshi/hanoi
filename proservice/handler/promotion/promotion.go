package promotion

import (
	"context"
	"fmt"
	//"log"
	"strconv"
	"time"
	"strings"
	"crypto/sha256"
    "encoding/hex"
	"hanoi/models"
	"hanoi/handler"
	"hanoi/handler/wallet"
	//"hanoi/handler/redisUtils"
	"github.com/shopspring/decimal"
   	"hanoi/repository"
   
	
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/go-redis/redis/v8"
)


var ctx = context.Background()


type PBody struct {
	UserID  string  `json:"userID"`
	ProID   string  `json:"proID"`
	Amount  float64 `json:"amount"`
}
var ttl time.Duration = 5 * time.Second
 
// ฟังก์ชันตรวจสอบการเชื่อมต่อ Redis
func checkRedisConnection(redisClient *redis.Client) error {
	_, err := redisClient.Ping(ctx).Result()
	return err
}
 
// func acquireLock(redisClient *redis.Client, lockKey string) (bool, error) {
// 	// ใช้ `HSet` เพื่อสร้างล็อก
// 	locked, err := redisClient.HSetNX(ctx, lockKey, "locked", true).Result()
// 	if err != nil {
// 		return false, err
// 	}
// 	if locked {
// 		// ตั้งเวลาหมดอายุ
// 		go func() {
// 			time.Sleep(10 * time.Second)
// 			handler.ReleaseLock(redisClient, lockKey)
// 		}()
// 	}
// 	return locked, nil
// }

// // ฟังก์ชันปลดล็อก
// func handler.ReleaseLock(redisClient *redis.Client, lockKey string) error {
// 	// ลบฟิลด์ที่เป็นล็อก
// 	_, err := redisClient.HDel(ctx, lockKey, "locked").Result()
// 	return err
// }

// ฟังก์ชันเช็คสถานะล็อก
func isLocked(redisClient *redis.Client, lockKey string) (bool, error) {
	locked, err := redisClient.HGet(ctx, lockKey, "locked").Result()
	if err != nil { 
	if err == redis.Nil {
			return false, nil
		}
		return false,err
	}
	return locked == "1", nil
}
// ฟังก์ชันเพื่อตรวจสอบว่าคีย์มีประเภทที่ถูกต้องหรือไม่
func checkKeyType(redisClient *redis.Client, key string, expectedType string) error {
	keyType, err := redisClient.Type(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error checking key type: %v", err)
	}
	if keyType != expectedType && expectedType != "any" {
		return fmt.Errorf("key %s is of type %s, expected %s", key, keyType, expectedType)
	}
	return nil
}

// ฟังก์ชันเพื่อสร้างคีย์ total_transactions ถ้ายังไม่มี
func ensureTotalTransactionsKey(redisClient *redis.Client, userID string) error {
	totalKey := fmt.Sprintf("%s:total_transactions", userID)

	// เช็คประเภทของคีย์
	keyType, err := redisClient.Type(ctx, totalKey).Result()
	if err != nil {
		return fmt.Errorf("error checking key type: %v", err)
	}

	// ถ้ายังไม่มีคีย์ หรือประเภทไม่ถูกต้อง ให้สร้างใหม่
	if keyType == "none" {
		// สร้างคีย์ใหม่ในประเภท hash
		if err := redisClient.HSet(ctx, totalKey, "total_amount", 0).Err(); err != nil {
			return fmt.Errorf("error initializing total transactions: %v", err)
		}
	}

	return nil
}

 


// ฟังก์ชันในการสร้าง UID
func generateUID(userID, proID, timestamp string) string {
    hash := sha256.New()
    hash.Write([]byte(userID + proID + timestamp))
    return hex.EncodeToString(hash.Sum(nil))
}
func selectPromotion(redisClient *redis.Client, userID string,balance float64, pro_setting map[string]interface{}) error {
    lockKey := fmt.Sprintf("%s:lock", userID)

    // พยายามล็อก
    if locked, err := handler.AcquireLock(redisClient, lockKey,ttl); err != nil || !locked {
        return fmt.Errorf("could not acquire lock: %v", err)
    }
    defer handler.ReleaseLock(redisClient, lockKey)

    // คีย์สำหรับโปรโมชั่นปัจจุบัน
    currentPromotionKey := fmt.Sprintf("%s:current_promotion", userID)

    // ดึงข้อมูลโปรโมชั่นปัจจุบัน
    promotionStatus, err := redisClient.HGet(ctx, currentPromotionKey, "status").Result()
    if err != nil && err != redis.Nil {
        return fmt.Errorf("error checking current promotion status: %v", err)
    }

    // ถ้าสถานะเป็น "" กำหนดให้เป็น "2" (ended)
    if promotionStatus == "" {
        promotionStatus = "2"
    }

    // ตรวจสอบสถานะ
    if promotionStatus != "0" && promotionStatus != "2" {
        return fmt.Errorf("โปรโมชั่นเดิม ยังไม่สิ้นสุด is %s", promotionStatus)
    }

	

	if promotionStatus == "2" && balance >= 1 {
		return fmt.Errorf("ยอดเงินคงเหลือมากกว่าศูนย์")
	}
	if balance < 1 {
		//return fmt.Errorf("ยอดเงินคงเหลือน้อยกว่าศูนย์")
		err := ClearRedisKey(redisClient, currentPromotionKey)

		if err != nil {
			fmt.Println("Error clearing Redis keys:", err)
		} else {
			fmt.Println("Successfully cleared Redis keys for user:", userID)
		}
	}
    // ดึง timestamp ปัจจุบัน
    now := time.Now()
    nowRFC3339 := now.Format(time.RFC3339)

	


    // สร้าง UID โดยการแฮช userID, proID และ timestamp
    uid := generateUID(userID, fmt.Sprintf("%d", pro_setting["Id"]), nowRFC3339)

 

	response := map[string]interface{}{
     "proID":         fmt.Sprintf("%d", pro_setting["Id"]),
    "status":        "0",
    "timestamp":     nowRFC3339,
    "uid":          uid,
    "Id":           fmt.Sprintf("%s", pro_setting["Userid"]), // แปลงเป็น string
    "Type":         fmt.Sprintf("%s", pro_setting["Type"]),
    "count":        fmt.Sprintf("%d", pro_setting["count"]), // แปลงเป็น string
    "MinTurnover":  fmt.Sprintf("%s", pro_setting["MinTurnover"]),
    "Example":      fmt.Sprintf("%s", pro_setting["Formular"]),
    "Name":         fmt.Sprintf("%s", pro_setting["Name"]),
    "TurnType":     fmt.Sprintf("%s", pro_setting["TurnType"]),
    "Week":         fmt.Sprintf("%s", pro_setting["Week"]),

    // ข้อมูลใหม่ที่เพิ่มเข้ามาจาก pro_setting
    "minDept":       fmt.Sprintf("%s", pro_setting["minDept"]), // แปลงเป็น string หากจำเป็น
    "maxDept":       fmt.Sprintf("%s", pro_setting["maxDept"]),
	 // แปลงเป็น string หากจำเป็น
    "Widthdrawmax":  fmt.Sprintf("%s", pro_setting["Widthdrawmax"]), // แปลงเป็น string
    "Widthdrawmin":  fmt.Sprintf("%s", pro_setting["Widthdrawmin"]), // แปลงเป็น string
    "MinSpendType":  fmt.Sprintf("%s", pro_setting["MinSpendType"]), // แปลงเป็น string
    "MinCreditType": fmt.Sprintf("%s", pro_setting["MinCreditType"]), // แปลงเป็น string
    "MaxWithdrawType": fmt.Sprintf("%s", pro_setting["MaxWithdrawType"]), // แปลงเป็น string
    "MinCredit":     fmt.Sprintf("%s", pro_setting["MinCredit"]), // แปลงเป็น string
    "Zerobalance":   fmt.Sprintf("%d", pro_setting["Zerobalance"]), // แปลงเป็น string
    "CreatedAt":     pro_setting["CreatedAt"], // ควรตรวจสอบประเภทให้แน่ใจว่าเป็นประเภทที่เหมาะสม
    "MaxUse":        fmt.Sprintf("%d", pro_setting["MaxUse"]), // แปลงเป็น string
    }

	  

    // บันทึกโปรโมชั่นใหม่โดยเขียนทับโปรโมชั่นเก่า
    if err := redisClient.HSet(ctx, currentPromotionKey, response).Err(); err != nil {
        return err
    }

    return nil
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

func deposit(redisClient *redis.Client, userID string, proID string, depositAmount float64) (map[string]string,error) {
    
	lockKey := fmt.Sprintf("%s:lock", userID)

    // ตรวจสอบการล็อก
    if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
        return nil,fmt.Errorf("could not acquire lock: %v", err)
    }
	defer handler.ReleaseLock(redisClient, lockKey)

    key := fmt.Sprintf("%s:%s", userID, proID)
 
    // ตรวจสอบให้แน่ใจว่ามี key สำหรับการทำธุรกรรมแต่ละครั้ง
    if err := ensureTotalTransactionsKey(redisClient, userID); err != nil {
        return nil,fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 297")
    }

	currentPromotionKey := fmt.Sprintf("%s:current_promotion", userID)

    // ดึงข้อมูลโปรโมชั่นปัจจุบัน
    currentStatus, err := redisClient.HGet(ctx, currentPromotionKey, "status").Result()
    if err != nil && err != redis.Nil {
		return nil,fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 347")
    }
    // // ตรวจสอบประเภทของ key
    // if err := checkKeyType(redisClient, key, "hash"); err != nil {
    //     return fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 302")
    // }

    // ดึงสถานะปัจจุบัน
    // currentStatus, err := redisClient.HGet(ctx, key, "status").Result()
    // if err != nil {
    //     return fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 308")
    // }
    //fmt.Println("currentStatus:",currentStatus)
    // แปลงสถานะเป็น int
    intStr, err := strconv.Atoi(currentStatus)
    if err != nil || intStr > 0 {
        return nil,fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 314")
    }

    // ดึงยอดเงินปัจจุบัน
    balanceStr, err := redisClient.HGet(ctx, key, "balance").Result()
    if err != nil && err == redis.Nil {
		//return nil,fmt.Errorf(err.Error())
		balanceStr = "0"
    }

    balance, err := strconv.ParseFloat(balanceStr, 64)
    if err != nil {
        return nil,fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 325")
    }
	
	minDeptStr, err := redisClient.HGet(ctx, currentPromotionKey, "minDept").Result()
	minDept,_ := strconv.ParseFloat(minDeptStr, 64)
    if err != nil && err != redis.Nil || depositAmount < minDept {
		return nil,fmt.Errorf("ยอดเงินฝาก น้อยกว่ายอดเงินฝากขั้นต่ำ")
    }
	example, err := redisClient.HGet(ctx, currentPromotionKey, "Example").Result()
    if err != nil && err != redis.Nil {
		return nil,fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 305")
    }
	Formular := strings.Replace(example, "deposit", fmt.Sprintf("%.2f", depositAmount), 1) // Convert deposit to string if necessary
	fmt.Printf(" %v ",Formular)
	// // Evaluate the expression (you may need to implement a function to evaluate the string expression)
	balanceIncrease, err := wallet.EvaluateExpression(Formular)
	balanceIncrease64,_ := balanceIncrease.Float64()
    // บันทึกยอดก่อนการฝาก
    if err := redisClient.HSet(ctx, key, "before_balance", balance).Err(); err != nil {
        return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
    }
	maxDeptStr,_ := redisClient.HGet(ctx, currentPromotionKey, "maxDept").Result()
	maxDeptFloat,_ := strconv.ParseFloat(maxDeptStr, 64)

	if balanceIncrease64 - depositAmount > maxDeptFloat {
		//BankStatement.Balance = users.Balance.Add(deposit.Add(pro_setting["maxDept"].(decimal.Decimal)))
		balanceIncrease64 = depositAmount + maxDeptFloat
	} //else {
		//fmt.Printf("wallet.go 453 line balanceIncrease: %v \n",balanceIncrease)
		//BankStatement.Balance = users.Balance.Add(balanceIncrease)
	//	BankStatement.Balance = balanceIncrease
	//}
 
    // เพิ่มจำนวนเงินฝาก
    if err := redisClient.HIncrByFloat(ctx, key, "amount", depositAmount).Err(); err != nil {
        return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
    }

	
	if err := redisClient.HSet(ctx, key, "bonus_amount",balanceIncrease64-depositAmount).Err(); err != nil {
        return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
    }
    // บันทึกยอดหลังการฝาก
    if err := redisClient.HSet(ctx, key, "after_balance",balanceIncrease64).Err(); err != nil {
        return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
    }

    // เปลี่ยนสถานะของโปรโมชั่นเป็น 1 (กำลังใช้งาน)
    if err := redisClient.HSet(ctx, key, "status", "1", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
        return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
    }

	  // ยอดรวมของ user
	  totalKey := fmt.Sprintf("%s:total_transactions", userID)
	  if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
		  return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
	  }
  
	  _, err = redisClient.HIncrByFloat(ctx, totalKey, "total_amount", balanceIncrease64).Result()
	  if err != nil {
		  return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
	  }


	turntype,err := redisClient.HGet(ctx, currentPromotionKey, "TurnType").Result()
	if err !=nil {
		return nil,fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
	}
	MinSpendType,err := redisClient.HGet(ctx, currentPromotionKey, "MinSpendType").Result()
	if err !=nil {
		return nil,fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
	}
	MinTurnover,_ := getRedisKV(redisClient,currentPromotionKey,"MinTurnover");
	switch turntype {
	case "turnover":
		var baseAmount float64
            if  MinSpendType == "deposit" {
                baseAmount = depositAmount
            } else {
                baseAmount = balanceIncrease64
            }
			fmt.Printf(" minTurnover: %v \n",MinTurnover)
			fmt.Printf(" MinSpendType: %v \n",MinSpendType)
			fmt.Printf(" baseAmount: %v \n",baseAmount)
			
			// fmt.Printf(" totalTurnover: %v \n",totalTurnover)
			// fmt.Printf(" userTurnover: %v \n",users.Turnover)
		//minTurnover,_ := decimal.NewFromString(MinTurnover)
		requiredTurnover, err := wallet.CalculateRequiredTurnover(MinTurnover, decimal.NewFromFloat(baseAmount))
    	if err !=nil {
			return nil,fmt.Errorf("turnover ผิดพลาด!")
		}
		requiredTurnover64,_ := requiredTurnover.Float64()
		fmt.Printf(" requiredTurnover: %v \n",requiredTurnover64)
		if err := redisClient.HSet(ctx, key, "requiredTurnover",requiredTurnover64).Err(); err != nil {
			return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
		}
	
	case "turncredit":
		



		minCreditStr,err := redisClient.HGet(ctx, currentPromotionKey, "MinCredit").Result()
		if err !=nil {
			return nil,fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
		}
		
		minCredit,err := decimal.NewFromString(minCreditStr)
		if err != nil {
			return nil,fmt.Errorf("ไม่สามารถแปลงค่ายอดเครดิตขั้นต่ำได้")
		}
		MinCreditType,err := redisClient.HGet(ctx, currentPromotionKey, "MinCreditType").Result()
		if err !=nil {
			return nil,fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
		}
		var baseAmount float64
		
		if MinCreditType == "deposit" {
			baseAmount = depositAmount
		} else {
			
			baseAmount = balanceIncrease64
		
		}
		totalKey := fmt.Sprintf("%s:total_transactions", userID)

		totalAmount, err := redisClient.HGet(ctx, totalKey, "total_amount").Result()
		if err != nil {
			return nil, err
		}
	
		var amount float64
		fmt.Sscanf(totalAmount, "%f", &amount)
		
		requiredCredit,_ := minCredit.Float64()
		requiredcredit := requiredCredit*baseAmount
		fmt.Printf(" minCredit: %v \n",minCredit)
		fmt.Printf(" amount: %v \n",amount)
		fmt.Printf(" baseAmount: %v \n",baseAmount)
		fmt.Printf(" requiredCredit: %v \n",requiredcredit)
		if err := redisClient.HSet(ctx, key, "requiredCredit",requiredcredit).Err(); err != nil {
			return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
		}
	
		
		// if amount < requiredCredit {
		// 	return nil,fmt.Errorf("ยอดเครดิต %v น้อยกว่ายอดเครดิตขั้นต่ำ %v ", amount, requiredcredit)
		// }
		
		 
	 }
  
	 WidthdrawMaxStr,err := redisClient.HGet(ctx, currentPromotionKey, "Widthdrawmax").Result()
	 if err !=nil {
		 return nil,fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
	 }
	 
	 WidthdrawMax, err := strconv.ParseFloat(WidthdrawMaxStr, 64)
	 if err != nil {
		 return nil, fmt.Errorf("ไม่สามารถแปลง Widthdrawmax เป็น float64: %v", err)
	 }
	 MaxWithdrawType,err := redisClient.HGet(ctx, currentPromotionKey, "MaxWithdrawType").Result()
	 if err !=nil {
		 return nil,fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
	 }
	 //deposit,_ := decimal.NewFromFloat(depositAmount)
	 if MaxWithdrawType == "deposit" {
		WidthdrawMax = WidthdrawMax*depositAmount
	 } else {
		WidthdrawMax = WidthdrawMax*balanceIncrease64
	 }

	 if err := redisClient.HSet(ctx, key, "WidthdrawMax",WidthdrawMax).Err(); err != nil {
		return nil,fmt.Errorf("ไม่สามารถฝากเงินเพิ่มได้ ")
	}

    // เรียกใช้ selectPromotion เพื่อเลือกโปรโมชั่น
    // คุณอาจต้องการทำการตั้งค่า proID ให้ถูกต้อง และหากต้องการส่งข้อมูลเพิ่มเติมให้กับ selectPromotion
    //promotion := Promotion{} // สร้างอ็อบเจ็กต์ Promotion ที่จำเป็นในที่นี้
    //ProItem := ProItem{}     // สร้างอ็อบเจ็กต์ ProItem ที่จำเป็น

    // เรียกใช้ ฟังก์ชัน selectPromotion
    // if err := selectPromotion(redisClient, userID, proID, promotion, ProItem); err != nil {
    //     return fmt.Errorf("failed to select promotion: %v", err)
    // }


	result, err := redisClient.HGetAll(ctx, currentPromotionKey).Result()
    if err != nil {
        fmt.Println("Error fetching all values:", err)
          // ต้อง exit ถ้ามีข้อผิดพลาด
    }

    // แสดงผลลัพธ์
    fmt.Println("All fields and values:", result)
	return result,nil
}
  
func withdraw(redisClient *redis.Client, prefix string,userID string, proID string, withdrawAmount float64,turnoverAmount float64) (map[string]string,error) {
	
	//var percentValue decimal.Decimal
	//var percentStr = ""
	key := fmt.Sprintf("%s%s",prefix,userID)
	lockKey := fmt.Sprintf("%s:lock", key)

    // ตรวจสอบการล็อก
    if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
        return nil,fmt.Errorf("could not acquire lock: %v", err)
    }
	defer handler.ReleaseLock(redisClient, lockKey)
	
    //key := fmt.Sprintf("%s:%s", userID, proID)
 
    // ตรวจสอบให้แน่ใจว่ามี key สำหรับการทำธุรกรรมแต่ละครั้ง
    if err := ensureTotalTransactionsKey(redisClient, key); err != nil {
        return nil,fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถถอนเงินเพิ่มได้ 297")
    }

	currentPromotionKey := fmt.Sprintf("%s:current_promotion", key)

	checkpro,_ := redisClient.HGetAll(ctx, currentPromotionKey).Result()
	if len(checkpro) > 0 && checkpro["status"] != "2" {
		
	
	
    // ดึงข้อมูลโปรโมชั่นปัจจุบัน
    currentStatus, err := redisClient.HGet(ctx, currentPromotionKey, "status").Result()
	

    if err != nil && err != redis.Nil {
		return nil,fmt.Errorf("ไม่พบข้อมูล โปรโมชั่น")
    }
	
	totalKey := fmt.Sprintf("%s:total_transactions", key)
	totalAmount, err := redisClient.HGet(ctx, totalKey, "total_amount").Result()
	if err != nil {
		return nil, err
	}
	var amount float64
	fmt.Sscanf(totalAmount, "%f", &amount)

	if amount < 1 {

		currentStatus = "2"
	}






	//fmt.Println(" currentStatus: ",currentStatus)
	
	
	// Widthdrawmin,err := redisClient.HGet(ctx, currentPromotionKey, "Widthdrawmin").Result()
	// if err !=nil {
	// 	return fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
	// }
	// widthdrawmin,_ := decimal.NewFromString(Widthdrawmin)
	// if widthdrawmin.GreaterThan(users.Balance) {
		 
	// 		return fmt.Errorf("ยอดคงเหลือน้อยกว่ายอดถอนขั้นต่ำของโปรโมชั่น (%v %v)", Widthdrawmin, users.Currency)
		 
	// }

	// response,err := getTotalTransactions(redisClient,fmt.Sprintf("%s%s","ckd",userID))
	// fmt.Printf("Response: %+v \n",response)

	// balanceStr, err := redisClient.HGet(ctx, currentPromotionKey, "after_balance").Result()
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(" balanceStr: ",balanceStr)
	// balance, err := strconv.ParseFloat(balanceStr, 64)
	// if err != nil {
	// 	return err
	// }

	// if balance < withdrawAmount {
	// 	return fmt.Errorf("ยอดคงเหลือไม่พอถอน!")
	// }
	

	/// check withdeaw conditionn by type

    // 1.check type turnover or turncredit
	// 2.turnover check maxWithdraw amount and check % or amount
	// 3.turncredit check type deposit or deposit+bonus and multiply with % or amount

	turntype,err := redisClient.HGet(ctx, currentPromotionKey, "TurnType").Result()
	if err !=nil {
		return nil,fmt.Errorf("ไม่พบข้อมูลโปรโมชั่น!")
	}

	requiredTurnOver,_ := getRedisKV(redisClient,currentPromotionKey,"requiredTurnover");
	requiredTurnover64,_ := strconv.ParseFloat(requiredTurnOver, 64)
	//fmt.Println("TurnType:",turntype)

	if turntype == "turnover" {

			

			

			//fmt.Printf("Total: %s \n",checkpro)
			fmt.Printf("RequiredTurnover %v \n",requiredTurnOver)
			fmt.Printf("turnoverAmount %v \n",turnoverAmount)

			if turnoverAmount < requiredTurnover64 {
				return nil,fmt.Errorf("ยอดเทิร์นโอเวอร์  น้อยกว่า %s ! ",requiredTurnOver)
			}
		
			
		} else {
			requiredCredit,_ := getRedisKV(redisClient,currentPromotionKey,"requiredCredit");
			requiredCredit64,_ := strconv.ParseFloat(requiredCredit, 64)
 
			if amount < requiredCredit64 {
				return nil,fmt.Errorf("ยอดเทิร์นเครดิต น้อยกว่า %s ! ",requiredCredit64)
			}
		}

		intStr,_ := strconv.Atoi(currentStatus)
		if turnoverAmount > requiredTurnover64 {
			intStr = 2
		}
		if err != nil || intStr == 1  {
			return nil,fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถถอนเงินเพิ่มได้ ")
		}

		WidthdrawMin,_ := getRedisKV(redisClient,currentPromotionKey,"Widthdrawmin");
		WidthdrawMin64,_ := strconv.ParseFloat(WidthdrawMin, 64)

		
		if amount <  WidthdrawMin64 {
			//withdrawAmount = WidthdrawMax64
			return nil,fmt.Errorf("ยอดถอน น้อยกว่าที่กำหนดขั้นต่ำ! ")
		}
	

		WidthdrawMax,_ := getRedisKV(redisClient,currentPromotionKey,"WidthdrawMax");
		WidthdrawMax64,_ := strconv.ParseFloat(WidthdrawMax, 64)

		
		if withdrawAmount <  WidthdrawMax64 {
			withdrawAmount = WidthdrawMax64
			//return nil,fmt.Errorf("ยอดเทิร์นโอเวอร์ ไม่ตรงเงื่อนไขการถอน! ")
		}
	


	// if err := redisClient.HSet(ctx, key, "before_balance", balance).Err(); err != nil {
	// 	return err
	// }

	// if err := redisClient.HIncrByFloat(ctx, key, "amount", -withdrawAmount).Err(); err != nil {
	// 	return err
	// }

	// // บันทึกยอดหลังการฝาก
	if err := redisClient.HSet(ctx, currentPromotionKey, "withdraw_amount", -withdrawAmount).Err(); err != nil {
		return nil,err
	}
	// if err := redisClient.HSet(ctx, key, "status", "2", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
	// 	return err
	// }
	result, err := redisClient.HGetAll(ctx, currentPromotionKey).Result()
    if err != nil {
        fmt.Println("Error fetching all values:", err)
          // ต้อง exit ถ้ามีข้อผิดพลาด
    }

    // แสดงผลลัพธ์
    //fmt.Println("All fields and values:", result)
	return result,nil
	} 
	 
	     
	// 	err := wallet.NormalTurnover(prefix,userID)
	// 	if err != nil {
	// 		fmt.Println("Error fetching all values:", err)
	// 		  // ต้อง exit ถ้ามีข้อผิดพลาด
	// 	}
	 return nil,nil
 

	// totalKey := fmt.Sprintf("%s:total_transactions", userID)
	// if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
	// 	return err
	// }
	// _, err := redisClient.HIncrByFloat(ctx, totalKey, "total_amount", -withdrawAmount).Result()
	// return err

	
	
}
 
// ฟังก์ชันการเล่นเกม
func playGame(redisClient *redis.Client, userID string, proID string, gameAmount float64) (map[string]float64, error) {
	lockKey := fmt.Sprintf("%s:lock", userID)

	if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
		return nil, fmt.Errorf("could not acquire lock: %v", err)
	}

	if locked, err := handler.AcquireLock(redisClient, lockKey,ttl); err != nil || !locked {
		return nil, fmt.Errorf("could not acquire lock: %v", err)
	}
	defer handler.ReleaseLock(redisClient, lockKey)

	key := fmt.Sprintf("%s:%s", userID, proID)

	if err := checkKeyType(redisClient, key, "hash"); err != nil {
		return nil, err
	}

	currentStatus, err := redisClient.HGet(ctx, key, "status").Result()
	intStr, _ := strconv.Atoi(currentStatus)
	if err != nil || intStr < 1 {
		if intStr == 0 {
			return nil, fmt.Errorf("โปรโมชั่นไม่อยู่ในสถานะ ใช้งาน!")
		}
		return nil, fmt.Errorf("โปรโมชั่นไม่สามารถใช้งานได้!")
	}

	// บันทึกยอดก่อนเล่นเกม
	balanceStr, err := redisClient.HGet(ctx, key, "balance").Result()
	if err != nil {
		return nil, err
	}
	balance, _ := strconv.ParseFloat(balanceStr, 64)

	if balance < gameAmount {
		return nil, fmt.Errorf("ยอดคงเหลือไม่พอเล่นเกม!")
	}

	// ดำเนินการเล่นเกม (การคำนวณผลที่ได้)
	netGain := gameAmount * 2 // ยกตัวอย่างการชนะ (ขึ้นอยู่กับกฎเกมที่คุณสร้าง)

	// อัปเดตยอด
	if err := redisClient.HIncrByFloat(ctx, key, "balance", netGain).Err(); err != nil {
		return nil, err
	}
	
	// ตรวจสอบยอดหลังเล่นเกม
	newBalance := balance + netGain // คำนวณยอดใหม่หลังจากเล่นเกม

	// หากยอดคงเหลือเหลือ 0 ให้ปรับสถานะโปรโมชั่นเป็น 2
	if newBalance == 0 {
		if err := redisClient.HSet(ctx, key, "status", "2", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
			return nil, err
		}
	} else {
		// หากไม่เป็น 0 ก็ให้ปรับสถานะเป็น 1 (กำลังใช้งาน)
		if err := redisClient.HSet(ctx, key, "status", "1", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
			return nil, err
		}
	}

	// คืนค่าผลลัพธ์
	result := map[string]float64{
		"before_balance": balance,
		"game_amount":    gameAmount,
		"after_balance":  newBalance,
	}

	return result, nil
}

func getTransactionsByUserIDAndDate(rdb *redis.Client, userID string, startDate, endDate time.Time) ([]string, error) {
   
	zKey := fmt.Sprintf("user_transactions:%s", userID)

    startScore := float64(startDate.Unix())
    endScore := float64(endDate.Unix())

    // ค้นหา transaction IDs ในช่วงวันที่กำหนด
    transactionIDs, err := rdb.ZRangeByScore(ctx, zKey, &redis.ZRangeBy{
        Min: fmt.Sprintf("%f", startScore),
        Max: fmt.Sprintf("%f", endScore),
    }).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to get transactions: %v", err)
    }

    return transactionIDs, nil
}
func getTransactionsBank(rdb *redis.Client, userID string, uid string) (map[string]string, error) {
    // สร้าง key ของธุรกรรม
    transactionKey := fmt.Sprintf("bank_statement:%s:%s", userID, uid)
	fmt.Printf("transactionKey: %s \n",transactionKey)
    // ดึงข้อมูลธุรกรรมทั้งหมดจาก Redis
    transactionData, err := rdb.HGetAll(ctx, transactionKey).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to get transactions from Redis: %v", err)
    }

    return transactionData, nil
}
// ฟังก์ชันเพื่อดูยอดธุรกรรมทั้งหมด
func getTotalTransactions(redisClient *redis.Client, userID string) (float64, error) {
	totalKey := fmt.Sprintf("%s:total_transactions", userID)
	if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
		return 0, err
	}

	totalAmount, err := redisClient.HGet(ctx, totalKey, "total_amount").Result()
	if err != nil {
		return 0, err
	}

	var amount float64
	fmt.Sscanf(totalAmount, "%f", &amount)

	return amount, nil
}
// ฟังก์ชันเพื่อดูสถานะโปรโมชั่น
func getPromotionStatus(redisClient *redis.Client, userID string, proID string) (map[string]string, error) {
	key := fmt.Sprintf("%s:%s", userID, proID)

	// ตรวจสอบประเภทของคีย์
	if err := checkKeyType(redisClient, key, "hash"); err != nil {
		return nil, err
	}

	// ดึงค่าที่เกี่ยวข้องจาก Redis
	status, err := redisClient.HGet(ctx, key, "status").Result()
	if err != nil {
		if err == redis.Nil {
			status = "not found"
		} else {
			return nil, err
		}
	}

	timestamp, err := redisClient.HGet(ctx, key, "timestamp").Result()
	if err != nil {
		if err == redis.Nil {
			timestamp = "not found"
		} else {
			return nil, err
		}
	}

	return map[string]string{
		"status":    status,
		"timestamp": timestamp,
	}, nil
}
func getPromotionStatusForUser(redisClient *redis.Client, userID string) (map[string]string, error) {
    // คีย์สำหรับโปรโมชั่นปัจจุบัน
    currentPromotionKey := fmt.Sprintf("%s:current_promotion", userID)

    // ตรวจสอบว่ามีโปรโมชั่นปัจจุบันอยู่หรือไม่
    if err := checkKeyType(redisClient, currentPromotionKey, "hash"); err != nil {
        return nil, err
    }

    // ดึงข้อมูลโปรโมชั่นปัจจุบัน
    promotionData, err := redisClient.HGetAll(ctx, currentPromotionKey).Result()
    if err != nil {
        return nil, err
    }

	fmt.Printf("\n PromotionData: %+v \n",promotionData)
    // ตรวจสอบว่าโปรโมชั่นมีข้อมูลหรือไม่
	//fmt.Println("PromotionData:",len(promotionData))
    if len(promotionData) == 0 {
        return nil, fmt.Errorf("no current promotion found for user %s", userID)
    }

    // แสดงผลข้อมูลโปรโมชั่น
    promotionData["is_expired"] = "false" // ตัวอย่างค่าที่กำลังคงอยู่

    return promotionData, nil
}
func getPromotionStatusForAllUsers(redisClient *redis.Client, proID string) (map[string]map[string]string, error) {
	promotionStatuses := make(map[string]map[string]string)

	// ค้นหาทุกคีย์ของผู้ใช้ที่มีโปรโมชั่น
	userKeys, err := redisClient.Keys(ctx, "*:*").Result()
	if err != nil {
		return nil, err
	}

	for _, userKey := range userKeys {
		parts := strings.Split(userKey, ":")
		if len(parts) != 2 {
			continue // ป้องกันการแตกในกรณีที่คีย์ไม่มีรูปแบบที่ถูกต้อง
		}

		userID := parts[0]
		currentProID := parts[1]

		if currentProID != proID {
			continue // ข้ามโปรโมชั่นที่ไม่ตรงกัน
		}

		if err := checkKeyType(redisClient, userKey, "hash"); err != nil {
			return nil, err
		}

		status, err := redisClient.HGet(ctx, userKey, "status").Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		if err == redis.Nil {
			status = "0" // กำหนดเป็น 0 หากไม่พบสถานะ
		}

		timestamp, err := redisClient.HGet(ctx, userKey, "timestamp").Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		if err == redis.Nil {
			timestamp = "0" // กำหนดเป็น 0 หากไม่พบ timestamp
		}

		beforeBalanceStr, err := redisClient.HGet(ctx, userKey, "before_balance").Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		var beforeBalance float64
		if err == redis.Nil {
			beforeBalance = 0 // กำหนดเป็น 0 หากไม่พบยอดก่อน
		} else {
			beforeBalance, _ = strconv.ParseFloat(beforeBalanceStr, 64)
		}

		totalKey := fmt.Sprintf("%s:total_transactions", userID)
		if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
			fmt.Println("Error:", err)
		}

		// ดึงยอดรวม
		totalAmountStr, err := redisClient.HGet(ctx, totalKey, "total_amount").Result()
		var totalAmount float64
		if err != nil {
			totalAmount = 0 // กำหนดเป็น 0 หากไม่พบยอดรวม
			fmt.Println("Error:", err)
		} else {
			totalAmount, _ = strconv.ParseFloat(totalAmountStr, 64)
		}

		// คำนวณ after_balance
		afterBalance := beforeBalance + totalAmount

		// บันทึกสถานะของโปรโมชั่นของผู้ใช้
		promotionStatuses[userID] = map[string]string{
			"status":         status,
			"timestamp":      timestamp,
			"before_balance": fmt.Sprintf("%.2f", beforeBalance),
			"amount":         fmt.Sprintf("%.2f", totalAmount),
			"after_balance":  fmt.Sprintf("%.2f", afterBalance),
		}
	}

	return promotionStatuses, nil
}
// Endpoint สำหรับเลือกรับโปรโมชั่น
func SelectPromotion(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง UserID จาก context
		userID := c.Locals("ID")
		//userID := c.Locals("ID")
		//proID := c.Params("proID")
		
		db, err := handler.GetDBFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
			}
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}

		body := PBody{}

		if err := c.BodyParser(&body); err != nil {
			//fmt.Println("Error parsing body:", err.Error())
			response := fiber.Map{
				"Status":  false,
				"Message": err.Error(),
			}
			return c.JSON(response)
		}
		var users models.Users
		err = db.Debug().Select("balance").Where("id= ?", userID).Find(&users).Error
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
			}
		pro_setting,err :=  handler.GetProdetail(db,body.ProID)
		fmt.Printf("Prosetting: %+v\n", pro_setting) 
		if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Status": false,
			"Message": err.Error(),
		})
		}


		pro_setting["Userid"] = userIDStr;
		balance,_ := users.Balance.Float64()
		if err = selectPromotion(redisClient, fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr),balance, pro_setting); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
		} else {
			uid,_ := getRedisKV(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr),"uid")
			updates := map[string]interface{}{
				"ProID": uid,
				"ProStatus": "0",
				//"Turnover": users.Turnover,
				//"ProStatus": users.ProStatus,
			}
			_err := repository.UpdateUserFields(db, userID.(int), updates) // อัปเดตยูสเซอร์ที่มี ID = 1
			if _err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Status": false,
				"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
				"Data": fiber.Map{ 
					"id": -1,
				}})
			}
		}

		

		return c.JSON(fiber.Map{
			"Status": true,
			//"Data": response,
			"Message": "Selected Promotion successfully",
		})
	}
}
func ClearRedisKey(redisClient *redis.Client, userID string) error {
	// คีย์สำหรับลบข้อมูลโปรโมชั่น
	promotionListKey := fmt.Sprintf("%s:current_promotion", userID)
	if err := redisClient.Del(ctx, promotionListKey).Err(); err != nil {
		return fmt.Errorf("could not delete promotions key: %v", err)
	}

	// ลบโปรโมชั่นแต่ละรายการ
	oldPromotionKeys, err := redisClient.LRange(ctx, promotionListKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, proID := range oldPromotionKeys {
		promotionKey := fmt.Sprintf("%s:%s", userID, proID)
		if err := redisClient.Del(ctx, promotionKey).Err(); err != nil {
			return fmt.Errorf("could not delete promotion key %s: %v", promotionKey, err)
		}
	}

	// สามารถลบคีย์อื่น ๆ ที่ต้องการที่เกี่ยวข้องกับ userID ได้ที่นี่
	totalKey := fmt.Sprintf("%s:total_transactions", userID)
	if err := redisClient.Del(ctx, totalKey).Err(); err != nil {
		return fmt.Errorf("could not delete total transactions key: %v", err)
	}

	return nil
}

func GetTransactionHandler(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
    // สมมุติว่าเรามีค่า userID และ uid
    //userID := c.Query("userID") // หรือจาก params / body ตามที่คุณต้องการ
    //uid := c.Query("uid") // หรือจาก params / body ตามที่คุณต้องการ
	

	userID := c.Locals("ID")
	prefix := c.Locals("prefix")
	uid := c.Locals("uid")
	var userIDStr string
	switch v := userID.(type) {
	case string:
		userIDStr = v
	case int:
		userIDStr = fmt.Sprintf("%d", v) 
	}
	
	userid := fmt.Sprintf("%s%s", prefix, userIDStr) // Prefix + userID
	uidStr := fmt.Sprintf("%s", uid) 

	//userid := fmt.Sprintf("%s%s", prefix.(string), userID.(string))
    key := fmt.Sprintf("bank_statement:%s:%s", userid, uidStr)
    fmt.Printf("Key: %s \n", key)
    // เรียกฟังก์ชัน getTransaction
    transactionData, err := getTransaction(redisClient, userid, uidStr)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status": false,
            "message": err.Error(),
            "data": nil,
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "Status": true,
        "Data": transactionData,
    })
}
}

func getTransaction(rdb *redis.Client, userID string, uid string) (map[string]string, error) {
    // สร้าง key ของธุรกรรม
    key := fmt.Sprintf("bank_statement:%s:%s", userID, uid)

    // ดึงข้อมูลธุรกรรมทั้งหมดจาก Redis
    transactionData, err := rdb.HGetAll(ctx, key).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to get transaction from Redis: %v", err)
    }

    return transactionData, nil
}

func GetAllPromotion(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {

		body := PBody{}

		if err := c.BodyParser(&body); err != nil {
			//fmt.Println("Error parsing body:", err.Error())
			response := fiber.Map{
				"Status":  false,
				"Message": err.Error(),
			}
			return c.JSON(response)
		}
		response,err := getPromotionStatusForAllUsers(redisClient, body.ProID);
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"Status": true,
			"Data": response,
			"Message": "Promotion selected successfully",
		})

	 }
}
func GetUserTotalPromotion(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("ID")
		//userID := c.Locals("ID")
		//proID := c.Params("proID")
		
		//db, _err := handler.GetDBFromContext(c)
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}
		response,err := getTotalTransactions(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"Status": true,
			"Data": response,
			"Message": "Promotion selected successfully",
		})
	}
}
func GetPromotionsUsersID(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("ID")
		
		//db, _err := handler.GetDBFromContext(c)
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}
		response,_ := getPromotionStatusForUser(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr))
		// if err != nil {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"Status": false,
		// 		"Message": err.Error(),
		// 	})
		// }

		//fmt.Printf("Response: %+v \n",response)
		return c.JSON(fiber.Map{
			"Status": true,
			//"Data": fiber.Map{
			//	"ProId": body.ProID,
			"Data":response,
			//},
			"Message": "Promotion selected successfully",
		})
	
	}
}
func GetPromotionStatus(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("ID")
		//userID := c.Locals("ID")
		//proID := c.Params("proID")
		body := PBody{}

		if err := c.BodyParser(&body); err != nil {
			fmt.Println("Error parsing body:", err.Error())
			response := fiber.Map{
				"Status":  false,
				"Message": err.Error(),
			}
			return c.JSON(response)
		}


		//db, _err := handler.GetDBFromContext(c)
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}
		response,err := getPromotionStatus(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr),body.ProID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"Status": true,
			//"Data": fiber.Map{
			//	"ProId": body.ProID,
			"Data":response,
			//},
			"Message": "Promotion selected successfully",
		})
	}
}
func Withdraw(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("ID")
		prefix := c.Locals("prefix")
	 	var percentValue decimal.Decimal
		var percentStr = ""

		body := PBody{}
		var users models.Users
		if err := c.BodyParser(&body); err != nil {
			fmt.Println("Error parsing body:", err.Error())
			response := fiber.Map{
				"Status":  false,
				"Message": err.Error(),
			}
			return c.JSON(response)
		}
		
		BankStatement := new(models.BankStatement)

		db, _err := handler.GetDBFromContext(c)
		if _err != nil {
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": _err,
				"Data": fiber.Map{ 
					"id": -1,
				}})
		}
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}
		BankStatement.Userid = userID.(int)
		BankStatement.Walletid = userID.(int)

		//fmt.Printf(" body : %+v \n",body)
		// amount, err := strconv.ParseFloat(body.Amount, 64)
		// if err != nil {
		// 	fmt.Printf("Error converting amount: %v\n", err)
		// 	return err
		// }



		if err_ := db.Where("id  = ? ", userIDStr).First(&users).Error; err_ != nil {
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": err_,
				"Data": fiber.Map{ 
					"id": -1,
				}})
		}

		var totalTurnover decimal.Decimal
		var lastdate models.BankStatement
		if err := db.Debug().Model(&models.BankStatement{}).
		Where("userid = ? and status='verified'",userIDStr).
		Select("updatedAt").
		Order("id desc").
		First(&lastdate).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Message": "ไม่สามารถคำนวณยอดเทิร์นได้ !",
				"Status": false,
				"Data": "เกิดข้อผิดพลาด!",
			})
		}

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

		totalTurnover64,_ := totalTurnover.Float64()

		resultx,err := withdraw(redisClient,prefix.(string),userIDStr,body.ProID,body.Amount,totalTurnover64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
		}
		deposit := decimal.NewFromFloat(body.Amount)

		if users.Balance.LessThan(deposit.Abs()) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": fmt.Sprintf("ยอดถอนมากกว่า  %v ยอดคงเหลือ %v %v !",deposit,users.Balance,users.Currency),
				"Data": fiber.Map{
					"id": -1,
				}})
		}
		BankStatement.Beforebalance = users.Balance
		BankStatement.Transactionamount =  deposit
	    BankStatement.Balance = users.Balance.Add(deposit)
		
		if resultx == nil  {

			BankStatement.Beforebalance = users.Balance
			BankStatement.Transactionamount =  deposit
		    BankStatement.Balance = users.Balance.Add(deposit)
			
			

			BankStatement.Turnover = totalTurnover

			var count_trans float64

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

			//totalTurnover64,_ := totalTurnover.Float64()

			

			userid := fmt.Sprintf("%s%d", users.Prefix, users.ID)
			key := fmt.Sprintf("bank_statement:%s:%s", userid, users.Uid)
			transaction_amount, err := getRedisKV(redisClient, key, "transaction_amount")

			if strings.Contains(users.MinTurnoverDef, "%") {
				percentStr = strings.TrimSuffix(users.MinTurnoverDef, "%")
				//fmt.Printf(" MinturnoverDef : %s %\n",percentStr)
				// แปลงเป็น float64
				percentValue, _ = decimal.NewFromString(percentStr)
		 
		
			// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
				if count_trans > 0 {
					percentValue = percentValue.Mul(decimal.NewFromFloat(count_trans)).Div(decimal.NewFromInt(100))
				} else {
					percentValue = percentValue.Div(decimal.NewFromInt(100))	
				}
				 
			} else {
				if count_trans > 0 {
				 	percentStr = users.MinTurnoverDef  
					percentValue, _ = decimal.NewFromString(percentStr)
					percentValue = percentValue.Mul(decimal.NewFromFloat(count_trans))
				} else {
					percentStr = users.MinTurnoverDef
					percentValue, _ = decimal.NewFromString(percentStr)
				}
			  
			// แปลงเป็น float64
				
		  
			}
			
			//fmt.Printf(" Minturnover : %s \n",percentStr)
			// แปลงเป็น float64
			//percentValue, _ := decimal.NewFromString(percentStr)
		 
		
			
			//err := redisClient.HSet(ctx, key, "transaction_amount", transactionsub.TransactionAmount.String()).Err()
			minresult := users.LastDeposit.Mul(percentValue)
			//total_deposit_amount, _ := getRedisKV(redisClient, key, "total_deposit_amount")
			//fmt.Printf(" total_deposit_amount: %v \n",total_deposit_amount)
			
			if err != nil {
				return fmt.Errorf("could not get transaction_amount: %v", err)
			}
			//fmt.Printf("transaction_amount: %v \n", transaction_amount)
			
			if transaction_amount == "" {
				minresult = users.Balance.Mul(percentValue)
				//fmt.Printf(" minresult: %v \n",minresult)
			} else {
				percentStr = users.MinTurnoverDef
				percentValue, _ = decimal.NewFromString(percentStr)
				minresult = users.LastDeposit.Mul(percentValue)
			}
			fmt.Printf(" count_trans : %v \n",count_trans)
			// แปลงเปอร์เซ็นต์เป็นค่าทศนิยม
			//percentValue = percentValue.Div(decimal.NewFromInt(100))//.Add(decimal.NewFromInt(1))
			fmt.Printf(" MinturnoverDef : %s \n",percentStr)
			// ใช้ในสูตรคำนวณ
			//baseValue := 500.0
			fmt.Printf(" PercentValue: %v \n",percentValue)
			//getTransactionsBank(redisClient,)
			fmt.Printf(" minresult: %v \n",minresult)
			
			if totalTurnover.LessThanOrEqual(minresult)  || totalTurnover.IsZero() {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Status": false,
					"Message": fmt.Sprintf("ยอดเทิร์นโอเวอร์น้อยกว่ายอดเทิร์นโอเวอร์ขั้นต่ำ %v %v = %v ของยอดฝากล่าสุด !",users.MinTurnoverDef,users.Currency,minresult.StringFixed(2)),
					"Data": fiber.Map{
						"id": -1,
					}})
			} 
			BankStatement.Turnover = totalTurnover
		} 
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		BankStatement.Userid = userID.(int)
		BankStatement.Walletid = userID.(int)
		BankStatement.Beforebalance = users.Balance
		BankStatement.ProStatus = users.ProStatus
		// ถ้ามีโปรโมชั่นให้ปรับเป็น 0
		BankStatement.Bankname = users.Bankname
		BankStatement.Accountno = users.Banknumber
		BankStatement.Transactionamount = deposit
		
	
		// บันทึกรายกา
		request := &wallet.PayInRequest{
			Ref:            users.Username,
			BankAccountName: users.Fullname,
			Amount:         deposit.Abs().String(),
			BankCode:       users.Bankname,
			BankAccountNo:  users.Banknumber,
			//MerchantURL:    "https://www.xn--9-twft5c6ayhzf2bxa.com/",
		}
		
		var result wallet.PayInResponse
		result, err = wallet.Payout(request) // เรียกใช้ฟังก์ชัน paying พร้อมส่ง request
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
		BankStatement.StatementType =  "Withdraw"
		
		// return c.JSON(fiber.Map{
		// 	"Status": true,
		// 	"Message": "ถอนเงินสำเร็จ",
		// 	"Data": fiber.Map{
		// 		"id": BankStatement.ID,
		// 		"beforebalance": BankStatement.Beforebalance,
		// 		"transactionamount": BankStatement.Transactionamount,
		// 		"balance": BankStatement.Balance,
		// 		"method": BankStatement.StatementType,
		// 	},
		// })


		if err := tx.Create(&BankStatement).Error; err != nil {
			tx.Rollback()
			return c.JSON(fiber.Map{
				"Status": false,
				"Message": "ไม่สามารถบันทึกรายการได้",
				"Data": fiber.Map{"id": -1},
			})
		}

		updates := map[string]interface{}{
			"Balance": decimal.Zero,
			"LastWithdraw": deposit,
		}
	 
		if users.ProStatus != "" {
			updates["ProStatus"] = ""
			if err := tx.Model(&users).Updates(updates).Error; err != nil {
				tx.Rollback()
				return c.JSON(fiber.Map{
					"Status": false,
					"Message": "ไม่สามารถอัพเดทข้อมูลผู้ใช้ได้",
					"Data": fiber.Map{"id": -1},
				})
			}
		} else {
			updates["Balance"] = BankStatement.Balance
			if err := tx.Model(&users).Updates(updates).Error; err != nil {
				tx.Rollback()
				return c.JSON(fiber.Map{
					"Status": false,
					"Message": "ไม่สามารถอัพเดทข้อมูลผู้ใช้ได้",
					"Data": fiber.Map{"id": -1},
				})
			}
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

	}}



func Deposit(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("ID")
		//userID := c.Locals("ID")
		//proID := c.Params("proID")
		//body := PBody{}
		BankStatement := new(models.BankStatement)
		if err := c.BodyParser(&BankStatement); err != nil {
			fmt.Println("Error parsing body:", err.Error())
			response := fiber.Map{
				"Status":  false,
				"Message": err.Error(),
			}
			return c.JSON(response)
		}


		db, _err := handler.GetDBFromContext(c)
		if _err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": _err.Error(),
			})
		}
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}

		// amount, err := strconv.ParseFloat(body.Amount, 64)
		// if err != nil {
		// 	fmt.Printf("Error converting amount: %v\n", err)
		// 	return err
		// }

		promotion,_ := getPromotionStatusForUser(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr))
		
		//fmt.Printf("promotion: %+v\n", promotion)

		if len(promotion) > 0 && promotion["status"]!="2"{
			
			transactionamount,_ := BankStatement.Transactionamount.Float64()

			result,err := deposit(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr),"current_promotion",transactionamount)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Status": false,
					"Message": err.Error(),
				})
			}
			
			// แสดงผลลัพธ์
			//fmt.Println("All fields and values:", result)
			 Result,err := wallet.Deposit(db,result,BankStatement); 
			 if err !=nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"Status": false,
					"Message": err.Error(),
				})
			}
	
			return c.Status(200).JSON(fiber.Map{
				"Status": true,
				"Data":  Result,
			})
			} else {

				Result,err := wallet.XDeposit(db,userIDStr,BankStatement); 
				if err !=nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"Status": false,
						"Message": err.Error(),
					})
				}
		
				return c.Status(200).JSON(fiber.Map{
					"Status": true,
					"Data":  Result,
				})
			}
		// if err != nil {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"Status": false,
		// 		"Message": fmt.Sprintf("1210 %s",err.Error()),
		// 	})
		// }  
		
		
	}
}
func Playgame(redisClient *redis.Client) fiber.Handler {
	 return func(c *fiber.Ctx) error {
		userID := c.Locals("ID")
		//userID := c.Locals("ID")
		//proID := c.Params("proID")
		body := PBody{}

		if err := c.BodyParser(&body); err != nil {
			fmt.Println("Error parsing body:", err.Error())
			response := fiber.Map{
				"Status":  false,
				"Message": err.Error(),
			}
			return c.JSON(response)
		}
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}

		// amount, err := strconv.ParseFloat(body.Amount, 64)
		// if err != nil {
		// 	fmt.Printf("Error converting amount: %v\n", err)
		// 	return err
		// }
		response,err := playGame(redisClient,userIDStr,body.ProID,body.Amount)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"Status": true,
			//"Data": fiber.Map{
			//	"ProId": body.ProID,
			"Data":response,
			//},
			"Message": "Promotion Deposit successfully",
		})
	 }
}
func ClearData(redisClient *redis.Client) fiber.Handler {
		return func(c *fiber.Ctx) error {

		userID := c.Locals("ID")
		prefix := c.Locals("prefix")
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}


		
		db, _err := handler.GetDBFromContext(c)
		if _err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Status": false,
				"Message": _err.Error(),
			})
		}
		key := fmt.Sprintf("%s%d", prefix, userID)
		currentPromotionKey := fmt.Sprintf("%s:current_promotion", key)
	
		//currentPromotionKey := fmt.Sprintf("%s:current_promotion", userID)
		status,_ := getRedisKV(redisClient,currentPromotionKey,"status");
		if status == "0" || status == "2" {
		// สมมุติว่าคุณมี redisClient และ userID ของผู้ใช้
		err := ClearRedisKey(redisClient, fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr))
 
		
		if err != nil {
			fmt.Println("Error clearing Redis keys:", err)
		} else {
			fmt.Println("Successfully cleared Redis keys for user:", userID)
		}

		

		updates := map[string]interface{}{
			"ProID": "",
			"ProStatus": "0",
			//"Turnover": users.Turnover,
			//"ProStatus": users.ProStatus,
		}
		_err := repository.UpdateUserFields(db, userID.(int), updates) // อัปเดตยูสเซอร์ที่มี ID = 1
		if _err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status": false,
			"Message":  "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
		}
		
		return c.JSON(fiber.Map{
			"Status": true,
			//"Data": fiber.Map{
			//	"ProId": body.ProID,
			//"Data":response,
			//},
			"Message": "Clear successfully",
		})
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status": false,
			"Message":  "สถานะโปรโมชั่น ไม่สามารถยกเลิกโปรโมชั่นได้!",
			"Data": fiber.Map{ 
				"id": -1,
			}})
		}
	
	}
}

 