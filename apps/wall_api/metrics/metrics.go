package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProcessedMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processed_messages",
		Help: "A counter for the number of messages received",
	})
	WSClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_clients",
		Help: "A Gauge for the number of connected WebSocket clients",
	})
)
