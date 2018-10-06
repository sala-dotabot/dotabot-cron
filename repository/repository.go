package repository

type SubscriptionRepository interface {
	GetLastKnownMatchId(subscription TelegramMatchSubscription) (int64, error)
	SaveLastKnownMatchId(subscription TelegramMatchSubscription, matchId uint64) error
	RemoveLastKnownMatchId(subscription TelegramMatchSubscription) error
	FindAll() ([]TelegramMatchSubscription, error)
}