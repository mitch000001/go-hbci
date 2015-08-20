package segment

import "github.com/mitch000001/go-hbci/element"

func NewEncryptedDataSegment(encryptedData []byte) *EncryptedDataSegment {
	e := &EncryptedDataSegment{
		Data: element.NewBinary(encryptedData, -1),
	}
	e.Segment = NewBasicSegment(999, e)
	return e
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment EncryptedDataSegment

type EncryptedDataSegment struct {
	Segment
	Data *element.BinaryDataElement
}

func (e *EncryptedDataSegment) Version() int         { return 1 }
func (e *EncryptedDataSegment) ID() string           { return "HNVSD" }
func (e *EncryptedDataSegment) referencedId() string { return "" }
func (e *EncryptedDataSegment) sender() string       { return senderBoth }

func (e *EncryptedDataSegment) elements() []element.DataElement {
	return []element.DataElement{
		e.Data,
	}
}
