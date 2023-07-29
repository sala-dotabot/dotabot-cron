package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
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
	url.RawQuery = query.Encode()

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

	if err != nil {
		return
	}
	if resp.StatusCode == 200 {
		var respStruct Response
		err = json.Unmarshal(respBuffer.Bytes(), &respStruct)
		if err != nil {
			return
		}
		result = respStruct.Write
	} else {
		var respStruct ErrorResponse
		err = json.Unmarshal(respBuffer.Bytes(), &respStruct)
		if err != nil {
			return
		}
		err = errors.New("Status Code: " + strconv.FormatInt(int64(resp.StatusCode), 10) + ", message: " + respStruct.Message)
	}
	return
}
