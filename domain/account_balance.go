package domain

import "time"

type AccountBalance struct {
	Account          AccountConnection
	ProductName      string
	Currency         string
	BookedBalance    Balance
	EarmarkedBalance *Balance
	CreditLimit      *Amount
	AvailableAmount  *Amount
	UsedAmount       *Amount
	BookingDate      *time.Time
	DueDate          *time.Time
}
