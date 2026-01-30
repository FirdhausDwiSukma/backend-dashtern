package main

import (
	"backend-dashboard/internal/config"
	"backend-dashboard/internal/domain"
	"backend-dashboard/pkg/database"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	log.Println("Dropping all tables...")

	// Drop all tables in reverse order (to respect foreign keys)
	db.Migrator().DropTable(
		&domain.AuditLog{},
		&domain.NineGridResult{},
		&domain.PotentialScore{},
		&domain.PerformanceScore{},
		&domain.MentorReview{},
		&domain.Attendance{},
		&domain.Task{},
		&domain.HRProfile{},
		&domain.PICProfile{},
		&domain.InternProfile{},
		&domain.User{},
		&domain.Role{},
	)

	log.Println("All tables dropped successfully")

	// Now run migration
	database.AutoMigrate(db)
	database.SeedRoles(db)
	database.SeedSuperAdmin(db)
	database.SeedSampleData(db)

	log.Println("Migration and seeding completed!")
}
