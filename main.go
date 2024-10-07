package main

import (
	"log"
	"os"
	"server/controllers"
	"server/middlewares"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	// Connect to PostgreSQL using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&models.User{}, &models.Message{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize controllers
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	authController := controllers.AuthController{
		DB:     db,
		JwtKey: jwtKey,
	}

	contactsController := controllers.ContactsController{
		DB: db,
	}

	messagingController := controllers.MessageController{
		DB: db,
	}

	// Start the message broadcasting handler
	go messagingController.HandleBroadcast()

	// Initialize Gin router
	router := gin.Default()

	// Authentication routes
	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware(jwtKey))
	{
		protected.POST("/check-contacts", contactsController.CheckRegisteredContacts)
		protected.GET("/messages", messagingController.GetMessages)
		protected.GET("/ws", func(c *gin.Context) {
			messagingController.HandleConnections(c, jwtKey)
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
