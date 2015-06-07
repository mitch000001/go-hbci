package hbci

func MakeCall() string {
	return ""
}

func InitializeDialog() {}

type EncryptedMessage struct {
	*Message
	EncryptionHeader *EncryptionHeader
	EncryptedData    *EncryptedData
}

type EncryptionHeader struct{}

type EncryptedData struct{}

type Message struct {
	Header *MessageHeader
	End    *MessageEnd
}

type BankMessage struct {
	*Message
	SignatureBegin         *SignatureHeader
	SignatureEnd           *SignatureEnd
	MessageAcknowledgement *MessageAcknowledgement
	SegmentAcknowledgement *SegmentAcknowledgement
	DataSegments           []*DataSegment
}

type MessageAcknowledgement struct{}
type SegmentAcknowledgement struct{}
type DataSegment struct{}

type ClientMessage struct {
	*Message
	SignatureBegin *SignatureHeader
	SignatureEnd   *SignatureEnd
	Jobs           []*Job
}

type MessageHeader struct {
	*Segment
}

type MessageEnd struct{}
type SignatureHeader struct{}
type SignatureEnd struct{}
type Job struct{}

type SegmentSequence struct{}
