package domain

import "time"

type AccountTransaction struct {
	Account              AccountConnection
	Amount               Amount
	ValutaDate           time.Time
	BookingDate          time.Time
	BankID               string
	AccountID            string
	Purpose              string
	Purpose2             string
	AccountBalanceBefore Balance
	AccountBalanceAfter  Balance
}
