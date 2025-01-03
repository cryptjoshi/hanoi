package main

import (
	
	"hanoi/route"
	"hanoi/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"log"
	"time"
	"github.com/swaggo/fiber-swagger"
	_ "hanoi/docs"
	
)

 

type PBody struct {
	UserID  string  `json:"userID"`
	ProID   string  `json:"proID"`
	DepositAmount string `json:"amount"`
}
 
  
// @title Backend API Document
// @version 1.0
// @description This is a sample API server.
// @host 152.42.185.164:4006
// @BasePath /
func main() {

	redisClient := handler.CreateRedisClient()
	defer redisClient.Close()
	if err := handler.CheckRedisConnection(redisClient); err != nil {
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

	 
	v2 := app.Group("/api/v2")
	route.SetupRoutes(v2,true,redisClient)
	// เริ่มต้นการทำงานของ Fiber
	log.Fatal(app.Listen(":8060"))
}