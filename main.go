package main

import (
	"dirwatcher/api"
	"dirwatcher/database"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPass := ""

	startTask := make(chan bool)
	stopTask := make(chan bool)
	fmt.Println("connecting DB...")

	// Initialize database
	db, err := database.InitDB(dbHost, dbPort, dbUser, dbPass)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// Register API routes
	api.RegisterAPIRoutes(router, db, startTask, stopTask)

	// Start the API server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
