package main

import (
	"context"
	// "fmt"
	// //"log"
	// "strconv"
	// "time"
	// "strings"

	"hanoi/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
	"github.com/swaggo/fiber-swagger"
	_ "hanoi/docs"
	//"github.com/shopspring/decimal"
)

var ctx = context.Background()
var redis_master_host = "redis" //os.Getenv("REDIS_HOST")
var redis_master_port = "6379" 

type PBody struct {
	UserID  string  `json:"userID"`
	ProID   string  `json:"proID"`
	DepositAmount string `json:"amount"`
}

// สร้าง Redis Client
// func createRedisClient() *redis.Client {
// 	return 	 redis.NewClient(&redis.Options{
// 		Addr:     redis_master_host + ":" + redis_master_port,
// 		Password: "", //redis_master_password,
// 		DB:       0,  // ใช้ database 0
// 	})
//  }

 func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: redis_master_host + ":" + redis_master_port,
	})
}
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
// 			releaseLock(redisClient, lockKey)
// 		}()
// 	}
// 	return locked, nil
// }

// // ฟังก์ชันปลดล็อก
// func releaseLock(redisClient *redis.Client, lockKey string) error {
// 	// ลบฟิลด์ที่เป็นล็อก
// 	_, err := redisClient.HDel(ctx, lockKey, "locked").Result()
// 	return err
// }

// // ฟังก์ชันเช็คสถานะล็อก
// func isLocked(redisClient *redis.Client, lockKey string) (bool, error) {
// 	locked, err := redisClient.HGet(ctx, lockKey, "locked").Result()
// 	if err != nil { 
// 	if err == redis.Nil {
// 			return false, nil
// 		}
// 		return false,err
// 	}
// 	return locked == "1", nil
// }
// // ฟังก์ชันเพื่อตรวจสอบว่าคีย์มีประเภทที่ถูกต้องหรือไม่
// func checkKeyType(redisClient *redis.Client, key string, expectedType string) error {
// 	keyType, err := redisClient.Type(ctx, key).Result()
// 	if err != nil {
// 		return fmt.Errorf("error checking key type: %v", err)
// 	}
// 	if keyType != expectedType && expectedType != "any" {
// 		return fmt.Errorf("key %s is of type %s, expected %s", key, keyType, expectedType)
// 	}
// 	return nil
// }

// // ฟังก์ชันเพื่อสร้างคีย์ total_transactions ถ้ายังไม่มี
// func ensureTotalTransactionsKey(redisClient *redis.Client, userID string) error {
// 	totalKey := fmt.Sprintf("%s:total_transactions", userID)

// 	// เช็คประเภทของคีย์
// 	keyType, err := redisClient.Type(ctx, totalKey).Result()
// 	if err != nil {
// 		return fmt.Errorf("error checking key type: %v", err)
// 	}

// 	// ถ้ายังไม่มีคีย์ หรือประเภทไม่ถูกต้อง ให้สร้างใหม่
// 	if keyType == "none" {
// 		// สร้างคีย์ใหม่ในประเภท hash
// 		if err := redisClient.HSet(ctx, totalKey, "total_amount", 0).Err(); err != nil {
// 			return fmt.Errorf("error initializing total transactions: %v", err)
// 		}
// 	}

// 	return nil
// }

// // ฟังก์ชันเลือกโปรโมชั่น
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

// 	var hasActivePromotion bool
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
// 			hasActivePromotion = true
// 			break
// 		}
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

// 	if balance != "0" {
// 		return fmt.Errorf("balance must be zero to select a new promotion")
// 	}

// 	if hasActivePromotion {
// 		return fmt.Errorf("cannot select new promotion while the previous promotion is still active")
// 	}

// 	// กำหนดสถานะโปรโมชั่นใหม่ พร้อม timestamp
// 	if err := redisClient.HSet(ctx, key, "status", "0", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // ฟังก์ชันฝากเงิน
// func depositSuccess(redisClient *redis.Client, userID string, proID string, depositAmount float64) error {
// 	lockKey := fmt.Sprintf("%s:lock", userID)

// 	if locked, err := isLocked(redisClient, lockKey); err != nil || locked {
// 		return fmt.Errorf("could not acquire lock: %v", err)
// 	}

// 	// if locked, err := acquireLock(redisClient, lockKey); err != nil || !locked {
// 	// 	return fmt.Errorf("could not acquire lock: %v", err)
// 	// }
// 	// defer releaseLock(redisClient, lockKey)

// 	// key := fmt.Sprintf("%s:%s", userID, proID)

// 	// if err := checkKeyType(redisClient, key, "hash"); err != nil {
// 	// 	return err
// 	// }
// 	key := fmt.Sprintf("%s:%s", userID, proID)

// 	if err := ensureTotalTransactionsKey(redisClient, userID); err != nil {
// 		return err
// 	}

// 	if err := checkKeyType(redisClient, key, "hash"); err != nil {
// 		return err
// 	}

// 	currentStatus, err := redisClient.HGet(ctx, key, "status").Result()
// 	if err != nil || currentStatus != "0" {
// 		return fmt.Errorf("promotion is either not selected or is not in a valid state for deposit")
// 	}

// 	if err := redisClient.HIncrByFloat(ctx, key, "balance", depositAmount).Err(); err != nil {
// 		return err
// 	}

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

// // ฟังก์ชันถอนเงิน
// func withdraw(redisClient *redis.Client, userID string, proID string, withdrawAmount float64) error {
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
// 	if err != nil || currentStatus != "1" {
// 		return fmt.Errorf("promotion is either not in a valid state for withdrawal")
// 	}

// 	balanceStr, err := redisClient.HGet(ctx, key, "balance").Result()
// 	if err != nil {
// 		return err
// 	}

// 	balance, err := strconv.ParseFloat(balanceStr, 64)
// 	if err != nil {
// 		return err
// 	}

// 	if balance < withdrawAmount {
// 		return fmt.Errorf("insufficient balance for withdrawal")
// 	}

// 	if err := redisClient.HIncrByFloat(ctx, key, "balance", -withdrawAmount).Err(); err != nil {
// 		return err
// 	}

// 	if err := redisClient.HSet(ctx, key, "status", "2", "timestamp", time.Now().Format(time.RFC3339)).Err(); err != nil {
// 		return err
// 	}

// 	totalKey := fmt.Sprintf("%s:total_transactions", userID)
// 	if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
// 		return err
// 	}
// 	_, err = redisClient.HIncrByFloat(ctx, totalKey, "total_amount", -withdrawAmount).Result()
// 	return err
// }

// // ฟังก์ชันเพื่อดูยอดธุรกรรมทั้งหมด
// func getTotalTransactions(redisClient *redis.Client, userID string) (float64, error) {
// 	totalKey := fmt.Sprintf("%s:total_transactions", userID)
// 	if err := checkKeyType(redisClient, totalKey, "hash"); err != nil {
// 		return 0, err
// 	}

// 	totalAmount, err := redisClient.HGet(ctx, totalKey, "total_amount").Result()
// 	if err != nil {
// 		return 0, err
// 	}

// 	var amount float64
// 	fmt.Sscanf(totalAmount, "%f", &amount)

// 	return amount, nil
// }

// // ฟังก์ชันเพื่อดูสถานะโปรโมชั่น
// func getPromotionStatus(redisClient *redis.Client, userID string, proID string) (map[string]string, error) {
// 	key := fmt.Sprintf("%s:%s", userID, proID)

// 	// ตรวจสอบประเภทของคีย์
// 	if err := checkKeyType(redisClient, key, "hash"); err != nil {
// 		return nil, err
// 	}

// 	// ดึงค่าที่เกี่ยวข้องจาก Redis
// 	status, err := redisClient.HGet(ctx, key, "status").Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			status = "not found"
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	timestamp, err := redisClient.HGet(ctx, key, "timestamp").Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			timestamp = "not found"
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	return map[string]string{
// 		"status":    status,
// 		"timestamp": timestamp,
// 	}, nil
// }


// // ฟังก์ชันเพื่อดูสถานะโปรโมชั่นสำหรับ proID ที่ระบุในทุกผู้ใช้
// func getPromotionStatusForAllUsers(redisClient *redis.Client, proID string) (map[string]map[string]string, error) {
// 	promotionStatuses := make(map[string]map[string]string)
	
// 	key := fmt.Sprintf("*:*")
// 	// ค้นหาทุกคีย์ของผู้ใช้ที่มีโปรโมชั่นที่ตรงกับ proID
// 	userKeys, err := redisClient.Keys(ctx, key).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, userKey := range userKeys {
		 
// 		parts := strings.Split(userKey, ":")
// 		if len(parts) != 2 {
// 			continue // ป้องกันการแตกในกรณีที่คีย์ไม่มีรูปแบบที่ถูกต้อง
// 		}

// 		userID := parts[0]
// 		currentProID := parts[1]

// 		fmt.Println("UserID:",userID)
// 		fmt.Println("currentProID:",currentProID)
	
// 		if currentProID != proID {
// 			continue // ข้ามโปรโมชั่นที่ไม่ตรงกัน
// 		}
// 		if err := checkKeyType(redisClient, userKey, "hash"); err != nil {
// 			return nil,err
// 		}
	
// 		status, err := redisClient.HGet(ctx, userKey, "status").Result()
// 		if err != nil && err != redis.Nil {
// 			return nil, err
// 		}
// 		if err == redis.Nil {
// 			status = "not found"
// 		}

// 		timestamp, err := redisClient.HGet(ctx, userKey, "timestamp").Result()
// 		if err != nil && err != redis.Nil {
// 			return nil, err
// 		}
// 		if err == redis.Nil {
// 			timestamp = "not found"
// 		}

// 		// บันทึกสถานะของโปรโมชั่นของผู้ใช้
// 		promotionStatuses[userID] = map[string]string{
// 			"status":    status,
// 			"timestamp": timestamp,
// 		}
// 	}

// 	return promotionStatuses, nil
// }

// @title User Management API
// @version 1.0
// @description This is a sample API server.
// @host http://152.42.185.164:8020
// @BasePath /
func main() {

	redisClient := createRedisClient()
	defer redisClient.Close()
	if err := checkRedisConnection(redisClient); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins: "*", // อนุญาตทุกโดเมน (ในโปรดักชันให้ระบุโดเมนที่จำเป็นเท่านั้น)
        AllowMethods: "GET,POST,PUT,DELETE",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))
	app.Use(compress.New())
	app.Use(func(c *fiber.Ctx) error {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		c.Locals("location", loc)
		return c.Next()
	})

	app.Use(logger.New())
	
 

    // Swagger route
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    // Example route
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Welcome to the User Management API!")
    })

	// // Endpoint สำหรับฝากเงิน
	// app.Post("/deposit", func(c *fiber.Ctx) error {
	// 	//userID := c.Params("userID")
	// 	//proID := c.Params("proID")
	// 	//depositAmount := c.FormValue("amount")

	// 	body := PBody{}

	// 	if err := c.BodyParser(&body); err != nil {
	// 		fmt.Println("Error parsing body:", err.Error())
	// 		response := fiber.Map{
	// 			"Status":  false,
	// 			"Message": err.Error(),
	// 		}
	// 		return c.JSON(response)
	// 	}


	// 	amount, err := strconv.ParseFloat(body.DepositAmount, 64)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "Invalid deposit amount",
	// 		})
	// 	}

	// 	if err := depositSuccess(redisClient, body.UserID, body.ProID,amount); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"message": "Deposit successful",
	// 	})
	// })

	// // Endpoint สำหรับถอนเงิน
	// app.Post("/withdraw/:userID/:proID", func(c *fiber.Ctx) error {
	// 	userID := c.Params("userID")
	// 	proID := c.Params("proID")
	// 	withdrawAmount := c.FormValue("amount")

	// 	amount, err := strconv.ParseFloat(withdrawAmount, 64)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "Invalid withdraw amount",
	// 		})
	// 	}

	// 	if err := withdraw(redisClient, userID, proID, amount); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"message": "Withdrawal successful",
	// 	})
	// })

	// app.Get("/promotions/status/:userID/:proID", func(c *fiber.Ctx) error {
	// 	userID := c.Params("userID")
	// 	proID := c.Params("proID")

	// 	status, err := getPromotionStatus(redisClient, userID, proID)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"userID":   userID,
	// 		"proID":    proID,
	// 		"status":   status["status"],
	// 		"timestamp": status["timestamp"],
	// 	})
	// })

	// // Endpoint สำหรับดูยอดธุรกรรมทั้งหมด
	// app.Get("/total/:userID", func(c *fiber.Ctx) error {
	// 	userID := c.Params("userID")

	// 	totalAmount, err := getTotalTransactions(redisClient, userID)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"userID":       userID,
	// 		"total_amount": totalAmount,
	// 	})
	// })

	// app.Post("/promotions/status/all", func(c *fiber.Ctx) error {
	// 	// proID := c.Params("proID")
	// 	// fmt.Println("proID:",proID)

	// 	body := PBody{}

	// 	if err := c.BodyParser(&body); err != nil {
	// 		fmt.Println("Error parsing body:", err.Error())
	// 		response := fiber.Map{
	// 			"Status":  false,
	// 			"Message": err.Error(),
	// 		}
	// 		return c.JSON(response)
	// 	}


	// 	statuses, err := getPromotionStatusForAllUsers(redisClient, body.ProID)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"proID":            body.ProID,
	// 		"promotion_status": statuses,
	// 	})
	// })

	v2 := app.Group("/api/v2")
	route.SetupRoutes(v2,true,redisClient)
	// เริ่มต้นการทำงานของ Fiber
	log.Fatal(app.Listen(":8020"))
}