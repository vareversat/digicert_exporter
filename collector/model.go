package collector

import (
	"fmt"
	"time"

	_ "time/tzdata"
)

type OrderList struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID           int          `json:"id"`
	Certificate  Certificate  `json:"certificate"`
	Organization Organization `json:"organization"`
}

type Certificate struct {
	ID         int    `json:"id"`
	CommonName string `json:"common_name"`
	ValidUntil string `json:"valid_till"`
}

type Organization struct {
	Name string `json:"name"`
}

func (o *Order) FormatDateTimestamp() time.Time {
	timestamp, _ := time.Parse(
		time.RFC3339,
		fmt.Sprintf("%sT00:00:00+00:00", o.Certificate.ValidUntil),
	)
	location, _ := time.LoadLocation("Europe/Paris")
	return timestamp.In(location)
}
