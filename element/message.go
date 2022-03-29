package element

import (
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
)

// NewReferencingMessage returns a new referencing message for the given
// dialogID and messageNumber
func NewReferencingMessage(dialogID string, messageNumber int) *ReferencingMessageDataElement {
	r := &ReferencingMessageDataElement{
		DialogID:      NewIdentification(dialogID),
		MessageNumber: NewNumber(messageNumber, 4),
	}
	r.DataElement = NewDataElementGroup(referenceMessageDEG, 2, r)
	return r
}

// ReferencingMessageDataElement represents a reference to another message for
// a given dialog
type ReferencingMessageDataElement struct {
	DataElement
	DialogID      *IdentificationDataElement
	MessageNumber *NumberDataElement
}

// Val returns the value of r as domain.ReferencingMessage
func (r *ReferencingMessageDataElement) Val() domain.MessageReference {
	return domain.MessageReference{
		DialogID:      r.DialogID.Val(),
		MessageNumber: r.MessageNumber.Val(),
	}
}

// IsValid returns true if the DialogID and MessageNumber are set and the
// underlying value is valid.
func (r *ReferencingMessageDataElement) IsValid() bool {
	if r.DialogID == nil || r.MessageNumber == nil {
		return false
	}
	return r.DataElement.IsValid()
}

// UnmarshalHBCI unmarshals value into r
func (r *ReferencingMessageDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if len(elements) != 2 {
		return fmt.Errorf("Malformed marshaled value")
	}
	dialogID := charset.ToUTF8(elements[0])
	num, err := strconv.Atoi(charset.ToUTF8(elements[1]))
	if err != nil {
		return fmt.Errorf("%T: Malformed message number: %v", r, err)
	}
	*r = *NewReferencingMessage(dialogID, num)
	return nil
}

// GroupDataElements returns the grouped DataElements
func (r *ReferencingMessageDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		r.DialogID,
		r.MessageNumber,
	}
}
