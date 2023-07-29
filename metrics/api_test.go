package metrics

import (
	"encoding/json"
	"log"
	"testing"
)

type IamToken string

func (this IamToken) SignedToken() (string, error) {
	return string(this), nil
}

func TestApi(t *testing.T) {
	metric := CreateSimpleMetric("name", DGAUGE, 1)

	buf, _ := json.Marshal(metric)

	log.Print(string(buf))
}
