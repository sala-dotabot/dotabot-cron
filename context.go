package main

import (
	"dotabot-cron/dota"
	"dotabot-cron/matches"
	"dotabot-cron/repository"
	"dotabot-cron/telegram"
	"os"

	"github.com/go-redis/redis"
)

type Context struct {
	DotaApi                dota.DotaApi
	TelegramApi            telegram.TelegramApi
	MatchSubscriber        matches.MatchSubscriber
	SubscriptionRepository repository.SubscriptionRepository
}

func InitContext() (context *Context, err error) {
	dotaApiBaseUrl := getDotaApiBaseUrl()
	dotaApiToken := getDotaApiToken()

	telegramApiBaseUrl := getTelegramApiBaseUrl()
	telegramApiToken := getTelegramApiToken()

	dotaApi, err := dota.NewDotaClient(dotaApiBaseUrl, dotaApiToken)
	if err != nil {
		return
	}

	telegramApi, err := telegram.CreateTelegramApiClient(telegramApiBaseUrl, telegramApiToken)
	if err != nil {
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: "",
		DB:       0,
	})

	var subscriptionRepository repository.SubscriptionRepository = repository.CreateRedisRepository(client)

	matchSubscriber := matches.CreateMatchSubscriber(dotaApi, subscriptionRepository, telegramApi)
	if err != nil {
		return
	}

	context = &Context{
		DotaApi:                dotaApi,
		TelegramApi:            telegramApi,
		MatchSubscriber:        matchSubscriber,
		SubscriptionRepository: subscriptionRepository}
	return
}

func getDotaApiBaseUrl() string {
	dotaApiBaseUrl := os.Getenv("DOTA_API_BASE_URL")
	if dotaApiBaseUrl != "" {
		return dotaApiBaseUrl
	} else {
		return "https://api.steampowered.com"
	}
}

func getDotaApiToken() string {
	return os.Getenv("DOTA_API_TOKEN")
}

func getTelegramApiBaseUrl() string {
	telegramApiBaseUrl := os.Getenv("TELEGRAM_API_BASE_URL")
	if telegramApiBaseUrl != "" {
		return telegramApiBaseUrl
	} else {
		return "https://api.telegram.org"
	}
}

func getTelegramApiToken() string {
	return os.Getenv("TELEGRAM_API_TOKEN")
}

func getRedisAddr() string {
	return os.Getenv("REDIS_ADDR")
}
