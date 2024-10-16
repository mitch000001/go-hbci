package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

// NewDecryptedMessage creates a new decrypted message from rawMessage
func NewDecryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, rawMessage []byte) (BankMessage, error) {
	unmarshaler := NewUnmarshaler(rawMessage)
	err := unmarshaler.Unmarshal()
	if err != nil {
		return nil, fmt.Errorf("Malformed decrypted message bytes: %v", err)
	}
	var acknowledgements []domain.Acknowledgement
	messageAcknowledgement, ok := unmarshaler.SegmentByID("HIRMG").(*segment.MessageAcknowledgement)
	if ok {
		messageAcknowledgement.SetReferencingMessage(header.ReferencingMessage())
		acknowledgements = append(acknowledgements, messageAcknowledgement.Acknowledgements()...)
	}
	for _, seg := range unmarshaler.SegmentsByID("HIRMS") {
		if segmentAcknowledgement, ok := seg.(*segment.SegmentAcknowledgement); ok {
			acknowledgements = append(acknowledgements, segmentAcknowledgement.Acknowledgements()...)
		} else {
			panic(fmt.Errorf("Error while unmarshaling segments"))
		}
	}
	version, ok := segment.SupportedHBCIVersions[header.HBCIVersion.Val()]
	if !ok {
		return nil, fmt.Errorf("Unsupported HBCI version: %d", header.HBCIVersion.Val())
	}
	message := &decryptedMessage{
		rawMessage:       rawMessage,
		acknowledgements: acknowledgements,
		unmarshaler:      unmarshaler,
		hbciVersion:      version,
	}
	// TODO: set hbci message appropriate, if possible
	message.message = NewBasicMessageWithHeaderAndEnd(header, end, message)
	return message, nil
}

// A decryptedMessage represents a message which was decrypted using a CryptoProvider
type decryptedMessage struct {
	rawMessage       []byte
	message          Message
	acknowledgements []domain.Acknowledgement
	unmarshaler      *Unmarshaler
	hbciVersion      segment.HBCIVersion
}

// MarshalHBCI marshals d to HBCI wire format
func (d *decryptedMessage) MarshalHBCI() ([]byte, error) {
	return d.rawMessage, nil
}

func (d *decryptedMessage) MessageHeader() *segment.MessageHeaderSegment {
	return d.message.MessageHeader()
}

func (d *decryptedMessage) MessageEnd() *segment.MessageEndSegment {
	return d.message.MessageEnd()
}

func (d *decryptedMessage) FindMarshaledSegment(segmentID string) []byte {
	return d.unmarshaler.MarshaledSegmentByID(segmentID)
}

func (d *decryptedMessage) FindMarshaledSegments(segmentID string) [][]byte {
	return d.unmarshaler.MarshaledSegmentsByID(segmentID)
}

func (d *decryptedMessage) MarshaledSegments() [][]byte {
	return d.unmarshaler.MarshaledSegments()
}

func (d *decryptedMessage) FindSegment(segmentID string) segment.Segment {
	return d.unmarshaler.SegmentByID(segmentID)
}

func (d *decryptedMessage) FindSegments(segmentID string) []segment.Segment {
	return d.unmarshaler.SegmentsByID(segmentID)
}

func (d *decryptedMessage) SegmentPosition(segmentID string) int {
	seg := d.unmarshaler.MarshaledSegmentByID(segmentID)
	if len(seg) == 0 {
		return -1
	}
	elements, err := segment.ExtractElements(seg)
	if err != nil {
		return -1
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return -1
	}
	return header.Position.Val()
}

func (d *decryptedMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{}
}

func (d *decryptedMessage) HBCIVersion() segment.HBCIVersion {
	return d.hbciVersion
}

func (d *decryptedMessage) Acknowledgements() []domain.Acknowledgement {
	return d.acknowledgements
}

func (d *decryptedMessage) SupportedSegments() []segment.VersionedSegment {
	var versionedSegments []segment.VersionedSegment
	for _, s := range d.unmarshaler.MarshaledSegments() {
		vid, err := extractVersionedSegmentIdentifier(s)
		if err != nil {
			continue
		}
		if len(vid.ID) <= 5 {
			continue
		}
		versionedSegments = append(versionedSegments, vid)
	}
	return versionedSegments
}
