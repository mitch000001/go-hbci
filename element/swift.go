package element

import (
	"github.com/mitch000001/go-hbci/swift"
)

// SwiftMT940DataElement represents a DataElement containing SWIFT MT940
// binary data
type SwiftMT940DataElement struct {
	*BinaryDataElement
	swiftMT940Messages *swift.MT940Messages
}

// UnmarshalHBCI unmarshals value into s
func (s *SwiftMT940DataElement) UnmarshalHBCI(value []byte) error {
	s.BinaryDataElement = &BinaryDataElement{}
	err := s.BinaryDataElement.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	s.swiftMT940Messages = swift.NewMT940Messages(s.BinaryDataElement.Val())
	return nil
}

// Val returns the embodied transactions as *swift.MT940Messages
func (s *SwiftMT940DataElement) Val() *swift.MT940Messages {
	return s.swiftMT940Messages
}
