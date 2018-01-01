package domain

import "time"

// Balance represents a blanace at a given date
type Balance struct {
	Amount           Amount
	TransmissionDate time.Time
	TransmissionTime *time.Time
}
