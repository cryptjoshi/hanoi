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

func getMaster(prefix string) (Setting,error) {

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

	setting,_ := getMaster(prefix)

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
	migrateNormal(db)
	//CheckAndCreateTable(db,models.BankStatement{})
	migrationAffiliate(db)
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
