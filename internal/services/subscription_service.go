package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"subscription-service/internal/models"
	"subscription-service/internal/repositories"
)

// SubscriptionService handles business logic for subscription operations
type SubscriptionService struct {
	repo *repositories.SubscriptionRepository
}

// NewSubscriptionService creates a new SubscriptionService
func NewSubscriptionService(repo *repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// validateDateFormat validates that the date string is in MM-YYYY format
func validateDateFormat(dateStr string) error {
	_, err := time.Parse("01-2006", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format '%s', expected MM-YYYY", dateStr)
	}
	return nil
}

// CreateSubscription creates a new subscription record
func (s *SubscriptionService) CreateSubscription(req *models.SubscriptionRequest) (*models.SubscriptionResponse, error) {
	logrus.Infof("Creating subscription for service '%s'", req.ServiceName)

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		logrus.Warnf("Invalid user ID format: %v", err)
		return nil, errors.New("invalid user ID format")
	}

	if err := validateDateFormat(req.StartDate); err != nil {
		logrus.Warnf("Invalid start_date: %v", err)
		return nil, err
	}

	if req.EndDate != nil {
		if err := validateDateFormat(*req.EndDate); err != nil {
			logrus.Warnf("Invalid end_date: %v", err)
			return nil, err
		}
	}

	subscription := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err = s.repo.Create(subscription)
	if err != nil {
		logrus.Errorf("Failed to create subscription: %v", err)
		return nil, err
	}

	response := toSubscriptionResponse(subscription)
	logrus.Infof("Created subscription %s for user %s", subscription.ID, userID)
	return response, nil
}

// GetSubscription retrieves a subscription by its ID
func (s *SubscriptionService) GetSubscription(id uuid.UUID) (*models.SubscriptionResponse, error) {
	logrus.Infof("Getting subscription %s", id)

	subscription, err := s.repo.GetByID(id)
	if err != nil {
		logrus.Errorf("Failed to get subscription %s: %v", id, err)
		return nil, err
	}

	if subscription == nil {
		logrus.Warnf("Subscription %s not found", id)
		return nil, nil
	}

	response := toSubscriptionResponse(subscription)
	logrus.Infof("Retrieved subscription %s", id)
	return response, nil
}

// GetAllSubscriptions retrieves all subscriptions
func (s *SubscriptionService) GetAllSubscriptions() ([]models.SubscriptionResponse, error) {
	logrus.Info("Getting all subscriptions")

	subscriptions, err := s.repo.GetAll()
	if err != nil {
		logrus.Errorf("Failed to get all subscriptions: %v", err)
		return nil, err
	}

	responses := make([]models.SubscriptionResponse, 0, len(subscriptions))
	for _, sub := range subscriptions {
		responses = append(responses, *toSubscriptionResponse(&sub))
	}

	logrus.Infof("Retrieved %d subscriptions", len(responses))
	return responses, nil
}

// UpdateSubscription updates a subscription by its ID
func (s *SubscriptionService) UpdateSubscription(id uuid.UUID, req *models.SubscriptionRequest) (*models.SubscriptionResponse, error) {
	logrus.Infof("Updating subscription %s", id)

	existing, err := s.repo.GetByID(id)
	if err != nil {
		logrus.Errorf("Failed to get subscription %s for update: %v", id, err)
		return nil, err
	}

	if existing == nil {
		logrus.Warnf("Subscription %s not found for update", id)
		return nil, errors.New("subscription not found")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		logrus.Warnf("Invalid user ID format: %v", err)
		return nil, errors.New("invalid user ID format")
	}

	if err := validateDateFormat(req.StartDate); err != nil {
		logrus.Warnf("Invalid start_date: %v", err)
		return nil, err
	}

	if req.EndDate != nil {
		if err := validateDateFormat(*req.EndDate); err != nil {
			logrus.Warnf("Invalid end_date: %v", err)
			return nil, err
		}
	}

	existing.ServiceName = req.ServiceName
	existing.Price = req.Price
	existing.UserID = userID
	existing.StartDate = req.StartDate
	existing.EndDate = req.EndDate

	err = s.repo.Update(existing)
	if err != nil {
		logrus.Errorf("Failed to update subscription %s: %v", id, err)
		return nil, err
	}

	response := toSubscriptionResponse(existing)
	logrus.Infof("Updated subscription %s", id)
	return response, nil
}

// DeleteSubscription deletes a subscription by its ID
func (s *SubscriptionService) DeleteSubscription(id uuid.UUID) error {
	logrus.Infof("Deleting subscription %s", id)

	err := s.repo.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			logrus.Warnf("Subscription %s not found for deletion", id)
			return errors.New("subscription not found")
		}
		logrus.Errorf("Failed to delete subscription %s: %v", id, err)
		return err
	}

	logrus.Infof("Deleted subscription %s", id)
	return nil
}

// AggregateTotalPrice calculates the total cost of subscriptions with optional filters
func (s *SubscriptionService) AggregateTotalPrice(req *models.AggregateRequest) (*models.AggregateResponse, error) {
	logrus.Info("Aggregating total subscription price")

	if req.UserID != nil && *req.UserID != "" {
		_, err := uuid.Parse(*req.UserID)
		if err != nil {
			logrus.Warnf("Invalid user ID format in aggregate request: %v", err)
			return nil, errors.New("invalid user ID format")
		}
	}

	if req.StartDate != nil && *req.StartDate != "" {
		if err := validateDateFormat(*req.StartDate); err != nil {
			logrus.Warnf("Invalid start_date in aggregate request: %v", err)
			return nil, err
		}
	}

	if req.EndDate != nil && *req.EndDate != "" {
		if err := validateDateFormat(*req.EndDate); err != nil {
			logrus.Warnf("Invalid end_date in aggregate request: %v", err)
			return nil, err
		}
	}

	totalPrice, err := s.repo.AggregateTotalPrice(
		req.UserID,
		req.ServiceName,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		logrus.Errorf("Failed to aggregate total price: %v", err)
		return nil, err
	}

	response := &models.AggregateResponse{
		TotalPrice: totalPrice,
	}

	logrus.Infof("Aggregated total price: %d", totalPrice)
	return response, nil
}

// toSubscriptionResponse converts a Subscription model to a SubscriptionResponse
func toSubscriptionResponse(sub *models.Subscription) *models.SubscriptionResponse {
	return &models.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}
