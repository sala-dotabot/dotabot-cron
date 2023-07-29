package main

import (
	"log"
	"time"

	"github.com/saladinkzn/dotabot-cron/metrics"
)

func main() {
	context, err := InitContext()
	if err != nil {
		log.Fatalf("Error while initializing context: %s", err)
	}

	subscriber := context.MatchSubscriber
	repository := context.SubscriptionRepository

	var metricList []metrics.Metric

	subscriptions, err := repository.FindAll()
	if err != nil {
		log.Print(err)
	}
	log.Printf("Found %d subscriptions", len(subscriptions))
	subscriptions_count_metric := metrics.CreateSimpleMetric("subscriptions_count", metrics.IGAUGE, float64(len(subscriptions)))
	metricList = append(metricList, subscriptions_count_metric)

	successfullyProcessed := 0
	errorProcessed := 0
	for _, subscription := range subscriptions {
		log.Printf("Processing subscription: %d", subscription.ChatId)
		err := subscriber.ProcessSubscription(subscription)
		if err != nil {
			errorProcessed += 1
			log.Print(err)
		} else {
			successfullyProcessed += 1
		}
	}

	metricList = append(metricList, metrics.CreateSimpleMetric("subscriptions_success_count", metrics.IGAUGE, float64(successfullyProcessed)))
	metricList = append(metricList, metrics.CreateSimpleMetric("subscriptions_failure_count", metrics.IGAUGE, float64(errorProcessed)))

	result, err := context.YandexMonitoringClient.Write(time.Now(), makeLabels(), metricList)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Successfully written %d metrics", result)
}

func makeLabels() map[string]string {
	labels := make(map[string]string)
	labels["app_name"] = "dotabot-cron"
	return labels
}
