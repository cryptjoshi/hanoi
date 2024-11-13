package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"os"
	"strings"
	"sync"
)

var Database *gorm.DB
var (
	dbConnections = make(map[string]*gorm.DB)
	mutex         sync.Mutex
)

//var DSN = 'root:helloworld@tcp(db:3306)/tsxbet_dev?tls=true'
//var DSN string = "root:helloworld@tcp(db:3306)/ckd_development?charset=utf8mb4&parseTime=True&loc=Local"
const baseDSN = "web:1688XdAs@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"

type Setting struct {
	Value string `gorm:"column:value"` // Adjust the struct according to your table schema
}
func Connect() error {
	var err error


	
	dsn := fmt.Sprintf(baseDSN)//, "ckd_development")
	//dsn := fmt.Sprintf("%s&parseTime=True", baseDSN)// os.Getenv("DSN"))

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
		return db, nil
	}
	// prefix = strings.ToLower(prefix)

	// // Check if the connection already exists
	// if db, exists := dbConnections[prefix]; exists {
	// 	return db, nil
	// }

	// // Read database prefixes and environment from environment variable
	// //prefixes := strings.Split(os.Getenv("DB_PREFIXES"), ",")
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
	// // Determine the database name based on the prefix
	// if contains(prefixes, prefix) {
	// 	dbName = fmt.Sprintf("%s_%s", prefix, suffix)
	// } else {
	// 	//return nil, fmt.Errorf("unknown prefix: %s", prefix)
	// 	dbName = fmt.Sprintf("%s", prefix)
	// }

	// Create the DSN for the selected database
	dsn := fmt.Sprintf(baseDSN, dbName)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
	})

	if err == nil {
		dbConnections[prefix] = db // Store the connection in the map
		fmt.Println("Successfully connected to DB:", dbName)
	} else {
		return nil, err // Return the error if connection fails
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
func ConnectToDBX(prefix string) (*gorm.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()

	prefix = strings.ToLower(prefix)

	//var setting Setting

	setting,_ := getMaster(prefix)

	// Determine the database name based on the retrieved value
	dbName := fmt.Sprintf("%s_%s", prefix, setting.Value)

	// Check if the connection already exists
	if db, exists := dbConnections[dbName]; exists {
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
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
	})

	if err == nil {
		dbConnections[dbName] = db // Store the connection in the map
		fmt.Println("Successfully connected to DB:", dbName)
	} else {
		return nil, err // Return the error if connection fails
	}
	//migrateNormal(db)
	return db, nil
}