package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"subscription-service/internal/config"
	"subscription-service/internal/models"
)

// NewDatabaseConnection creates a new database connection and runs auto-migration
func NewDatabaseConnection(cfg *config.Config) *gorm.DB {
	logrus.Info("Connecting to database...")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	logrus.Debugf("Database DSN: host=%s port=%s dbname=%s user=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.User, cfg.Database.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	logrus.Info("Successfully connected to database")

	logrus.Info("Running auto-migration...")
	err = db.AutoMigrate(&models.Subscription{})
	if err != nil {
		logrus.Fatalf("Failed to auto-migrate database: %v", err)
	}
	logrus.Info("Auto-migration completed")

	return db
}
