package metrics

import (
	"encoding/json"
	"log"
	"testing"
)

func TestApi(t *testing.T) {
	metric := CreateSimpleMetric("name", DGAUGE, 1)

	buf, _ := json.Marshal(metric)

	log.Print(string(buf))
}
