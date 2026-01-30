package main

import (
	"backend-dashboard/internal/config"
	"backend-dashboard/internal/delivery/http"
	"backend-dashboard/internal/delivery/http/middleware"
	"backend-dashboard/internal/repository"
	"backend-dashboard/internal/usecase"
	"backend-dashboard/pkg/database"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect to Database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// 3. Migration & Seeding
	database.AutoMigrate(db)
	database.SeedRoles(db)
	database.SeedSuperAdmin(db)
	database.SeedSampleData(db)

	// 4. Init Layers
	userRepo := repository.NewPostgresUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)

	// 5. Setup Router
	r := gin.Default()

	// CORS Middleware (Simple version for development)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Rate Limit: 5 requests every 1 minute
	// rate.Every(1 * time.Minute / 5) = 1 token every 12 seconds
	// Burst 5 = allow 5 requests immediately
	loginRateLimiter := middleware.RateLimitMiddleware(rate.Every(1*time.Minute/5), 5)

	// User Handler
	userHandler := http.NewUserHandler(r, userUsecase)

	// Public routes
	r.POST("/login", loginRateLimiter, userHandler.Login)

	// API routes (in production, these should be protected with JWT middleware)
	api := r.Group("/api")
	{
		// User management
		api.GET("/users", userHandler.GetUsers)
		api.GET("/users/:id", userHandler.GetUser)
		api.POST("/users", userHandler.CreateUser)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeactivateUser)

		// Hard delete - only for super_admin (TODO: protect with JWT middleware)
		api.DELETE("/users/:id/permanent", userHandler.HardDeleteUser)
	}

	// 6. Run Server
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
