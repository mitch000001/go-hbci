package domain

import "time"

type Balance struct {
	Amount           Amount
	TransmissionDate time.Time
	TransmissionTime *time.Time
}
