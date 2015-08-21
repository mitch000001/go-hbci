package domain

type Amount struct {
	Amount   float64
	Currency string
}

type BusinessTransaction struct {
	ID               string
	NeededSignatures int
	Limit            *AccountLimit
}
