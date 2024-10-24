package handler

import (
	"context"
	"fmt"
	"hanoi/models"
	"math/rand"

	"github.com/amalfra/etag"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	//"github.com/golang-jwt/jwt"
	jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	"encoding/json"
	"hanoi/database"
	"hanoi/repository"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var Words = []string{
	"apple", "banana", "cherry", "date", "elderberry",
	"fig", "grape", "honeydew", "kiwi", "lemon",
	"mango", "orange", "papaya", "quince", "raspberry",
	"strawberry", "tangerine", "watermelon", "blueberry", "blackberry",
	"apricot", "cranberry", "pineapple", "pomegranate", "grapefruit",
	"avocado", "coconut", "guava", "lime", "passionfruit",
	"lychee", "nectarine", "plum", "apricot", "kiwifruit",
	"boysenberry", "cantaloupe", "rambutan", "starfruit", "persimmon",
	"currant", "dragonfruit", "gooseberry", "papaya", "ugli fruit",
	"quince", "ackee", "durian", "jackfruit", "kumquat",
	"litchi", "mulberry", "olive", "rhubarb", "tamarind",
	"tomato", "coconut", "cucumber", "eggplant", "zucchini",
	"potato", "carrot", "onion", "garlic", "broccoli",
	"cauliflower", "spinach", "kale", "lettuce", "cabbage",
	"brussels sprouts", "artichoke", "asparagus", "celery", "green bean",
	"peas", "corn", "radish", "beet", "turnip",
	"rutabaga", "pars"}

var ctx = context.Background()
var amqp_connection *amqp.Connection
var amqp_channel *amqp.Channel
var queue amqp.Queue = amqp.Queue{}
var is_connection = false
var has_channel = false
var has_queue = false

var redis_master_host = "redis" //os.Getenv("REDIS_HOST")
var redis_master_port = "6379"  //os.Getenv("REDIS_PORT")
var redis_master_password = os.Getenv("REDIS_PASSWORD")
var redis_slave_host = os.Getenv("REDIS_SLAVE_HOST")
var redis_slave_port = os.Getenv("REDIS_SLAVE_PORT")
var redis_slave_password = os.Getenv("REDIS_SLAVE_PASSWORD")
var queue_name = "wallet"                   //os.Getenv("QUEUE_NAME")
var exchange_name = "wallet"                //os.Getenv("EXCHANGE_NAME")
var rabbit_mq = "amqp://128.199.92.45:5672" //os.Getenv("RABBIT_MQ") @rabbitmq:5672/wallet
var connection_timeout = os.Getenv("CONNECTION_TIMEOUT")
var redis_database = getEnv("REDIS_DATABASE", "0")
var go_pixel_log = os.Getenv("GO_PIXEL_LOG")
var mysql_host = os.Getenv("MYSQL_HOST")
var mysql_user = os.Getenv("MYSQL_ROOT_USER")
var mysql_pass = os.Getenv("MYSQL_ROOT_PASSWORD")
var jwtSecret = os.Getenv("PASSWORD_SECRET")

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

type ExchangeRate struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
}

type ExchangeRateResponse struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func InitAMQP() {
	fmt.Println("Init AMQP RABBIT")
	fmt.Println("channel")
	fmt.Println(amqp_channel)
	fmt.Println(connection_timeout)

	conn, err := amqp.DialConfig(rabbit_mq, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			conn_timeout, _ := strconv.ParseInt(connection_timeout, 10, 32)
			conn_timeout = conn_timeout * 365 * 24 * 60
			return net.DialTimeout(network, addr, time.Duration(conn_timeout)*time.Second)
		},
	})
	failOnError(err, "Failed to connect to RabbitMQ")
	defer fmt.Println("conenction close")
	// defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer fmt.Println("channel close")
	// defer ch.Close()

	if err == nil {
		amqp_channel = ch
		amqp_connection = conn
		is_connection = true
	}
}

func get_channel() *amqp.Channel {
	return amqp_channel
}

func FastPost(url string, referrer string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	req.SetRequestURI(url)
	req.Header.Add("Referer", referrer)
	req.Header.SetMethod("POST")

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetPixelTrack(c *fiber.Ctx) error {
	if !is_connection {
		InitAMQP()
	}

	fmt.Println("is close")
	fmt.Println(amqp_connection.IsClosed())

	if amqp_connection.IsClosed() {
		is_connection = false
		InitAMQP()
	}

	if !amqp_connection.IsClosed() {
		body := c.OriginalURL() + "&HTTP_REFERER=" + string(c.Request().Header.Referer())
		err := amqp_channel.Publish(
			exchange_name, // exchange
			"",            // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		is_log, _ := strconv.ParseBool(go_pixel_log)
		if is_log {
			log.Printf(" [x] Sent %s", body)
			log.Printf("Exchange name %s", exchange_name)
		} else {
			log.Printf(" [x] Sent Data")
		}

		failOnError(err, "Failed to publish a message")

		message := fmt.Sprintf("Success")
		return c.SendString(message)
	}

	return c.SendString("404 not found")
}

func GetPixelPath(c *fiber.Ctx) error {
	log.Println("==============================================")
	log.Println("Get Pixel Path")
	// InitAMQP()
	c.Append("Cache-Control", "public, max-age=300")
	c.Append("content-type", "text/javascript")
	c.Append("Accept-Encoding", "gzip, deflate, brotli")
	c.Append("Expires", time.Now().AddDate(0, 0, 1).Format(http.TimeFormat))

	key := c.Params("key")
	etagKey := key + "_etag"
	db, _ := strconv.Atoi(redis_database)
	rdbWrite := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: redis_master_password,
		DB:       db,
	})
	rdbRead := redis.NewClient(&redis.Options{
		Addr:     redis_slave_host + ":" + redis_slave_port,
		Password: redis_slave_password,
		DB:       db,
	})
	pixel, err := rdbRead.Get(ctx, key).Result()

	etagPixel, errEtag := rdbRead.Get(ctx, etagKey).Result()
	reqEtag := string(c.Request().Header.Peek("If-None-Match"))

	fmt.Println("If non match : " + reqEtag)

	if errEtag != redis.Nil && err != redis.Nil {
		if etagPixel == reqEtag {
			fmt.Println("is etag match")
			return c.SendStatus(304)
		}
	}

	if err == redis.Nil || len(pixel) == 0 {
		fmt.Println("use rest client")
		url := os.Getenv("PHP_URL") + "/pixel/" + key
		fmt.Println(url)
		resp, resp_err := FastPost(url, string(c.Request().Header.Referer()))

		if resp_err == nil {
			m := minify.New()
			m.AddFunc("text/javascript", js.Minify)
			// pixel_resp, _ := m.String("text/javascript", string(resp.Body()))
			pixel_resp := string(resp.Body())

			fmt.Println(pixel_resp)

			eTag := etag.Generate(pixel_resp, false)
			redis_err := rdbWrite.Set(ctx, key, pixel_resp, 0).Err()

			if redis_err != nil {
				fmt.Println(redis_err)
			}

			etag_err := rdbWrite.Set(ctx, etagKey, eTag, 0).Err()

			if etag_err != nil {
				fmt.Println(etag_err)
			}

			c.Append("ETag", eTag)

			return c.SendString(pixel_resp)
		}

		return c.SendString("404 not found")

	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("use redis")
		eTag := etag.Generate(pixel, false)
		redis_err := rdbWrite.Set(ctx, etagKey, eTag, 0).Err()

		if redis_err != nil {
			fmt.Println(redis_err)
		}

		c.Append("ETag", eTag)
		return c.SendString(pixel)
	}

	return c.SendString("404 not found")
}

func GetRoot(c *fiber.Ctx) error {
	// if !is_connection {
	// 	InitAMQP()
	// }

	// fmt.Println("is close")
	// fmt.Println(amqp_connection.IsClosed())

	// if amqp_connection.IsClosed() {
	// 	is_connection = false
	// 	InitAMQP()
	// }
	// tokenString := getToken(c.Request())
	// if tokenString == "" {
	// 	return c.SendStatus(http.StatusUnauthorized)
	// }
	// token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(jwtSecret), nil
	//   })

	//   if err != nil {
	// 	return c.SendStatus(http.StatusUnauthorized)
	//   }

	//   claims := token.Claims.(*MyCustomClaims)
	//   query := "SELECT * FROM USERS WHERE username = ?"
	//   row := db.QueryRow(query, claims.username)
	//   var user User
	//   err2 := row.Scan(&user.username, &user.id, &user.role, &user.active)
	//   if err2 != nil {
	// 	return c.SendStatus(http.StatusNotFound)
	//   }

	//   return c.JSON(user)

	message := fmt.Sprintf("I am GOPRO!")
	return c.SendString(message)
}

//  authorized handler

func Signup(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	message := fmt.Sprintf("I am GROOT!")
	return c.SendString(message)
}

// Protected route
func Protected(c *fiber.Ctx) error {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	username := claims["Username"].(string)
	favPhrase := claims["PartnersKey"].(string)

	return c.SendString("Bye Bye 👋" + username + " " + favPhrase)

}

func Logout(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	// message := fmt.Sprintf("I am GROOT!")
	// return c.SendString(message)
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	username := claims["Username"].(string)
	//favPhrase := claims["PartnersKey"].(string)
	//id := claims["ID"]
	prefix := username[:3]
	db, _ := database.ConnectToDB(prefix)
	//if claims != nil {
	updates := map[string]interface{}{
		"Token": "",
	}

	// อัปเดตข้อมูลยูสเซอร์
	repository.UpdateFieldsUserString(db, username, updates)

	response := fiber.Map{
		"Message": "ออกจากระบบสำเร็จ!",
		"Status":  true,
		"Data": fiber.Map{
			"id": -1,
		},
	}
	return c.JSON(response)
	// } else {
	// 	response := fiber.Map{
	// 		"Message": "มีข้อผิดพลาด กรุณาออกจากระบบอีกครั้ง!",
	// 		"Status":  false,
	// 		"Data": fiber.Map{
	// 			"id": -1,
	// 		},
	// 	}
	// 	return c.JSON(response)
	// }

}

func GenerateSeedPhrase(length int) string {
	seedPhrase := ""
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		randomInt := rand.Intn(40)
		seedPhrase = fmt.Sprintf("%s %s", seedPhrase, Words[randomInt])
	}

	return seedPhrase

}

func GetDBFromContext(c *fiber.Ctx) (*gorm.DB, error) {
	dbInterface := c.Locals("db")
	if dbInterface == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "No database connection found")
	}

	// แปลง interface{} ให้เป็น *gorm.DB
	db, ok := dbInterface.(*gorm.DB)
	if !ok {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Invalid database connection")
	}

	return db, nil
}
func handleError(err error) {
	log.Fatal(err)
}
func migrateNormal(db *gorm.DB) {

	if err := db.AutoMigrate(&models.Product{}, &models.BanksAccount{}, &models.Users{}, &models.TransactionSub{}, &models.BankStatement{}, &models.BuyInOut{}); err != nil {
		handleError(err)
	}

	fmt.Println("Migrations Normal Tables executed successfully")
}
func migrateAdmin(db *gorm.DB) {

	if err := db.AutoMigrate(&models.TsxAdmin{}, &models.Provider{}); err != nil {
		handleError(err)
	}
	fmt.Println("Migrations Admin Tables executed successfully")
}

// database

type Dbstruct struct {
	DBName   string   `json:"dbname"`
	Prefix   string   `json:"prefix"`
	Username string   `json:"username"`
	Dbnames  []string `json:"dbnames"`
}

func createDatabase(dbName string) *gorm.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host)

	// Connect to MySQL without a specific database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", dbName)
	if err := db.Exec(createDBQuery).Error; err != nil {

		return nil
	}
	newDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host, dbName)
	newDB, err := gorm.Open(mysql.Open(newDsn), &gorm.Config{})

	return newDB
}

// Function to connect and create a database with a specific prefix and name
func CreateDatabase(c *fiber.Ctx) error {

	dbstruct := new(Dbstruct)

	if err := c.BodyParser(dbstruct); err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host)

	// Connect to MySQL without a specific database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}

	for _, dbname := range dbstruct.Dbnames {
		// Create the database with the provided prefix and name
		//fullDBName := dbstruct.Prefix + "_" + dbname
		createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", dbname)
		if err := db.Exec(createDBQuery).Error; err != nil {
			response := fiber.Map{
				"Message": "สร้างฐานข้อมูลไม่สำเร็จ!!",
				"Status":  false,
			}
			return c.JSON(response)
		}

		// Switch to the new database
		newDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host, dbname)
		newDB, err := gorm.Open(mysql.Open(newDsn), &gorm.Config{})
		if err != nil {
			response := fiber.Map{
				"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
				"Status":  false,
			}
			return c.JSON(response)
		}
		migrateAdmin(newDB)
		migrateNormal(newDB)
	}

	response := fiber.Map{
		"Message": "สร้างฐานข้อมูลสำเร็จ!",
		"Status":  true,
	}
	return c.JSON(response)
}
func GetDatabaseList(c *fiber.Ctx) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host)

	// Connect to MySQL without a specific database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}
	type DatabaseInfo struct {
		Prefix string   `json:"prefix"`
		Names  []string `json:"names"`
	}
	// Query to get all databases
	groupedDatabases := make(map[string][]string)

	rows, err := db.Raw("SHOW DATABASES").Rows()
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดในการดึงรายชื่อฐานข้อมูล",
			"Status":  false,
		}
		return c.JSON(response)
	}
	defer rows.Close()

	systemDatabases := map[string]bool{
		"information_schema": true,
		"mysql":              true,
		"performance_schema": true,
		"sys":                true,
	}

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			continue
		}
		//fmt.Println(dbName)
		//fmt.Println(systemDatabases[dbName])
		// Include databases with underscore in their names and exclude system databases
		//
		if strings.Contains(dbName, "_") && !systemDatabases[dbName] {
			parts := strings.SplitN(dbName, "_", 2)
			if len(parts) == 2 {
				prefix := parts[0]
				if _, exists := groupedDatabases[prefix]; !exists {
					groupedDatabases[prefix] = []string{}
				}
				groupedDatabases[prefix] = append(groupedDatabases[prefix], dbName)
			}
		}
	}
	var databaseList []DatabaseInfo
	for prefix, names := range groupedDatabases {
		databaseList = append(databaseList, DatabaseInfo{
			Prefix: prefix,
			Names:  names,
		})
	}

	response := fiber.Map{
		"Message":   "ดึงรายชื่อฐานข้อมูลสำเร็จ",
		"Status":    true,
		"Databases": groupedDatabases,
	}
	return c.JSON(response)
}
func GetDatabaseByPrefix(c *fiber.Ctx) error {
	loginRequest := new(Dbstruct)

	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}

	rows, err := db.Raw("SHOW DATABASES").Rows()
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดในการดึงรายชื่อฐานข้อมูล",
			"Status":  false,
		}
		return c.JSON(response)
	}
	defer rows.Close()

	systemDatabases := map[string]bool{
		"information_schema": true,
		"mysql":              true,
		"performance_schema": true,
		"sys":                true,
	}

	var databaseNames []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			continue
		}

		if strings.HasPrefix(dbName, loginRequest.Prefix+"_") && !systemDatabases[dbName] {
			databaseNames = append(databaseNames, dbName)
		}
	}

	response := fiber.Map{
		"Message":   "ดึงรายชื่อฐานข้อมูลสำเร็จ",
		"Status":    true,
		"Databases": databaseNames,
		// map[string][]string{
		//     loginRequest.Prefix: databaseNames,
		// },
	}
	return c.JSON(response)
}
func RootLogin(c *fiber.Ctx) error {

	loginRequest := new(models.Users)

	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// เชื่อมต่อกับฐานข้อมูลระบบ 'mysql'
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/mysql?charset=utf8mb4&parseTime=True&loc=Local", loginRequest.Username, loginRequest.Password, mysql_host)
	// fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		response := fiber.Map{
			"Message": "เข้าสู่ระบบไม่สำเร็จ!!",
			"Status":  false,
		}
		return c.JSON(response)
		//	return err
		// return fmt.Errorf("Failed to connect to MySQL: %v", err)
	}
	c.Locals("db", db)
	// fmt.Println("Login successful")

	response := fiber.Map{
		"Message": "เข้าสู่ระบบสำเร็จ!!",
		"Status":  true,
	}
	return c.JSON(response)

}

// ฟังก์ชันสำหรับสร้างฐานข้อมูลใหม่
func CreateDB(db *gorm.DB, dbName string) error {
	// ตรวจสอบว่าฐานข้อมูลที่ต้องการสร้างมีอยู่หรือยัง ถ้าไม่มีให้สร้าง
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)

	// รันคำสั่ง SQL เพื่อสร้างฐานข้อมูล
	if err := db.Exec(createDBQuery).Error; err != nil {
		return fmt.Errorf("Failed to create database: %v", err)
	}

	fmt.Printf("Database %s created successfully\n", dbName)
	return nil
}

// promotion

type promotiondata struct {
	Prefix string `json:"prefix"`
	Body   struct {
		Name               string              `json:"name"`
		Description        string              `json:"description"`
		PercentDiscount    decimal.NullDecimal `json:"percentDiscount"`
		StartDate          string              `json:"startDate"`
		EndDate            string              `json:"endDate"`
		MaxDiscount        decimal.NullDecimal `json:"maxDiscount"`
		UsageLimit         int                 `json:"usageLimit"`
		SpecificTime       string              `json:"specificTime"`
		PaymentMethod      string              `json:"paymentMethod"`
		MinSpend           decimal.NullDecimal `json:"minSpend"`
		MaxSpend           decimal.NullDecimal `json:"maxSpend"`
		TermsAndConditions string              `json:"termsAndConditions"`
		Status             int                 `json:"status"`
		Includegames       string              `json:"includegames"`
		Excludegames       string              `json:"excludegames"`
	} `json:"body"`

	PromotionId int `json:"promotionId"`
}

func CreatePromotion(c *fiber.Ctx) error {

	var data promotiondata

	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err)
		return c.JSON(fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!" + err.Error(),
			"Status":  false,
			"Data":    err.Error(),
		})
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: data.Prefix + "_development",
		production:  data.Prefix + "_production",
	}
	db, err := database.ConnectToDB(prefixs.development)
	if db == nil {

		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)

	}
	if err != nil {
		log.Fatal(err)
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)

	}
	// database.CheckAndCreateTable(db, models.Promotion{})

	promotion := models.Promotion{
		Name:               data.Body.Name,
		Description:        data.Body.Description,
		PercentDiscount:    data.Body.PercentDiscount.Decimal,
		StartDate:          data.Body.StartDate,
		EndDate:            data.Body.EndDate,
		MaxDiscount:        data.Body.MaxDiscount.Decimal,
		UsageLimit:         data.Body.UsageLimit,
		SpecificTime:       data.Body.SpecificTime,
		PaymentMethod:      data.Body.PaymentMethod,
		MinSpend:           data.Body.MinSpend.Decimal,
		MaxSpend:           data.Body.MaxSpend.Decimal,
		TermsAndConditions: data.Body.TermsAndConditions,
		Status:             data.Body.Status,
		Includegames:       data.Body.Includegames,
		Excludegames:       data.Body.Excludegames,
	}

	err = db.Create(&promotion).Error

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	response := fiber.Map{
		"Message": "สร้างข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    promotion,
	}

	return c.JSON(response)
}

func GetPromotion(c *fiber.Ctx) error {
	body := new(Dbstruct)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}
	//fmt.Printf("prefixs: %s",prefixs)
	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {

		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}
	//database.CheckAndCreateTable(db, models.Promotion{})
	err = db.AutoMigrate(&models.Promotion{})
	if err != nil {
		fmt.Println("err:", err)
	}
	promotions := []models.Promotion{}

	err = db.Debug().Where("status=1 and enddate>?", time.Now().Format("2006-01-02")).Find(&promotions).Error

	// fmt.Println(promotions)
	if err != nil {

		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}

	return c.JSON(fiber.Map{
		"Message": "ดึงข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    promotions,
	})

	//	return c.JSON(promotions)
}

func GetPromotionById(c *fiber.Ctx) error {
	body := new(promotiondata)
	if err := c.BodyParser(body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
			"error":   err.Error(),
		})
	}
	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}
	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {

		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}
	promotion := models.Promotion{}
	err = db.Debug().First(&promotion, body.PromotionId).Error
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "ดึงข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    promotion,
	}
	return c.JSON(response)

}

func UpdatePromotion(c *fiber.Ctx) error {
	data := new(promotiondata)
	if err := c.BodyParser(data); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	promotion := models.Promotion{
		Name:               data.Body.Name,
		Description:        data.Body.Description,
		PercentDiscount:    data.Body.PercentDiscount.Decimal,
		StartDate:          data.Body.StartDate,
		EndDate:            data.Body.EndDate,
		MaxDiscount:        data.Body.MaxDiscount.Decimal,
		UsageLimit:         data.Body.UsageLimit,
		SpecificTime:       data.Body.SpecificTime,
		PaymentMethod:      data.Body.PaymentMethod,
		MinSpend:           data.Body.MinSpend.Decimal,
		MaxSpend:           data.Body.MaxSpend.Decimal,
		TermsAndConditions: data.Body.TermsAndConditions,
		Status:             data.Body.Status,
		Includegames:       data.Body.Includegames,
		Excludegames:       data.Body.Excludegames,
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: data.Prefix + "_development",
		production:  data.Prefix + "_production",
	}
	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//err = db.AutoMigrate(&models.Promotion{})
	//db.AutoMigrate(&models.Promotion{});
	//AutoMigrate(&models.TsxAdmin{},&models.Provider{},&models.Promotion{});
	//promotion = models.Promotion{}
	err = db.Debug().Model(&promotion).Where("id = ?", data.PromotionId).Updates(promotion).Error
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!" + err.Error(),
			"Status":  false,
		}
		return c.JSON(response)
	}

	response := fiber.Map{
		"Message": "อัปเดตข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    promotion,
	}
	return c.JSON(response)
}

func DeletePromotion(c *fiber.Ctx) error {
	body := new(promotiondata)
	if err := c.BodyParser(body); err != nil {
		fmt.Println(err)
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!" + err.Error(),
			"Status":  false,
		}
		return c.JSON(response)
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}
	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}
	//db.AutoMigrate(&models.Promotion{})
	//fmt.Println(body)
	err = db.Debug().Delete(&models.Promotion{}, "id = ?", body.PromotionId).Error
	if err != nil {
		response := fiber.Map{
			"Message": "มีข้อผิดพลาดเกิดขึ้น!!",
			"Status":  false,
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "ลบข้อมูลสำเร็จ",
		"Status":  true,
	}
	return c.JSON(response)
}

// game
type gameData struct {
	Prefix string `json:"prefix"`
	ID     int    `json:"id"`
	Body   struct {
		ProductCode string `json:"productcode"`
		Product     string `json:"product"`
		GameType    string `json:"gameType"`
		Active      int    `json:"active"`
		Remark      string `json:"remark"`
		Position    string `json:"position"`
		Urlimage    string `json:"urlimage"`
		Name        string `json:"name"`
		Status      string `json:"status"`
	} `json:"body"`
}

func CreateGame(c *fiber.Ctx) error {
	data := new(gameData)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	game := models.Games{
		ProductCode: data.Body.ProductCode,
		Product:     data.Body.Product,
		GameType:    data.Body.GameType,
		Active:      data.Body.Active,
		Remark:      data.Body.Remark,
		Position:    data.Body.Position,
		Urlimage:    data.Body.Urlimage,
		Name:        data.Body.Name,
		Status:      data.Body.Status,
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: data.Prefix + "_development",
		production:  data.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//database.CheckAndCreateTable(db, models.Games{})
	err = db.AutoMigrate(&models.Games{})
	if err != nil {
		fmt.Println("err:", err)
	}
	err = db.Create(&game).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	response := fiber.Map{
		"Message": "สร้างข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    game,
	}
	return c.JSON(response)
}

func GetGameList(c *fiber.Ctx) error {
	body := new(Dbstruct)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//err = database.CheckAndCreateTable(db, models.Games{})
	err = db.AutoMigrate(&models.Games{})
	if err != nil {
		fmt.Println("err:", err)
	}

	games := []models.Games{}
	err = db.Debug().Find(&games).Error
	if err != nil {
		response := fiber.Map{
			"Message": "ดึงข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	if len(games) == 0 {

		sql := `
		INSERT INTO Games (product, productCode, gametype, active, status, remark, position, urlImage) VALUES
		(1017, 'TF Gaming', NULL, 1, '{"id":"13","name":"Esport"}', NULL, 'OK', NULL),
		(1009, 'CQ9', NULL, 1, '{"id":"8","name":"Fishing"}', NULL, 'OK', NULL),
		(1091, 'Jili', NULL, 1, '{"id":"8","name":"Fishing"}', NULL, 'OK', NULL),
		(1002, 'Evolution Gaming', NULL, 2, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1003, 'All Bet', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1004, 'Big Gaming', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1011, 'Play Tech', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1020, 'WM Casino', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1022, 'Sexy Gaming', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1052, 'Dream Gaming', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1077, 'SkyWind', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1053, 'Nexus 4D', NULL, 1, '{"id":"5","name":"Lottery"}', NULL, 'OK', NULL),
		(1074, 'HKGP Lottery', NULL, 1, '{"id":"5","name":"Lottery"}', NULL, 'OK', NULL),
		(1076, 'AMB Poker', NULL, 1, '{"id":"7","name":"p2p"}', NULL, 'OK', NULL),
		(1006, 'Pragmatic Play', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1009, 'CQ9', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1011, 'Play Tech', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1013, 'Joker', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1048, 'Reevo', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1049, 'EvoPlay', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1050, 'PlayStar', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1075, 'SlotXo', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1077, 'SkyWind', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1085, 'JDB', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1091, 'Jili', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1046, 'IBC', NULL, 1, '{"id":"3","name":"Sport Book"}', NULL, 'OK', NULL),
		(1081, 'BTI', NULL, 1, '{"id":"3","name":"Sport Book"}', NULL, 'OK', NULL),
		(1105, 'Royal Slot Gaming', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1110, 'Red Tiger', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1012, 'SBO', NULL, 1, '{"id":"3","name":"Sport Book"}', NULL, 'OK', NULL),
		(9999, 'GCLUB', NULL, 1, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(8888, 'PGSoft', NULL, 1, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1017, 'TF Gaming', NULL, 1, '{"id":"13","name":"Esport"}', NULL, 'OK', NULL),
		(1009, 'CQ9', NULL, 1, '{"id":"8","name":"Fishing"}', NULL, 'OK', NULL),
		(1091, 'Jili', NULL, 1, '{"id":"8","name":"Fishing"}', NULL, 'OK', NULL),
		(1002, 'Evolution Gaming', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1003, 'All Bet', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1004, 'Big Gaming', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1005, 'SA Gaming', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1011, 'Play Tech', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1020, 'WM Casino', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1022, 'Sexy Gaming', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1038, 'King 855', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1052, 'Dream Gaming', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1077, 'SkyWind', NULL, 0, '{"id":"2","name":"Live Casino"}', NULL, 'OK', NULL),
		(1053, 'Nexus 4D', NULL, 0, '{"id":"5","name":"Lottery"}', NULL, 'OK', NULL),
		(1074, 'HKGP Lottery', NULL, 0, '{"id":"5","name":"Lottery"}', NULL, 'OK', NULL),
		(1076, 'AMB Poker', NULL, 0, '{"id":"2","name":"Pp"}', NULL, 'OK', NULL),
		(1006, 'Pragmatic Play', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1009, 'CQ9', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1011, 'Play Tech', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1013, 'Joker', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1014, 'Dragon Soft', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1039, 'AMAYA', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1041, 'Habanero', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1048, 'Reevo', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1049, 'EvoPlay', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1050, 'PlayStar', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1075, 'SlotXo', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1077, 'SkyWind', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1084, 'Advant Play', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1085, 'JDB', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1091, 'Jili', NULL, 0, '{"id":"1","name":"Slot"}', NULL, 'OK', NULL),
		(1046, 'IBC', NULL, 0, '{"id":"3","name":"Sport Book"}', NULL, 'OK', NULL),
		(1081, 'BTI', NULL, 0, '{"id":"3","name":"Sport Book"}', NULL, 'OK', NULL)
		`

		err = db.Exec(sql).Error
		if err != nil {
			response := fiber.Map{
				"Message": "สร้างข้อมูลผิดพลาด",
				"Status":  false,
				"Data":    err.Error(),
			}
			return c.JSON(response)
		}
	}
	games = []models.Games{}
	err = db.Debug().Find(&games).Error
	if err != nil {
		response := fiber.Map{
			"Message": "ดึงข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "ดึงข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    games,
	}
	return c.JSON(response)
}

func GetGameById(c *fiber.Ctx) error {
	body := new(gameData)
	if err := c.BodyParser(body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//database.CheckAndCreateTable(db, models.Games{})
	err = db.AutoMigrate(&models.Games{})
	if err != nil {
		fmt.Println("err:", err)
	}
	game := models.Games{}
	err = db.Debug().First(&game, body.ID).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	response := fiber.Map{
		"Message": "ดึงข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    game,
	}
	return c.JSON(response)
}

func UpdateGame(c *fiber.Ctx) error {
	body := new(gameData)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//fmt.Println(body)
	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	database.CheckAndCreateTable(db, models.Games{})

	game := models.Games{}

	game.ProductCode = body.Body.ProductCode
	game.Product = body.Body.Product
	game.GameType = body.Body.GameType
	game.Active = body.Body.Active
	game.Remark = body.Body.Remark
	game.Position = body.Body.Position
	game.Urlimage = body.Body.Urlimage
	game.Name = body.Body.Name
	game.Status = body.Body.Status

	err = db.Debug().Model(&models.Games{}).Where("id = ?", body.ID).Updates(body.Body).Error
	if err != nil {
		response := fiber.Map{
			"Message": "อัปเดตข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "อัปเดตข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    body.Body,
	}
	return c.JSON(response)
}

func GetGameStatus(c *fiber.Ctx) error {

	type GameStatus struct {
		ProductCode string `json:"productCode"`
		Status      string `json:"status"`
	}

	body := new(gameData)
	if err := c.BodyParser(body); err != nil {
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
		}
		return c.JSON(response)
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check Redis for cached game status
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", // redis_master_password,
		DB:       0,  // Use database 0
	})

	cachedStatus, err := rdb.Get(ctx, "game_status").Result()
	if err == nil {
		// If cached data is found, return it
		type Status struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		// Define the main struct that includes the status
		type Product struct {
			ProductCode string `json:"productCode"`
			Status      Status `json:"status"`
		}
		var products []Product
		var tempProducts []struct {
			ProductCode string `json:"productCode"`
			Status      string `json:"status"` // Keep status as string for initial parsing
		}

		// Unmarshal the main JSON
		if err := json.Unmarshal([]byte(cachedStatus), &tempProducts); err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		// Step 2: Iterate through the temporary products and unmarshal the status
		for _, item := range tempProducts {
			var status Status
			if err := json.Unmarshal([]byte(item.Status), &status); err != nil {
				log.Fatalf("Error unmarshalling status JSON: %v", err)
			}
			products = append(products, Product{
				ProductCode: item.ProductCode,
				Status:      status,
			})
		}

		// Print the result
		// for _, product := range products {
		// 	fmt.Printf("Product Code: %s, Status: %+v\n", product.ProductCode, product.Status)
		// }

	}

	// If no cached data, query the database
	gameStatus := []GameStatus{}
	err = db.Debug().Table("Games").Select("DISTINCT status").Find(&gameStatus).Error
	if err != nil {
		response := fiber.Map{
			"Message": "ดึงข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

	// Cache the result in Redis with a 1-day expiration
	statusJSON, err := json.Marshal(gameStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error converting to JSON")
	}

	err = rdb.Set(ctx, "game_status", statusJSON, 24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error caching game status")
	}

	response := fiber.Map{
		"Message": "ดึงข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    gameStatus,
	}
	return c.JSON(response)
}

func GetMemberList(c *fiber.Ctx) error {
	body := new(Dbstruct)
	if err := c.BodyParser(body); err != nil {
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		response := fiber.Map{
			"Message": "ติดต่อฐานข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	//database.CheckAndCreateTable(db, models.Users{})
	games := []models.Users{}
	err = db.Debug().Find(&games).Error
	if err != nil {
		response := fiber.Map{
			"Message": "ดึงข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "ดึงข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    games,
	}
	return c.JSON(response)
}

type MemberBody struct {
	Prefix string `json:"prefix"`
	ID     int    `json:"id"`
	Body   struct {
		Fullname   string `json:"fullname"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Status     int    `json:"status"`
		Bankname   string `json:"bankname"`
		Banknumber string `json:"banknumber"`
		ProStatus  string `json:"prostatus"`
	}
}

func CreateMember(c *fiber.Ctx) error {

	body := new(MemberBody)
	if err := c.BodyParser(body); err != nil {
		fmt.Println(err)
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		response := fiber.Map{
			"Message": "ติดต่อฐานข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

	member := models.Users{}
	member.Username = body.Body.Username
	member.Password = body.Body.Password
	member.Status = body.Body.Status
	member.Role = "user"

	err = db.Debug().Create(&member).Error
	if err != nil {
		response := fiber.Map{
			"Message": "สร้างข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "สร้างข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    body.Body,
	}
	return c.JSON(response)
}

func GetMemberById(c *fiber.Ctx) error {
	body := new(MemberBody)
	if err := c.BodyParser(body); err != nil {
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		response := fiber.Map{
			"Message": "ติดต่อฐานข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	//database.CheckAndCreateTable(db, models.Users{})
	user := models.Users{}
	err = db.Debug().First(&user, body.ID).Error
	if err != nil {
		response := fiber.Map{
			"Message": "ดึงข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "ดึงข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    user,
	}
	return c.JSON(response)
}

func UpdateMember(c *fiber.Ctx) error {

	body := new(MemberBody)
	if err := c.BodyParser(body); err != nil {
		fmt.Println(err)
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	var prefixs = struct {
		development string
		production  string
	}{
		development: body.Prefix + "_development",
		production:  body.Prefix + "_production",
	}

	db, err := database.ConnectToDB(prefixs.development)
	if err != nil {
		response := fiber.Map{
			"Message": "ติดต่อฐานข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

	member := models.Users{}
	member.Username = body.Body.Username
	member.Password = body.Body.Password
	member.Status = body.Body.Status
	member.Bankname = body.Body.Bankname
	member.Banknumber = body.Body.Banknumber
	member.ProStatus = body.Body.ProStatus

	err = db.Debug().Model(&member).Where("id = ?", body.ID).Updates(body.Body).Error
	if err != nil {
		response := fiber.Map{
			"Message": "อัพเดตข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	response := fiber.Map{
		"Message": "อัพเดตข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    body.Body,
	}
	return c.JSON(response)

}

func GetExchangeRates(c *fiber.Ctx) error {
	// เชื่อมต่อ Redis
	type ExchangeRateBody struct {
		Currency string `json:"currency"`
	}

	cbody := new(ExchangeRateBody)
	if err := c.BodyParser(cbody); err != nil {
		fmt.Println(err)
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", //redis_master_password,
		DB:       0,  // ใช้ database 0
	})

	// ทดสอบการเชื่อมต่อ Redis
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error connecting to Redis")
		// อาจจะ return error หรือจัดการข้อผิดพลาดตามที่คุณต้องการ
	} else {
		fmt.Println("เชื่อมต่อ Redis สำเร็จ:", pong)
	}

	// ตรวจสอบว่ามีข้อมูลใน Redis หรือไม่
	cachedRates, err := rdb.Get(ctx, "exchange_rates").Result()
	if err == nil {
		// ถ้ามีข้อมูลใน Redis ส่งกลับเลย
		fmt.Println("ถ้ามีข้อมูลใน Redis ส่งกลับเลย")
		return c.SendString(cachedRates)
	}
	//fmt.Println(cachedRates)
	fmt.Println("ถ้าไม่มีข้อมูลใน Redis ดึงข้อมูลจาก API")
	// ถ้าไม่มีข้อมูลใน Redis ดึงข้อมูลจาก API
	resp, err := http.Get("https://api.exchangerate-api.com/v4/latest/" + cbody.Currency)
	if err != nil {

		// ถ้าไม่สามารถเข้าถึง API ได้ ลองดึงข้อมูลเก่าจาก Redis
		oldRates, err := rdb.Get(ctx, "old_exchange_rates").Result()
		if err == nil {
			fmt.Println("ถ้าไม่สามารถเข้าถึง API ได้ ลองดึงข้อมูลเก่าจาก Redis")
			return c.SendString(oldRates)
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching exchange rates")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading response body")
	}
	defer resp.Body.Close()

	var exchangeRates ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRates)
	if err != nil {
		// Log the error and the response body for debugging
		fmt.Printf("Error parsing JSON: %v\n", err)
		fmt.Printf("Response body: %s\n", string(body))
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing JSON")
	}

	// แปลงข้อมูลกลับเป็น JSON string
	ratesJSON, err := json.Marshal(exchangeRates)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).SendString("Error converting to JSON")
	}

	// บันทึกข้อมูลลง Redis พร้อมกำหนดเวลาหมดอายุ 1 ชั่วโมง
	err = rdb.Set(ctx, "exchange_rates", ratesJSON, 24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error caching exchange rates")
	}

	// บันทึกข้อมูลเดียวกันไว้เป็นข้อมูลสำรองโดยไม่มีกำหนดเวลาหมดอายุ
	err = rdb.Set(ctx, "old_exchange_rates", ratesJSON, 0).Err()
	if err != nil {
		// Log error, but don't return it to the user
		// log.Printf("Error caching old exchange rates: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error caching old exchange rates")
	}

	// ส่งข้อมูลกลับไปยังผู้ใช้
	return c.SendString(string(ratesJSON))
}

func UpdateMaster(c *fiber.Ctx) error {

	type MasterBody struct {
		Prefix string          `json:"prefix"`
		ID     int             `json:"id"`
		Body   models.Settings `json:"body"`
	}
	body := new(MasterBody)
	if err := c.BodyParser(body); err != nil {
		response := fiber.Map{
			"Message": "รับข้อมูลผิดพลาด",
			"Status":  false,
			"Data":    err.Error(),
		}
		return c.JSON(response)
	}
	var settings models.Settings

	db := createDatabase("master")
	if db == nil {
		response := fiber.Map{
			"Message": "ติดต่อฐานข้อมูลผิดพลาด",
			"Status":  false,
		}

		return c.JSON(response)
	}

	db.AutoMigrate(&models.Settings{})
	var RowsAffected int64
	db.Debug().Model(&settings).Select("id").Scan(&settings).Count(&RowsAffected)
	if RowsAffected == 0 {
		db.Debug().Create(&body.Body)
	} else {
		db.Debug().Model(&settings).Updates(body.Body)
	}
	response := fiber.Map{
		"Message": "อัพเดตข้อมูลสำเร็จ",
		"Status":  true,
		"Data":    body.Body,
	}
	return c.JSON(response)
}
