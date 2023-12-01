package collector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDateTimestamp(t *testing.T) {
	order := Order{
		ID:           0,
		Certificate:  Certificate{ID: 0, CommonName: "example.com", ValidUntil: "2023-11-23"},
		Organization: Organization{Name: "My Org."},
	}

	expectedTimestamp, _ := time.Parse(time.RFC3339, "2023-11-23T00:00:00+00:00")
	location, _ := time.LoadLocation("Europe/Paris")

	timestamp := order.FormatDateTimestamp()

	assert.Equal(t, expectedTimestamp.In(location), timestamp)
}
