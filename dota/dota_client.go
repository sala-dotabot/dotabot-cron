package dota

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type DotaApi interface {
	GetMatchesHistory(account_id string) (MatchesResult, error)
	GetMatchDetails(match_id uint64) (MatchResult, error)
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
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func (this DotaClient) GetMatchDetails(match_id uint64) (result MatchResult, err error) {
	url, err := url.Parse(this.baseUrl + "/IDOTA2Match_570/GetMatchDetails/v001/")
	if err != nil {
		return
	}

	query := url.Query()
	query.Set("match_id", strconv.FormatUint(match_id, 10))
	query.Set("key", this.key)
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}
