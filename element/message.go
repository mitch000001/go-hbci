package element

func NewReferenceMessage(dialogId string, messageNumber int) *ReferenceMessage {
	r := &ReferenceMessage{
		DialogID:      NewIdentification(dialogId),
		MessageNumber: NewNumber(messageNumber, 4),
	}
	r.DataElement = NewDataElementGroup(ReferenceMessageDEG, 2, r)
	return r
}

type ReferenceMessage struct {
	DataElement
	DialogID      *IdentificationDataElement
	MessageNumber *NumberDataElement
}

func (r *ReferenceMessage) IsValid() bool {
	if r.DialogID == nil || r.MessageNumber == nil {
		return false
	} else {
		return r.DataElement.IsValid()
	}
}

func (r *ReferenceMessage) Value() interface{} {
	return r
}

func (r *ReferenceMessage) GroupDataElements() []DataElement {
	return []DataElement{
		r.DialogID,
		r.MessageNumber,
	}
}
