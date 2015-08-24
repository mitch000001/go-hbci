package domain

type AccountConnection struct {
	AccountID                 string
	SubAccountCharacteristics string
	CountryCode               int
	BankID                    string
}

type InternationalAccountConnection struct {
	IBAN                      string
	BIC                       string
	AccountID                 string
	SubAccountCharacteristics string
	BankID                    BankId
}
