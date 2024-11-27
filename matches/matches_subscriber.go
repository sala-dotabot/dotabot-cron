package matches

import (
	"fmt"
	"log"
	"strconv"

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
		matchDetailsResult, err := this.dotaApi.GetMatchDetails(uint64(matchId))
		if err != nil {
			return err
		}

		message := GetMessage(matchDetailsResult.Result, accountId)

		err = this.telegramApi.SendMessage(subscription.ChatId, message)
		if err != nil {
			return err
		}

		this.subscriptionRepository.SaveLastKnownMatchId(subscription, uint64(matchId))
	}

	return nil
}

const SHORT_FORMAT = "New match: https://www.dotabuff.com/matches/%d"
const LONG_FORMAT = `New match: https://www.dotabuff.com/matches/%d
Player (%d) %s. KDA: %d/%d/%d
`

func GetMessage(match dota.MatchDetails, accountIdStr string) string {
	accountId, _ := strconv.ParseInt(accountIdStr, 10, 64)

	if len(match.Players) == 0 {
		return fmt.Sprintf(SHORT_FORMAT, match.MatchId)
	}
	thePlayer := match.Players[0]
	players := match.Players
	for _, player := range players {
		if player.AccountId == accountId {
			thePlayer = player
		}
	}

	isPlayerRadiant := thePlayer.TeamNumber == 0
	radiantWin := match.RadiantWin

	var winString string
	if isPlayerRadiant && radiantWin || !isPlayerRadiant && !radiantWin {
		winString = "WIN"
	} else {
		winString = "LOST"
	}
	return fmt.Sprintf(LONG_FORMAT, match.MatchId, accountId, winString, thePlayer.Kills, thePlayer.Deaths, thePlayer.Assists)
}
