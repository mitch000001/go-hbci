package domain

// BankParameterData represent metadata prvided by a bank institute that
// reflect limitations and limits when talking to that institute
type BankParameterData struct {
	Version                    int
	BankID                     BankID          `yaml:",inline"`
	BankName                   string          `yaml:"bankName"`
	MaxTransactionsPerMessage  int             `yaml:"maxTransactionsPerMessage"`
	MaxMessageSize             int             `yaml:"maxMessageSize"`
	MinTimeout                 int             `yaml:"minTimeout"`
	MaxTimeout                 int             `yaml:"maxTimeout"`
	PinTanBusinessTransactions map[string]bool `yaml:"pinTanBusinessTransactions"`
}

// PinTanBusinessTransaction provides information about whether a given Segment
// needs a TAN or not.
type PinTanBusinessTransaction struct {
	SegmentID string `yaml:"segmentID"`
	NeedsTan  bool   `yaml:"needsTan"`
}

type PinTanBusinessTransactions []PinTanBusinessTransaction

func (p PinTanBusinessTransactions) NeedsTan(segmentID string) (bool, bool) {
	var transaction *PinTanBusinessTransaction
	for _, tr := range p {
		if tr.SegmentID == segmentID {
			transaction = &tr
			break
		}
	}
	if transaction == nil {
		return false, false
	}
	return transaction.NeedsTan, true
}
