package gc

import 
(
	"github.com/gofiber/fiber/v2"
	"pkd/models"
	"pkd/database"
	//"fmt"
)
type Response struct {
    Message string      `json:"message"`
    Status  bool        `json:"status"`
    Data    interface{} `json:"data"` // ใช้ interface{} เพื่อรองรับข้อมูลหลายประเภทใน field data
}
// ฟังก์ชันตัวอย่างใน gclub.go

func Index(c *fiber.Ctx) error {

	var user []models.Users
	database.Database.Find(&user)
	response := Response{
		Message: "Welcome to GClub!!",
		Status:  true,
		Data: fiber.Map{ 
			"users":user,
		}, 
	}
	 
	return c.JSON(response)
 
 
   
}
