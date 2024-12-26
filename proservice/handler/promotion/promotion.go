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

 
// ฟังก์ชันตรวจสอบการเชื่อมต่อ Redis
func checkRedisConnection(redisClient *redis.Client) error {
	_, err := redisClient.Ping(ctx).Result()
	return err
}
 
func acquireLock(redisClient *redis.Client, lockKey string) (bool, error) {
	// ใช้ `HSet` เพื่อสร้างล็อก
	locked, err := redisClient.HSetNX(ctx, lockKey, "locked", true).Result()
	if err != nil {
		return false, err
	}
	if locked {
		// ตั้งเวลาหมดอายุ
		go func() {
			time.Sleep(10 * time.Second)
			releaseLock(redisClient, lockKey)
		}()
	}
	return locked, nil
}

// ฟังก์ชันปลดล็อก
func releaseLock(redisClient *redis.Client, lockKey string) error {
	// ลบฟิลด์ที่เป็นล็อก
	_, err := redisClient.HDel(ctx, lockKey, "locked").Result()
	return err
}

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

// ฟังก์ชันเลือกโปรโมชั่น
// func selectPromotion(redisClient *redis.Client, userID string, proID string) error {
// 	lockKey := fmt.Sprintf("%s:lock", userID)

// 	// พยายามล็อก
// 	if locked, err := acquireLock(redisClient, lockKey); err != nil || !locked {
// 		return fmt.Errorf("could not acquire lock: %v", err)
// 	}
// 	defer releaseLock(redisClient, lockKey)

// 	key := fmt.Sprintf("%s:%s", userID, proID)

// 	oldPromotionKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s:*", userID)).Result()
// 	if err != nil {
// 		return err
// 	}

// 	var previousPromotionStatus string

// 	for _, oldKey := range oldPromotionKeys {
// 		if err := checkKeyType(redisClient, oldKey, "hash"); err != nil {
// 			return err
// 		}

// 		oldStatus, err := redisClient.HGet(ctx, oldKey, "status").Result()
// 		if err != nil {
// 			if err == redis.Nil {
// 				oldStatus = "0" // หรือค่าเริ่มต้นอื่น ๆ ที่อาจหมายถึงไม่มีกระบวนการ
// 			} else {
// 				return err
// 			}
// 		}

// 		if oldStatus == "1" {
// 			return fmt.Errorf("ไม่สามารถเลือกโปรโมชั่นใหม่ได้ โปรโมชั่นปัจจุบันทำงาน")
// 		}

// 		previousPromotionStatus = oldStatus // บันทึกสถานะโปรโมชั่นก่อนหน้า
// 	}

// 	// ตรวจสอบยอดเงินก่อนหน้านี้
// 	totalKey := fmt.Sprintf("%s:total_transactions", userID)
// 	if err := ensureTotalTransactionsKey(redisClient, userID); err != nil {
// 		return err
// 	}
// 	if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
// 		return err
// 	}

// 	balance, err := redisClient.HGet(ctx, totalKey, "total_amount").Result()
// 	if err != nil && err != redis.Nil {
// 		return err
// 	}
	
	 

// 	balanceStr, err := strconv.ParseFloat(balance, 64)
// 	if err != nil {
// 		return fmt.Errorf("ไม่สามารถแปลงยอดคงเหลือเป็นตัวเลขได้: %v", err)
// 	}

// 	if balanceStr > 0 {
// 		return fmt.Errorf("ยอดคงเหลือมากกว่าศูนย์")
// 	}

// 	// ระบบจัดการสถานะ
// 	if previousPromotionStatus == "0" {
// 		// เปลี่ยนสถานะโปรโมชั่นก่อนหน้าเป็น -1
// 		for _, oldKey := range oldPromotionKeys {
// 			if err := redisClient.HSet(ctx, oldKey, "status", "-1").Err(); err != nil {
// 				return err
// 			}
// 		}
// 	} else if previousPromotionStatus != "-1" {
// 		// ถ้าสถานะก่อนหน้าเป็นอย่างอื่นที่ไม่ใช่ 2 และไม่ใช่ 0 ไม่มีการเปลี่ยนแปลงใด ๆ
// 		return fmt.Errorf("โปรโมชั่นปัจจุบันสถานะ  %s", previousPromotionStatus)
// 	}

// 	// กำหนดสถานะโปรโมชั่นใหม่ พร้อม timestamp
// 	if err := redisClient.HSet(ctx, key, "status", "0", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }


// ฟังก์ชันในการสร้าง UID
func generateUID(userID, proID, timestamp string) string {
    hash := sha256.New()
    hash.Write([]byte(userID + proID + timestamp))
    return hex.EncodeToString(hash.Sum(nil))
}
func selectPromotion(redisClient *redis.Client, userID string,balance float64, pro_setting map[string]interface{}) error {
    lockKey := fmt.Sprintf("%s:lock", userID)

    // พยายามล็อก
    if locked, err := acquireLock(redisClient, lockKey); err != nil || !locked {
        return fmt.Errorf("could not acquire lock: %v", err)
    }
    defer releaseLock(redisClient, lockKey)

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
		return fmt.Errorf("โปรโมชั่นเดิม ยังไม่สิ้นสุด is %s", promotionStatus)
	}

    // ดึง timestamp ปัจจุบัน
    now := time.Now()
    nowRFC3339 := now.Format(time.RFC3339)

	


    // สร้าง UID โดยการแฮช userID, proID และ timestamp
    uid := generateUID(userID, fmt.Sprintf("%d", pro_setting["Id"]), nowRFC3339)

	// response["Id"] = promotion.ID
	// response["Type"] = ProItem.ProType.Type
	// response["count"] = ProItem.UsageLimit
	// response["MinTurnover"] = promotion.MinSpend
	// response["Formular"] = promotion.Example
	// response["Name"] = promotion.Name
	// response["TurnType"]=promotion.TurnType
	// response["Example"]= promotion.Example
	// if ProItem.ProType.Type == "weekly" {
	// 	response["Week"] = ProItem.ProType.DaysOfWeek
	// } 


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
    "maxDept":       fmt.Sprintf("%s", pro_setting["maxDept"]), // แปลงเป็น string หากจำเป็น
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

	//fmt.Printf("Response: %+v \n",response)
	// response := map[string]interface{}{
	// 	"proID":         fmt.Sprintf("%d", pro_setting["Id"]),
	// 	"status":        "0",
	// 	"timestamp":     nowRFC3339,
	// 	"uid":          uid,
	// 	"Id":           fmt.Sprintf("%d", pro_setting["Id"]), // แปลงเป็น string
	// 	"Type":         fmt.Sprintf("%s", pro_setting["Type"]),
	// 	"count":        fmt.Sprintf("%d", pro_setting["UsageLimit"]), // แปลงเป็น string
	// 	"MinTurnover":  fmt.Sprintf("%s", pro_setting["MinTurnover"]),
	// 	"MinCredit":    fmt.Sprintf("%s", pro_setting["MinCredit"]), // แปลงเป็น string หากต้องการ
	// 	"Example":      fmt.Sprintf("%s", pro_setting["Formular"]),
	// 	"Name":         fmt.Sprintf("%s", pro_setting["Name"]),
	// 	"TurnType":     fmt.Sprintf("%s", pro_setting["TurnType"]),
	// 	"Week":         fmt.Sprintf("%s", pro_setting["Week"]),
	
	// 	// ข้อมูลใหม่ที่เพิ่มเข้ามาจาก pro_setting
	// 	"minDept":       fmt.Sprintf("%s", pro_setting["MinDept"]), // แปลงเป็น string หากจำเป็น
	// 	"maxDept":       fmt.Sprintf("%s", pro_setting["MaxDiscount"]), // แปลงเป็น string หากจำเป็น
	// 	"Widthdrawmax":  fmt.Sprintf("%s", pro_setting["MaxSpend"]), // แปลงเป็น string
	// 	"Widthdrawmin":  fmt.Sprintf("%s", pro_setting["Widthdrawmin"]), // แปลงเป็น string
	// 	"MinSpendType":  fmt.Sprintf("%s", pro_setting["MinSpendType"]), // แปลงเป็น string
	// 	"MinCreditType": fmt.Sprintf("%s", pro_setting["MinCreditType"]), // แปลงเป็น string
	// 	"MaxWithdrawType": fmt.Sprintf("%s", pro_setting["MaxWithdrawType"]), // แปลงเป็น string
	// 	// แปลงเป็น string
	// 	"Zerobalance":   fmt.Sprintf("%s", pro_setting["Zerobalance"]), // แปลงเป็น string
	// 	"CreatedAt":     pro_setting["CreatedAt"], // ควรตรวจสอบประเภทให้แน่ใจว่าเป็นประเภทที่เหมาะสม
	// 	"MaxUse":        fmt.Sprintf("%d", pro_setting["UsageLimit"]), // แปลงเป็น string
	// }
     


    // บันทึกโปรโมชั่นใหม่โดยเขียนทับโปรโมชั่นเก่า
    if err := redisClient.HSet(ctx, currentPromotionKey, response).Err(); err != nil {
        return err
    }

    return nil
}

func deposit(redisClient *redis.Client, userID string, proID string, depositAmount float64) (map[string]string,error) {
    
	lockKey := fmt.Sprintf("%s:lock", userID)

    // ตรวจสอบการล็อก
    if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
        return nil,fmt.Errorf("could not acquire lock: %v", err)
    }
	defer releaseLock(redisClient, lockKey)

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

    _, err = redisClient.HIncrByFloat(ctx, totalKey, "total_amount", depositAmount).Result()
    if err != nil {
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
 
// func deposit(redisClient *redis.Client, userID string, proID string, depositAmount float64) error {
// 	lockKey := fmt.Sprintf("%s:lock", userID)

// 	if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
// 		return fmt.Errorf("could not acquire lock: %v", err)
// 	}

// 	key := fmt.Sprintf("%s:%s", userID, proID)

// 	if err := ensureTotalTransactionsKey(redisClient, userID); err != nil {
// 		return err
// 	}

// 	if err := checkKeyType(redisClient, key, "hash"); err != nil {
// 		return err
// 	}

// 	currentStatus, err := redisClient.HGet(ctx, key, "status").Result()
// 	if err != nil {
// 		return err
// 	}
// 	intStr, _ := strconv.Atoi(currentStatus)
// 	if intStr > 0 {
// 		return fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้")
// 	}

// 	balanceStr, err := redisClient.HGet(ctx, key, "balance").Result()
// 	if err != nil {
// 		return err
// 	}

// 	balance, err := strconv.ParseFloat(balanceStr, 64)
// 	if err != nil {
// 		return err
// 	}

// 	// บันทึกยอดก่อนการฝาก
// 	if err := redisClient.HSet(ctx, key, "before_balance", balance).Err(); err != nil {
// 		return err
// 	}

// 	// เพิ่มจำนวนเงินฝาก
// 	if err := redisClient.HIncrByFloat(ctx, key, "amount", depositAmount).Err(); err != nil {
// 		return err
// 	}

// 	// บันทึกยอดหลังการฝาก
// 	if err := redisClient.HSet(ctx, key, "after_balance", balance+depositAmount).Err(); err != nil {
// 		return err
// 	}

// 	// เปลี่ยนสถานะของโปรโมชั่นเป็น 1 (กำลังใช้งาน)
// 	if err := redisClient.HSet(ctx, key, "status", "1", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
// 		return err
// 	}

// 	totalKey := fmt.Sprintf("%s:total_transactions", userID)
// 	if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
// 		return err
// 	}

// 	_, err = redisClient.HIncrByFloat(ctx, totalKey, "total_amount", depositAmount).Result()
// 	return err
// }
// ฟังก์ชันถอนเงิน
func withdraw(redisClient *redis.Client, userID string, proID string, withdrawAmount float64) error {
	lockKey := fmt.Sprintf("%s:lock", userID)

    // ตรวจสอบการล็อก
    if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
        return fmt.Errorf("could not acquire lock: %v", err)
    }
	defer releaseLock(redisClient, lockKey)
	
    key := fmt.Sprintf("%s:%s", userID, proID)
 
    // ตรวจสอบให้แน่ใจว่ามี key สำหรับการทำธุรกรรมแต่ละครั้ง
    if err := ensureTotalTransactionsKey(redisClient, userID); err != nil {
        return fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 297")
    }

	currentPromotionKey := fmt.Sprintf("%s:current_promotion", userID)

	checkpro,_ := redisClient.HGetAll(ctx, currentPromotionKey).Result()
	if len(checkpro) > 0 {
		
	
	
    // ดึงข้อมูลโปรโมชั่นปัจจุบัน
    currentStatus, err := redisClient.HGet(ctx, currentPromotionKey, "status").Result()
	

    if err != nil && err != redis.Nil {
		return fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถฝากเงินเพิ่มได้ 347")
    }
	//fmt.Println(" currentStatus: ",currentStatus)
	intStr,_ := strconv.Atoi(currentStatus)
	if err != nil || intStr > 0 {
        return fmt.Errorf("โปรโมชั่น ยังมีสถานะทำงาน ไม่สามารถถอนเงินเพิ่มได้ ")
    }

	balanceStr, err := redisClient.HGet(ctx, key, "balance").Result()
	if err != nil {
		return err
	}

	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return err
	}

	if balance < withdrawAmount {
		return fmt.Errorf("ยอดคงเหลือไม่พอถอน!")
	}

	if err := redisClient.HSet(ctx, key, "before_balance", balance).Err(); err != nil {
		return err
	}

	if err := redisClient.HIncrByFloat(ctx, key, "amount", -withdrawAmount).Err(); err != nil {
		return err
	}

	// บันทึกยอดหลังการฝาก
	if err := redisClient.HSet(ctx, key, "after_balance", balance + withdrawAmount).Err(); err != nil {
		return err
	}
	if err := redisClient.HSet(ctx, key, "status", "2", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
		return err
	}

	totalKey := fmt.Sprintf("%s:total_transactions", userID)
	if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
		return err
	}
	_, err = redisClient.HIncrByFloat(ctx, totalKey, "total_amount", -withdrawAmount).Result()
	return err
	} else {
	
		fmt.Printf("checkpro: %+v \n",checkpro)
		return nil
	}
}

// ฟังก์ชันการเล่นเกม
// func playGame(redisClient *redis.Client, userID string, proID string, gameAmount float64) error {
// 	lockKey := fmt.Sprintf("%s:lock", userID)

// 	if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
// 		return fmt.Errorf("could not acquire lock: %v", err)
// 	}

// 	if locked, err := acquireLock(redisClient, lockKey); err != nil || !locked {
// 		return fmt.Errorf("could not acquire lock: %v", err)
// 	}
// 	defer releaseLock(redisClient, lockKey)

// 	key := fmt.Sprintf("%s:%s", userID, proID)

// 	if err := checkKeyType(redisClient, key, "hash"); err != nil {
// 		return err
// 	}

// 	currentStatus, err := redisClient.HGet(ctx, key, "status").Result()
// 	intStr, _ := strconv.Atoi(currentStatus)
// 	if err != nil || intStr < 1 {
// 		if intStr == 0 {
// 			return fmt.Errorf("โปรโมชั่นไม่อยู่ในสถานะ ใช้งาน!")
// 		}
// 		return fmt.Errorf("โปรโมชั่นไม่สามารถใช้งานได้!")
// 	}

// 	// บันทึกยอดก่อนเล่นเกม
// 	balanceStr, err := redisClient.HGet(ctx, key, "balance").Result()
// 	if err != nil {
// 		return err
// 	}
// 	balance, _ := strconv.ParseFloat(balanceStr, 64)

// 	if balance < gameAmount {
// 		return fmt.Errorf("ยอดคงเหลือไม่พอเล่นเกม!")
// 	}

// 	// ดำเนินการเล่นเกม (การคำนวณผลที่ได้)
// 	netGain := gameAmount * 2 // ยกตัวอย่างการชนะ (ขึ้นอยู่กับกฎเกมที่คุณสร้าง)

// 	// อัปเดตยอด
// 	if err := redisClient.HIncrByFloat(ctx, key, "balance", netGain).Err(); err != nil {
// 		return err
// 	}
// 	// ตรวจสอบยอดหลังเล่นเกม
// 	balance += netGain // คำนวณยอดใหม่หลังจากเล่นเกม

// 	// หากยอดคงเหลือเหลือ 0 ให้ปรับสถานะโปรโมชั่นเป็น 2
// 	if balance == 0 {
// 		if err := redisClient.HSet(ctx, key, "status", "2", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
// 			return err
// 		}
// 	} else {
// 		// หากไม่เป็น 0 ก็ให้ปรับสถานะเป็น 1 (กำลังใช้งาน)
// 		if err := redisClient.HSet(ctx, key, "status", "1", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// ฟังก์ชันการเล่นเกม
func playGame(redisClient *redis.Client, userID string, proID string, gameAmount float64) (map[string]float64, error) {
	lockKey := fmt.Sprintf("%s:lock", userID)

	if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
		return nil, fmt.Errorf("could not acquire lock: %v", err)
	}

	if locked, err := acquireLock(redisClient, lockKey); err != nil || !locked {
		return nil, fmt.Errorf("could not acquire lock: %v", err)
	}
	defer releaseLock(redisClient, lockKey)

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

    // ตรวจสอบว่าโปรโมชั่นมีข้อมูลหรือไม่
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
		}

		

		return c.JSON(fiber.Map{
			"Status": true,
			//"Data": response,
			"Message": "Selected Promotion successfully",
		})
	}
}
func clearRedisKey(redisClient *redis.Client, userID string) error {
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
		response,err := getPromotionStatusForUser(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr))
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

		// amount, err := strconv.ParseFloat(body.Amount, 64)
		// if err != nil {
		// 	fmt.Printf("Error converting amount: %v\n", err)
		// 	return err
		// }
		
		err := withdraw(redisClient,fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr),body.ProID,body.Amount)
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
			//"Data":response,
			//},
			"Message": "Promotion withdraw successfully",
		})
	}
}
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
		
		if len(promotion) > 0 {
			
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
		var userIDStr string
		switch v := userID.(type) {
		case string:
			userIDStr = v
		case int:
			userIDStr = fmt.Sprintf("%d", v) 
		}
		// สมมุติว่าคุณมี redisClient และ userID ของผู้ใช้
		err := clearRedisKey(redisClient, fmt.Sprintf("%s%s",c.Locals("prefix"),userIDStr))
		if err != nil {
			fmt.Println("Error clearing Redis keys:", err)
		} else {
			fmt.Println("Successfully cleared Redis keys for user:", userID)
		}
		return c.JSON(fiber.Map{
			"Status": true,
			//"Data": fiber.Map{
			//	"ProId": body.ProID,
			//"Data":response,
			//},
			"Message": "Clear successfully",
		})
	}
}