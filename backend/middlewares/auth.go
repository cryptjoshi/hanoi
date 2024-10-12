package middlewares

import (
 "github.com/gofiber/fiber/v2"
 jwtware "github.com/gofiber/jwt/v3"
 "github.com/golang-jwt/jwt/v4"
 "os"
)

var jwtKey  = []byte(os.Getenv("PASSWORD_SECRET"))
// Middleware JWT function
func NewAuthMiddleware(secret string) fiber.Handler {
 return jwtware.New(jwtware.Config{
  SigningKey: []byte(secret),
 })
}

func JwtMiddleware(c *fiber.Ctx) error {
    tokenString := c.Get("Authorization")
    if tokenString == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
    }

    token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtKey), nil
    })

    

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // ดึงข้อมูล claims เช่น id, role หรือ prefix
        c.Locals("Walletid", claims["walletid"])
        c.Locals("ID", claims["id"])
        c.Locals("username", claims["username"])
        c.Locals("PartnersKey",claims["partnersKey"])
        c.Locals("role", claims["role"])
        c.Locals("prefix", claims["prefix"])
        return c.Next()
    } else {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
    }
}