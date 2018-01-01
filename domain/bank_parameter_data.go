package domain

// BankParameterData represent metadata prvided by a bank institute that
// reflect limitations and limits when talking to that institute
type BankParameterData struct {
	Version                    int
	BankID                     BankID
	BankName                   string
	MaxTransactionsPerMessage  int
	MaxMessageSize             int
	MinTimeout                 int
	MaxTimeout                 int
	PinTanBusinessTransactions map[string]bool
}

// PinTanBusinessTransaction provides information about whether a given Segment
// needs a TAN or not.
type PinTanBusinessTransaction struct {
	SegmentID string
	NeedsTan  bool
}
