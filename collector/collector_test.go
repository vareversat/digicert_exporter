package collector

import (
	"strings"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewDigicertCollectorWithoutMockedData(t *testing.T) {

	logger := log.NewNopLogger()
	digicertURL := "https://digicert.com"
	digicertAPIKey := "API_KEY"

	// Create a new DigicertCollector
	collector, err := NewDigicertCollector(logger, digicertURL, digicertAPIKey, false, true)

	assert.NoError(t, err)
	assert.False(t, collector.sandboxMode)
}

func TestNewDigicertCollectorWithMockedData(t *testing.T) {

	logger := log.NewNopLogger()
	digicertURL := ""
	digicertAPIKey := ""

	// Create a new DigicertCollector
	collector, err := NewDigicertCollector(logger, digicertURL, digicertAPIKey, true, true)

	assert.NoError(t, err)
	assert.True(t, collector.sandboxMode)
}

func TestUpdateMetricsApiDown(t *testing.T) {
	logger := log.NewNopLogger()
	digicertURL := "https://digicert.com"
	digicertAPIKey := "test"

	// Create a new DigicertCollector
	collector, err := NewDigicertCollector(logger, digicertURL, digicertAPIKey, false, true)
	assert.NoError(t, err)

	promChan := make(chan prometheus.Metric)

	go collector.UpdateMetrics(promChan)

	expected := strings.NewReader(
		`# HELP digicert_api_up Was the last Digicert API query successful.
# TYPE digicert_api_up gauge
digicert_api_up 0
`)

	err = testutil.CollectAndCompare(
		collector,
		expected,
		"digicert_api_up",
	)
}

func TestUpdateMetricsUse(t *testing.T) {
	logger := log.NewNopLogger()
	digicertURL := ""
	digicertAPIKey := ""

	// Create a new DigicertCollector
	collector, err := NewDigicertCollector(logger, digicertURL, digicertAPIKey, true, true)
	assert.NoError(t, err)
	assert.True(t, collector.sandboxMode)

	promChan := make(chan prometheus.Metric)

	go collector.UpdateMetrics(promChan)

	expected := strings.NewReader(
		`# HELP digicert_api_up Was the last Digicert API query successful.
# TYPE digicert_api_up gauge
digicert_api_up 0
`)

	err = testutil.CollectAndCompare(
		collector,
		expected,
		"digicert_api_up",
	)
}
