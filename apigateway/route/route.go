package route

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"apigateway/handler"
	"apigateway/middlewares"
	"apigateway/handler/ef" 
	"apigateway/handler/gc" 
	"apigateway/handler/pg" 
	"os"
)
var jwtSecret = os.Getenv("PASSWORD_SECRET")
func ProviderMiddleware(c *fiber.Ctx) error {
	username := c.FormValue("username") // Assuming username is part of the request
	if len(username) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username",
		})
	}
	
	// Extract prefix (first 3 characters)
	prefix := username[:3]

	// Determine provider based on prefix or logic
	switch prefix {
	case "EFI":
		return ef.GetBalance(c) // EFinity
	case "PGS":
		return pg.GetBalance(c) // PGSoft
	case "GCL":
		return gc.GetBalance(c) // GClub
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid provider or prefix",
		})
	}
}

func SetupRoutes(app *fiber.App) {
	// app.Use(etag.New())
	app.Use(compress.New())

	jwt := middlewares.NewAuthMiddleware(jwtSecret)

	// app.Static("/css", "./css")
	// app.Static("/js/libraries", "./js")
	//app.Get("/",handler.GetRoot)

	
	//auth
	// app.Post("api/auth/register",handler.AddUser)
	// app.Post("api/auth/login",handler.Login)
	// app.Delete("api/auth/logout",jwt,handler.Logout)

	app.Get("/protected", jwt, handler.Protected)
	app.Post("/api/gateway/getBalance",jwt, handler.GetBalance)

	// Define individual routes for each provider (if needed)
	app.Post("/callback/Seamless/GetBalance", ef.GetBalance)
	app.Post("/callback/pgsoft/checkBalance", pg.GetBalance)
	app.Post("/api/Wallet/Balance", gc.GetBalance)
	// user
	// app.Post("api/user/all",handler.GetUserAll)
	// app.Post("api/user/userinfo",jwt,handler.GetUserByID)
	// app.Get("api/user/balance",jwt,handler.GetBalance)
	// app.Post("api/user/statement",jwt,handler.GetUserStatement)


	// wallet
	// app.Get("api/user/withdraw",handler.WithDraw)
	// app.Get("api/user/deposit",handler.Deposit)
	// app.Post("api/statement/update",handler.UpdateStatement)
	// app.Post("api/statement/add",handler.AddStatement)

	// dashboard
	// app.Post("/api/bank/statement",handler.GetBankStatement)
	// app.Post("/api/first/statement",handler.GetFirstUsers)
	// app.Post("/api/user/all/statement",handler.GetUserAllStatement)
	// app.Post("/api/user/statement",jwt,handler.GetUserStatement)
	// app.Post("/api/user/statement/id",handler.GetIdStatement)
	// app.Post("/api/user/sum/statement",handler.GetUserSumStatement)


	// app.Post("api/user/token",jwt,handler.UpdateToken)
	// app.Post("api/user/byid",handler.GetBalanceFromID)

	 
}