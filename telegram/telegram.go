package telegram

import (
	"time"
	"errors"
	"log"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TelegramApi interface {
	SendMessage(chat_id int64, message string) error
}

type TelegramApiClient struct {
	baseUrl string
	token   string
	proxy   *url.URL
}

func CreateTelegramApiClient(baseUrl string, token string, proxy string) (result TelegramApi, err error) {
	if baseUrl == "" {
		err = errors.New("baseUrl cannot be empty")
		return
	}

	if token == "" {
		err = errors.New("token cannot be null")
		return
	}

	var proxyUrl *url.URL = nil
	if proxy != "" {
		log.Printf("Using proxy: %s", proxy)
		proxyUrl, err = url.Parse(proxy)
		if err != nil {
			return
		}
	}

	result = &TelegramApiClient{baseUrl: baseUrl, token: token, proxy: proxyUrl}
	return
}

// Sends message to telegram
func (this TelegramApiClient) SendMessage(chat_id int64, message string) error {
	urlTemplate, err := url.Parse(fmt.Sprintf("%s/bot%s/sendMessage", this.baseUrl, this.token))
	if err != nil {
		return err
	}

	var tr *http.Transport
	if this.proxy != nil  {
		log.Printf("Using proxy %s", this.proxy)
		tr = &http.Transport{
			Proxy: http.ProxyURL(this.proxy),
			DisableKeepAlives: false, 
		}
	} else {
		tr = nil
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: tr,
	}

	q := urlTemplate.Query()
	q.Set("chat_id", strconv.FormatInt(chat_id, 10))
	q.Set("text", message)
	urlTemplate.RawQuery = q.Encode()

	log.Printf("Sending message: %s", urlTemplate.String())
	// TODO: parse response
	_, err = client.Get(urlTemplate.String())
	if err != nil {
		return err
	} else {
		return nil
	}
}
