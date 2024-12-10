package repository

import (
 "errors"
 "hanoi/models"
 //"hanoi/database"
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
 
 
func UpdateFieldsUserString(db *gorm.DB,username string, updates map[string]interface{}) error {

    // ดึงข้อมูลของยูสเซอร์ที่ต้องการแก้ไขจากฐานข้อมูล
    var user models.Users

    if err := db.Debug().Where("username=?",username).First(&user).Error; err != nil {
        return errors.New("มีข้อผิดพลาด")
    }

    // ทำการอัปเดตเฉพาะฟิลด์ที่ต้องการ
    if err := db.Debug().Model(&user).Updates(updates).Error; err != nil {
        return errors.New("มีข้อผิดพลาด")
    }
    return nil
}

func UpdateUserFields(db *gorm.DB, userID int, updates map[string]interface{}) error {
    // ดึงข้อมูลของยูสเซอร์ที่ต้องการแก้ไขจากฐานข้อมูล
    var user models.Users

	if err := db.Debug().Where("walletid = ? ", userID).First(&user).Error; err != nil {
    //if err := db.First(&user, userID).Error; err != nil {
        return errors.New("มีข้อผิดพลาด")
    }

    // ทำการอัปเดตเฉพาะฟิลด์ที่ต้องการ
    if err := db.Debug().Model(&user).Updates(updates).Error; err != nil {
        return errors.New("มีข้อผิดพลาด")
    }
    return nil
}