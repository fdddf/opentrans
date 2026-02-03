package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/fdddf/opentrans/internal/database"
	"gorm.io/gorm"
)

// SubscriptionService handles subscription-related operations
type SubscriptionService struct {
	DB *database.Database
}

// Subscription tiers with their limits
const (
	FreeTier    = "free"
	BasicTier   = "basic"
	PremiumTier = "premium"
)

// Default limits per subscription tier
var SubscriptionLimits = map[string]struct {
	MaxApps         int
	MaxTranslations int
}{
	FreeTier:    {MaxApps: 1, MaxTranslations: 1000},
	BasicTier:   {MaxApps: 5, MaxTranslations: 10000},
	PremiumTier: {MaxApps: 20, MaxTranslations: 50000},
}

// CreateSubscription creates a new subscription for a user
func (s *SubscriptionService) CreateSubscription(userID uint, stripeSubscriptionID, stripeCustomerID, subscriptionType string, currentPeriodStart, currentPeriodEnd time.Time, trialEnd *time.Time, cancelAtPeriodEnd bool) (*database.Subscription, error) {
	// Check if user already has an active subscription
	var existingSubscription database.Subscription
	result := s.DB.Where("user_id = ? AND subscription_status != ?", userID, "canceled").First(&existingSubscription)
	if result.Error == nil {
		return nil, errors.New("user already has an active subscription")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing subscription: %v", result.Error)
	}

	subscription := &database.Subscription{
		UserID:               userID,
		StripeSubscriptionID: stripeSubscriptionID,
		StripeCustomerID:     stripeCustomerID,
		SubscriptionType:     subscriptionType,
		SubscriptionStatus:   "active", // Default to active
		CurrentPeriodStart:   currentPeriodStart,
		CurrentPeriodEnd:     currentPeriodEnd,
		TrialEnd:             trialEnd,
		CancelAtPeriodEnd:    cancelAtPeriodEnd,
	}

	result = s.DB.Create(subscription)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create subscription: %v", result.Error)
	}

	// Update user's subscription info
	limits := SubscriptionLimits[subscriptionType]
	err := s.DB.Model(&database.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_subscribed":     true,
		"subscription_type": subscriptionType,
		"max_apps":          limits.MaxApps,
		"max_translations":  limits.MaxTranslations,
		"current_usage":     0, // Reset usage to 0
		"subscription_end":  &currentPeriodEnd,
	}).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update user subscription info: %v", err)
	}

	return subscription, nil
}

// GetSubscription retrieves a subscription by ID
func (s *SubscriptionService) GetSubscription(subscriptionID uint) (*database.Subscription, error) {
	var subscription database.Subscription
	result := s.DB.Preload("User").First(&subscription, subscriptionID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("subscription not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &subscription, nil
}

// GetSubscriptionByUser retrieves a user's active subscription
func (s *SubscriptionService) GetSubscriptionByUser(userID uint) (*database.Subscription, error) {
	var subscription database.Subscription
	result := s.DB.Where("user_id = ? AND subscription_status != ?", userID, "canceled").First(&subscription)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("no active subscription found for user")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &subscription, nil
}

// GetSubscriptionByStripeID retrieves a subscription by Stripe subscription ID
func (s *SubscriptionService) GetSubscriptionByStripeID(stripeSubscriptionID string) (*database.Subscription, error) {
	var subscription database.Subscription
	result := s.DB.Preload("User").Where("stripe_subscription_id = ?", stripeSubscriptionID).First(&subscription)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("subscription not found")
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &subscription, nil
}

// UpdateSubscription updates an existing subscription
func (s *SubscriptionService) UpdateSubscription(subscriptionID uint, updates map[string]interface{}) error {
	result := s.DB.Model(&database.Subscription{}).Where("id = ?", subscriptionID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update subscription: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// UpdateSubscriptionByStripeID updates a subscription by Stripe subscription ID
func (s *SubscriptionService) UpdateSubscriptionByStripeID(stripeSubscriptionID string, updates map[string]interface{}) error {
	result := s.DB.Model(&database.Subscription{}).Where("stripe_subscription_id = ?", stripeSubscriptionID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update subscription: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// CancelSubscription cancels a subscription
func (s *SubscriptionService) CancelSubscription(subscriptionID uint) error {
	subscription, err := s.GetSubscription(subscriptionID)
	if err != nil {
		return err
	}

	result := s.DB.Model(&database.Subscription{}).Where("id = ?", subscriptionID).Updates(map[string]interface{}{
		"subscription_status": "canceled",
	})
	if result.Error != nil {
		return fmt.Errorf("failed to cancel subscription: %v", result.Error)
	}

	// If it was the active subscription, update user's subscription info
	if subscription.SubscriptionStatus != "canceled" {
		// Reset to free tier
		err = s.DB.Model(&database.User{}).Where("id = ?", subscription.UserID).Updates(map[string]interface{}{
			"is_subscribed":     false,
			"subscription_type": FreeTier,
			"max_apps":          SubscriptionLimits[FreeTier].MaxApps,
			"max_translations":  SubscriptionLimits[FreeTier].MaxTranslations,
			"subscription_end":  nil,
		}).Error

		if err != nil {
			return fmt.Errorf("failed to update user subscription info: %v", err)
		}
	}

	return nil
}

// CheckUserUsage checks if a user has exceeded their translation limit
func (s *SubscriptionService) CheckUserUsage(userID uint) (bool, int, int, error) {
	var user database.User
	result := s.DB.First(&user, userID)
	if result.Error != nil {
		return false, 0, 0, fmt.Errorf("failed to get user: %v", result.Error)
	}

	return user.CurrentUsage >= user.MaxTranslations, user.CurrentUsage, user.MaxTranslations, nil
}

// RequireActiveSubscriptionAndQuota ensures user has active subscription/quota before proceeding
func (s *SubscriptionService) RequireActiveSubscriptionAndQuota(userID uint) error {
	over, usage, limit, err := s.CheckUserUsage(userID)
	if err != nil {
		return err
	}
	if over {
		return fmt.Errorf("usage exceeded: %d/%d", usage, limit)
	}
	return nil
}

// ResetMonthlyUsage resets the monthly usage for all users (typically called by a cron job)
func (s *SubscriptionService) ResetMonthlyUsage() error {
	err := s.DB.Model(&database.User{}).UpdateColumn("current_usage", 0).Error
	if err != nil {
		return fmt.Errorf("failed to reset monthly usage: %v", err)
	}

	return nil
}

// GetUserSubscriptionInfo returns detailed subscription information for a user
func (s *SubscriptionService) GetUserSubscriptionInfo(userID uint) (*database.User, error) {
	var user database.User
	result := s.DB.First(&user, userID)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %v", result.Error)
	}

	return &user, nil
}

// ProcessSubscriptionWebhook processes a subscription update from Stripe webhook
func (s *SubscriptionService) ProcessSubscriptionWebhook(stripeSubscriptionID, status, subscriptionType string, currentPeriodStart, currentPeriodEnd time.Time, cancelAtPeriodEnd bool) error {
	subscription, err := s.GetSubscriptionByStripeID(stripeSubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %v", err)
	}

	// Update the subscription
	updates := map[string]interface{}{
		"SubscriptionStatus": status,
		"SubscriptionType":   subscriptionType,
		"CurrentPeriodStart": currentPeriodStart,
		"CurrentPeriodEnd":   currentPeriodEnd,
		"CancelAtPeriodEnd":  cancelAtPeriodEnd,
	}

	err = s.UpdateSubscription(subscription.ID, updates)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %v", err)
	}

	// Update user's subscription info based on the new subscription type
	limits := SubscriptionLimits[subscriptionType]
	userUpdates := map[string]interface{}{
		"subscription_type": subscriptionType,
		"max_apps":          limits.MaxApps,
		"max_translations":  limits.MaxTranslations,
		"is_subscribed":     status == "active",
		"subscription_end":  &currentPeriodEnd,
	}

	err = s.DB.Model(&database.User{}).Where("id = ?", subscription.UserID).Updates(userUpdates).Error
	if err != nil {
		return fmt.Errorf("failed to update user subscription info: %v", err)
	}

	return nil
}

// Global functions for backward compatibility
var subscriptionServiceInstance *SubscriptionService

func SetSubscriptionService(db *database.Database) {
	subscriptionServiceInstance = &SubscriptionService{DB: db}
}

func CreateSubscription(userID uint, stripeSubscriptionID, stripeCustomerID, subscriptionType string, currentPeriodStart, currentPeriodEnd time.Time, trialEnd *time.Time, cancelAtPeriodEnd bool) (*database.Subscription, error) {
	return subscriptionServiceInstance.CreateSubscription(userID, stripeSubscriptionID, stripeCustomerID, subscriptionType, currentPeriodStart, currentPeriodEnd, trialEnd, cancelAtPeriodEnd)
}

func GetSubscription(subscriptionID uint) (*database.Subscription, error) {
	return subscriptionServiceInstance.GetSubscription(subscriptionID)
}

func GetSubscriptionByUser(userID uint) (*database.Subscription, error) {
	return subscriptionServiceInstance.GetSubscriptionByUser(userID)
}

func GetSubscriptionByStripeID(stripeSubscriptionID string) (*database.Subscription, error) {
	return subscriptionServiceInstance.GetSubscriptionByStripeID(stripeSubscriptionID)
}

func UpdateSubscription(subscriptionID uint, updates map[string]interface{}) error {
	return subscriptionServiceInstance.UpdateSubscription(subscriptionID, updates)
}

func UpdateSubscriptionByStripeID(stripeSubscriptionID string, updates map[string]interface{}) error {
	return subscriptionServiceInstance.UpdateSubscriptionByStripeID(stripeSubscriptionID, updates)
}

func CancelSubscription(subscriptionID uint) error {
	return subscriptionServiceInstance.CancelSubscription(subscriptionID)
}

func CheckUserUsage(userID uint) (bool, int, int, error) {
	return subscriptionServiceInstance.CheckUserUsage(userID)
}

func RequireActiveSubscriptionAndQuota(userID uint) error {
	return subscriptionServiceInstance.RequireActiveSubscriptionAndQuota(userID)
}

func ResetMonthlyUsage() error {
	return subscriptionServiceInstance.ResetMonthlyUsage()
}

func GetUserSubscriptionInfo(userID uint) (*database.User, error) {
	return subscriptionServiceInstance.GetUserSubscriptionInfo(userID)
}

func ProcessSubscriptionWebhook(stripeSubscriptionID, status, subscriptionType string, currentPeriodStart, currentPeriodEnd time.Time, cancelAtPeriodEnd bool) error {
	return subscriptionServiceInstance.ProcessSubscriptionWebhook(stripeSubscriptionID, status, subscriptionType, currentPeriodStart, currentPeriodEnd, cancelAtPeriodEnd)
}
