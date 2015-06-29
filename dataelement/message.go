package dataelement

func NewReferenceMessage(dialogId string, messageNumber int) *ReferenceMessage {
	r := &ReferenceMessage{
		DialogID:      NewIdentificationDataElement(dialogId),
		MessageNumber: NewNumberDataElement(messageNumber, 4),
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
