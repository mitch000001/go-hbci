package hbci

type basicMessage struct {
	Header *MessageHeaderSegment
	End    *MessageEndSegment
}

type ClientMessage interface {
	Jobs() SegmentSequence
}

type BankMessage interface {
	DataSegments() SegmentSequence
}

type basicBankMessage struct {
	*basicMessage
	BankMessage
	SignatureBegin          *SignatureHeaderSegment
	SignatureEnd            *SignatureEndSegment
	MessageAcknowledgements *MessageAcknowledgement
	SegmentAcknowledgements *SegmentAcknowledgement
}

type basicClientMessage struct {
	*basicMessage
	SignatureBegin *SignatureHeaderSegment
	SignatureEnd   *SignatureEndSegment
}

type SegmentSequence []Segment

func NewDialogCancellationMessage(messageAcknowledgement *MessageAcknowledgement) *DialogCancellationMessage {
	d := &DialogCancellationMessage{
		MessageAcknowledgements: messageAcknowledgement,
	}
	return d
}

type DialogCancellationMessage struct {
	*basicMessage
	MessageAcknowledgements *MessageAcknowledgement
}

var validHBCIVersions = []int{201, 210, 220}

func NewReferencingMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int, referencedMessage *ReferenceMessage) *MessageHeaderSegment {
	m := NewMessageHeaderSegment(size, hbciVersion, dialogId, number)
	m.Ref = referencedMessage
	return m
}

func NewMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int) *MessageHeaderSegment {
	m := &MessageHeaderSegment{
		Size:        NewDigitDataElement(size, 12),
		HBCIVersion: NewNumberDataElement(hbciVersion, 3),
		DialogID:    NewIdentificationDataElement(dialogId),
		Number:      NewNumberDataElement(number, 4),
	}
	m.basicSegment = NewBasicSegment("HNHBK", 1, 3, m)
	return m
}

type MessageHeaderSegment struct {
	*basicSegment
	Size        *DigitDataElement
	HBCIVersion *NumberDataElement
	DialogID    *IdentificationDataElement
	Number      *NumberDataElement
	Ref         *ReferenceMessage
}

func (m *MessageHeaderSegment) elements() []DataElement {
	return []DataElement{
		m.Size,
		m.HBCIVersion,
		m.DialogID,
		m.Number,
		m.Ref,
	}
}

func (m *MessageHeaderSegment) SetSize(size int) {
	m.Size = NewDigitDataElement(size, 12)
}

func NewMessageEndSegment(segmentNumber, messageNumber int) *MessageEndSegment {
	end := &MessageEndSegment{
		Number: NewNumberDataElement(messageNumber, 4),
	}
	end.basicSegment = NewBasicSegment("HNHBS", segmentNumber, 1, end)
	return end
}

type MessageEndSegment struct {
	*basicSegment
	Number *NumberDataElement
}

func (m *MessageEndSegment) elements() []DataElement {
	return []DataElement{
		m.Number,
	}
}

func NewReferenceMessage(dialogId string, messageNumber int) *ReferenceMessage {
	r := &ReferenceMessage{
		DialogID:      NewIdentificationDataElement(dialogId),
		MessageNumber: NewNumberDataElement(messageNumber, 4),
	}
	r.elementGroup = NewDataElementGroup(ReferenceMessageDEG, 2, r)
	return r
}

type ReferenceMessage struct {
	*elementGroup
	DialogID      *IdentificationDataElement
	MessageNumber *NumberDataElement
}

func (r *ReferenceMessage) IsValid() bool {
	if r.DialogID == nil || r.MessageNumber == nil {
		return false
	} else {
		return r.elementGroup.IsValid()
	}
}

func (r *ReferenceMessage) Value() interface{} {
	return r
}

func (r *ReferenceMessage) groupDataElements() []DataElement {
	return []DataElement{
		r.DialogID,
		r.MessageNumber,
	}
}
