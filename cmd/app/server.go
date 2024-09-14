package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"
)

// Singleton
var (
	db   *gorm.DB
	once sync.Once
)

func redisConnection() *redis.Client {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379"
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Error parsing Redis URL: %v", err)
	}

	client := redis.NewClient(opts)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	fmt.Println("Connected successfully to Redis!")
	return client
}

func newDB() (*gorm.DB, error) {
	//Singleton
	once.Do(func() {
		var err error
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		dsn := os.Getenv("MYSQL_URL")

		if dsn == "" {
			log.Fatal("Environment variable DB_URL is not set")
		}

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic("failed to connect database")
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get database connection pool: %v", err)
		}

		// Cấu hình connection pool chung
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(40 * time.Minute)

	})
	return db, nil
}
