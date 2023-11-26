package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormatDateTimestamp(t *testing.T) {
	strDate := "2023-11-23"
	expectedTimestamp, _ := time.Parse(time.RFC3339, "2023-11-23T00:00:00+00:00")
	location, _ := time.LoadLocation("Europe/Paris")

	timestamp := formatDateTimestamp(strDate)

	assert.Equal(t, expectedTimestamp.In(location), timestamp)
}
