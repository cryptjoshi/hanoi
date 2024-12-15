package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
	"github.com/go-redis/redis/v8" 
	"context" // เพิ่มการนำเข้า Redis
	//"golang.org/x/net/context" // เพิ่มการนำเข้าคอนเท็กซ์
	
	
	
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

var mysql_host = os.Getenv("MYSQL_HOST")
var mysql_user = os.Getenv("MYSQL_ROOT_USER")
var mysql_pass = os.Getenv("MYSQL_ROOT_PASSWORD")
var jwtSecret = os.Getenv("PASSWORD_SECRET")

var ctx = context.Background() 
var Database *gorm.DB
var (
	dbConnections = make(map[string]*gorm.DB)
	mutex         sync.Mutex
)
var redisClient *redis.Client // สร้างตัวแปรสำหรับ Redis client

type DatabaseInfo struct {
	Prefix string   `json:"prefix"`
	Names  []string `json:"names"`
}


	// ฟังก์ชันสำหรับการเชื่อมต่อกับ Redis
func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // เปลี่ยนที่อยู่ตามที่ตั้งของ Redis
	})
}
//var DSN = 'root:helloworld@tcp(db:3306)/tsxbet_dev?tls=true'
//var DSN string = "root:helloworld@tcp(db:3306)/ckd_development?charset=utf8mb4&parseTime=True&loc=Local"
const baseDSN = "web:1688XdAs@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FBangkok"

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
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
		// Close the existing connection before opening a new one
		// sqlDB, _ := db.DB()
		// sqlDB.Close()
		// delete(dbConnections, dbName) // Remove the old connection from the map
		return db,nil
	}
	
	dsn := fmt.Sprintf(baseDSN, dbName)
	
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", //redis_master_password,
		DB:       0,  // ใช้ database 0
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return Setting{}, err
		// อาจจะ return error หรือจัดการข้อผิดพลาดตามที่คุณต้องการ
	} else {
		fmt.Println("เชื่อมต่อ Redis สำเร็จ:", pong)
	}
	
	// สร้าง context สำหรับ Redis

	// ตรวจสอบค่าจาก Redis
	val, err := rdb.Get(ctx, prefix).Result()
	if err == nil {
		// ถ้ามีค่าใน Redis ให้คืนค่า Setting
		return Setting{Value: val}, nil
	}
	// ตรวจสอบว่ามีข้อมูลใน Redis หรือไม่
	

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

	// เก็บค่าลง Redis
	err = rdb.Set(ctx, prefix, setting.Value, 0).Err() // 0 หมายถึงไม่มีการหมดอายุ
	if err != nil {
		fmt.Printf("Error setting value in Redis: %s\n", err)
	}

	 //masterDB.Close()
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



func getDataList() ([]DatabaseInfo,error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", mysql_user, mysql_pass, mysql_host)
	fmt.Printf(" DSN: %s \n",dsn)
	// Connect to MySQL without a specific database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil,err
	}
	
	// Query to get all databases
	groupedDatabases := make(map[string][]string)

	rows, err := db.Raw("SHOW DATABASES").Rows()
	fmt.Printf("Rows: %v \n",rows)
	fmt.Printf("Err: %v \n",err)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	systemDatabases := map[string]bool{
		"information_schema": true,
		"mysql":              true,
		"performance_schema": true,
		"sys":                true,
	}

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			continue
		}
		fmt.Printf("DBName: %v",dbName)
		fmt.Printf("SysDBName: %v",systemDatabases[dbName])
		// Include databases with underscore in their names and exclude system databases
		//
		if strings.Contains(dbName, "_") && !systemDatabases[dbName] {
			parts := strings.SplitN(dbName, "_", 2)
			if len(parts) == 2 {
				prefix := parts[0]
				if _, exists := groupedDatabases[prefix]; !exists {
					groupedDatabases[prefix] = []string{}
				}
				groupedDatabases[prefix] = append(groupedDatabases[prefix], dbName)
			}
		}
	}
	var databaseList []DatabaseInfo
	for prefix, names := range groupedDatabases {
		databaseList = append(databaseList, DatabaseInfo{
			Prefix: prefix,
			Names:  names,
		})
	}
	return databaseList,nil

}

func GetDBName(username string) (map[string]interface{},error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", //redis_master_password,
		DB:       0,  // ใช้ database 0
	})

	loginInfo, err := rdb.HMGet(ctx, username, "login_date", "prefix", "database").Result()
	if err != nil {
		fmt.Println("Error retrieving data from Redis:", err)
		return  map[string]interface{}{},err
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
	if err != nil {
		fmt.Printf(" dbdata Err: %s \n", err)
		return nil, err
	}
	
	prefix, ok := dbdata["prefix"].(string)
	if !ok {
		return nil, fmt.Errorf("prefix not found or not a string")
	}
	fmt.Printf("Prefix: %s \n",prefix)
	db, err := ConnectToDB(prefix)
	if err != nil {
		fmt.Printf(" Connection Error Err: %s \n", err)
		return nil, err
	}
	
	return db, nil
}