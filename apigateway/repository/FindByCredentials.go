package repository

import (
 "errors"
 "apigateway/models"
 "apigateway/database"
  
)
// Simulate a database call
func FindUser(username, password string) (*models.Users, error) {
    var user models.Users

    // ดึงข้อมูลโดยใช้ Username และ Password
    prefix := username[:3] // Extract prefix

    // Connect to the database based on the prefix
    db, err := database.ConnectToDB(prefix)
    if err = db.Where("preferredname = ? AND password = ?", username, password).First(&user).Error; err != nil {
        return &user,  errors.New("user not found")
    }
    return &user, nil
 
}

