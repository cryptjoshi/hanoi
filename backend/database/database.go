package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"os"
)

var Database *gorm.DB
//var DSN = 'root:helloworld@tcp(db:3306)/tsxbet_dev?tls=true'
var DSN string = "root:helloworld@tcp(db:3306)/ckd_dev?charset=utf8mb4&parseTime=True&loc=Local"
func Connect() error {
	var err error
	dsn := fmt.Sprintf("%s&parseTime=True", DSN)// os.Getenv("DSN"))

	Database, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, SkipDefaultTransaction:true,
			PrepareStmt:true},
	)

	if err == nil {
		fmt.Println("Successfully connected to DB!")
	}

	return err
}