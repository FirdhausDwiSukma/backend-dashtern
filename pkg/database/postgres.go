package database

import (
	"backend-dashboard/internal/config"
	"backend-dashboard/internal/domain"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Try to connect to the specific database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Attempt to create database if it doesn't exist (this requires connecting to default 'postgres' db first)
		// For simplicity in this task, we will assume the DB exists or print a clear message.
		// Detailed auto-creation logic can be complex due to permissions.
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
}

func SeedSuperAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("role = ?", "super_admin").Count(&count)

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		admin := domain.User{
			Username: "admin",
			Password: string(hashedPassword),
			Role:     "super_admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to seed super admin: %v", err)
		}
		log.Println("Super Admin seeded successfully: username=admin, password=password123")
	} else {
		log.Println("Super Admin already exists, skipping seed.")
	}
}
