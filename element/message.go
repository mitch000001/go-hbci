package element

import (
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
)

func NewReferencingMessage(dialogId string, messageNumber int) *ReferencingMessageDataElement {
	r := &ReferencingMessageDataElement{
		DialogID:      NewIdentification(dialogId),
		MessageNumber: NewNumber(messageNumber, 4),
	}
	r.DataElement = NewDataElementGroup(ReferenceMessageDEG, 2, r)
	return r
}

type ReferencingMessageDataElement struct {
	DataElement
	DialogID      *IdentificationDataElement
	MessageNumber *NumberDataElement
}

func (r *ReferencingMessageDataElement) Val() domain.ReferencingMessage {
	return domain.ReferencingMessage{
		DialogID:      r.DialogID.Val(),
		MessageNumber: r.MessageNumber.Val(),
	}
}

func (r *ReferencingMessageDataElement) IsValid() bool {
	if r.DialogID == nil || r.MessageNumber == nil {
		return false
	} else {
		return r.DataElement.IsValid()
	}
}

func (r *ReferencingMessageDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if len(elements) != 2 {
		return fmt.Errorf("Malformed marshaled value")
	}
	dialogId := charset.ToUTF8(elements[0])
	num, err := strconv.Atoi(charset.ToUTF8(elements[1]))
	if err != nil {
		return fmt.Errorf("%T: Malformed message number: %v", r, err)
	}
	*r = *NewReferencingMessage(dialogId, num)
	return nil
}

func (r *ReferencingMessageDataElement) Value() interface{} {
	return r
}

func (r *ReferencingMessageDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		r.DialogID,
		r.MessageNumber,
	}
}
