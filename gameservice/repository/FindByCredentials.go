package repository

import (
 "errors"
 "pkd/models"
 "pkd/database"
  
)
// Simulate a database call
func FindUser(username, password string) (*models.Users, error) {
    var user models.Users

    // ดึงข้อมูลโดยใช้ Username และ Password
    if err := database.Database.Where("preferredname = ? AND password = ?", username, password).First(&user).Error; err != nil {
        return &user,  errors.New("user not found")
    }
    return &user, nil
 
}

