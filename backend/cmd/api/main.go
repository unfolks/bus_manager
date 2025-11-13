package main

import (
	"log"
	"os"

	"bus-manager/internal/database"
	"bus-manager/internal/handlers"
	"bus-manager/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis
	redisClient := database.InitRedis()

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORS())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, redisClient)
	gameHandler := handlers.NewGameHandler(db, redisClient)

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "bus-manager-api",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// Game routes (protected)
		game := api.Group("/game")
		game.Use(middleware.AuthMiddleware(redisClient))
		{
			game.GET("/company", gameHandler.GetCompany)
			game.POST("/company", gameHandler.CreateCompany)
			game.GET("/depots", gameHandler.GetDepots)
			game.POST("/depots", gameHandler.CreateDepot)
			game.GET("/buses", gameHandler.GetBuses)
			game.POST("/buses", gameHandler.CreateBus)
			game.GET("/routes", gameHandler.GetRoutes)
			game.POST("/trips", gameHandler.CreateTrip)
			game.GET("/trips/active", gameHandler.GetActiveTrips)
		}
	}

	// WebSocket route for real-time updates
	r.GET("/ws/trips", handlers.HandleWebSocket)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
