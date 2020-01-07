package element

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/swift"
	"time"
)

// SwiftMT940DataElement represents a DataElement containing SWIFT MT940
// binary data
type SwiftMT940DataElement struct {
	*BinaryDataElement
	swiftMT940Elements []*swift.MT940
}

// UnmarshalHBCI unmarshals value into s
func (s *SwiftMT940DataElement) UnmarshalHBCI(value []byte) error {
	s.BinaryDataElement = &BinaryDataElement{}
	err := s.BinaryDataElement.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	messageExtractor := swift.NewMessageExtractor(s.BinaryDataElement.Val())
	messages, err := messageExtractor.Extract()
	if err != nil {
		return err
	}
	for _, message := range messages {
		tr := &swift.MT940{}
		tr.ReferenceDate = time.Now()
		err = tr.Unmarshal(message)
		if err != nil {
			return err
		}
		s.swiftMT940Elements = append(s.swiftMT940Elements, tr)
	}
	return nil
}

// Val returns the embodied transactions as []domain.AccountTransaction
func (s *SwiftMT940DataElement) Val() []domain.AccountTransaction {
	var transactions []domain.AccountTransaction
	for _, mt940 := range s.swiftMT940Elements {
		transactions = append(transactions, mt940.AccountTransactions()...)
	}
	return transactions
}
