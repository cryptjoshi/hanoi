package route

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/etag"
	// "fmt"
	// "strconv"
	// "strings"
	"hanoi/handler"
	"hanoi/handler/users"
	// "hanoi/handler/partner"
	"hanoi/handler/wallet"
	//"hanoi/middlewares"
	// "hanoi/handler/ef"
	// "hanoi/handler/gc"
	jwt "hanoi/handler/jwtn"
	// "hanoi/handler/pg"
	"github.com/go-redis/redis/v8"
	pro "hanoi/handler/promotion"
	"os"
)

var jwtSecret = os.Getenv("PASSWORD_SECRET")

 
func SetupRoutes(app fiber.Router, isV1 bool,redisClient *redis.Client) {
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
	app.Post("/users/info",jwt.JwtMiddleware,users.GetUser(redisClient))
	app.Post("/users/info/username",users.GetUserByUsername)
	app.Post("/users/statement",jwt.JwtMiddleware,users.GetUserStatement)
	app.Post("/users/transactions",jwt.JwtMiddleware,users.GetUserTransaction)
	app.Post("/users/update",jwt.JwtMiddleware,users.UpdateUser)
	//app.Post("/users/update/pro",jwt.JwtMiddleware,users.UpdateUserPro)
	app.Post("/users/commission",jwt.JwtMiddleware,handler.GetUserCommission)
	//app.Post("/users/promotion",jwt.JwtMiddleware,users.GetPromotionByUser)
	
    

	
	// app.Post("/db/create",handler.CreateDatabase)
	// app.Post("/db/register",handler.Register)
	// app.Post("/db/login",handler.RootLogin)
	// app.Post("/db/list",handler.GetDatabaseList)
	// app.Post("/db/prefix",handler.GetDatabaseByPrefix)
	// app.Post("/db/update",handler.UpdateDatabase)
	// app.Post("/db/setting",handler.GetMasterSetting)
	

	///   promotion

	//app.Post("/db/promotion/all",pro.GetAllPromotion)
	app.Post("/db/promotion/byuser",jwt.JwtMiddleware,handler.GetPromotionByUser)
	//app.Post("/users/promotions/webhook",jwt.JwtMiddleware,pro.Withdraw(redisClient))
	app.Post("/users/promotions/withdraw",jwt.JwtMiddleware,pro.Withdraw(redisClient))
	app.Post("/users/promotions/deposit",jwt.JwtMiddleware,pro.Deposit(redisClient))
	app.Post("/users/promotions/playgame",jwt.JwtMiddleware,pro.Playgame(redisClient))
	app.Post("/users/promotions/byuser",jwt.JwtMiddleware,pro.GetPromotionsUsersID(redisClient))
	app.Post("/users/promotions/status",jwt.JwtMiddleware,pro.GetPromotionStatus(redisClient))
	app.Post("/users/promotions/total",jwt.JwtMiddleware,pro.GetUserTotalPromotion(redisClient))
	app.Post("/users/promotions/all",pro.GetAllPromotion(redisClient))
	app.Post("/users/promotions/banks",jwt.JwtMiddleware,pro.GetTransactionHandler(redisClient))
	app.Post("/users/promotions",jwt.JwtMiddleware,pro.SelectPromotion(redisClient))
	app.Post("/users/promotions/clear",jwt.JwtMiddleware,pro.ClearData(redisClient))
	//app.Post("/db/promotion/update",pro.UpdatePromotion)
	//app.Post("/db/promotion/delete",pro.DeletePromotion)
	
	


	app.Post("/db/game/all",handler.GetGameList)
	app.Post("/db/game/byid",handler.GetGameById)
	app.Post("/db/game/bytype",jwt.JwtMiddleware,handler.GetGameByType)
	app.Post("/db/game/status",handler.GetGameStatus)
	app.Post("/db/game/create",handler.CreateGame)
	app.Post("/db/game/update",handler.UpdateGame)
 
	// app.Post("/db/game/all", handler.GetGameList)
	// app.Post("/db/game/byid", handler.GetGameById)
	app.Post("/db/game/status", handler.GetGameStatus)
	// app.Post("/db/game/create", handler.CreateGame)
	// app.Post("/db/game/update", handler.UpdateGame)

	// app.Post("/db/member/create", handler.CreateMember)
	// app.Post("/db/member/all", handler.GetMemberList)
	// app.Post("/db/member/byid", handler.GetMemberById)
	// app.Post("/db/member/update", handler.UpdateMember)
	// app.Post("/db/member/bypartner",jwt.JwtPMiddleware,handler.GetMemberByPartner)
	// app.Post("/db/member/bypartnerid",handler.GetMemberByPartnerId)

	// app.Post("/db/master/update", handler.UpdateMaster)
	// app.Post("/db/master/commission", handler.GetCommission)

	// app.Post("/db/exchange/rate", handler.GetExchangeRates)

	// app.Delete("/users/logout", jwt.JwtMiddleware, users.Logout)

	// app.Get("/protected", jwt.JwtMiddleware, handler.Protected)
	 
	app.Post("/statement/all",jwt.JwtMiddleware,wallet.GetStatement)
	app.Post("/statement/update", wallet.UpdateStatement)
	app.Post("/statement/withdraw",jwt.JwtMiddleware, pro.Withdraw(redisClient))
	app.Post("/statement/deposit",jwt.JwtMiddleware, pro.Deposit(redisClient))
	app.Post("/statement/webhook",wallet.Webhook(redisClient))
	
	app.Post("/transaction/add",jwt.JwtMiddleware,handler.AddTransactions)
	app.Post("/transaction/all",jwt.JwtMiddleware,handler.GetAllTransaction)
	 
	// app.Post("/db/partner",jwt.JwtPMiddleware, partner.GetPartner)
	// app.Post("/db/partner/all", partner.GetPartners)
	// app.Post("/db/partner/byid", partner.GetPartnerById)
	// app.Post("/db/partner/login", partner.Login)
	// app.Post("/db/partner/create",partner.Register)
	// app.Post("/db/partner/update",partner.UpdatePartner)
	// app.Post("/db/partner/checkseed",partner.GetSeed)
	// app.Post("/db/partner/overview",jwt.JwtPMiddleware,partner.Overview)
 
}
