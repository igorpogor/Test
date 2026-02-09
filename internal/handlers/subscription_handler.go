package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"subscription-service/internal/models"
	"subscription-service/internal/services"
)

type SubscriptionHandler struct {
	service *services.SubscriptionService
}

func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req models.SubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	response, err := h.service.CreateSubscription(&req)
	if err != nil {
		logrus.Errorf("Failed to create subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Invalid subscription ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID format"})
		return
	}

	response, err := h.service.GetSubscription(id)
	if err != nil {
		logrus.Errorf("Failed to get subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	responses, err := h.service.GetAllSubscriptions()
	if err != nil {
		logrus.Errorf("Failed to get all subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Invalid subscription ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID format"})
		return
	}

	var req models.SubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	response, err := h.service.UpdateSubscription(id, &req)
	if err != nil {
		logrus.Errorf("Failed to update subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Invalid subscription ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID format"})
		return
	}

	err = h.service.DeleteSubscription(id)
	if err != nil {
		logrus.Errorf("Failed to delete subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) AggregateTotalPrice(c *gin.Context) {
	var req models.AggregateRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		logrus.Errorf("Invalid query parameters: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	response, err := h.service.AggregateTotalPrice(&req)
	if err != nil {
		logrus.Errorf("Failed to aggregate total price: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SubscriptionHandler) RegisterRoutes(router *gin.Engine) {
	subscriptionsGroup := router.Group("/subscriptions")
	{
		subscriptionsGroup.POST("", h.CreateSubscription)
		subscriptionsGroup.GET("/:id", h.GetSubscription)
		subscriptionsGroup.GET("", h.GetAllSubscriptions)
		subscriptionsGroup.PUT("/:id", h.UpdateSubscription)
		subscriptionsGroup.DELETE("/:id", h.DeleteSubscription)
		subscriptionsGroup.GET("/aggregate", h.AggregateTotalPrice)
	}
}
