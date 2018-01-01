package domain

// AccountLimit represents a limit for an account with a possible timespan
type AccountLimit struct {
	Kind   string
	Amount Amount
	Days   int
}
