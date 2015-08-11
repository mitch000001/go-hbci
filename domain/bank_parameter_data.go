package domain

type BankParameterData struct {
	Version                    int
	BankID                     BankId
	BankName                   string
	MaxTransactionsPerMessage  int
	PinTanBusinessTransactions map[string]bool
}

type PinTanBusinessTransaction struct {
	SegmentID string
	NeedsTan  bool
}
