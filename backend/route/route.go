package route

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/etag"

	"hanoi/handler"
	"hanoi/handler/users"
	"hanoi/handler/partner"
	"hanoi/handler/wallet"
	//"hanoi/middlewares"
	"hanoi/handler/ef"
	"hanoi/handler/gc"
	"hanoi/handler/jwtn"
	"hanoi/handler/pg"
	//"github.com/swaggo/fiber-swagger"
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

func SetupRoutes(app fiber.Router) {
	// app.Use(etag.New())

	//jwtm := middlewares.NewAuthMiddleware(jwtSecret)

	// app.Static("/css", "./css")
	// app.Static("/js/libraries", "./js")
	//app.Get("/",handler.GetRoot)

	// เส้นทาง Swagger
	//app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// เส้นทาง API สำหรับดึงข้อมูลผู้ใช้งาน
	// user
	app.Post("/users/all", jwt.JwtMiddleware, users.GetUsers)
	app.Post("/users/login", users.Login)
	app.Post("/users/register",users.Register)
	app.Post("/users/balance",jwt.JwtMiddleware,users.GetBalance)
	app.Post("/users/sum/balance",jwt.JwtMiddleware,users.GetBalanceSum)
	app.Post("/users/info",jwt.JwtMiddleware,users.GetUser)
	app.Post("/users/info/username",users.GetUserByUsername)
	app.Post("/users/statement",jwt.JwtMiddleware,users.GetUserStatement)
	app.Post("/users/transactions",jwt.JwtMiddleware,users.GetUserTransaction)
	app.Post("/users/update",jwt.JwtMiddleware,users.UpdateUser)
	app.Post("/users/update/pro",jwt.JwtMiddleware,users.UpdateUserPro)
	app.Post("/users/commission",jwt.JwtMiddleware,handler.GetUserCommission)
	
    
	app.Post("/db/create",handler.CreateDatabase)
	app.Post("/db/register",handler.Register)
	app.Post("/db/login",handler.RootLogin)
	app.Post("/db/list",handler.GetDatabaseList)
	app.Post("/db/prefix",handler.GetDatabaseByPrefix)
	app.Post("/db/update",handler.UpdateDatabase)
	app.Post("/db/setting",handler.GetMasterSetting)
	
	app.Post("/db/promotion/all",handler.GetAllPromotion)
	app.Post("/db/promotion/byuser",jwt.JwtMiddleware,handler.GetPromotionByUser)
	app.Post("/db/promotion/byid",handler.GetPromotionById)
	app.Post("/db/promotion/create",handler.CreatePromotion)
	app.Post("/db/promotion/update",handler.UpdatePromotion)
	app.Post("/db/promotion/delete",handler.DeletePromotion)
	//app.Post("/db/promotion/delete/:id",handler.DeletePromotion)	
	app.Post("/db/game/all",handler.GetGameList)
	app.Post("/db/game/byid",handler.GetGameById)
	app.Post("/db/game/bytype",jwt.JwtMiddleware,handler.GetGameByType)
	app.Post("/db/game/status",handler.GetGameStatus)
	app.Post("/db/game/create",handler.CreateGame)
	app.Post("/db/game/update",handler.UpdateGame)

	// app.Post("/db/create", handler.CreateDatabase)
	// app.Post("/db/login", handler.RootLogin)
	// app.Post("/db/list", handler.GetDatabaseList)
	// app.Post("/db/prefix", handler.GetDatabaseByPrefix)

	// app.Post("/db/promotion/all", handler.GetPromotion)
	// app.Post("/db/promotion/user", handler.GetPromotionByUser)
	// app.Post("/db/promotion/byid", handler.GetPromotionById)
	// app.Post("/db/promotion/create", handler.CreatePromotion)
	// app.Post("/db/promotion/update", handler.UpdatePromotion)
	// app.Post("/db/promotion/delete", handler.DeletePromotion)
	//app.Post("/db/promotion/delete/:id",handler.DeletePromotion)
	app.Post("/db/game/all", handler.GetGameList)
	app.Post("/db/game/byid", handler.GetGameById)
	app.Post("/db/game/status", handler.GetGameStatus)
	app.Post("/db/game/create", handler.CreateGame)
	app.Post("/db/game/update", handler.UpdateGame)

	app.Post("/db/member/create", handler.CreateMember)
	app.Post("/db/member/all", handler.GetMemberList)
	app.Post("/db/member/byid", handler.GetMemberById)
	app.Post("/db/member/update", handler.UpdateMember)

	app.Post("/db/master/update", handler.UpdateMaster)
	app.Post("/db/master/commission", handler.GetCommission)

	app.Post("/db/exchange/rate", handler.GetExchangeRates)

	app.Delete("/users/logout", jwt.JwtMiddleware, users.Logout)

	app.Get("/protected", jwt.JwtMiddleware, handler.Protected)
	// app.Post("/api/gateway/getBalance",jwt.JwtMiddleware, users.GetBalance)

	// Define individual routes for each provider (if needed)
	// app.Post("/callback/Seamless/GetBalance", ef.GetBalance)
	// app.Post("/callback/pgsoft/checkBalance", pg.GetBalance)
	// app.Post("/api/Wallet/Balance", gc.GetBalance)

	// wallet
	// app.Post("/wallet/withdraw",wallet.WithDraw)
	// app.Post("/wallet/deposit",wallet.AddStatement)
	app.Post("/statement/all",jwt.JwtMiddleware,wallet.GetStatement)
	app.Post("/statement/update", wallet.UpdateStatement)
	app.Post("/statement/withdraw",jwt.JwtMiddleware, wallet.Withdraw)
	app.Post("/statement/deposit",jwt.JwtMiddleware, wallet.Deposit)
	app.Post("/statement/webhook",wallet.Webhook)
	// app.Post("/transaction/add",handler.AddTransactions)
	// app.Post("/transaction/update",handler.UpdateTransactions)
	app.Post("/transaction/add",jwt.JwtMiddleware,handler.AddTransactions)
	app.Post("/transaction/all",jwt.JwtMiddleware,handler.GetAllTransaction)
	// dashboard
	// app.Post("/api/bank/statement",handler.GetBankStatement)
	// app.Post("/api/first/statement",handler.GetFirstUsers)
	// app.Post("/api/user/all/statement",handler.GetUserAllStatement)
	// app.Post("/api/user/statement",jwt,handler.GetUserStatement)
	// app.Post("/api/user/statement/id",handler.GetIdStatement)
	// app.Post("/api/user/sum/statement",handler.GetUserSumStatement)

	// app.Post("api/user/token",jwt,handler.UpdateToken)
	// app.Post("api/user/byid",handler.GetBalanceFromID)
	app.Post("/db/partner/all", partner.GetPartners)
	app.Post("/db/partner/login", partner.Login)
	app.Post("/db/partner/register",partner.Register)

	app.Post("/db/partner/checkseed",partner.GetSeed)
	// app.Post("/users/balance",jwt.JwtMiddleware,users.GetBalance)
	// app.Post("/users/sum/balance",jwt.JwtMiddleware,users.GetBalanceSum)
	// app.Post("/users/info",jwt.JwtMiddleware,users.GetUser)
	// app.Post("/users/info/username",users.GetUserByUsername)
	// app.Post("/users/statement",jwt.JwtMiddleware,users.GetUserStatement)
	// app.Post("/users/transactions",jwt.JwtMiddleware,users.GetUserTransaction)
	// app.Post("/users/update",jwt.JwtMiddleware,users.UpdateUser)
	// app.Post("/users/update/pro",jwt.JwtMiddleware,users.UpdateUserPro)
	// app.Post("/users/commission",jwt.JwtMiddleware,handler.GetUserCommission)
}
