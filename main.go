package main

import (
	"log"
)

func main() {
	context, err := InitContext()
	if err != nil {
		log.Fatalf("Error while initializing context: %s", err)
	}

	subscriber := context.MatchSubscriber
	repository := context.SubscriptionRepository

	subscriptions, err := repository.FindAll()
	if err != nil {
		log.Print(err)
	}
	log.Printf("Found %d subscriptions", len(subscriptions))

	for _, subscription := range subscriptions {
		log.Printf("Processing subscription: %d", subscription.ChatId)
		err := subscriber.ProcessSubscription(subscription)
		if err != nil {
			log.Print(err)
		}
	}
}
