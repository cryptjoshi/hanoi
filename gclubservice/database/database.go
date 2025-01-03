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
	"github.com/go-redis/redis/v8" 
	"context" // เพิ่มการนำเข้า Redis
	"os"
)
var redis_master_host = "redis" //os.Getenv("REDIS_HOST")
var redis_master_port = "6379"  //os.Getenv("REDIS_PORT")
var redis_master_password = os.Getenv("REDIS_PASSWORD")
var redis_slave_host = os.Getenv("REDIS_SLAVE_HOST")
var redis_slave_port = os.Getenv("REDIS_SLAVE_PORT")
var redis_slave_password = os.Getenv("REDIS_SLAVE_PASSWORD")
// var queue_name = "wallet"                   //os.Getenv("QUEUE_NAME")
// var exchange_name = "wallet"                //os.Getenv("EXCHANGE_NAME")
// var rabbit_mq = "amqp://128.199.92.45:5672" //os.Getenv("RABBIT_MQ") @rabbitmq:5672/wallet
var connection_timeout = os.Getenv("CONNECTION_TIMEOUT")
var redis_database = getEnv("REDIS_DATABASE", "0")
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
var ctx = context.Background() 
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
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
	 
	fmt.Println("Migrations Normal Tables executed successfully")
}

func migrationPromotion(db *gorm.DB){
	if err := db.AutoMigrate(&models.PromotionLog{});err != nil {
		fmt.Errorf("Tables schema migration not successfully\n")
	}
	fmt.Println("Migrations Promotion Tables executed successfully")
}
func migrationAffiliate(db *gorm.DB){
	if err := db.AutoMigrate(&models.Referral{},&models.Partner{},&models.Affiliate{},&models.AffiliateTracking{},&models.Users{},&models.AffiliateLog{},&models.Promotion{});err != nil {
		fmt.Errorf("Tables schema migration not successfully\n")
	}
	fmt.Println("Migrations Affiliate Tables executed successfully")
}

type Setting struct {
	Value string `gorm:"column:value"` // Adjust the struct according to your table schema
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

	//var setting Setting

	setting,_ := GetMaster(prefix)

	// Determine the database name based on the retrieved value
	dbName := fmt.Sprintf("%s_%s", prefix, setting.Value)

	// Check if the connection already exists
	if db, exists := dbConnections[dbName]; exists {
		migrationPromotion(db)
		migrationAffiliate(db)
		return db, nil
	}

	// Read database prefixes and environment from environment variable
	//prefixes := strings.Split(os.Getenv("DB_PREFIXES"), ",")
	// env := os.Getenv("ENVIRONMENT") // Read the environment variable
	// var dbName string
	// suffix := "development" // Default to dev

	// if env == "production" {
	// 	suffix = "production"
	// }

	// if strings.Contains(prefix, suffix) {
	// 	//suffix = "production"
	// 	dbName = fmt.Sprintf("%s", prefix)
	// } else {
	// 	dbName = fmt.Sprintf("%s_%s", prefix, suffix)
	// }
 
	// Connect to the master database to read settings

	//defer masterDB.Close() // Ensure the masterDB connection is closed

	// Read the value from the settings table


	// Create the DSN for the selected database
	dsn := fmt.Sprintf(baseDSN, dbName)
	//fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
	})
	//migrationPromotion(db)
	if err == nil {
		dbConnections[dbName] = db // Store the connection in the map
		fmt.Println("Successfully connected to DB:", dbName)
	} else {
		return nil, err // Return the error if connection fails
	}
	if err := db.Exec("SET time_zone = 'Asia/Bangkok'"); err != nil {
		fmt.Printf("Warning: Failed to set timezone: %v\n", err)
	}
	//migrateNormal(db)
	//CheckAndCreateTable(db,models.BankStatement{})
	//migrationAffiliate(db)
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

 
func GetDBName(username string) (map[string]interface{},error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", //redis_master_password,
		DB:       0,  // ใช้ database 0
	})
	
	// keys, err := rdb.Keys(ctx, "*").Result() // ดึงทุกคีย์ใน Redis
	// if err != nil {
	// 	fmt.Println("Error retrieving keys from Redis:", err)
	// 	return  nil,err
	// }

	// // แสดงผลลัพธ์
	// fmt.Println("Keys in Redis:", keys)

	loginInfo, err := rdb.HMGet(ctx, username, "login_date", "prefix", "database").Result()
	if err != nil {
		fmt.Println("Error retrieving data from Redis:", err)
		return  map[string]interface{}{},err
	}
	if loginInfo[0] == nil && loginInfo[1] == nil && loginInfo[2] == nil {
		return nil, fmt.Errorf("loginInfo is empty: %v", loginInfo)
	}

	 

// แสดงผลลัพธ์
	response := map[string]interface{}{
			"login_date": loginInfo[0], // login_date
			"prefix":   loginInfo[1], // prefix
			"database":   loginInfo[2], // database
	}
	return response,nil
}
 
func GetDatabaseConnection(memberName string) (*gorm.DB, error) {
	dbdata, err := GetDBName(memberName)
	//fmt.Printf(" dbdata : %+v \n",dbdata)
	if err != nil {
		fmt.Printf(" dbdata Err: %s \n", err)
		return nil, err
	}
	
	prefix, ok := dbdata["prefix"].(string)
	if !ok {
		return nil, fmt.Errorf("prefix not found or not a string")
	}
	//fmt.Printf("Prefix: %s \n",prefix)
	db, err := ConnectToDB(prefix)
	if err != nil {
		fmt.Printf(" Connection Error Err: %s \n", err)
		return nil, err
	}
	
	return db, nil
}