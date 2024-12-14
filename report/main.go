package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/compress"
	//"github.com/gin-gonic/gin"
	"log"
	//"os"
	"fmt"
	//"hanoi/rabbitmq"
	//"hanoi/handler"
	//"hanoi/users"
	"hanoi/route"
	// "hanoi/database"
	// "hanoi/models"
	//"hanoi/handler/njwt"
	//"gorm.io/gorm"
	"time"
	//"github.com/swaggo/gin-swagger"
	//"github.com/swaggo/fiber-swagger"
	//"github.com/swaggo/files"
	 _ "hanoi/docs" // สำหรับเอกสาร Swagger

	
)

// func loadDatabase() {
// 	if err := database.Connect(); err != nil {
// 		handleError(err)
// 	}

// }

// func DropTable () {

// 	database.Database.Migrator().DropTable(&models.TransactionSub{})
// 	database.Database.Migrator().DropTable(&models.BuyInOut{})

// }

// func migrateNormal(db *gorm.DB) {

// 	if err := db.AutoMigrate(&models.Product{},&models.BanksAccount{},&models.Users{},&models.TransactionSub{},&models.BankStatement{},&models.BuyInOut{}); err != nil {
// 		handleError(err)
// 	}
	 
// 	fmt.Println("Migrations Normal Tables executed successfully")
// }
// func migrateAdmin() {
 
// 	if err := database.Database.AutoMigrate(&models.TsxAdmin{},&models.Provider{},&models.Promotion{}); err != nil {
// 		handleError(err)
// 	}
// 	fmt.Println("Migrations Admin Tables executed successfully")
// }



// @title Api Goteway in Go
// @version 1.0
// @description Api Goteway in Go.
// @host localhost:4006
// @BasePath /api/v1


func handleError(err error) {
	log.Fatal(err)
}
func printMessage(message string) {
    for i := 0; i < 5; i++ {
        fmt.Println(message)
        time.Sleep(1 * time.Second)
    }
}
func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
        AllowOrigins: "*", // อนุญาตทุกโดเมน (ในโปรดักชันให้ระบุโดเมนที่จำเป็นเท่านั้น)
        AllowMethods: "GET,POST,PUT,DELETE",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))
	app.Use(compress.New())

	// app.Use(func(c *fiber.Ctx) error {
	// 	// ดึง prefix จาก token
	// 	prefix, err := jwt.ExtractPrefixFromToken(c)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	// 	}

	// 	// เชื่อมต่อฐานข้อมูลตาม prefix
	// 	db, err := database.ConnectToDB(prefix)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to database"})
	// 	}

	// 	// เก็บการเชื่อมต่อใน context เพื่อให้ endpoint อื่นๆ ใช้งานได้
	// 	c.Locals("db", db)

	// 	// ไปยัง handler ต่อไป
	// 	return c.Next()
	// })

	app.Use(func(c *fiber.Ctx) error {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		c.Locals("location", loc)
		return c.Next()
	})

	app.Use(logger.New())
	
	// go printMessage("Hello from Goroutine")

    // // ฟังก์ชันใน main จะทำงานต่อไป
    // printMessage("Hello from main")

    // // Sleep เพื่อรอให้ goroutine ทำงานเสร็จ (ตัวอย่างเพื่อให้เห็นผล)
    // time.Sleep(6 * time.Second)
	//migrateAdmin()

	 v1 := app.Group("/api/v1")
	 route.SetupRoutes(v1)
 
    // เรียกใช้ฟังก์ชันจาก efinity.go
	log.Fatal(app.Listen(":8050"))
	 
	
}