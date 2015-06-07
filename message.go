package hbci

func NewReferencingMessageHeaderSegment(header *SegmentHeader, size int, hbciVersion int, dialogId string, number int, referencedMessage *ReferenceMessage) *MessageHeaderSegment {
	return &MessageHeaderSegment{
		Header:      header,
		Size:        NewDigitDataElement(size, 12),
		HBCIVersion: NewNumberDataElement(hbciVersion, 3),
		DialogID:    NewAlphaNumericDataElement(dialogId, 30),
		Number:      NewNumberDataElement(number, 4),
		Ref:         referencedMessage,
	}
}

func NewMessageHeaderSegment(header *SegmentHeader, size int, hbciVersion int, dialogId string, number int) *MessageHeaderSegment {
	return &MessageHeaderSegment{
		Header:      header,
		Size:        NewDigitDataElement(size, 12),
		HBCIVersion: NewNumberDataElement(hbciVersion, 3),
		DialogID:    NewAlphaNumericDataElement(dialogId, 30),
		Number:      NewNumberDataElement(number, 4),
	}
}

type MessageHeaderSegment struct {
	Header      *SegmentHeader
	Size        DataElement
	HBCIVersion DataElement
	DialogID    DataElement
	Number      DataElement
	Ref         *ReferenceMessage
}

type ReferenceMessage struct{}
