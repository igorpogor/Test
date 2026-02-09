package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"subscription-service/internal/models"
	"subscription-service/internal/repositories"
)

type SubscriptionService struct {
	repo *repositories.SubscriptionRepository
}

func NewSubscriptionService(repo *repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(req *models.SubscriptionRequest) (*models.SubscriptionResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		logrus.Errorf("Invalid user ID format: %v", err)
		return nil, errors.New("invalid user ID format")
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

	response := &models.SubscriptionResponse{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   subscription.StartDate,
		EndDate:     subscription.EndDate,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}

	logrus.Infof("Created subscription %s for user %s", subscription.ID, userID)
	return response, nil
}

func (s *SubscriptionService) GetSubscription(id uuid.UUID) (*models.SubscriptionResponse, error) {
	subscription, err := s.repo.GetByID(id)
	if err != nil {
		logrus.Errorf("Failed to get subscription %s: %v", id, err)
		return nil, err
	}

	if subscription == nil {
		logrus.Warnf("Subscription %s not found", id)
		return nil, nil
	}

	response := &models.SubscriptionResponse{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   subscription.StartDate,
		EndDate:     subscription.EndDate,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}

	logrus.Infof("Retrieved subscription %s", id)
	return response, nil
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.SubscriptionResponse, error) {
	subscriptions, err := s.repo.GetAll()
	if err != nil {
		logrus.Errorf("Failed to get all subscriptions: %v", err)
		return nil, err
	}

	var responses []models.SubscriptionResponse
	for _, sub := range subscriptions {
		responses = append(responses, models.SubscriptionResponse{
			ID:          sub.ID,
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserID:      sub.UserID,
			StartDate:   sub.StartDate,
			EndDate:     sub.EndDate,
			CreatedAt:   sub.CreatedAt,
			UpdatedAt:   sub.UpdatedAt,
		})
	}

	logrus.Infof("Retrieved %d subscriptions", len(responses))
	return responses, nil
}

func (s *SubscriptionService) UpdateSubscription(id uuid.UUID, req *models.SubscriptionRequest) (*models.SubscriptionResponse, error) {
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
		logrus.Errorf("Invalid user ID format: %v", err)
		return nil, errors.New("invalid user ID format")
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

	response := &models.SubscriptionResponse{
		ID:          existing.ID,
		ServiceName: existing.ServiceName,
		Price:       existing.Price,
		UserID:      existing.UserID,
		StartDate:   existing.StartDate,
		EndDate:     existing.EndDate,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   existing.UpdatedAt,
	}

	logrus.Infof("Updated subscription %s", id)
	return response, nil
}

func (s *SubscriptionService) DeleteSubscription(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		logrus.Errorf("Failed to delete subscription %s: %v", id, err)
		return err
	}

	logrus.Infof("Deleted subscription %s", id)
	return nil
}

func (s *SubscriptionService) AggregateTotalPrice(req *models.AggregateRequest) (*models.AggregateResponse, error) {
	var userIDStr *string
	if req.UserID != nil {
		_, err := uuid.Parse(*req.UserID)
		if err != nil {
			logrus.Errorf("Invalid user ID format in aggregate request: %v", err)
			return nil, errors.New("invalid user ID format")
		}
		userIDStr = req.UserID
	}

	totalPrice, err := s.repo.AggregateTotalPrice(
		userIDStr,
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
