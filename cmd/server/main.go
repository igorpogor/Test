package main

import (
	"github.com/sirupsen/logrus"

	_ "subscription-service/docs"
	"subscription-service/internal/config"
	"subscription-service/internal/database"
	"subscription-service/internal/handlers"
	"subscription-service/internal/repositories"
	"subscription-service/internal/services"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Subscription Service API
// @version         1.0
// @description     REST API service for aggregating user online subscription data.
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.LoadConfig()

	logrus.Info("Initializing database connection...")
	db := database.NewDatabaseConnection(cfg)

	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	router := gin.Default()

	subscriptionHandler.RegisterRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logrus.Info("Swagger documentation available at /swagger/index.html")

	logrus.Infof("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
