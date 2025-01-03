// redisUtils.go
package handler

import (
	"context"
	"fmt"
	"time"
	"os"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redis_master_host = os.Getenv("REDIS_HOST")
var redis_master_port = os.Getenv("REDIS_PORT")
// Try to acquire a lock
func AcquireLock(redisClient *redis.Client, lockKey string, ttl time.Duration) (bool, error) {
	success, err := redisClient.SetNX(ctx, lockKey, "locked", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("error acquiring lock: %v", err)
	}
	return success, nil
}

// Release the lock
func ReleaseLock(redisClient *redis.Client, lockKey string) error {
	_, err := redisClient.Del(ctx, lockKey).Result()
	if err != nil {
		return fmt.Errorf("error releasing lock: %v", err)
	}
	return nil
}

func CreateRedisClient() *redis.Client {
	return 	 redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: "", //redis_master_password,
		DB:       0,  // ใช้ database 0
	})
 }