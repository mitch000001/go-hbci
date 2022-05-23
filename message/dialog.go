package message

import "github.com/mitch000001/go-hbci/segment"

// NewDialogInitializationClientMessage creates a basic client message used for DialogInitialization
func NewDialogInitializationClientMessage(hbciVersion segment.HBCIVersion) *DialogInitializationClientMessage {
	d := &DialogInitializationClientMessage{
		hbciVersion: hbciVersion,
	}
	d.BasicMessage = NewBasicMessage(d)
	return d
}

// DialogInitializationClientMessage represents a client message used to initialize a dialog
type DialogInitializationClientMessage struct {
	*BasicMessage
	Identification             *segment.IdentificationSegment
	ProcessingPreparation      *segment.ProcessingPreparationSegmentV3
	TanRequest                 *segment.TanRequestSegment
	PublicSigningKeyRequest    *segment.PublicKeyRequestSegment
	PublicEncryptionKeyRequest *segment.PublicKeyRequestSegment
	PublicKeyRequest           *segment.PublicKeyRequestSegment
	hbciVersion                segment.HBCIVersion
}

// HBCIVersion returns the version used for this message
func (d *DialogInitializationClientMessage) HBCIVersion() segment.HBCIVersion {
	return d.hbciVersion
}

// HBCISegments returns all segment from this message
func (d *DialogInitializationClientMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{
		d.Identification,
		d.ProcessingPreparation,
		d.TanRequest,
		d.PublicSigningKeyRequest,
		d.PublicEncryptionKeyRequest,
		d.PublicKeyRequest,
	}
}

func (d *DialogInitializationClientMessage) jobs() []segment.Segment {
	return []segment.Segment{
		d.Identification,
		d.ProcessingPreparation,
		d.TanRequest,
		d.PublicSigningKeyRequest,
		d.PublicEncryptionKeyRequest,
	}
}

// NewDialogFinishingMessage creates a message used to finish a dialog
func NewDialogFinishingMessage(hbciVersion segment.HBCIVersion, dialogID string) *DialogFinishingMessage {
	d := &DialogFinishingMessage{
		DialogEnd:   segment.NewDialogEndSegment(dialogID),
		hbciVersion: hbciVersion,
	}
	d.BasicMessage = NewBasicMessage(d)
	return d
}

// DialogFinishingMessage represents a message used to finish a dialog
type DialogFinishingMessage struct {
	*BasicMessage
	DialogEnd   *segment.DialogEndSegment
	hbciVersion segment.HBCIVersion
}

// HBCIVersion returns the version used for this message
func (d *DialogFinishingMessage) HBCIVersion() segment.HBCIVersion {
	return d.hbciVersion
}

// HBCISegments returns all segment from this message
func (d *DialogFinishingMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{
		d.DialogEnd,
	}
}

func (d *DialogFinishingMessage) jobs() []segment.ClientSegment {
	return []segment.ClientSegment{
		d.DialogEnd,
	}
}

// NewDialogCancellationMessage creates a message to cancel a dialog
func NewDialogCancellationMessage(messageAcknowledgement *segment.MessageAcknowledgement) *DialogCancellationMessage {
	d := &DialogCancellationMessage{
		MessageAcknowledgements: messageAcknowledgement,
	}
	return d
}

// DialogCancellationMessage represents a message used to cancel a dialog
type DialogCancellationMessage struct {
	*BasicMessage
	MessageAcknowledgements *segment.MessageAcknowledgement
}

// AnonymousDialogMessage represents a message used by anonymous dialogs
type AnonymousDialogMessage struct {
	*BasicMessage
	Identification        *segment.IdentificationSegment
	ProcessingPreparation *segment.ProcessingPreparationSegmentV3
}
