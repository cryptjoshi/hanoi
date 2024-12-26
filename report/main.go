package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/go-redis/redis/v8"
	//"github.com/gorilla/websocket"
	//"github.com/gofiber/contrib/websocket"
	//"github.com/valyala/fasthttp"
	//"github.com/gin-gonic/gin"
	//"github.com/google/uuid" 
	"log"
	"os"
	"fmt"
	//"hanoi/rabbitmq"
	//"hanoi/handler"
	//"hanoi/users"
	"hanoi/handler/partner"
	"hanoi/route"

	//jwtn "hanoi/handler/jwtn"
	// "hanoi/database"
	// "hanoi/models"
	//"hanoi/handler/njwt"
	//"gorm.io/gorm"
	"time"
	//"github.com/swaggo/gin-swagger"
	//"github.com/swaggo/fiber-swagger"
	//"github.com/swaggo/files"
	 _ "hanoi/docs" // สำหรับเอกสาร Swagger
	//"crypto/sha256"
	//"github.com/gofiber/contrib/websocket"
	socketio "github.com/doquangtan/gofiber-socket.io"
 	"encoding/json"
	
)


type MessageObject struct {
    Data  string `json:"data"`
    From  string `json:"from"`
    Event string `json:"event"`
    To    string `json:"to"`
}
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

type Message struct {
    ID      string `json:"id"`
    Message string `json:"message"`
}



var redis_master_host = os.Getenv("REDIS_HOST")
var redis_master_port = os.Getenv("REDIS_PORT")
var redis_master_password = os.Getenv("REDIS_PASSWORD")
var redis_slave_host = os.Getenv("REDIS_SLAVE_HOST")
var redis_slave_port = os.Getenv("REDIS_SLAVE_PORT")
var redis_slave_password = os.Getenv("REDIS_SLAVE_PASSWORD")
var redis_database = getEnv("REDIS_DATABASE", "0")
var secretKey = os.Getenv("PASSWORD_SECRET")
var io *socketio.Io
 
//var hub = make(map[*websocket.Conn]bool)
// type Client struct {
// 	ID   string
// 	Conn *websocket.Conn
// }

var rdb *redis.Client
var ctx = context.Background()
 
// @title Api Goteway in Go
// @version 1.0
// @description Api Goteway in Go.
// @host localhost:4006
// @BasePath /api/v1

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func handleError(err error) {
	log.Fatal(err)
}

func printMessage(message string) {
    for i := 0; i < 5; i++ {
        fmt.Println(message)
        time.Sleep(1 * time.Second)
    }
}

func subscribeMessages(client *redis.Client, channel string,io *socketio.Io) {
	sub := client.Subscribe(ctx, channel)
	defer sub.Close()

	// Wait for confirmation that subscription is received
	_, err := sub.Receive(ctx)
	if err != nil {
		log.Fatalf("Could not subscribe: %v", err)
	}
	// key := []byte(secretKey) 
	// hashedKey := sha256.Sum256(key)
	currentTime := time.Now()

	ch := sub.Channel()
	for msg := range ch {

		user := partner.OBody{}


		user.Startdate = currentTime.Format("2006-01-02")
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
			fmt.Printf("Unmarshal: %v", err)
		}
		//payload := 
		//user.Prefix = 
		
		fmt.Printf("User: %+v\n",user)
		fmt.Printf("Result: %+v\n",result)
		user.Prefix =  result["Prefix"].(string)

		payload,err := partner.GetOverview(user)
		if err != nil {
			log.Fatalf("Could not decrypted: %v", err)
		}
		fmt.Printf("Payload: %+v\n",payload)
		// คีย์ที่ใช้งานเป็นต้องมีขนาด 16, 24 หรือ 32 bytes
		// decryptedData,err := jwtn.Decrypt(msg.Payload,hashedKey[:])
		// if err != nil {
		// 	log.Fatalf("Could not decrypted: %v", err)
		// }
		// decompressedData,err := jwtn.DecompressData(decryptedData)
		// if err != nil {
		// 	log.Fatalf("Could not decompress: %v", err)
		// }
		// count := len(hub) // นับจำนวนการเชื่อมต่อ websocket
		// fmt.Printf("Current number of websocket connections: %d\n", count)

		// for _, client := range hub {   // hub ควรเป็นแผนที่ของการเชื่อมต่อ websocket
		// 	if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
		// 		log.Println("Error sending message:", err)
		// 	}
		// }
		//fmt.Println("payload ",msg.Payload)
		//fmt.Println("Decrypted and Decompress ",string(decompressedData))
	//	io := io.(*socketio.Io)
		io.Emit("message", msg.Payload)
		//fmt.Printf("Received and emitted message: %s\n", string(decompressedData))

	}
}
 
func socketIoRoute(app fiber.Router) {
	io = socketio.New()

	io.OnAuthorization(func(params map[string]string) bool {
		// auth, ok := params["Authorization"]
		// if !ok {
		// 	return false
		// }
		return true
	})

	io.OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)
		socket.Join("demo")
		io.To("demo").Emit("hello", socket.Id+" join us room...", "server message")


		socket.On("demo.hello", func(event *socketio.EventPayload) {
			//var messageData Message
			//json.Unmarshal([]byte(msg), &messageData)
			// ใช้ messageData.ID และ messageData.Message ต่อไป
			if data, ok := event.Data[0].(map[string]interface{}); ok {
				id, idOk := data["id"].(string)
				message, msgOk := data["message"].(string)
				
				if !idOk || !msgOk {
					fmt.Println("ข้อมูลบางส่วนไม่ถูกต้อง")
					return
				}
				
				fmt.Printf("ID: %s\nMessage: %s\n", id, message)
			}
		})

		socket.On("test", func(event *socketio.EventPayload) {
			socket.Emit("test", event.Data...)
		})

		socket.On("join-room", func(event *socketio.EventPayload) {
			if len(event.Data) > 0 && event.Data[0] != nil {
				socket.Join(event.Data[0].(string))
			}
		})

		socket.On("leave-room", func(event *socketio.EventPayload) {
			socket.Leave("demo")
			io.To("demo").Emit("hello", socket.Id+" leave us room...", "server message")
		})

		socket.On("room-emit", func(event *socketio.EventPayload) {
			socket.To("demo").Emit("hello", socket.Id, event.Data)
		})

		socket.On("disconnecting", func(event *socketio.EventPayload) {
			println("disconnecting", socket.Nps, socket.Id)
		})

		socket.On("disconnect", func(event *socketio.EventPayload) {
			println("disconnect", socket.Nps, socket.Id)
		})
	})

	io.Of("/admin").OnConnection(func(socket *socketio.Socket) {
		println("connect", socket.Nps, socket.Id)
	})

	app.Use("/", io.Middleware)
	app.Route("/socket.io", io.Server)
}
func main() {

	//init()
	// สร้างเซิร์ฟเวอร์ Socket.IO
	//clients := make(map[string]string)

	rdb = redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", // redis_master_password,
		DB:       0,  // Use database 0
	})

	
	//go handleMessages(ctx)

	 


	app := fiber.New()
	app.Route("/", socketIoRoute)
	go subscribeMessages(rdb, "balance_update_channel",io)
	// app.Use(func(c *fiber.Ctx) error {
	// 	// IsWebSocketUpgrade returns true if the client
	// 	// requested upgrade to the WebSocket protocol.
	// 	if websocket.IsWebSocketUpgrade(c) {
	// 		c.Locals("allowed", true)
	// 		return c.Next()
	// 	}
	// 	return fiber.ErrUpgradeRequired
	// })
	// app.Use(cors.New(cors.Config{
    //     AllowOrigins: "*", // อนุญาตทุกโดเมน (ในโปรดักชันให้ระบุโดเมนที่จำเป็นเท่านั้น)
    //     AllowMethods: "GET,POST,PUT,DELETE",
    //     AllowHeaders: "Origin, Content-Type, Accept",
    // }))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions",
		ExposeHeaders:    "Content-Length",
		//AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))
	app.Use(compress.New())
	// Setup the middleware to retrieve the data sent in first GET request

	//app.Use("/socket.io/", socketServer)

	
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
	//app.Use("/ws", websocket.New(websocketHandler))

	v1 := app.Group("/api/v1")
	 route.SetupRoutes(v1,true)
 
	 // เพิ่มกลุ่มใหม่สำหรับ /api
	 api := app.Group("/api")
	 route.SetupRoutes(api,false) // เรียกใช้ฟังก์ชัน SetupRoutes สำหรับกลุ่มนี้
     
	

	 app.Get("/test", func(c *fiber.Ctx) error {
		io := c.Locals("io").(*socketio.Io)

		io.Emit("message", "Hello from socket.io server")

		io.Of("/admin").Emit("event", map[string]interface{}{
			"Ok": 1,
		})

		return c.SendStatus(200)
	})
	 

    // เรียกใช้ฟังก์ชันจาก efinity.go
	log.Fatal(app.Listen(":8030"))
	 
	
	
}