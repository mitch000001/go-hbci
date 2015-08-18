package domain

type BankParameterData struct {
	Version                    int
	BankID                     BankId
	BankName                   string
	MaxTransactionsPerMessage  int
	MaxMessageSize             int
	MinTimeout                 int
	MaxTimeout                 int
	PinTanBusinessTransactions map[string]bool
}

type PinTanBusinessTransaction struct {
	SegmentID string
	NeedsTan  bool
}
