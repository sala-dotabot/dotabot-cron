package main

import (
	"dotabot-cron/repository"
	"log"
)

func getSubscriptions() []repository.TelegramMatchSubscription {
	moja := repository.TelegramMatchSubscription{ChatId: 151904085, DotaAccountId: "70766996"}
	return []repository.TelegramMatchSubscription{moja}
}

func main() {
	context, err := InitContext()
	if err != nil {
		log.Fatal(err)
	}

	subscriber := context.MatchSubscriber
	repository := context.SubscriptionRepository

	for _, subscription := range getSubscriptions() {
		repository.SaveLastKnownMatchId(subscription, 0)
	}

	subscriptions := repository.FindAll()
	for _, subscription := range subscriptions {
		log.Printf("Processing subscription: %d", subscription.ChatId)
		err := subscriber.ProcessSubscription(subscription)
		if err != nil {
			log.Print(err)
		}
	}
}
