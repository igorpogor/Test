package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"subscription-service/internal/models"
)

// SubscriptionRepository handles database operations for subscriptions
type SubscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new SubscriptionRepository
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create inserts a new subscription into the database
func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	logrus.Debug("Creating subscription in database")
	return r.db.Create(subscription).Error
}

// GetByID retrieves a subscription by its ID
func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	logrus.Debugf("Fetching subscription by ID: %s", id)
	var subscription models.Subscription
	err := r.db.First(&subscription, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}

// GetAll retrieves all subscriptions from the database
func (r *SubscriptionRepository) GetAll() ([]models.Subscription, error) {
	logrus.Debug("Fetching all subscriptions from database")
	var subscriptions []models.Subscription
	err := r.db.Find(&subscriptions).Error
	return subscriptions, err
}

// Update saves updated subscription data to the database
func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	logrus.Debugf("Updating subscription %s in database", subscription.ID)
	return r.db.Save(subscription).Error
}

// Delete removes a subscription from the database by ID
func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	logrus.Debugf("Deleting subscription %s from database", id)
	result := r.db.Delete(&models.Subscription{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AggregateTotalPrice calculates the total price of subscriptions matching the given filters
func (r *SubscriptionRepository) AggregateTotalPrice(userID *string, serviceName *string, startDate *string, endDate *string) (int, error) {
	logrus.Debug("Aggregating total price with filters")
	query := r.db.Model(&models.Subscription{})

	if userID != nil && *userID != "" {
		logrus.Debugf("Filtering by user_id: %s", *userID)
		query = query.Where("user_id = ?", *userID)
	}
	if serviceName != nil && *serviceName != "" {
		logrus.Debugf("Filtering by service_name: %s", *serviceName)
		query = query.Where("service_name = ?", *serviceName)
	}

	if startDate != nil && *startDate != "" && endDate != nil && *endDate != "" {
		logrus.Debugf("Filtering by period: %s - %s", *startDate, *endDate)
		query = query.Where(
			"(end_date IS NULL AND start_date <= ?) OR "+
				"(end_date IS NOT NULL AND start_date <= ? AND end_date >= ?)",
			*endDate, *endDate, *startDate)
	} else if startDate != nil && *startDate != "" {
		logrus.Debugf("Filtering by start_date >= %s", *startDate)
		query = query.Where(
			"(end_date IS NULL AND start_date >= ?) OR "+
				"(end_date IS NOT NULL AND end_date >= ?)",
			*startDate, *startDate)
	} else if endDate != nil && *endDate != "" {
		logrus.Debugf("Filtering by end_date <= %s", *endDate)
		query = query.Where("start_date <= ?", *endDate)
	}

	var totalPrice sql.NullInt64
	err := query.Select("COALESCE(SUM(price), 0)").Row().Scan(&totalPrice)
	if err != nil {
		logrus.Errorf("Failed to aggregate total price: %v", err)
		return 0, err
	}

	if totalPrice.Valid {
		return int(totalPrice.Int64), nil
	}
	return 0, nil
}
