package main

import (
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


func loadDatabase() {
	if err := database.Connect(); err != nil {
		handleError(err)
	}

}

func DropTable () {

	database.Database.Migrator().DropTable(&models.TransactionSub{})
	database.Database.Migrator().DropTable(&models.BuyInOut{})

}

func migrateNormal() {
	
	if err := database.Database.AutoMigrate(&models.Product{},&models.BanksAccount{},&models.Users{},&models.TransactionSub{},&models.BankStatement{},&models.BuyInOut{}); err != nil {
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

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		c.Locals("location", loc)
		return c.Next()
	})

	//rabbitmq.Init()
	
	 loadDatabase()
	//  DropTable()
	 //migrateNormal()
	//  migrateAdmin()

	app.Use(logger.New())

	route.SetupRoutes(app)
 

    // เรียกใช้ฟังก์ชันจาก efinity.go
  
	log.Fatal(app.Listen(":8070"))
	//log.Fatal(app.Listen(":3006"))
	
}