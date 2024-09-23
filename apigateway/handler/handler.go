package handler

import (
	//"log"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
)
 
func Protected(c *fiber.Ctx) error {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	
	username := claims["Username"].(string)
	favPhrase := claims["PartnersKey"].(string)

	return c.SendString("Bye Bye 👋" + username + " " + favPhrase)
   	
}


func Signup(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	message := fmt.Sprintf("I am GROOT!")
	return c.SendString(message)
}

func Login(c *fiber.Ctx) error {
	 
 
	loginRequest := new(models.Users)

	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	 
	user,err := repository.FindUser(loginRequest.Preferredname,loginRequest.Password)
	if err != nil {
		response := Response{
			Message: "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
	
		return c.JSON(response)
	}

	
	//day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID": user.ID,
		"Walletid": user.Walletid,
		"Username": user.Username,
		"Role": user.Role,
		"PartnersKey": user.PartnersKey,
		//"exp":   time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256,claims)

	t,err_ := token.SignedString([]byte(jwtSecret))
	
	
	if err_ != nil {
		response := Response{
			Message: "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	}
	updates := map[string]interface{}{
		"Token": t,
			}

	// อัปเดตข้อมูลยูสเซอร์
	_err := repository.UpdateFieldsUser(user.ID, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	if _err != nil {
		response := Response{
			Message: "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	} else {
		response := Response{
			Message: "เข้าระบบสำเร็จ!",
			Status:  true,
			Data: models.LoginResponse{  
				Token: t, 
				},
		}
		return c.JSON(response)
	}
 
}

 

func Logout(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	// message := fmt.Sprintf("I am GROOT!")
	// return c.SendString(message)
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	
	username := claims["Username"].(string)
	//favPhrase := claims["PartnersKey"].(string)
	//id := claims["ID"]
	 
	//if claims != nil {
		updates := map[string]interface{}{
			"Token": "",
				}
	
		// อัปเดตข้อมูลยูสเซอร์
		 repository.UpdateFieldsUserString(username, updates) 

		response := Response{
			Message: "ออกจากระบบสำเร็จ!",
			Status:  true,
			Data: fiber.Map{ 
				"id": -1,
			},
		}
		return c.JSON(response)
	// } else {
	// 	response := Response{
	// 		Message: "มีข้อผิดพลาด กรุณาออกจากระบบอีกครั้ง!",
	// 		Status:  false,
	// 		Data: fiber.Map{ 
	// 			"id": -1,
	// 		},
	// 	}
	// 	return c.JSON(response)
	// }
	 
}




