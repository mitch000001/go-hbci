package hbci

type EncryptedMessage struct {
	*Message
	EncryptionHeader *EncryptionHeaderSegment
	EncryptedData    *EncryptedDataSegment
}

type Message struct {
	Header *MessageHeaderSegment
	End    *MessageEndSegment
}

type BankMessage struct {
	*Message
	SignatureBegin         *SignatureHeader
	SignatureEnd           *SignatureEnd
	MessageAcknowledgement *MessageAcknowledgement
	SegmentAcknowledgement *SegmentAcknowledgement
	DataSegments           []*DataSegment
}

type DataSegment struct{}

type ClientMessage struct {
	*Message
	SignatureBegin *SignatureHeader
	SignatureEnd   *SignatureEnd
	Jobs           []*Job
}

type SignatureHeader struct{}
type SignatureEnd struct{}
type Job struct{}

type SegmentSequence struct{}

func NewReferencingMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int, referencedMessage *ReferenceMessage) *MessageHeaderSegment {
	segmentHeader := NewSegmentHeader("HNHBK", 1, 3)
	return &MessageHeaderSegment{
		Header:      segmentHeader,
		Size:        NewDigitDataElement(size, 12),
		HBCIVersion: NewNumberDataElement(hbciVersion, 3),
		DialogID:    NewIdentificationDataElement(dialogId),
		Number:      NewNumberDataElement(number, 4),
		Ref:         referencedMessage,
	}
}

func NewMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int) *MessageHeaderSegment {
	segmentHeader := NewSegmentHeader("HNHBK", 1, 3)
	return &MessageHeaderSegment{
		Header:      segmentHeader,
		Size:        NewDigitDataElement(size, 12),
		HBCIVersion: NewNumberDataElement(hbciVersion, 3),
		DialogID:    NewIdentificationDataElement(dialogId),
		Number:      NewNumberDataElement(number, 4),
	}
}

type MessageHeaderSegment struct {
	*segment
	Header      *SegmentHeader
	Size        *DigitDataElement
	HBCIVersion *NumberDataElement
	DialogID    *IdentificationDataElement
	Number      *NumberDataElement
	Ref         *ReferenceMessage
}

func (m *MessageHeaderSegment) DataElements() []DataElement {
	return []DataElement{
		m.Size,
		m.HBCIVersion,
		m.DialogID,
		m.Number,
		m.Ref,
	}
}

func NewMessageEndSegment(segmentNumber, messageNumber int) *MessageEndSegment {
	segmentHeader := NewSegmentHeader("HNHBS", segmentNumber, 1)
	end := &MessageEndSegment{
		Number: NewNumberDataElement(messageNumber, 4),
	}
	end.segment = NewSegment(segmentHeader, end)
	return end
}

type MessageEndSegment struct {
	*segment
	Number *NumberDataElement
}

func (m *MessageEndSegment) DataElements() []DataElement {
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

func (r *ReferenceMessage) Valid() bool {
	if r.DialogID == nil || r.MessageNumber == nil {
		return false
	} else {
		return r.elementGroup.Valid()
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
