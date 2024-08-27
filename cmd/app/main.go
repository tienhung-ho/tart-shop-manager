package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	routerv1 "tart-shop-manager/api/router/v1"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := NewDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	fmt.Println("Database connection successful", db)
	rdb := redisConnection()
	fmt.Println("Redis connection successful", rdb)
	// Thực hiện các thao tác với cơ sở dữ liệu ở đây

	port := os.Getenv("PORT")
	r := routerv1.NewRouter(db, rdb)
	if err := r.Run(port); err != nil {
		return
	} // listen and serve (for windows "localhost:3000")
}
