package message

import "github.com/mitch000001/go-hbci/segment"

func NewDialogInitializationClientMessage() *DialogInitializationClientMessage {
	d := &DialogInitializationClientMessage{}
	d.BasicClientMessage = NewBasicClientMessage(d)
	return d
}

type DialogInitializationClientMessage struct {
	*BasicClientMessage
	Identification             *segment.IdentificationSegment
	ProcessingPreparation      *segment.ProcessingPreparationSegment
	PublicSigningKeyRequest    *segment.PublicKeyRequestSegment
	PublicEncryptionKeyRequest *segment.PublicKeyRequestSegment
}

func (d *DialogInitializationClientMessage) jobs() []segment.Segment {
	return []segment.Segment{
		d.Identification,
		d.ProcessingPreparation,
		d.PublicSigningKeyRequest,
		d.PublicEncryptionKeyRequest,
	}
}

type DialogInitializationBankMessage struct {
	*basicBankMessage
	BankParams            segment.SegmentSequence
	UserParams            segment.SegmentSequence
	PublicKeyTransmission *segment.PublicKeyTransmissionSegment
	Announcement          *segment.BankAnnouncementSegment
}

type DialogFinishingMessage struct {
	*BasicClientMessage
	DialogEnd *segment.DialogEndSegment
}

func (d *DialogFinishingMessage) jobs() []segment.Segment {
	return []segment.Segment{
		d.DialogEnd,
	}
}

func NewDialogCancellationMessage(messageAcknowledgement *segment.MessageAcknowledgement) *DialogCancellationMessage {
	d := &DialogCancellationMessage{
		MessageAcknowledgements: messageAcknowledgement,
	}
	return d
}

type DialogCancellationMessage struct {
	*BasicMessage
	MessageAcknowledgements *segment.MessageAcknowledgement
}

type AnonymousDialogMessage struct {
	*BasicMessage
	Identification        *segment.IdentificationSegment
	ProcessingPreparation *segment.ProcessingPreparationSegment
}
