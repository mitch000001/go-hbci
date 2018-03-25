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

// ToAccountConnection transforms i into an AccountConnection
func (i InternationalAccountConnection) ToAccountConnection() AccountConnection {
	return AccountConnection{
		AccountID:                 i.AccountID,
		SubAccountCharacteristics: i.SubAccountCharacteristics,
		CountryCode:               i.BankID.CountryCode,
		BankID:                    i.BankID.ID,
	}
}
