package domain

type BankParameterData struct {
	Version                   int
	BankID                    BankId
	BankName                  string
	MaxTransactionsPerMessage int
}
