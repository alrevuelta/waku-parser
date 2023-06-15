package metrics

import (
	"fmt"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	SentMessages = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "parser",
			Name:      "parser_sent_messages",
			Help:      "Amount of messages that were sent",
		},
		[]string{
			"container",
		},
	)

	ReceivedMessages = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "parser",
			Name:      "parser_received_messages",
			Help:      "Amount of messages that were received",
		},
		[]string{
			"container",
		},
	)

	AverageDelay = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "parser",
			Name:      "parser_average_delay_usec",
			Help:      "Average delay between message sent and received",
		},
		[]string{
			"container",
		},
	)
)

func RunMetrics(port int) {
	go func() {
		log.Info("Prometheus server started on port: ", port)
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}
