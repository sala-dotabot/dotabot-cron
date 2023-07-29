package metrics

import "time"

type MetricType string

const (
	DGAUGE  MetricType = "DGAUGE"
	IGAUGE  MetricType = "IGAUGE"
	COUNTER MetricType = "COUNTER"
	RATE    MetricType = "RATE"
)

type Point struct {
	Timestamp string  `json:"ts"`
	Value     float64 `json:"value"`
}

type Metric struct {
	Name       string
	Labels     map[string]string `json:"labels,omitempty"`
	MetricType string            `json:"type"`
	Timestamp  string            `json:"ts,omitempty"`
	Value      float64           `json:"value"`
	Timeseries []Point           `json:"timeseries,omitempty"`
}

type Payload struct {
	Timestamp string            `json:"ts"`
	Labels    map[string]string `json:"labels"`
	Metrics   []Metric          `json:"metrics"`
}

type Response struct {
	Write bool `json:"write"`
}

func CreatePayload(timestamp time.Time, labels map[string]string, metrics []Metric) Payload {
	return Payload{
		Timestamp: timestamp.Format(time.RFC3339),
		Labels:    labels,
		Metrics:   metrics,
	}
}

func CreatePoint(timestamp time.Time, value float64) Point {
	return Point{
		Timestamp: timestamp.Format(time.RFC3339),
		Value:     value,
	}
}

func CreateSimpleMetric(name string, metricType MetricType, value float64) Metric {
	return Metric{
		Name:       name,
		Labels:     make(map[string]string, 0),
		MetricType: string(metricType),
		Timestamp:  "",
		Value:      value,
		Timeseries: make([]Point, 0),
	}
}

func CreateMetric(name string, labels map[string]string, metricType MetricType, timestamp time.Time, value float64, points []Point) Metric {
	return Metric{
		Name:       name,
		Labels:     labels,
		MetricType: string(metricType),
		Timestamp:  timestamp.Format(time.RFC3339),
		Value:      value,
		Timeseries: points,
	}
}
