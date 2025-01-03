package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	//"os"
	"fmt"
	//"pkd/rabbitmq"
	"pkd/route"
	"pkd/database"
	"pkd/models"
	"time"

	// นำเข้า package game
	//"gorm.io/driver/mysql"
    //"gorm.io/gorm"
	
)

var ctx = context.Background()
var redis_master_host = "redis" //os.Getenv("REDIS_HOST")
var redis_master_port = "6379" 




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

func loadDatabase() {
	if err := database.Connect(); err != nil {
		handleError(err)
	}

}

 

func DropTable () {

	database.Database.Migrator().DropTable(&models.TransactionSub{})
	//database.Database.Migrator().DropTable(&models.BuyInOut{})

}

func migrateNormal() {
	
	//if err := database.Database.AutoMigrate(&models.Product{},&models.BanksAccount{},&models.Users{},&models.TransactionSub{},&models.BankStatement{},&models.BuyInOut{}); err != nil {
	if err := database.Database.AutoMigrate(&models.TransactionSub{}); err != nil {
		handleError(err)
	}
	 
	fmt.Println("Migrations Normal Tables executed successfully")
}
func migrateAdmin() {
 
	if err := database.Database.AutoMigrate(&models.TsxAdmin{},&models.Provider{}); err != nil {
		handleError(err)
	}
	fmt.Println("Migrations Admin Tables executed successfully")
}
func handleError(err error) {
	log.Fatal(err)
}

func main() {

	
	redisClient := createRedisClient()
	defer redisClient.Close()
	if err := checkRedisConnection(redisClient); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		c.Locals("location", loc)
		return c.Next()
	})

	//rabbitmq.Init()
	
	//loadDatabase()
    // DropTable()
	//migrateNormal()
	//  migrateAdmin()

	app.Use(logger.New())

	route.SetupRoutes(app,redisClient)
 

    // เรียกใช้ฟังก์ชันจาก efinity.go
  
	log.Fatal(app.Listen(":8070"))
	//log.Fatal(app.Listen(":3006"))
	
}