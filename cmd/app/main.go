package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	routerv1 "tart-shop-manager/api/router/v1"
	policiesutil "tart-shop-manager/internal/util/policies"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := newDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	fmt.Println("Database connection successful", db)
	rdb := redisConnection()
	fmt.Println("Redis connection successful", rdb)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	// // Define model and policy paths
	modelPath := filepath.Join(cwd, "config/casbin", "rbac_model.conf")

	_, err = policiesutil.InitEnforcer(db, modelPath)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	r := routerv1.NewRouter(db, rdb)
	if err := r.Run(port); err != nil {
		return
	} // listen and serve (for windows "localhost:3000")
}
