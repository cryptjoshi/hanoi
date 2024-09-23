package repository

import (
 "errors"
 "apigateway/models"
 "apigateway/database"
 "gorm.io/gorm"
  
)

// ฟังก์ชันสำหรับแก้ไขข้อมูลยูสเซอร์และบันทึกการเปลี่ยนแปลง
func UpdateUser(db *gorm.DB, user *models.Users) error {
    // ทำการบันทึกข้อมูลใหม่ของยูสเซอร์
    if err := db.Save(user).Error; err != nil {
        return errors.New("มีข้อผิดพลาด")
    }
    return nil
}
// func UpdateFieldsUser(userID int, updates map[string]interface{}) error {
//     // ดึงข้อมูลของยูสเซอร์ที่ต้องการแก้ไขจากฐานข้อมูล
//     var user models.Users
//     prefix := userID[:3] // Extract prefix

//     // Connect to the database based on the prefix
//     db, err := database.Connect(prefix)
//     if err = db.First(&user, userID).Error; err != nil {
//         return errors.New("มีข้อผิดพลาด")
//     }

//     // ทำการอัปเดตเฉพาะฟิลด์ที่ต้องการ
//     if err := database.Database.Model(&user).Updates(updates).Error; err != nil {
//         return errors.New("มีข้อผิดพลาด")
//     }
//     return nil
// }

func UpdateFieldsUserString(username string, updates map[string]interface{}) error {
    // ดึงข้อมูลของยูสเซอร์ที่ต้องการแก้ไขจากฐานข้อมูล
    var user models.Users

    prefix := username[:3] // Extract prefix

    // Connect to the database based on the prefix
    db, err := database.ConnectToDB(prefix)

    if err = db.Where("username=?",username).First(&user).Error; err != nil {
        return errors.New("มีข้อผิดพลาด")
    }

    // ทำการอัปเดตเฉพาะฟิลด์ที่ต้องการ
    if err = db.Model(&user).Updates(updates).Error; err != nil {
        return errors.New("มีข้อผิดพลาด")
    }
    return nil
}

// func UpdateUserFields(db *gorm.DB, userID int, updates map[string]interface{}) error {
//     // ดึงข้อมูลของยูสเซอร์ที่ต้องการแก้ไขจากฐานข้อมูล
//     var user models.Users
//     prefix := userID[:3] // Extract prefix

//     // Connect to the database based on the prefix
//     db, err := database.Connect(prefix)
// 	if err = db.Where("walletid = ? ", userID).First(&user).Error; err != nil {
//     //if err := db.First(&user, userID).Error; err != nil {
//         return errors.New("มีข้อผิดพลาด")
//     }

//     // ทำการอัปเดตเฉพาะฟิลด์ที่ต้องการ
//     if err = db.Model(&user).Updates(updates).Error; err != nil {
//         return errors.New("มีข้อผิดพลาด")
//     }
//     return nil
// }