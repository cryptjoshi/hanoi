package database

import (
	"fmt"
	"os"
 	"strings"
	"sync"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
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

const baseDSN = "web:1688XdAs@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
//const baseDSN = "root:1688XdAs@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
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
// Connect function to establish a database connection based on the prefix
func ConnectToDB(prefix string) (*gorm.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()
	
	prefix = strings.ToLower(prefix)
 
	
	// Check if the connection already exists
	if db, exists := dbConnections[prefix]; exists {
		return db, nil
	}

	// Read database prefixes and environment from environment variable
	//prefixes := strings.Split(os.Getenv("DB_PREFIXES"), ",")
	env := os.Getenv("ENVIRONMENT") // Read the environment variable
	var dbName string
	suffix := "development" // Default to dev

	if env == "production" {
		suffix = "production"
	}
	// fmt.Printf("%s %s",prefixes,suffix)
	// // Determine the database name based on the prefix
	// if contains(prefixes, prefix) {
	// 	dbName = fmt.Sprintf("%s_%s", prefix, suffix)
	// } else {
	// 	//return nil, fmt.Errorf("unknown prefix: %s", prefix)
	// 	dbName = fmt.Sprintf("%s", prefix)
	// }
	//fmt.Println("Prefix:"+prefix)
	if strings.Contains(prefix, suffix)  {
		//suffix = "production"
		dbName = fmt.Sprintf("%s", prefix)
	} else {
		dbName = fmt.Sprintf("%s_%s", prefix, suffix)
	}
	// Create the DSN for the selected database
	dsn := fmt.Sprintf(baseDSN, dbName)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction: true,
		PrepareStmt: true,
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
