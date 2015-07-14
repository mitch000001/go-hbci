package domain

type AccountInformation struct {
	AccountConnection           *AccountConnection
	UserID                      string
	Currency                    string
	Name1                       string
	Name2                       string
	ProductID                   string
	Limit                       *AccountLimit
	AllowedBusinessTransactions []BusinessTransaction
}
