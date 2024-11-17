package handler

import (
	"context"
	"fmt"
	"github.com/amalfra/etag"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
	"github.com/valyala/fasthttp"
	_ "github.com/go-sql-driver/mysql"
	"pkd/models"
	//"github.com/golang-jwt/jwt"
	jtoken "github.com/golang-jwt/jwt/v4"
	"github.com/shopspring/decimal"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	"pkd/database"
	"pkd/repository"
	"pkd/common" 
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"encoding/base64"
    "strings"
	"encoding/json"
)

var ctx = context.Background()
var amqp_connection *amqp.Connection
var amqp_channel *amqp.Channel
var queue amqp.Queue = amqp.Queue{}
var is_connection = false
var has_channel = false
var has_queue = false

var redis_master_host = os.Getenv("REDIS_HOST")
var redis_master_port = os.Getenv("REDIS_PORT")
var redis_master_password = os.Getenv("REDIS_PASSWORD")
var redis_slave_host = os.Getenv("REDIS_SLAVE_HOST")
var redis_slave_port = os.Getenv("REDIS_SLAVE_PORT")
var redis_slave_password = os.Getenv("REDIS_SLAVE_PASSWORD")
var queue_name = "wallet" //os.Getenv("QUEUE_NAME")
var exchange_name = "wallet" //os.Getenv("EXCHANGE_NAME")
var rabbit_mq = "amqp://128.199.92.45:5672" //os.Getenv("RABBIT_MQ") @rabbitmq:5672/wallet
var connection_timeout = os.Getenv("CONNECTION_TIMEOUT")
var redis_database = getEnv("REDIS_DATABASE", "0")
var go_pixel_log = os.Getenv("GO_PIXEL_LOG")

var jwtSecret = os.Getenv("PASSWORD_SECRET")


func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
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

// func MakeGPostRequest(url string,jsonData string) (*fasthttp.Response, error) {
	
	
// 	//var encrypt = GEncrypt(jsonString)

// 	// สร้าง Request โดยใช้ fasthttp
// 	req := fasthttp.AcquireRequest()
// 	defer fasthttp.ReleaseRequest(req)

// 	// ตั้งค่า HTTP method เป็น POST
// 	req.Header.SetMethod("POST")
// 	req.SetRequestURI(url)
	
// 	//encrypt := GEncrypt(jsonData)

// 	// กำหนด headers (X-API-ClientID, X-API-Signature, X-API-Timestamp)
// 	req.Header.Set("X-API-ClientID", CLIENT_ID) // กำหนด Client ID ของคุณ
// 	req.Header.Set("X-API-Signature", encrypt.Enc) // กำหนด Signature ของคุณ
// 	req.Header.Set("X-API-Timestamp", string(encrypt.Unx)) // กำหนดค่า Timestamp ปัจจุบัน
// 	req.Header.Set("Content-Type", "application/json")

// 	// Optional: เพิ่ม referrer หากต้องการ
// 	// if referrer != "" {
// 	// 	req.Header.Set("Referer", referrer)
// 	// }

// 	// ตั้งค่า body ของ request เป็นค่า encryptedData
// 	req.SetBody([]byte(encrypt.Des))

// 	// สร้าง Response
// 	resp := fasthttp.AcquireResponse()

// 	// ส่ง Request โดยใช้ fasthttp client
// 	client := &fasthttp.Client{}
// 	err := client.Do(req, resp)
// 	if err != nil {
// 		// คืนค่า nil และ error หากมีปัญหาในการส่ง request
// 		return nil, fmt.Errorf("error sending request: %v", err)
// 	}

// 	// คืนค่า response และ nil (ไม่มี error)
// 	return resp, nil
// }

type EFResponse struct {
	ErrorCode  int  `json:"errorcode"`
    ErrorMessage string `json:"errormessage"`
    Balance  decimal.Decimal `json:"balance"`
    BeforeBalance decimal.Decimal `json:"beforebalance"`
}

type ResponseBalance struct {
	BetAmount decimal.Decimal `json:"betamount"`
	BeforeBalance decimal.Decimal `json:"beforebalance"`
	Balance decimal.Decimal `json:"balance"`
}
func fastPost(url, headstr string, encoded *fiber.Map) (*fasthttp.Response, error) {
	if encoded == nil {
		return nil, fmt.Errorf("encoded map cannot be nil")
	}

	method := "POST"

	// สร้าง request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	// ตั้งค่า headers
	req.Header.Set("Authorization", headstr) // ตั้งค่า Authorization header
	req.Header.Set("Content-Type", "application/json") // ตั้งค่า Content-Type header

	// แปลง encoded เป็น JSON
	body, err := json.Marshal(encoded) // แปลงเป็น JSON
	if err != nil {
		return nil, fmt.Errorf("error marshalling encoded map: %v", err)
	}
	req.SetBody(body) // ตั้งค่า body ของ request

	// สร้าง response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// ส่ง request
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	return resp, nil
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

	//req.SetBody([]byte(encryptedData))

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
	if !is_connection {
		InitAMQP()
	}

	fmt.Println("is close")
	fmt.Println(amqp_connection.IsClosed())

	if amqp_connection.IsClosed() {
		is_connection = false
		InitAMQP()
	}
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

func Signup(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	message := fmt.Sprintf("I am GROOT!")
	return c.SendString(message)
}

func Login(c *fiber.Ctx) error {
	 



	loginRequest := new(models.Users)

	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	 
	user,err := repository.FindUser(loginRequest.Preferredname,loginRequest.Password)
	if err != nil {
		response := Response{
			Message: "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
	
		return c.JSON(response)
	}

	
	//day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID": user.ID,
		"Walletid": user.Walletid,
		"Username": user.Username,
		"Role": user.Role,
		"PartnersKey": user.PartnersKey,
		//"exp":   time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256,claims)

	t,err_ := token.SignedString([]byte(jwtSecret))
	
	
	if err_ != nil {
		response := Response{
			Message: "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	}
	updates := map[string]interface{}{
		"Token": t,
			}

	// อัปเดตข้อมูลยูสเซอร์
	_err := repository.UpdateFieldsUser(user.ID, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	if _err != nil {
		response := Response{
			Message: "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	} else {
		response := Response{
			Message: "เข้าระบบสำเร็จ!",
			Status:  true,
			Data: models.LoginResponse{  
				Token: t, 
				},
		}
		return c.JSON(response)
	}

	
	

	

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
	 
	//if claims != nil {
		updates := map[string]interface{}{
			"Token": "",
				}
	
		// อัปเดตข้อมูลยูสเซอร์
		 repository.UpdateFieldsUserString(username, updates) 

		response := Response{
			Message: "ออกจากระบบสำเร็จ!",
			Status:  true,
			Data: fiber.Map{ 
				"id": -1,
			},
		}
		return c.JSON(response)
	// } else {
	// 	response := Response{
	// 		Message: "มีข้อผิดพลาด กรุณาออกจากระบบอีกครั้ง!",
	// 		Status:  false,
	// 		Data: fiber.Map{ 
	// 			"id": -1,
	// 		},
	// 	}
	// 	return c.JSON(response)
	// }
	 
}
type TransactionRequest struct {
    Status          int                 `json:"status"`
    GameProvide     string                  `json:"gameProvide"`
    MemberName      string                  `json:"memberName"`
    TransactionAmount decimal.Decimal       `json:"transactionAmount"`
	PayoutAmount      decimal.Decimal       `json:"payoutAmount"`
	BetAmount         decimal.Decimal       `json:"betAmount"`
    ProductID         int64                     `json:"productID"`
    BeforeBalance   string                  `json:"beforeBalance"`
    Balance         decimal.Decimal                  `json:"balance"`
    AfterBalance    string                  `json:"afterBalance"`
}

func createBasicAuthHeader(operatorCode, secretKey string) string {
    // Combine the operator code and secret key
    credentials := fmt.Sprintf("%s:%s", operatorCode, secretKey)

    // Encode the credentials to base64
    encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

    // Return the Basic Auth header
    return "basic " + encodedCredentials
}

// var SECRET_KEY = os.Getenv("PASSWORD_SECRET")
// var pg_prod_code = os.Getenv("PG_PRODUCT_ID")

// var OPERATOR_CODE = "sunshinetest" //"sunshinepgthb"//"sunshinetest",
// var SECRET_API_KEY = os.Getenv("PG_API_KEY") //"9dc857f4-2225-45ef-bf0f-665bcf7d4a1b" //os.Getenv("PG_API_KEY")
// var PG_PROD_CODE= os.Getenv("PG_PRODUCT_ID")
// var PG_API_URL = "https://test.ambsuperapi.com"//os.Getenv("PG_API_URL") //"https://prod_md.9977997.com"
// var PG_PROD_URL = "https://api.hentory.io" 


func makePostRequest(url string, bodyData interface{}) (*fasthttp.Response, error) {
	// Marshal requestData struct เป็น JSON
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// สร้าง Request และ Response
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// ตั้งค่า URL, Method, และ Body
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	authHeader := createBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
	req.Header.Add("Authorization", authHeader)
	req.SetBody(jsonData)

	// ส่ง request
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)
	
	return resp, nil
}
func makeGetRequest(url string) (*fasthttp.Response, error) {
	// Marshal requestData struct เป็น JSON
	// jsonData, err := json.Marshal(bodyData)
	// if err != nil {
	// 	return nil, fmt.Errorf("error marshaling JSON: %v", err)
	// }

	// สร้าง Request และ Response
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// ตั้งค่า URL, Method, และ Body
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	authHeader := createBasicAuthHeader(common.OPERATOR_CODE, common.SECRET_API_KEY)
	req.Header.Add("Authorization", authHeader)
	//req.SetBody(jsonData)

	// ส่ง request
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)
	
	return resp, nil
}




func GetGameList(c *fiber.Ctx) error {
	type BodyGame struct {
		ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	}

	type PgRequest struct {
		Id string `json:"id"`
		TimestampMillis int `json:"timestampmillis"`
		ProductID string `json:"productid`
		Currency string `json:"currency"`
		Username string `json:"username"`
		SessionToken string `json:"sessiontoken"`
		StatusCode  int  `json:"statuscode"`
		Balance  decimal.Decimal `json:"balance"`
		//ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	//	Txns []TxnsRequest `json:"txns"`
	}
	type CResponse struct {
		Message string      `json:"message"`
		Status  bool        `json:"status"`
		Data    interface{} `json:"data"`  
	}

	type CEResponse struct {
		Message interface{}  `json:"message"`
		Status  bool        `json:"status"`
		Data    interface{} `json:"data"`  
	}

	var response CEResponse
	// bodyRequest := new(BodyGame)

	// if err := c.BodyParser(&bodyRequest); err != nil {
	// 	fmt.Printf(" %s ", err.Error())
	// 	response := fiber.Map{
	// 		"Status":  false,
	// 		"Message": err.Error(),
	// 	}
	// 	return c.JSON(response)
	// }

	// fmt.Printf("Body: %s",bodyRequest.Body)
	//var tokenString := c.Get("Authorization")[7:]
	request := new(PgRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//var users models.Users
	//users = ValidateJWTReturn(request.SessionToken);

	//fmt.Printf("users: %v ",users)
	//fmt.Printf("request: %s ",request.SessionToken)

	// var args = fiber.Map{
	// 	"username": strings.ToLower(users.Username),//user.data.username,
	// 	"productId":pg_prod_code,
	// 	"gameCode": request.ProductID,
	// 	"isMobileLogin": true,
	// 	"sessionToken": request.SessionToken,
	// 	//"betLimit": [],
	// 	"callbackUrl":"https://www.โชคดี789.com/lobby/slot/game?id=8888&type=1", //`${req.protocol}://${req.get('host')}${req.originalUrl}`
	// }
	
	
	
	resp,err := makeGetRequest(common.PG_API_URL+"/seamless/games?productId=PGSOFT2")		
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	resultBytes := resp.Body()
	resultString := string(resultBytes)
	// แสดงผล string ที่ได้
	//fmt.Println("Response body as string:", resultString)

	

	err = json.Unmarshal([]byte(resultString), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	//url := fmt.Sprintf(PG_PROD_URL,"/seamless/login")

	respon := fiber.Map{
		"Status":  true,
		"Message": response.Message,
		"Data": response.Data,
	}
	return c.JSON(respon)
}


func LaunchGame(c *fiber.Ctx) error {
	type BodyGame struct {
		ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	}

	type PgRequest struct {
		Id string `json:"id"`
		TimestampMillis int `json:"timestampmillis"`
		ProductID string `json:"productid`
		Currency string `json:"currency"`
		Username string `json:"username"`
		SessionToken string `json:"sessiontoken"`
		StatusCode  int  `json:"statuscode"`
		Balance  decimal.Decimal `json:"balance"`
		//ProductID string  `json:"productid"`
		LanguageCode string `json:"languagecode"`
		Platform string `json:"platform"`
		GameID string `json:"gameid"`
		GameType string `json:"gametype"`
		callbackUrl string `json:"callbackurl"`
	//	Txns []TxnsRequest `json:"txns"`
	}
	type CResponse struct {
		Message string      `json:"message"`
		Status  bool        `json:"status"`
		Data    interface{} `json:"data"`  
	}

	var response CResponse
	// bodyRequest := new(BodyGame)

	// if err := c.BodyParser(&bodyRequest); err != nil {
	// 	fmt.Printf(" %s ", err.Error())
	// 	response := fiber.Map{
	// 		"Status":  false,
	// 		"Message": err.Error(),
	// 	}
	// 	return c.JSON(response)
	// }

	// fmt.Printf("Body: %s",bodyRequest.Body)
	//var tokenString := c.Get("Authorization")[7:]
	request := new(PgRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	var users models.Users
	users = ValidateJWTReturn(request.SessionToken);

	//fmt.Printf("users: %v ",users)
	//fmt.Printf("request: %s ",request.SessionToken)
	// efargs = {
	// 	"OperatorCode": OPERATOR_CODE,
	// 	"MemberName": req.body.username,
	// 	"Password":   response.data.uid,
	// 	"ProductID": ProductID,
	// 	"GameType": GameType,
	// 	"GameID": GameID,
	// 	"LanguageCode": LanguageCode,
	// 	"Platform": Platform,
	// 	"Sign": hashSignature("LaunchGame",RequestTime),
	// 	"RequestTime": RequestTime
	// 	}
	var args = fiber.Map{
		"username": strings.ToLower(users.Username),//user.data.username,
		"productId":common.PG_PROD_CODE,
		"gameCode": request.ProductID,
		"isMobileLogin": true,
		"sessionToken": request.SessionToken,
		//"betLimit": [],
		"callbackUrl":"https://www.โชคดี789.com/lobby/slot/game?id=8888&type=1", //`${req.protocol}://${req.get('host')}${req.originalUrl}`
	}
	
	//fmt.Printf(" args : %s ",args)
	
	resp,err := makePostRequest(common.PG_API_URL+"/seamless/login",args)		
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	resultBytes := resp.Body()
	resultString := string(resultBytes)
	// แสดงผล string ที่ได้
	fmt.Println("Response body as string:", resultString)

	err = json.Unmarshal([]byte(resultString), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	//url := fmt.Sprintf(PG_PROD_URL,"/seamless/login")
// 	fmt.Printf("args % v \n",args)
// 	c.Append("content-type", "text/javascript")
// 	//OPERATOR_CODE := "your_operator_code"
//     //SECRET_KEY := "your_secret_key"
// 	c.Append("maxBodyLength","Infinity")
//     authHeader := createBasicAuthHeader(OPERATOR_CODE, SECRET_API_KEY)
//     fmt.Printf("authheader: %s \n",authHeader)
// 	c.Append("Authorization",authHeader )

// 	headers := c.Request().Header

// 	var headerString string
// 	headers.VisitAll(func(key, value []byte) {
// 		headerString += string(key) + ": " + string(value) + "\n"
// 	})

// //Basic c3Vuc2hpbmVwZ3RoYjo5ZGM4NTdmNC0yMjI1LTQ1ZWYtYmYwZi02NjViY2Y3ZDRhMWI=
// //basic c3Vuc2hpbmVwZ3RoYjo5ZGM4NTdmNC0yMjI1LTQ1ZWYtYmYwZi02NjViY2Y3ZDRhMWI=

// 	resp, err := fastPost(PG_API_URL+"/seamless/login", headerString,&args)
// 	fmt.Println(PG_API_URL+"/seamless/login")
// 	fmt.Printf("err: %v ",err)
// 	if err != nil {
// 		response := fiber.Map{
// 			"Status":  false,
// 			"Message": err.Error(),
// 		}
// 		return c.JSON(response)
// 	}
	// const gameurl = await fetch(`${PG_PROD_URL}/seamless/login`,
	// {
	// 	method: 'POST',
	// 	maxBodyLength: Infinity,
	// 	body: JSON.stringify(args),
	// 	headers: {'Content-Type': 'application/json','Authorization': 'basic ' + new Buffer.from(`${OPERATOR_CODE}:${SECRET_KEY}`).toString('base64')}
	// })


   // logger.info('args:'+JSON.stringify(args))
	//const gameurl = await fetch(`${PG_PROD_URL}/seamless/login`,
	respon := fiber.Map{
		"Status":  true,
		"Message": response.Message,
		"Data": response.Data,
	}
	return c.JSON(respon)
}

func AddTransactions(c *fiber.Ctx) error {

 

	type bodyRequest struct {
		Body TransactionRequest
	}
	
	//fmt.Println("transactionsub:",transactionsub)
	//Body := make(map[string]interface{})
	transactionRequest := new(bodyRequest)

	if err := c.BodyParser(&transactionRequest); err != nil {
		fmt.Printf(" %s ", err.Error())
		response := fiber.Map{
			"Status":  false,
			"Message": err.Error(),
		}
		return c.JSON(response)
	}

	fmt.Printf("Body: %s",transactionRequest.Body)

	
	// transactionRequest := new(TransactionRequest)
	// if err := c.BodyParser(body); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"Status": false,
	// 		"Message": "กรุณาตรวจสอบข้อมูลอีกครั้ง!",
	// 		"Data": map[string]interface{}{
	// 			"id": -1,
	// 		},
	// 	})
	// }

	transactionsub := models.TransactionSub{}

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
    if err_ := database.Database.Where("username = ? ", transactionRequest.Body.MemberName).First(&users).Error; err_ != nil {
		response = Response{
			Status: false,
			Message: "ไม่พบข้อมูล",
			Data:map[string]interface{}{
				"id": -1,
			},
    	}
	}

	

    transactionsub.Status = transactionRequest.Body.Status
    transactionsub.GameProvide = transactionRequest.Body.GameProvide
    transactionsub.MemberName = transactionRequest.Body.MemberName
	transactionsub.ProductID = transactionRequest.Body.ProductID
	transactionsub.BetAmount = transactionRequest.Body.BetAmount
	transactionsub.PayoutAmount = transactionRequest.Body.PayoutAmount
	transactionsub.BeforeBalance = users.Balance
	transactionsub.Balance = transactionRequest.Body.Balance

	//fmt.Println(transactionsub)

	result := database.Database.Debug().Create(&transactionsub); 
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
	
		 
		  repository.UpdateFieldsUserString(users.Username, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
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
				BetAmount: transactionsub.BetAmount,
			},
		}
		
	}
 
	return c.JSON(response)
}



