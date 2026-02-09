package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"subscription-service/internal/models"
	"subscription-service/internal/services"
)

// SubscriptionHandler handles HTTP requests for subscription operations
type SubscriptionHandler struct {
	service *services.SubscriptionService
}

// NewSubscriptionHandler creates a new SubscriptionHandler
func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscription godoc
// @Summary      Create a new subscription
// @Description  Create a new subscription record
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      models.SubscriptionRequest  true  "Subscription data"
// @Success      201           {object}  models.SubscriptionResponse
// @Failure      400           {object}  models.ErrorResponse
// @Failure      500           {object}  models.ErrorResponse
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	logrus.Info("Handling POST /subscriptions request")
	var req models.SubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("Invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request format: " + err.Error()})
		return
	}

	response, err := h.service.CreateSubscription(&req)
	if err != nil {
		logrus.Errorf("Failed to create subscription: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	logrus.Infof("Subscription created successfully: %s", response.ID)
	c.JSON(http.StatusCreated, response)
}

// GetSubscription godoc
// @Summary      Get a subscription by ID
// @Description  Get subscription details by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "Subscription ID (UUID)"
// @Success      200  {object}  models.SubscriptionResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	logrus.Infof("Handling GET /subscriptions/%s request", c.Param("id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Warnf("Invalid subscription ID format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid subscription ID format"})
		return
	}

	response, err := h.service.GetSubscription(id)
	if err != nil {
		logrus.Errorf("Failed to get subscription: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if response == nil {
		logrus.Warnf("Subscription %s not found", id)
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "subscription not found"})
		return
	}

	logrus.Infof("Subscription %s retrieved successfully", id)
	c.JSON(http.StatusOK, response)
}

// GetAllSubscriptions godoc
// @Summary      Get all subscriptions
// @Description  Get a list of all subscriptions
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   models.SubscriptionResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	logrus.Info("Handling GET /subscriptions request")
	responses, err := h.service.GetAllSubscriptions()
	if err != nil {
		logrus.Errorf("Failed to get all subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	logrus.Infof("Retrieved %d subscriptions", len(responses))
	c.JSON(http.StatusOK, responses)
}

// UpdateSubscription godoc
// @Summary      Update a subscription
// @Description  Update subscription details by ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      string                      true  "Subscription ID (UUID)"
// @Param        subscription  body      models.SubscriptionRequest  true  "Updated subscription data"
// @Success      200           {object}  models.SubscriptionResponse
// @Failure      400           {object}  models.ErrorResponse
// @Failure      404           {object}  models.ErrorResponse
// @Failure      500           {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	logrus.Infof("Handling PUT /subscriptions/%s request", c.Param("id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Warnf("Invalid subscription ID format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid subscription ID format"})
		return
	}

	var req models.SubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("Invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request format: " + err.Error()})
		return
	}

	response, err := h.service.UpdateSubscription(id, &req)
	if err != nil {
		if err.Error() == "subscription not found" {
			logrus.Warnf("Subscription %s not found for update", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		logrus.Errorf("Failed to update subscription: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	logrus.Infof("Subscription %s updated successfully", id)
	c.JSON(http.StatusOK, response)
}

// DeleteSubscription godoc
// @Summary      Delete a subscription
// @Description  Delete a subscription by ID
// @Tags         subscriptions
// @Param        id   path  string  true  "Subscription ID (UUID)"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	logrus.Infof("Handling DELETE /subscriptions/%s request", c.Param("id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Warnf("Invalid subscription ID format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid subscription ID format"})
		return
	}

	err = h.service.DeleteSubscription(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			logrus.Warnf("Subscription %s not found for deletion", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		logrus.Errorf("Failed to delete subscription: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	logrus.Infof("Subscription %s deleted successfully", id)
	c.Status(http.StatusNoContent)
}

// AggregateTotalPrice godoc
// @Summary      Aggregate total price of subscriptions
// @Description  Calculate total price of subscriptions for a selected period with optional filtering by user ID and service name
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  false  "User ID (UUID) to filter by"
// @Param        service_name  query     string  false  "Service name to filter by"
// @Param        start_date    query     string  false  "Start date (MM-YYYY) of the period"
// @Param        end_date      query     string  false  "End date (MM-YYYY) of the period"
// @Success      200           {object}  models.AggregateResponse
// @Failure      400           {object}  models.ErrorResponse
// @Failure      500           {object}  models.ErrorResponse
// @Router       /subscriptions/aggregate [get]
func (h *SubscriptionHandler) AggregateTotalPrice(c *gin.Context) {
	logrus.Info("Handling GET /subscriptions/aggregate request")
	var req models.AggregateRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		logrus.Warnf("Invalid query parameters: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid query parameters: " + err.Error()})
		return
	}

	response, err := h.service.AggregateTotalPrice(&req)
	if err != nil {
		logrus.Errorf("Failed to aggregate total price: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	logrus.Infof("Aggregated total price: %d", response.TotalPrice)
	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers all subscription routes
func (h *SubscriptionHandler) RegisterRoutes(router *gin.Engine) {
	logrus.Info("Registering subscription routes")
	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("", h.CreateSubscription)
		subscriptions.GET("", h.GetAllSubscriptions)
		// IMPORTANT: /aggregate must be registered before /:id to avoid route conflicts
		subscriptions.GET("/aggregate", h.AggregateTotalPrice)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.DELETE("/:id", h.DeleteSubscription)
	}
}
