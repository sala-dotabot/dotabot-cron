package main

import (
	"log"
	"os"

	"github.com/saladinkzn/dotabot-cron/dota"
	"github.com/saladinkzn/dotabot-cron/matches"
	"github.com/saladinkzn/dotabot-cron/metrics"
	"github.com/saladinkzn/dotabot-cron/repository"
	"github.com/saladinkzn/dotabot-cron/telegram"

	"github.com/go-redis/redis"
)

type Context struct {
	DotaApi                dota.DotaApi
	TelegramApi            telegram.TelegramApi
	MatchSubscriber        matches.MatchSubscriber
	SubscriptionRepository repository.SubscriptionRepository
	YandexMonitoringClient metrics.YandexMonitoringClient
}

func InitContext() (context *Context, err error) {
	dotaApiBaseUrl := getDotaApiBaseUrl()
	dotaApiToken := getDotaApiToken()
	log.Printf("Loaded dota settings %s, %s", dotaApiBaseUrl, dotaApiToken)

	telegramApiBaseUrl := getTelegramApiBaseUrl()
	telegramApiToken := getTelegramApiToken()
	telegramProxyUrl := getTelegramProxyUrl()
	log.Printf("Loaded telegram settings: %s, %s, %s", telegramApiBaseUrl, telegramApiToken, telegramProxyUrl)

	dotaApi, err := dota.NewDotaClient(dotaApiBaseUrl, dotaApiToken)
	if err != nil {
		return
	}

	telegramApi, err := telegram.CreateTelegramApiClient(telegramApiBaseUrl, telegramApiToken, telegramProxyUrl)
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

	monitoringClient := metrics.MakeYandexMonitoringClientImpl(getMonitoringUrl(), getFolderId(), getIamContext())

	context = &Context{
		DotaApi:                dotaApi,
		TelegramApi:            telegramApi,
		MatchSubscriber:        matchSubscriber,
		SubscriptionRepository: subscriptionRepository,
		YandexMonitoringClient: monitoringClient}
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

func getTelegramProxyUrl() string {
	return os.Getenv("TELEGRAM_PROXY_URL")
}

func getMonitoringUrl() string {
	return os.Getenv("MONITORING_URL")
}

func getFolderId() string {
	return os.Getenv("FOLDER_ID")
}

func getServiceAccountId() string {
	return os.Getenv("SERVICE_ACCOUNT_ID")
}

func getKeyId() string {
	return os.Getenv("KEY_ID")
}

func getKeyFileName() string {
	return os.Getenv("KEY_FILE_NAME")
}

func getIamContext() metrics.IamContext {
	return metrics.CreateIamContext(
		getKeyId(),
		getServiceAccountId(),
		getKeyFileName(),
	)
}
