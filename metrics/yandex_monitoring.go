package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type YandexMonitoringClient interface {
	Write(timestamp time.Time, labels map[string]string, metrics []Metric) (result int64, err error)
}

type YandexMonitoringClientImpl struct {
	baseUrl    string
	folderId   string
	iamContext IamContext
}

func MakeYandexMonitoringClientImpl(baseUrl string, folderId string, iamContext IamContext) YandexMonitoringClient {
	return &YandexMonitoringClientImpl{
		baseUrl:    baseUrl,
		folderId:   folderId,
		iamContext: iamContext,
	}
}

func (this *YandexMonitoringClientImpl) Write(timestamp time.Time, labels map[string]string, metrics []Metric) (result int64, err error) {
	url, err := url.Parse(this.baseUrl + "/data/write")
	if err != nil {
		return
	}

	query := url.Query()
	query.Set("folderId", this.folderId)
	query.Set("service", "custom")

	payload := Payload{
		Timestamp: timestamp.Format(time.RFC3339),
		Labels:    labels,
		Metrics:   metrics,
	}

	req, err := json.Marshal(payload)
	if err != nil {
		return
	}

	body := bytes.NewBuffer(req)

	request, err := http.NewRequest("POST", url.String(), body)
	if err != nil {
		return
	}
	iamToken, err := this.iamContext.SignedToken()
	if err != nil {
		return
	}
	request.Header.Add("Authorization", "Bearer "+iamToken)
	request.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	respBuffer := bytes.Buffer{}
	respBuffer.ReadFrom(resp.Body)

	var respStruct Response
	json.Unmarshal(respBuffer.Bytes(), &respStruct)
	if resp.StatusCode == 200 {
		result = respStruct.Write
	} else {
		err = errors.New(respStruct.ErrorMessage)
	}
	return
}
