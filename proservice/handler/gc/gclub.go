package gc

import 
(
	"github.com/gofiber/fiber/v2"
	"hanoi/models"
	"hanoi/database"
	//"pkd/handler"
	"hanoi/repository"
	"hanoi/encrypt"
	//"crypto/md5"
	//"crypto/des"
	//"crypto/cipher"
	//"bytes"
	//"encoding/base64"
	"encoding/json"
	//"encoding/hex"
	//"encoding/json"
	//"pkd/repository"
	//"github.com/shopspring/decimal"
	//jtoken "github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4"
	"fmt"
	"time"
	"log"
	//"strconv"
	//"os"
	"strings"
	"github.com/valyala/fasthttp"
)
var (
	CLIENT_ID     = "6342e1be-fa03-456f-8d2d-8e1c9513c351" 
	CLIENT_SECRET = "6d83ac42"
	SYSTEMCODE    = "ckthb"
	WEBID         = "ckdthb"
	DESKEY 		  = "9c62a148"
	DESIV 		  =	"8e014099"
)
type Response struct {
    Message string      `json:"message"`
    Status  bool        `json:"status"`
    Data    interface{} `json:"data"` // ใช้ interface{} เพื่อรองรับข้อมูลหลายประเภทใน field data
}
type GResponse struct {
	MsgID    int    `json:"msgId"`
	Message  string `json:"message"`
	Data     GResponseData   `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

// Struct สำหรับข้อมูลใน field "data"
type GResponseData struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}
type GRequestData struct  {
	SystemCode string `json:"systemcode"`
	WebId      string  `json:"webid"`
	DataList   []interface{} `json:"datalist"`
}

type RequestEncrypt struct {
	Data  string `json:"data"`
	Key string `json:"key"`
	Iv string `json:"iv"`
}
 
type GResult struct {
	MsgId int `json:"msgid"`
	Message string `json:"message"`
	Data GResponseData `json:"data"`
	Timestamp int `json:"timestamp"`
}
type GClaims struct {
	SystemCode    string `json:"systemcode"`
	WebID         string `json:"webid"`
	MemberAccount string `json:"memberaccount"`
	TokenType     string `json:"tokentype"`
	jwt.RegisteredClaims
}

type GcRequest struct {
	Systemcode string `json:"systemcode"`
	Webid string `json:"webid"`
	Account string `json:"account"`
	Requestid string `json:"requestid"`
	Token string `json:"token"`
	Username string `json:"username"`
	
	// Id string `json:"id"`
	// TimestampMillis int `json:"timestampmillis"`
	// ProductID string `json:"productid`
	// Currency string `json:"currency"`
	// Username string `json:"username"`
	// SessionToken string `json:"sessiontoken"`
	// StatusCode  int  `json:"statuscode"`
	// Balance  decimal.Decimal `json:"balance"`
	//Txns []TxnsRequest `json:"txns"`
}

// ฟังก์ชันตัวอย่างใน gclub.go
type LoginRquest struct {
	BackUrl        string `json:"BackUrl"`
	GroupLimitID   string `json:"GroupLimitID"`
	ItemNo         string `json:"ItemNo"`
	Lang           string `json:"Lang"`
	MemberAccount  string `json:"MemberAccount"`
	SystemCode     string `json:"SystemCode"`
	WebId          string `json:"WebId"`
}



var API_URL_G = "http://rcgapiv2.rcg666.com/"
var API_URL_PROXY = "http://api.tsxbet.info:8001"




func Index(c *fiber.Ctx) error {

	//var user []models.Users
	
	response := Response{
		Status: true,
		Message: "OK",
		Data:   []interface{}{},
		} 
	 
	// database.Database.Find(&user)
	// response = GResponse{
	// 	MsgId: 0,
	// 	Message: "OK",
	// 	Data: GData{
	// 	  SystemCode: "DocDemoSystem",
	// 	  WebId: "DocDemoWeb",
	// 	  DataList: []interface{}{},
	// 	},
	// 	TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
	// }
	 
	return c.JSON(response)
   
}


 
func CheckUser(c *fiber.Ctx) error {
	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}

	var authToken,_ = Gsign(request.Account,24,"AuthToken")
	var sessionToken,_ = Gsign(request.Account,1,"SessionToken")
	
	var users models.Users
	//users = handler.ValidateJWTReturn(request.SessionToken);
	db, _ := database.ConnectToDB(request.Account)
	var rowsAffected = db.Debug().Where("username = ? AND g_token = ?", strings.ToUpper(request.Account),request.Token).First(&users).RowsAffected
  
	if rowsAffected == 0 {
		var response = fiber.Map{
			"msgId": 4,
			"message": "Invalid GameToken",
			"data": fiber.Map{
				"requestId": request.Requestid,
				"account": request.Account,
				"token": request.Token,
			},
			"timestamp": time.Now().Unix(),
		  }
		  return c.Status(400).JSON(response)
	} else {
		var response = fiber.Map{
			"msgId": 0,
			"message": "OK",
			"data": fiber.Map{
				"requestId": request.Requestid,
				"account": request.Account,
				"authToken": authToken,
				"sessionToken": sessionToken,
			},
			"timestamp": time.Now().Unix(),
			}
		return c.Status(200).JSON(response)

	}	
}
func GetBalance(c *fiber.Ctx) error {

	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//var users models.Users

	tokenString := c.Get("Authorization")[7:] 
	claims, _ :=  Gverify(tokenString)	
	fmt.Println(claims)
	//users,_err = handler.Gverify(tokenString);
	var user models.Users
	db, _ := database.ConnectToDB(claims.MemberAccount)
	db.Where("username = ?", claims.MemberAccount).First(&user)
    balanceFloat, _ := user.Balance.Float64()
	var response = fiber.Map{
		"msgId": 0,
		"message": "OK",
		"data": fiber.Map{
			"status": 0,
			"requestId": request.Requestid,
			"account": claims.MemberAccount,
			"balance": balanceFloat,
		},
		"timestamp": time.Now().Unix(),
		}

	return c.Status(200).JSON(response)
	
}
func Test(c* fiber.Ctx) error {

	 // "ak+xb8pip08kqqijH/vcAYZ56//9nZWqm/Tu7E2ZpjL4zaHQo91QP+F6wbsZfEhgAH02smpi470="

	 
	data := encrypt.Data{
		Token:  "VBBF",
		Amount: 100.0,
		TranID: "T1_1700154",
	}

	// แปลง struct เป็น JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}



	encode,_err_ := encrypt.EncryptDESAndMD5(string(jsonData),"12345678","98765432","CF7861C7-556F-499A-890C-F9C7C4190266","p@ssw0rd")
	if _err_ != nil {
		fmt.Println(_err_)
	}
	
	//encode := &encrypt.ECResult{} // สมมุติว่า encode เป็น pointer

	result := *encode 
	
	fmt.Printf("Encrypted Data: %+v\n", result)
	//eresult := encrypt.ECResult(encode)
	
	decode,_err := encrypt.DecryptDES(result.Des,"12345678","98765432")
	if _err != nil {
		fmt.Printf("Error Data: %+v\n", _err)
	} else {
		fmt.Printf("Decrypted Data: %+v\n", decode)
	}
	return c.JSON("Test")
}
func Login(c *fiber.Ctx) (error) {
	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "โทเคนหมดอายุ หรือไม่ถูกต้อง",
		})
	}
	var users models.Users

	tokenString := c.Get("Authorization")[7:] 
	claims, _ :=  Gverify(tokenString)
	
	if claims != nil {
	fmt.Println(claims)
	}
	db, _ := database.ConnectToDB(request.Account)
	db.Where("username = ?", strings.ToUpper(request.Account)).First(&users)
	//fmt.Println(users)
	loginResponse,err := loging(request.Account)

	//loginResponse, err := parseLoginResponse(responseString)
	if err != nil {
		
		log.Fatal("Post Error",err)
	 	 
	}

	// แสดงผล
	// fmt.Printf("MsgID: %s\n", loginResponse.MsgID)
	// fmt.Printf("Message: %s\n", loginResponse.Message)
	// fmt.Printf("Data: %s\n", loginResponse.Data)
	// fmt.Printf("TimeStamp: %s\n", loginResponse.Timestamp)
	// fmt.Printf("GroupLimitID: %s\n", loginResponse.GroupLimitID)
	// fmt.Printf("ItemNo: %s\n", loginResponse.ItemNo)
	// fmt.Printf("Lang: %s\n", loginResponse.Lang)
	// fmt.Printf("MemberAccount: %s\n", loginResponse.MemberAccount)
	// fmt.Printf("SystemCode: %s\n", loginResponse.SystemCode)
	// fmt.Printf("WebId: %s\n", loginResponse.WebId)
 
	//  if err != nil {
	// 	log.Fatal("Post Error",err)
	 	 
	//  	} else {
	// 	str_resp := string(resp.Body())
	// 	desenc_str,_xerr := encrypt.DecryptDES(str_resp,handler.DESKEY,handler.DESIV)
		
	// 	if _xerr != nil {
	// 		log.Fatal("Post Error",_xerr)
	// 	}

	// 	fmt.Println(desenc_str)
	// 	//return  str_resp
	//  }

	//var user models.Users
	//database.Database.Where("username = ?", strings.ToUpper(request.Account)).First(&user)
	return c.JSON(loginResponse)
}
func LaunchGame(c *fiber.Ctx) error {

	request := new(GcRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	// //var users models.Users
	type CResponse struct {
		Message string      `json:"message"`
		Status  bool        `json:"status"`
		Data    GResponseData `json:"data"`  
	}
	response := CResponse{
		Status:false,
		Message:"",
		Data: GResponseData{},
	}

	db, _ := database.ConnectToDB(request.Username)

	strdata := fiber.Map{
		"account": request.Account,
	}
	resp,err := makePostRequest("http://gservice:9003/LaunchGame",strdata)
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	bodyBytes := resp.Body()
	bodyString := string(bodyBytes)
	

	 

 	err = json.Unmarshal([]byte(bodyString), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	//var users models.Users
	updates := map[string]interface{}{
		"g_token": response.Data.Token,
		}

	repository.UpdateFieldsUserString(db,request.Account, updates) 

	//fmt.Printf("URL: %s\n", response.Data.Url)
	//fmt.Printf("TOKEN: %s\n", response)
	
	 
	return c.JSON(response)
}
func loging(account string) (GResponse,error) {
	
	data := fiber.Map{
		"SystemCode": SYSTEMCODE,
		"WebId": WEBID,
		"MemberAccount": account,
		"ItemNo": "1",
		"BackUrl": "https://tsx.bet/",
		"GroupLimitID": "1,4,12",
		"Lang": "th-TH",
	}

	//fmt.Println(data)

	var loginResponse GResponse

	jsonData, err := json.Marshal(data)

	// if err != nil {
	// 	log.Fatal("Error marshalling JSON:", err)
	// }

	// des,en_err := encrypt.EncryptDES(string(jsonData),DESKEY,DESIV)
	// if en_err != nil {
	// 	log.Fatal("Encode Error:",en_err)
	// }
	// //ecresult := encrypt.ECResult{}
	// _,ecresult := encrypt.CreateSignature(CLIENT_ID,CLIENT_SECRET, des)

	strdata := RequestEncrypt{
		Data: string(jsonData),
		Key: DESKEY,
		Iv: DESIV,
	}
	
	type Response struct {
		Data encrypt.ECResult `json:"data"`
	}

	resp,err := makePostRequest("http://gservice:9003/encryption",strdata)		
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	bodyBytes := resp.Body()
	bodyString := string(bodyBytes)

	// แสดงผล string ที่ได้
	//fmt.Println("Response body as string:", bodyString)

	var response Response

 	err = json.Unmarshal([]byte(bodyString), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return loginResponse,err
	}

	fmt.Printf("Enc: %s\n", response.Data.Enc)
	fmt.Printf("Des: %s\n", response.Data.Des)
	fmt.Printf("Unx: %d\n", response.Data.Unx)

	 
	

	 resp,err = GPostRequest(API_URL_PROXY+"/api/Player/Login",CLIENT_ID,&response.Data)
	
	 bodyBytes = resp.Body()
	 bodyString = string(bodyBytes)

	
	
	strdata = RequestEncrypt{
		Data: bodyString,
		Key: DESKEY,
		Iv: DESIV,
	}
	resp,err = makePostRequest("http://gservice:9003/decryption",strdata)		
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	resultBytes := resp.Body()
	resultString := string(resultBytes)
	// แสดงผล string ที่ได้
	fmt.Println("Response body as string:", resultString)

	return loginResponse, nil

}
func GetUserOnline(c *fiber.Ctx) (error){
	// var data = fiber.Map{
	// 	"SystemCode": handler.SYSTEMCODE,
	// 	"WebId": handler.WEBID,
	// }
	response := Response{
		Status:false,
		Data: encrypt.GData{},
	}
	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	log.Fatal("Error marshalling JSON:", err)
	// }

	// encode,_enerr := encrypt.EncryptDES([]byte(jsonData),[]byte(handler.DESKEY),[]byte(handler.DESIV))
	// //EncryptDESAndMD5(string(jsonData),handler.DESKEY,handler.DESIV,handler.CLIENT_ID,handler.CLIENT_SECRET)
	
	// if _enerr != nil {
	// 	return _enerr
	// }
	// resultx,err_ := encrypt.DecryptDES([]byte(encode.Des),[]byte(handler.DESKEY),[]byte(handler.DESIV))
	// if err_ != nil {
	// 	return err_
	// }
	// //result := *encode 
	//  fmt.Println(resultx)
	//  resp,err := MakePostRequest(API_URL_PROXY+"/api/Player/GetPlayerOnlineList",encode)	
	//  str_resp := string(resp.Body())

	// //  decode,err := encrypt.DecryptDES(str_resp,[]byte(handler.DESKEY),[]byte(handler.DESIV))
	// //  if err !=nil {
	// // 	 response = Response{
	// // 		 Status:false,
	// // 		 Data: err,
	// // 	 }
	// //  }
	  
 

	return c.JSON(response)
}
func createOrUpdate(account string,name string)(*fasthttp.Response,error) {

	var data = fiber.Map{
			"SystemCode": SYSTEMCODE,
			"WebId":  WEBID,
			"MemberAccount": account,
			"MemberName": name,
			"StopBalance": -1,
			"BetLimitGroup": "1,4,12",
			"Currency": "THB",
			"Language": "th-TH",
			"OpenGameList": "ALL",
		}
 
 
		jsonData, _ := json.Marshal(data)
	
		des,_ := encrypt.EncryptDES(string(jsonData),CLIENT_ID,CLIENT_SECRET)
	 
		_,ecresult := encrypt.CreateSignature(CLIENT_ID,CLIENT_SECRET, des)
		
		// ecresult := encrypt.ECResult{}
		 
		 
		
		fmt.Println("Des:",ecresult.Des)
		fmt.Println("Unx:",ecresult.Unx)
		fmt.Println("Enc:",ecresult.Enc)
		
		
		
		dex,_ := encrypt.DecryptDES(des,DESKEY,DESIV)
	
		fmt.Println("Decrpt Text:",dex)
		
		resp,_ := fastPost(API_URL_PROXY+"/api/Player/Login",CLIENT_ID,&ecresult)
		 
	return resp,nil

}
func fastPost(url,clienid string,encoded *encrypt.ECResult) (*fasthttp.Response,error) {

	//url = "http://api.tsxbet.info:8001/api/Player/Login"
	method := "POST"

	// สร้าง request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	req.Header.Set("X-API-ClientID",clienid )
	req.Header.Set("X-API-Signature", encoded.Enc)
	req.Header.Set("X-API-Timestamp",string(encoded.Unx))

	// สร้าง response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	req.SetBody([]byte(encoded.Des))
	// ส่ง request
	 
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)

	return resp, nil
}
func Gsign(account string, expire int64, tokenType string) (string, error) {
	// ถ้า expire คือ ชั่วโมง (h) ให้แปลงเป็นวินาที
	expire = expire * 3600

	// กำหนดเวลาสำหรับ nbf, iat, exp
	now := time.Now().Unix()
	expirationTime := now + expire

	claims := &GClaims{
		SystemCode:   SYSTEMCODE,
		WebID:          WEBID,
		MemberAccount: account,
		TokenType:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Unix(now, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(now, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := fmt.Sprintf("%s%s", CLIENT_ID, CLIENT_SECRET)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func Gverify(tokenString string) (*GClaims, error) {
	secret := fmt.Sprintf("%s%s", CLIENT_ID, CLIENT_SECRET)

	token, err := jwt.ParseWithClaims(tokenString, &GClaims{}, func(token *jwt.Token) (interface{}, error) {
		 
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า Token เป็นของ Claims ที่เรากำหนดไว้
	if claims, ok := token.Claims.(*GClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
func parseLoginResponse(responseString string) (GResponse, error) {
	var loginResponse GResponse

	// แปลง JSON string เป็น struct
	err := json.Unmarshal([]byte(responseString), &loginResponse)
	if err != nil {
		return GResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	return loginResponse, nil
}
func responseToString(resp *fasthttp.Response) string {
	// อ่าน body ของ response
	body := resp.Body()

	// แปลงเป็น string และคืนค่า
	return string(body)
}
func parseResponseToECResult(resp *fasthttp.Response) (*encrypt.ECResult, error) {
	// อ่าน response body
	body := resp.Body()

	// สร้าง struct ECResult เปล่า
	var result encrypt.ECResult

	// Unmarshal JSON response ลงใน struct
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	
	return &result, nil
}
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
func GPostRequest(url ,clienid string,ecresult *encrypt.ECResult) (*fasthttp.Response, error) {
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
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set("X-API-ClientID",clienid )
	req.Header.Set("X-API-Signature", ecresult.Enc)
	req.Header.Set("X-API-Timestamp",string(ecresult.Unx))
	req.SetBody([]byte(ecresult.Des))

	// ส่ง request
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	// ปล่อย Request (เนื่องจาก fasthttp ใช้ memory pool)
	fasthttp.ReleaseRequest(req)
	
	return resp, nil
}