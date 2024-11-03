package route

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/compress"
	//"github.com/golang-jwt/jwt/v4"
	"pkd/handler"
	//"pkd/middlewares"
	"pkd/handler/ef" 
	"pkd/handler/gc" 
	"pkd/handler/pg" 
	"os"
)
var jwtSecret = os.Getenv("PASSWORD_SECRET")
func SetupRoutes(app *fiber.App) {
	// app.Use(etag.New())
	app.Use(compress.New())
	
	//jwt := middlewares.NewAuthMiddleware(jwtSecret)
	// app.Static("/css", "./css")
	// app.Static("/js/libraries", "./js")

	// EFinity
	app.Get("/callback/Seamless",ef.Index)
	app.Post("/callback/Seamless/GetBalance",ef.GetBalance)
	app.Post("/callback/Seamless/PlaceBet",ef.PlaceBet)
	app.Post("/callback/Seamless/GameResult",ef.GameResult)
	app.Post("/callback/Seamless/RollBack",ef.RollBack)
	app.Post("/callback/Seamless/CancelBet",ef.CancelBet)
	app.Post("/callback/Seamless/Bonus",ef.Bonus)
	app.Post("/callback/Seamless/Jackpot",ef.Jackpot)
	app.Post("/callback/Seamless/BuyIn",ef.BuyIn)
	app.Post("/callback/Seamless/BuyOut",ef.BuyOut)
	app.Post("/callback/Seamless/PushBet",ef.PushBet)
	app.Post("/callback/Seamless/MobileLogin",ef.MobileLogin)
	   

	// PGSOFT
	app.Get("/callback/pgsoft",pg.Index)
	app.Post("/callback/pgsoft/checkBalance",pg.GetBalance)
	app.Post("/callback/pgsoft/settleBets",pg.PlaceBet)
	
	
	// GCLUB
	app.Get("/api/Auth",gc.Index)
	app.Get("/api/Player/HelloWorld",gc.Index)
	app.Post("/api/Auth/CheckUser",gc.CheckUser)
	app.Post("/api/Auth/LaunchGame",gc.LaunchGame)
	//app.Post("/api/Auth/Login",gc.Login)
	app.Post("/api/Wallet/Balance",gc.GetBalance)
	// app.Post("/api/Auth/RequestExtendToken",gc.GetBalance)
	// app.Post("/api/Wallet/Debit",gc.GetBalance)
	// app.Post("/api/Wallet/Credit",gc.GetBalance)
	// app.Post("/api/Wallet/Cancel",gc.GetBalance)

	app.Get("/api/UserOnline",gc.GetUserOnline)
	// app.post("/api/LaunchGame",gc.launchGame)
	// app.get("/api/LobbyGame",gc.getGameList)
	// app.post("/api/GetGame",gc.getGame)
	// app.get("/api/GetAllGames",gc.getGames)

	
	// app.Get("/signup",handler.SignupPage)
	// app.Get("/login",handler.LoginPage)
	
	//auth
	// app.Get("api/auth/all",jwt,handler.GetUser)
	// app.Post("api/auth/signup",handler.Signup)
	
	// app.Post("api/auth/register",handler.AddUser)
	// app.Post("api/auth/login",handler.Login)
	// app.Delete("api/auth/logout",jwt,handler.Logout)

	// app.Get("/protected", jwt, handler.Protected)
	
	// // user
	// app.Post("api/user/all/balance",handler.GetUserAll)
	//app.Get("api/user/balance",jwt,handler.GetBalance)
	// app.Post("api/user/token",jwt,handler.UpdateToken)
	// app.Post("api/user/byid",handler.GetBalanceFromID)
	// app.Post("/api/user/all/statement",handler.GetUserAllStatement)
	// app.Post("/api/user/statement",jwt,handler.GetUserStatement)
	// app.Post("/api/user/statement/id",handler.GetIdStatement)
	// app.Post("/api/user/sum/statement",handler.GetUserSumStatement)
	
	
	// //BankStatement
	// app.Post("/api/status/statement",handler.UpdateStatement)
	// app.Post("/api/statement",handler.AddStatement)


	// //Transactions
	// //app.Get("/api/transaction/all",handler.GetAllTransaction)
	// //app.Post("/api/status/statement",handler.UpdateStatement)
	app.Post("/api/v1/transaction/add",handler.AddTransactions)


	// // dashboard
	// app.Post("/api/bank/statement",handler.GetBankStatement)
	// app.Post("/api/first/statement",handler.GetFirstUsers)

	// // game
	// // GClub
	//  app.Get("/api/game/gc",gc.Index)
	 
	// // PGSoft 
	//  app.Get("/api/game/pg",pg.Index)
	 
	
	
	 //  Seamless.post('/GameResult',app.gameResult)
	//  Seamless.post('/Rollback',app.rollBack)
	//  Seamless.post('/CancelBet',app.cancelBet)
	//  Seamless.post('/Bonus',app.bonus)
	//  Seamless.post('/Jackpot',app.jackpot)
	//  Seamless.post('/MobileLogin',app.mobileLogin)
	//  Seamless.post('/BuyIn',app.buyIn)
	//  Seamless.post('/BuyOut',app.buyOut)
	//  Seamless.post('/PushBet',app.pushBet)

	// app.Get("/pixel-track", handler.GetPixelTrack)
	// app.Get("/pixel/:key", handler.GetPixelPath)
}