package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "digicert"

type DigicertCollector struct {
	digicertAPIEndpoint     string
	digicertAPIKey          string
	showExpiredCertificates bool
	up                      *prometheus.Desc
	scrapeDuration          *prometheus.Desc
	certificateExpire       *prometheus.Desc
	logger                  log.Logger
}

func (c *DigicertCollector) Collect(ch chan<- prometheus.Metric) {
	c.HitDigicertAPIAndUpdateMetrics(ch)
}

func (c *DigicertCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
	ch <- c.certificateExpire
}

func NewDigicertCollector(logger log.Logger,
	digicertURL string,
	digicertAPIKey string,
	digicertShowExpiredCertificates bool) (*DigicertCollector, error) {

	// Ping Digicert API
	var client = http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("HEAD", digicertURL, nil)
	req.Header.Add("X-DC-DEVKEY", digicertAPIKey)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, fmt.Errorf("failed to create http.NewRequest: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to reach Digicert API: %s, status code %d",
			err,
			resp.StatusCode,
		)
	}
	resp.Body.Close()

	level.Info(logger).Log("msg", "Digicert API is reachable")

	// Build the collector
	c := &DigicertCollector{
		digicertAPIEndpoint:     digicertURL,
		digicertAPIKey:          digicertAPIKey,
		showExpiredCertificates: digicertShowExpiredCertificates,
		logger:                  logger,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "api", "up"),
			"Was the last Digicert API query successful.",
			nil, nil,
		),
		scrapeDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape_duration", "seconds"),
			"Exporter scrape duration in seconds.",
			nil, nil,
		),
		certificateExpire: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "certificate_expire", "timestamp_seconds"),
			"Certificate expiration date.",
			[]string{"order_id", "certificate_id", "certificate_common_name", "organization"}, nil,
		),
	}

	return c, nil
}

func (c *DigicertCollector) HitDigicertAPIAndUpdateMetrics(ch chan<- prometheus.Metric) {

	start := time.Now()

	var orderList OrderList
	// Load channel stats
	req, err := http.NewRequest("GET", c.digicertAPIEndpoint, nil)
	if err != nil {
		level.Error(c.logger).Log("msg", err)
	}

	// This one line implements the authentication required for the task.
	req.Header.Add("X-DC-DEVKEY", c.digicertAPIKey)
	req.Header.Add("Content-Type", "application/json")
	// Make request and show output.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			c.up, prometheus.GaugeValue, 0,
		)
		level.Error(c.logger).Log("msg", err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		level.Error(c.logger).Log("msg", err)
	}

	err = json.Unmarshal(body, &orderList)
	if err != nil {
		level.Error(c.logger).Log("msg", err)
	}

	seenCertificationCommonName := make(map[string]Order)

	for i := 0; i < len(orderList.Orders); i++ {
		certificateCommonName := orderList.Orders[i].Certificate.CommonName
		certificateExpireDate := formatDateTimestamp(orderList.Orders[i].Certificate.ValidUntil)

		// A valid date must be in the future, or show all if showExpiredCertificates = true
		if certificateExpireDate.After(time.Now()) || c.showExpiredCertificates {
			// Test if the exporter already encounter this cert common name
			if formatDateTimestamp(
				seenCertificationCommonName[certificateCommonName].Certificate.ValidUntil,
			).IsZero() {
				// If no, insert into the map
				seenCertificationCommonName[certificateCommonName] = orderList.Orders[i]
			} else {
				// If yes AND the new date is after the current one, replace it
				if certificateExpireDate.After(formatDateTimestamp(seenCertificationCommonName[certificateCommonName].Certificate.ValidUntil)) {
					seenCertificationCommonName[certificateCommonName] = orderList.Orders[i]
				}
			}
		}
	}

	for name, order := range seenCertificationCommonName {
		if err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.certificateExpire,
				prometheus.UntypedValue,
				float64(formatDateTimestamp(order.Certificate.ValidUntil).Unix()),
				strconv.Itoa(order.ID),
				strconv.Itoa(order.Certificate.ID),
				name,
				order.Organization.Name,
			)
		}
	}

	end := time.Now()

	ch <- prometheus.MustNewConstMetric(
		c.scrapeDuration, prometheus.GaugeValue, end.Sub(start).Seconds(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.up, prometheus.GaugeValue, 1,
	)
}
