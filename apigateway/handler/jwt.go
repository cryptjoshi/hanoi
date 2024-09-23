package handler

import (
	//"bytes"
	//"crypto/cipher"
	//"crypto/des"
	//"encoding/base64"
	"apigateway/models"
	"apigateway/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	//"github.com/golang-jwt/jwt"
	//jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	// "pkd/repository"
	//"log"
	// "net"
	// "net/http"
	"encoding/json"
	"os"
	//"strconv"
	"time"
	//"strings"
	"fmt"
	"errors"
)
var jwtKey  = []byte(os.Getenv("PASSWORD_SECRET"))
//var CLIENT_ID = "6342e1be-fa03-456f-8d2d-8e1c9513c351" //[]byte(os.Getenv("CLIENT_ID"))
//var CLIENT_SECRET = "6d83ac42" //[]byte(os.Getenv("CLIENT_SECRET"))
//var DESKEY = "9c62a148"
//var DESIV =	"8e014099"
//var SYSTEMCODE = "tsxthb"
//var WEBID = "tsxthb"



// Struct สำหรับ JWT Claims

// type ECResult struct {
// 	Enc string `json:"enc"`
// 	Unx int64 `json:"unx"` // ค่า unx คุณสามารถกำหนดเอง
// 	Des string `json:"des"` // ค่า dex คุณสามารถกำหนดเอง
// }




type Claims struct {
    Username string `json:"username"`
	Id int `json:"id"`
	Role string `json:"role"`
	Prefix string `json:"prefix"`
	Walletid int `json:"walletid"`
	Checker string `json:"checker"`
    jwt.RegisteredClaims
}

 


func MapToJSONString(data fiber.Map) (string, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}


func ValidateJWTReturn(tokenString string) models.Users {
	claims := &Claims{}
	//dbClaims := &Claims{}
	//tokenString := c.Get("Authorization")[7:]
	// token, claims, err := ValidateJWT(tokenString) // เรียกใช้ฟังก์ชันจาก utils
	// if err != nil || !token.Valid {
	// 	return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	// }
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
	username := claims.Username
	 
	prefix := username[:3] // Extract prefix

	// Connect to the database based on the prefix
	db, err := database.ConnectToDB(prefix)
	//checkerFromRequest := claims.Checker
	var user models.Users
	//fmt.Println(err)
	if err==nil {
		db.Select("id,username,balance,Token").Where("username = ?", username).First(&user)
	}

	 
	return user
	 
}



func ValidateJWT(tokenString string) (error) {
	claims := &Claims{}
	//dbClaims := &Claims{}
	//tokenString := c.Get("Authorization")[7:]
	// token, claims, err := ValidateJWT(tokenString) // เรียกใช้ฟังก์ชันจาก utils
	// if err != nil || !token.Valid {
	// 	return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	// }
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
	username := claims.Username
	//checkerFromRequest := claims.Checker

	

	// ดึง JWT Token ที่เก็บไว้ในฐานข้อมูลสำหรับผู้ใช้ที่เกี่ยวข้อง
	var user models.Users
	prefix := username[:3] // Extract prefix

	// Connect to the database based on the prefix
	db, err := database.ConnectToDB(prefix)
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return result.Error
	}
	
	utoken := user.Token
	fmt.Println(utoken)
	//fmt.Println(tokenString)
	// ตรวจสอบและเปรียบเทียบค่า checker
	// _,err_ := jwt.ParseWithClaims(utoken, dbClaims, func(token *jwt.Token) (interface{}, error) {
    //     return jwtKey, nil
    // })
	
	
	// if err_!= nil {
	// 	fmt.Println("77")
	// 	fmt.Println(err_)
	// 	return err_
	// }
	// checkerFromDB := dbClaims.Checker
	// fmt.Println(&dbClaims)
	// fmt.Println(&claims)
	// แสดงค่า checker จาก request และจากฐานข้อมูล
	//fmt.Printf("Checker from request token: %s\n", checkerFromRequest)
	//fmt.Printf("Checker from DB token: %s\n", checkerFromDB)

	// เปรียบเทียบค่า checker
	if utoken != tokenString {
		//return c.Status(fiber.StatusUnauthorized).SendString("Checker mismatch")
		return errors.New("checker ไม่ตรง!")
	}
	return err
}

// ฟังก์ชันสำหรับตรวจสอบและแยก JWT Token
func validatedJWT(tokenString string) (error) {
    claims := &Claims{}
    
    _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
	
	var user models.Users
	prefix := claims.Username[:3] // Extract prefix

	// Connect to the database based on the prefix
	db, err := database.ConnectToDB(prefix)
	result := db.Debug().Where("username = ?", claims.Username).First(&user)
	if result.Error != nil {
		//http.Error(w, "User not found", http.StatusUnauthorized)
		return errors.New("มีข้อผิดพลาด")
	}
	fmt.Println("------------")
	fmt.Println(claims.Checker)
	fmt.Println("------------")
	fmt.Println(tokenString)
	fmt.Println("------------")

	// ตรวจสอบว่า token ที่ส่งมาไม่ตรงกับ token ที่เก็บในฐานข้อมูล
	if user.Token != tokenString {
		//http.Error(w, "Token ไม่ตรง", http.StatusUnauthorized)
		return errors.New("มีข้อผิดพลาด")
	}

	// หาก token ถูกต้องและตรงกัน
	//fmt.Fprintf(w, "Hello, %s", claims.Username)

    return err
}

// ฟังก์ชันสำหรับสร้าง JWT Token (เพื่อใช้ทดสอบ)
func createJWT(username string) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}


// func pad(data []byte, blockSize int) []byte {
// 	padding := blockSize - len(data)%blockSize
// 	padText := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(data, padText...)
// }

// func encryptDesCbc(data, key, iv []byte) ([]byte, error) {
// 	// สร้าง block สำหรับ DES
// 	block, err := des.NewCipher(key)
// 	if err != nil {
// 		return  nil, err
// 	}

// 	// เติม padding ให้ข้อมูลตาม block size
// 	data = pad(data, block.BlockSize())

// 	// สร้าง Cipher Block Mode (CBC)
// 	mode := cipher.NewCBCEncrypter(block, iv)

// 	// เข้ารหัสข้อมูล
// 	encrypted := make([]byte, len(data))
// 	mode.CryptBlocks(encrypted, data)

// 	return encrypted,err
// 	// เข้ารหัส base64 ก่อนส่งออก
// 	//return base64.StdEncoding.EncodeToString(encrypted), nil
// }

 



// Helper function to check if string is empty and return default value if so
// func ifEmpty(value, defaultValue string) string {
// 	if value == "" {
// 		return defaultValue
// 	}
// 	return value
// }

 
