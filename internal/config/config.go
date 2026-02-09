package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Logging struct {
		Level string
	}
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Warnf("Could not load .env file: %v (using environment variables)", err)
	} else {
		logrus.Info("Loaded configuration from .env file")
	}

	logLevel := getEnv("LOG_LEVEL", "info")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
		logrus.Warnf("Invalid log level '%s', defaulting to 'info'", logLevel)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg := &Config{
		Server: struct {
			Port string
		}{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
			SSLMode  string
		}{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "subscription_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Logging: struct {
			Level string
		}{
			Level: logLevel,
		},
	}

	logrus.Infof("Configuration loaded: server_port=%s, db_host=%s, db_port=%s, db_name=%s, log_level=%s",
		cfg.Server.Port, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Logging.Level)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
