package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// RedisClient is the connection handle you’ll reuse everywhere.
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis wires up the client and pings to verify connectivity.
func InitRedis() {
	db := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			log.Fatalf("invalid REDIS_DB value: %v", err)
		}
	}

	fmt.Printf("Connecting to Redis at %s, DB %d...\n", os.Getenv("REDIS_ADDR"), db)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"),     // e.g. "localhost:6379"
		Username:     os.Getenv("REDIS_USER"),     // "" if none
		Password:     os.Getenv("REDIS_PASSWORD"), // "" if none
		DB:           db,                          // 0 is default DB
		ReadTimeout:  3 * time.Second,             // optional hardening
		WriteTimeout: 3 * time.Second,
		PoolSize:     10, // tune per workload
	})

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
}
