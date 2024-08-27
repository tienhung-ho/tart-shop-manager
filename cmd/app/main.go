package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	//port := os.Getenv("PORT")

	//rdb := redisConnection()

	db, err := NewDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	fmt.Println("Database connection successful")
	fmt.Printf("%+v\n", db)
	rdb := redisConnection()
	fmt.Println("Redis connection successful", rdb)
	// Thực hiện các thao tác với cơ sở dữ liệu ở đây
}
