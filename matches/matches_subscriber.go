package matches

import (
	"fmt"
	"log"

	"github.com/saladinkzn/dotabot-cron/dota"
	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"
)

type MatchSubscriber interface {
	ProcessSubscription(subscription repository.TelegramMatchSubscription) error
}

type MatchSubscriberImpl struct {
	dotaApi                dota.DotaApi
	subscriptionRepository repository.SubscriptionRepository
	telegramApi            telegram.TelegramApi
}

func CreateMatchSubscriber(dotaApi dota.DotaApi, subscriptionRepository repository.SubscriptionRepository, telegramApi telegram.TelegramApi) MatchSubscriber {
	return &MatchSubscriberImpl{dotaApi: dotaApi,
		subscriptionRepository: subscriptionRepository,
		telegramApi:            telegramApi}
}

func (this MatchSubscriberImpl) ProcessSubscription(subscription repository.TelegramMatchSubscription) error {
	accountId := subscription.DotaAccountId

	wrapper, err := this.dotaApi.GetMatchesHistory(accountId)
	if err != nil {
		return err
	}
	log.Printf("Found matches: %d", len(wrapper.Result.Matches))

	var matches = wrapper.Result.Matches
	var matchId int64
	if len(matches) > 0 {
		matchId = matches[0].MatchId
	} else {
		matchId = -1
	}

	lastKnownId, err := this.subscriptionRepository.GetLastKnownMatchId(subscription)
	if err != nil {
		return err
	}
	log.Printf("matchId: %d, lastKnownId: %d", matchId, lastKnownId)

	if matchId > lastKnownId {
		log.Printf("Sending message for match %d", matchId)
		err := this.telegramApi.SendMessage(subscription.ChatId, fmt.Sprintf("New match: https://www.dotabuff.com/matches/%d", matchId))
		if err != nil {
			return err
		}

		this.subscriptionRepository.SaveLastKnownMatchId(subscription, uint64(matchId))
	}

	return nil
}
