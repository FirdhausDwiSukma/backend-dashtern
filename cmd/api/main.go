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

	internRepo := repository.NewInternRepository(db)
	internUsecase := usecase.NewInternUsecase(internRepo, userRepo)

	// 5. Setup Router
	r := gin.Default()

	// CORS Middleware (Simple version for development)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

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

	// Middlewares
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	superAdminOnly := middleware.RoleMiddleware(1) // role_id 1 = super_admin
	hrOrAbove := middleware.RoleMiddleware(1, 2)   // role_id 1,2 = super_admin, hr
	// picOrAbove := middleware.RoleMiddleware(1, 2, 3)   // role_id 1,2,3 = super_admin, hr, pic (for future use)

	// Handlers
	userHandler := http.NewUserHandler(r, userUsecase)
	internHandler := http.NewInternHandler(internUsecase)
	profileHandler := http.NewProfileHandler(userUsecase)

	// Public routes
	r.POST("/login", loginRateLimiter, userHandler.Login)

	// Protected API routes
	api := r.Group("/api")
	api.Use(authMiddleware) // All API routes require authentication
	{
		// User management (HR or above)
		users := api.Group("/users")
		users.Use(hrOrAbove)
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeactivateUser)

			// Hard delete - only for super_admin
			users.DELETE("/:id/permanent", superAdminOnly, userHandler.HardDeleteUser)
		}

		// Intern management (HR or above can create, all authenticated users can view)
		interns := api.Group("/interns")
		{
			interns.POST("", hrOrAbove, internHandler.CreateIntern)
			interns.GET("", internHandler.GetInterns)
			interns.GET("/:id", internHandler.GetIntern)
		}

		// Profile management (all authenticated users)
		profile := api.Group("/profile")
		{
			profile.GET("", profileHandler.GetProfile)
			profile.PUT("", profileHandler.UpdateProfile)
			profile.PUT("/password", profileHandler.UpdatePassword)
			profile.POST("/avatar", profileHandler.UpdateAvatar)
		}
	}

	// 6. Run Server
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
