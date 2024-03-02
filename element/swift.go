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

// SwiftMT942DataElement represents a DataElement containing SWIFT MT942
// binary data
type SwiftMT942DataElement struct {
	*BinaryDataElement
	swiftMT942Messages *swift.MT942Messages
}

// UnmarshalHBCI unmarshals value into s
func (s *SwiftMT942DataElement) UnmarshalHBCI(value []byte) error {
	s.BinaryDataElement = &BinaryDataElement{}
	err := s.BinaryDataElement.UnmarshalHBCI(value)
	if err != nil {
		return err
	}
	s.swiftMT942Messages = swift.NewMT942Messages(s.BinaryDataElement.Val())
	return nil
}

// Val returns the embodied transactions as *swift.MT942Messages
func (s *SwiftMT942DataElement) Val() *swift.MT942Messages {
	return s.swiftMT942Messages
}
