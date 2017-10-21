package dota

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type DotaApi interface {
	GetMatchesHistory(account_id string) (MatchesResult, error)
}

type DotaClient struct {
	baseUrl string
	key     string
}

func NewDotaClient(baseUrl string, key string) (client DotaApi, err error) {
	if baseUrl == "" {
		err = errors.New("baseUrl cannot be empty")
		return
	}

	if key == "" {
		err = errors.New("key cannot be empty")
		return
	}

	client = &DotaClient{baseUrl: baseUrl, key: key}
	return
}

func (this DotaClient) GetMatchesHistory(account_id string) (result MatchesResult, err error) {
	url, err := url.Parse(this.baseUrl + "/IDOTA2Match_570/GetMatchHistory/V001/")
	if err != nil {
		return
	}

	query := url.Query()
	query.Set("account_id", account_id)
	query.Set("key", this.key)
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return
	}

	b := bytes.Buffer{}
	b.ReadFrom(resp.Body)

	err = json.Unmarshal(b.Bytes(), &result)
	return
}
