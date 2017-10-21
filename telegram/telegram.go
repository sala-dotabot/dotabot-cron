package telegram

import (
	"errors"
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
}

func CreateTelegramApiClient(baseUrl string, token string) (result TelegramApi, err error) {
	if baseUrl == "" {
		err = errors.New("baseUrl cannot be empty")
		return
	}

	if token == "" {
		err = errors.New("token cannot be null")
		return
	}

	result = &TelegramApiClient{baseUrl: baseUrl, token: token}
	return
}

// Отправляет сообщение в телеграм
func (this TelegramApiClient) SendMessage(chat_id int64, message string) error {
	urlTemplate, err := url.Parse(fmt.Sprintf("%s/bot%s/sendMessage", this.baseUrl, this.token))
	if err != nil {
		return err
	}

	q := urlTemplate.Query()
	q.Set("chat_id", strconv.FormatInt(chat_id, 10))
	q.Set("text", message)
	urlTemplate.RawQuery = q.Encode()

	// TODO: parse response
	_, err = http.Get(urlTemplate.String())
	if err != nil {
		return err
	} else {
		return nil
	}
}
