package database

import (
	"fmt"
	"log"
	//"os"
	"strings"
	"sync"
	"hanoi/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// var Database *gorm.DB
// //var DSN = 'root:helloworld@tcp(db:3306)/tsxbet_dev?tls=true'
// var DSN string = "root:helloworld@tcp(db:3306)/ckd_dev?charset=utf8mb4&parseTime=True&loc=Local"
// func Connect() error {
// 	var err error
// 	dsn := fmt.Sprintf("%s&parseTime=True", DSN)// os.Getenv("DSN"))

// 	Database, err = gorm.Open(
// 		mysql.Open(dsn),
// 		&gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, SkipDefaultTransaction:true,
// 			PrepareStmt:true},
// 	)

// 	if err == nil {
// 		fmt.Println("Successfully connected to DB!")
// 	}

// 	return err
// }

var (
	dbConnections = make(map[string]*gorm.DB)
	mutex         sync.Mutex
)

const baseDSN = "web:1688XdAs@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FBangkok"

// const baseDSN = "root:1688XdAs@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
func handleError(err error) {
	log.Fatal(err)
}

func CheckAndCreateTable(db *gorm.DB, model interface{}) error {
	migrator := db.Migrator()
	// ใช้ db.Model(model) เพื่อกำหนดค่า Statement.Table
	tableName := db.Model(model).Statement.Table
	if tableName == "" {
		return fmt.Errorf("table name is empty, please ensure the model is correctly defined")
	}

	if !migrator.HasTable(tableName) {
		fmt.Printf("Table '%s' does not exist. Creating...\n", tableName)
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to create table '%s': %v", tableName, err)
		}
		fmt.Printf("Table '%s' created successfully\n", tableName)
	} else {
		fmt.Printf("Table '%s' already exists. Updating schema...\n", tableName)
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to update table '%s': %v", tableName, err)
		}
		fmt.Printf("Table '%s' schema updated successfully\n", tableName)
	}

	return nil
}
func migrateNormal(db *gorm.DB) {

	if err := db.AutoMigrate(&models.Referral{},&models.Partner{},&models.Affiliate{},&models.AffiliateLog{},&models.Product{},&models.BanksAccount{},&models.Users{},&models.TransactionSub{},
		&models.BankStatement{},&models.BuyInOut{},&models.PromotionLog{},&models.Games{},&models.Promotion{},&models.Provider{},&models.TsxAdmin{}); err != nil {
		fmt.Errorf("Tables schema migration not successfully\n")
	}
	 
	//fmt.Println("Migrations Normal Tables executed successfully")
}

func MigrationPromotion(db *gorm.DB){
	if err := db.AutoMigrate(&models.Promotion{},&models.PromotionLog{},&models.BankStatement{});err != nil {
		fmt.Errorf("Tables schema migration not successfully\n")
	}
	//fmt.Println("Migrations Promotion Tables executed successfully")
}
func migrationAffiliate(db *gorm.DB){
	if err := db.AutoMigrate(&models.Referral{},&models.Partner{},&models.Affiliate{},&models.AffiliateTracking{},&models.Users{},&models.AffiliateLog{},&models.Promotion{});err != nil {
		fmt.Errorf("Tables schema migration not successfully\n")
	}
	//fmt.Println("Migrations Affiliate Tables executed successfully")
}

type Setting struct {
	Value string `gorm:"column:value"` // Adjust the struct according to your table schema
}

func GetDBName(prefix string) (string,error) {

	masterDSN := fmt.Sprintf(baseDSN, "master") // Assuming 'master' is the database name for settings
	masterDB, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
	})
	if err != nil {
		fmt.Errorf("Error %s ",err)
		return "", err// Return the error if connection to master fails
	}
	defer func() {
		sqlDB, _ := masterDB.DB() // รับการเชื่อมต่อ SQL
		sqlDB.Close()              // ปิดการเชื่อมต่อ
	}()

	var setting Setting
	
	if err := masterDB.Table("Settings").Where("`key` = ?", prefix).First(&setting).Error; err != nil {
		return "", fmt.Errorf("failed to read setting for prefix '%s': %v", prefix, err)
	}
	dbName := fmt.Sprintf("%s_%s", prefix, setting.Value) 
	return dbName,nil
}


func GetMaster(prefix string) (Setting,error) {

	masterDSN := fmt.Sprintf(baseDSN, "master") // Assuming 'master' is the database name for settings
	masterDB, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
	})
	if err != nil {
		fmt.Errorf("Error %s ",err)
		return Setting{}, err// Return the error if connection to master fails
	}
	defer func() {
		sqlDB, _ := masterDB.DB() // รับการเชื่อมต่อ SQL
		sqlDB.Close()              // ปิดการเชื่อมต่อ
	}()

	var setting Setting
	
	if err := masterDB.Table("Settings").Where("`key` = ?", prefix).First(&setting).Error; err != nil {
		return Setting{}, fmt.Errorf("failed to read setting for prefix '%s': %v", prefix, err)
	}
	 
	return setting,nil
}

// Connect function to establish a database connection based on the prefix
func ConnectToDB(prefix string) (*gorm.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()

	prefix = strings.ToLower(prefix)
	setting, _ := GetMaster(prefix)
	dbName := fmt.Sprintf("%s_%s", prefix, setting.Value)

	// ตรวจสอบการเชื่อมต่อที่มีอยู่
	if db, exists := dbConnections[dbName]; exists {
		MigrationPromotion(db)
		migrationAffiliate(db)
		return db, nil
	}

	// สร้าง DSN พร้อมกำหนด timezone
	dsn := fmt.Sprintf(baseDSN, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                             true,
	})

	if err != nil {
		return nil, err
	}

	// ตั้งค่า timezone หลังจากเชื่อมต่อสำเร็จ
	if sqlDB, err := db.DB(); err == nil {
		// ตั้งค่า timezone
		if _, err := sqlDB.Exec("SET time_zone = 'Asia/Bangkok'"); err != nil {
			fmt.Printf("Warning: Failed to set timezone: %v\n", err)
		}

		// ตั้งค่าการเชื่อมต่อเพิ่มเติม
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	if err == nil {
		dbConnections[dbName] = db
		fmt.Println("Successfully connected to DB:", dbName)
	}

	return db, nil
}

// Helper function to check if a prefix is valid
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

 