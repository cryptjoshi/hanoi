package route

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"pkd/handler"
	"pkd/middlewares"
	//"pkd/handler/ef" 
	// "pkd/handler/gc" 
	// "pkd/handler/pg" 
	"os"
)
var jwtSecret = os.Getenv("PASSWORD_SECRET")
func SetupRoutes(app *fiber.App) {
	// app.Use(etag.New())
	app.Use(compress.New())

	jwt := middlewares.NewAuthMiddleware(jwtSecret)

	// app.Static("/css", "./css")
	// app.Static("/js/libraries", "./js")
	app.Get("/",handler.GetRoot)

	
	//auth
	app.Post("api/auth/register",handler.AddUser)
	app.Post("api/auth/login",handler.Login)
	app.Delete("api/auth/logout",jwt,handler.Logout)

	app.Get("/protected", jwt, handler.Protected)
	
	// user
	app.Post("api/user/all",handler.GetUserAll)
	app.Post("api/user/userinfo",jwt,handler.GetUserByID)
	app.Get("api/user/balance",jwt,handler.GetBalance)
	app.Post("api/user/statement",jwt,handler.GetUserStatement)


	// wallet
	// app.Get("api/user/withdraw",handler.WithDraw)
	// app.Get("api/user/deposit",handler.Deposit)
	// app.Post("api/statement/update",handler.UpdateStatement)
	// app.Post("api/statement/add",handler.AddStatement)

	// dashboard
	app.Post("/api/bank/statement",handler.GetBankStatement)
	app.Post("/api/first/statement",handler.GetFirstUsers)
	app.Post("/api/user/all/statement",handler.GetUserAllStatement)
	app.Post("/api/user/statement",jwt,handler.GetUserStatement)
	app.Post("/api/user/statement/id",handler.GetIdStatement)
	app.Post("/api/user/sum/statement",handler.GetUserSumStatement)


	// app.Post("api/user/token",jwt,handler.UpdateToken)
	// app.Post("api/user/byid",handler.GetBalanceFromID)

	 
}