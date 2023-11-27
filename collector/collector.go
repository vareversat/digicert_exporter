package collector

import (
	"github.com/go-kit/log/level"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "digicert"

type DigicertCollector struct {
	digicertAPIEndpoint     string
	digicertAPIKey          string
	showExpiredCertificates bool
	digicertMock            bool
	up                      *prometheus.Desc
	scrapeDuration          *prometheus.Desc
	certificateExpire       *prometheus.Desc
	logger                  log.Logger
}

func (c *DigicertCollector) Collect(ch chan<- prometheus.Metric) {
	c.UpdateMetrics(ch)
}

func (c *DigicertCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
	ch <- c.certificateExpire
}

func NewDigicertCollector(logger log.Logger,
	digicertURL string,
	digicertAPIKey string,
	digicertShowExpiredCertificates bool,
	digicertMock bool) (*DigicertCollector, error) {

	// Build the collector
	c := &DigicertCollector{
		digicertAPIEndpoint:     digicertURL,
		digicertAPIKey:          digicertAPIKey,
		showExpiredCertificates: digicertShowExpiredCertificates,
		digicertMock:            digicertMock,
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

	level.Info(logger).Log("msg", "Exporter started correctly")

	return c, nil
}

func (c *DigicertCollector) UpdateMetrics(ch chan<- prometheus.Metric) {

	start := time.Now()

	orderList, err := c.FetchDigicertData()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			c.up, prometheus.GaugeValue, 0,
		)
		return
	}

	seenCertificationCommonName := make(map[string]Order)

	for i := 0; i < len(orderList.Orders); i++ {
		order := orderList.Orders[i]
		certificateCommonName := order.Certificate.CommonName
		certificateExpireDate := order.FormatDateTimestamp()

		// A valid date must be in the future, or show all if showExpiredCertificates = true
		if certificateExpireDate.After(time.Now()) || c.showExpiredCertificates {
			seenOrder := seenCertificationCommonName[certificateCommonName]
			// Test if the collector already encounter this cert common name
			if seenOrder.FormatDateTimestamp().IsZero() {
				// If no, insert into the map
				seenCertificationCommonName[certificateCommonName] = orderList.Orders[i]
			} else {
				// If yes AND the new date is after the current one, replace it
				if certificateExpireDate.After(seenOrder.FormatDateTimestamp()) {
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
				float64(order.FormatDateTimestamp().Unix()),
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
