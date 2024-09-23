package handler

import {
	"apigateway/database"
	"apigateway/models"

}

func AddUser(c *fiber.Ctx) error {
	
	var currency =  os.Getenv("CURRENCY")
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(200).SendString(err.Error())
	}
	//user.Walletid = user.ID
	//user.Username = user.Prefix + user.Username + currency
	db, _ := database.ConnectToDB(user.Prefix)
	result := db.Create(&user); 

	
	

	// ส่ง response เป็น JSON
   


	if result.Error != nil {
		response := Response{
			Message: "เกิดข้อผิดพลาดไม่สามารถเพิ่มข้อมูลได้!",
			Status:  false,
			Data:    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
		} else {

			updates := map[string]interface{}{
				"Walletid":user.ID,
				"Preferredname": user.Username,
				"Username":user.Prefix + user.Username + currency,
			}
			if err := db.Model(&user).Updates(updates).Error; err != nil {
				return errors.New("มีข้อผิดพลาด")
			}
		
		response := Response{
			Message: "เพิ่มยูสเซอร์สำเร็จ!",
			Status:  true,
			Data:    fiber.Map{ 
				"id": user.ID,
				"walletid":user.ID,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	}

	 
}

func GetBalance(c *fiber.Ctx) error {
	
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	
	var users models.Users
	
	db, _ := database.ConnectToDB(users.Prefix)

	db.Debug().Where("id= ?",claims["ID"]).Find(&users)
	 
	
	if users == (models.Users{}) {
	 
			response := Response{
				Message: "ไม่พบรหัสผู้ใช้งาน!!",
				Status:  false,
			}
			return c.JSON(response)
	}

	tokenString := c.Get("Authorization")[7:] 
	
	_err := validateJWT(tokenString);
	//fmt.Println(_err)
	 if _err != nil {
	   
		  response := Response{
			Message: "โทเคนไม่ถูกต้อง!!",
			Status:  false,
		}
		return c.JSON(response)
	}else {
	//fmt.Println('')
		response := Response{
			Status: true,
			Message: "สำเร็จ",
			Data: fiber.Map{ 
			"balance":users.Balance,
		}}
		return c.JSON(response)
	}
	
	
}