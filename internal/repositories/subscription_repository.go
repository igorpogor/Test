package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"subscription-service/internal/models"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.First(&subscription, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subscription, err
}

func (r *SubscriptionRepository) GetAll() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Subscription{}, "id = ?", id).Error
}

func (r *SubscriptionRepository) AggregateTotalPrice(userID *string, serviceName *string, startDate *string, endDate *string) (int, error) {
	query := r.db.Model(&models.Subscription{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if serviceName != nil {
		query = query.Where("service_name = ?", *serviceName)
	}

	if startDate != nil || endDate != nil {
		if startDate != nil && endDate != nil {
			query = query.Where(
				"(end_date IS NULL AND start_date <= ?) OR "+
					"(end_date IS NOT NULL AND start_date <= ? AND end_date >= ?)",
				*endDate, *endDate, *startDate)
		} else if startDate != nil {
			query = query.Where(
				"(end_date IS NULL AND start_date <= ?) OR "+
					"(end_date IS NOT NULL AND end_date >= ?)",
				*startDate, *startDate)
		} else if endDate != nil {
			query = query.Where(
				"(end_date IS NULL AND start_date <= ?) OR "+
					"(end_date IS NOT NULL AND start_date <= ?)",
				*endDate, *endDate)
		}
	}

	var totalPrice int
	err := query.Select("SUM(price)").Row().Scan(&totalPrice)
	if err != nil {
		return 0, err
	}

	return totalPrice, nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("01-2006", dateStr)
}

func isActiveDuringPeriod(subscription *models.Subscription, periodStart, periodEnd *time.Time) bool {
	startDate, err := parseDate(subscription.StartDate)
	if err != nil {
		return false
	}

	var endDate *time.Time
	if subscription.EndDate != nil {
		parsedEndDate, err := parseDate(*subscription.EndDate)
		if err != nil {
			return false
		}
		endDate = &parsedEndDate
	}

	if endDate == nil {
		if periodEnd == nil {
			return !startDate.After(*periodStart)
		}
		return !startDate.After(*periodEnd) && !startDate.Before(*periodStart)
	}

	if periodEnd == nil {
		return !startDate.After(*periodStart) && !endDate.Before(*periodStart)
	}

	return !startDate.After(*periodEnd) && !endDate.Before(*periodStart)
}
