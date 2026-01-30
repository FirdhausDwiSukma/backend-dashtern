package database

import (
	"backend-dashboard/internal/config"
	"backend-dashboard/internal/domain"
	"fmt"
	"log"
	"time"

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
	// Migrate all models in correct order (dependencies first)
	err := db.AutoMigrate(
		&domain.Role{},
		&domain.User{},
		&domain.InternProfile{},
		&domain.PICProfile{},
		&domain.HRProfile{},
		&domain.Task{},
		&domain.Attendance{},
		&domain.MentorReview{},
		&domain.PerformanceScore{},
		&domain.PotentialScore{},
		&domain.NineGridResult{},
		&domain.AuditLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
}

func SeedRoles(db *gorm.DB) {
	roles := []domain.Role{
		{Name: "super_admin", Description: "Super Administrator with full system access"},
		{Name: "hr", Description: "Human Resources staff"},
		{Name: "pic", Description: "Person In Charge / Mentor for interns"},
		{Name: "intern", Description: "Internship participant"},
	}

	for _, role := range roles {
		var count int64
		db.Model(&domain.Role{}).Where("name = ?", role.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to seed role %s: %v", role.Name, err)
			} else {
				log.Printf("Role '%s' seeded successfully", role.Name)
			}
		}
	}
}

func SeedSuperAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ?", "super_admin").Count(&count)

	if count == 0 {
		// Get super_admin role
		var superAdminRole domain.Role
		if err := db.Where("name = ?", "super_admin").First(&superAdminRole).Error; err != nil {
			log.Fatalf("Failed to find super_admin role: %v", err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		now := time.Now()
		admin := domain.User{
			FullName:     "Super Administrator",
			Username:     "admin",
			Email:        "admin@dashtern.com",
			PasswordHash: string(hashedPassword),
			RoleID:       superAdminRole.ID,
			Status:       "active",
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to seed super admin: %v", err)
		}
		log.Println("Super Admin seeded successfully:")
		log.Println("  Username: admin")
		log.Println("  Password: password123")
		log.Println("  Email: admin@dashtern.com")
	} else {
		log.Println("Super Admin already exists, skipping seed.")
	}
}

func SeedSampleData(db *gorm.DB) {
	// Check if sample data already exists
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	if userCount > 1 { // More than just super admin
		log.Println("Sample data already exists, skipping seed.")
		return
	}

	log.Println("Seeding sample data...")

	// Get roles
	var hrRole, picRole, internRole domain.Role
	db.Where("name = ?", "hr").First(&hrRole)
	db.Where("name = ?", "pic").First(&picRole)
	db.Where("name = ?", "intern").First(&internRole)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	now := time.Now()

	// Create HR user
	hrUser := domain.User{
		FullName:     "Jane Doe",
		Username:     "hr_jane",
		Email:        "jane.hr@dashtern.com",
		PasswordHash: string(hashedPassword),
		RoleID:       hrRole.ID,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	db.Create(&hrUser)

	hrProfile := domain.HRProfile{
		UserID:     hrUser.ID,
		Department: "Human Resources",
		CreatedAt:  now,
	}
	db.Create(&hrProfile)

	// Create PIC users
	pic1 := domain.User{
		FullName:     "John Smith",
		Username:     "pic_john",
		Email:        "john.pic@dashtern.com",
		PasswordHash: string(hashedPassword),
		RoleID:       picRole.ID,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	db.Create(&pic1)

	picProfile1 := domain.PICProfile{
		UserID:    pic1.ID,
		Position:  "Senior Developer",
		Division:  "Engineering",
		Expertise: "Backend Development, Golang",
		CreatedAt: now,
	}
	db.Create(&picProfile1)

	pic2 := domain.User{
		FullName:     "Alice Williams",
		Username:     "pic_alice",
		Email:        "alice.pic@dashtern.com",
		PasswordHash: string(hashedPassword),
		RoleID:       picRole.ID,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	db.Create(&pic2)

	picProfile2 := domain.PICProfile{
		UserID:    pic2.ID,
		Position:  "UI/UX Lead",
		Division:  "Design",
		Expertise: "Frontend Development, React",
		CreatedAt: now,
	}
	db.Create(&picProfile2)

	// Create Intern users
	interns := []struct {
		fullName  string
		username  string
		email     string
		pic       domain.User
		batch     string
		division  string
		education string
	}{
		{"Bob Johnson", "intern_bob", "bob@example.com", pic1, "2026-01", "Engineering", "Computer Science"},
		{"Charlie Brown", "intern_charlie", "charlie@example.com", pic1, "2026-01", "Engineering", "Software Engineering"},
		{"Diana Prince", "intern_diana", "diana@example.com", pic2, "2026-01", "Design", "Graphic Design"},
	}

	for _, intern := range interns {
		user := domain.User{
			FullName:     intern.fullName,
			Username:     intern.username,
			Email:        intern.email,
			PasswordHash: string(hashedPassword),
			RoleID:       internRole.ID,
			Status:       "active",
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		db.Create(&user)

		profile := domain.InternProfile{
			UserID:    user.ID,
			PICID:     intern.pic.ID,
			Batch:     intern.batch,
			Division:  intern.division,
			StartDate: now.AddDate(0, -1, 0), // Started 1 month ago
			EndDate:   now.AddDate(0, 2, 0),  // Ends in 2 months
			Education: intern.education,
			CreatedAt: now,
		}
		db.Create(&profile)
	}

	log.Println("Sample data seeded successfully")
	log.Println("  HR: hr_jane / password123")
	log.Println("  PIC: pic_john / password123")
	log.Println("  PIC: pic_alice / password123")
	log.Println("  Intern: intern_bob / password123")
	log.Println("  Intern: intern_charlie / password123")
	log.Println("  Intern: intern_diana / password123")
}
