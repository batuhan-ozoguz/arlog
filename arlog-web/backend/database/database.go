package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"arlog/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds the database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect establishes a connection to the PostgreSQL database
func Connect(config Config) error {
	// Build connection string (DSN)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// If in production, set logger to Error level only
	if os.Getenv("ENVIRONMENT") == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying SQL database
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection established successfully")

	return nil
}

// Migrate runs database migrations for all models
func Migrate() error {
	log.Println("🔄 Running database migrations...")

	err := DB.AutoMigrate(
		&models.Team{},
		&models.Permission{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// SeedDatabase populates the database with initial test data
// This is useful for development and testing
func SeedDatabase() error {
	log.Println("🌱 Seeding database with test data...")

	// Check if data already exists
	var count int64
	DB.Model(&models.Team{}).Count(&count)
	if count > 0 {
		log.Println("ℹ️  Database already contains data, skipping seed")
		return nil
	}

	// Create test teams
	cosmosTeam := models.Team{
		TeamName:    "Cosmos Team",
		OktaGroupID: "cosmos-team-okta-group",
	}

	jupiterTeam := models.Team{
		TeamName:    "Jupiter Team",
		OktaGroupID: "jupiter-team-okta-group",
	}

	if err := DB.Create(&cosmosTeam).Error; err != nil {
		return fmt.Errorf("failed to create cosmos team: %w", err)
	}

	if err := DB.Create(&jupiterTeam).Error; err != nil {
		return fmt.Errorf("failed to create jupiter team: %w", err)
	}

	// Create test permissions for Cosmos Team
	cosmosPermissions := []models.Permission{
		{
			TeamID:              cosmosTeam.ID,
			ClusterName:         "dev-cluster",
			Namespace:           "cosmos-namespace",
			ServiceAccountToken: "dummy-token-cosmos-dev",
		},
		{
			TeamID:              cosmosTeam.ID,
			ClusterName:         "test-cluster",
			Namespace:           "cosmos-namespace",
			ServiceAccountToken: "dummy-token-cosmos-test",
		},
	}

	// Create test permissions for Jupiter Team
	jupiterPermissions := []models.Permission{
		{
			TeamID:              jupiterTeam.ID,
			ClusterName:         "dev-cluster",
			Namespace:           "jupiter-namespace",
			ServiceAccountToken: "dummy-token-jupiter-dev",
		},
	}

	// Insert permissions
	for _, perm := range cosmosPermissions {
		if err := DB.Create(&perm).Error; err != nil {
			return fmt.Errorf("failed to create permission: %w", err)
		}
	}

	for _, perm := range jupiterPermissions {
		if err := DB.Create(&perm).Error; err != nil {
			return fmt.Errorf("failed to create permission: %w", err)
		}
	}

	log.Println("✅ Database seeded successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}


