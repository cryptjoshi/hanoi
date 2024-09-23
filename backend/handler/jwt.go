package handler

import (
	// "context"
	// "fmt"
	// "github.com/amalfra/etag"
	// "github.com/go-redis/redis/v8"
	//"github.com/gofiber/fiber/v2"
	//"github.com/shopspring/decimal"
	// "github.com/streadway/amqp"
	// "github.com/tdewolff/minify/v2"
	// "github.com/tdewolff/minify/v2/js"
	// "github.com/valyala/fasthttp"
	// _ "github.com/go-sql-driver/mysql"
	"pkd/models"
	"pkd/database"
	 "github.com/golang-jwt/jwt/v4"
	//"github.com/golang-jwt/jwt"
	//jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	// "pkd/repository"
	// "log"
	// "net"
	// "net/http"
	"os"
	// "strconv"
	"time"
	//"strings"
	"fmt"
	"errors"
)
var jwtKey  = []byte(os.Getenv("PASSWORD_SECRET"))

// Struct สำหรับ JWT Claims
type Claims struct {
    Username string `json:"username"`
	Id int `json:"id"`
	Role string `json:"role"`
	Prefix string `json:"prefix"`
	Walletid int `json:"walletid"`
	Checker string `json:"checker"`
    jwt.RegisteredClaims
}


func validateJWT(tokenString string) (error) {
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
	result := database.Database.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return result.Error
	}
	
	utoken := user.Token
	//fmt.Println(utoken)
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
	result := database.Database.Debug().Where("username = ?", claims.Username).First(&user)
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