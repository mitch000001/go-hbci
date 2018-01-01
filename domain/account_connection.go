package domain

// AccountConnection represents an identification for a bank account
type AccountConnection struct {
	AccountID                 string
	SubAccountCharacteristics string
	CountryCode               int
	BankID                    string
}

// InternationalAccountConnection represents an international identification
// for a bank account
type InternationalAccountConnection struct {
	IBAN                      string
	BIC                       string
	AccountID                 string
	SubAccountCharacteristics string
	BankID                    BankID
}
