package domain

type AccountLimit struct {
	Kind   string
	Amount Amount
	Days   int
}
