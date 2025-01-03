// promotionUtils.go
package handler

import (
	"fmt"
	"time"

	"crypto/sha256"
	"encoding/hex"

	"github.com/go-redis/redis/v8"
)

// Generate a unique promotion ID
func GenerateUID(userID, proID string) string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := sha256.New()
	hash.Write([]byte(userID + proID + timestamp))
	return hex.EncodeToString(hash.Sum(nil))
}

// Validate if a promotion can be applied
func ValidatePromotion(redisClient *redis.Client, userID, promotionKey string) error {
	status, err := redisClient.HGet(ctx, promotionKey, "status").Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("error fetching promotion status: %v", err)
	}
	if status != "0" && status != "2" {
		return fmt.Errorf("promotion is not eligible: status %s", status)
	}
	return nil
}
