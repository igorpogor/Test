package main

import (
	"log"

	_ "subscription-service/docs"
	"subscription-service/internal/config"
	"subscription-service/internal/database"
	"subscription-service/internal/handlers"
	"subscription-service/internal/repositories"
	"subscription-service/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewDatabaseConnection(cfg)

	subscriptionRepo := repositories.NewSubscriptionRepository(db)

	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	router := gin.Default()

	router.Use(cors.Default())

	subscriptionHandler.RegisterRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
