package message

import (
	"fmt"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

func NewDecryptedMessage(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, rawMessage []byte) (*DecryptedMessage, error) {
	unmarshaler := NewUnmarshaler(rawMessage)
	err := unmarshaler.Unmarshal()
	if err != nil {
		return nil, fmt.Errorf("Malformed decrypted message bytes: %v", err)
	}
	messageAcknowledgement := unmarshaler.SegmentById("HIRMG").(*segment.MessageAcknowledgement)
	messageAcknowledgement.SetReferencingMessage(header.ReferencingMessage())
	acknowledgements := messageAcknowledgement.Acknowledgements()
	for _, seg := range unmarshaler.SegmentsById("HIRMS") {
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
	decryptedMessage := &DecryptedMessage{
		rawMessage:       rawMessage,
		acknowledgements: acknowledgements,
		unmarshaler:      unmarshaler,
		hbciVersion:      version,
	}
	// TODO: set hbci message appropriate, if possible
	decryptedMessage.message = NewBasicMessageWithHeaderAndEnd(header, end, decryptedMessage)
	return decryptedMessage, nil
}

type DecryptedMessage struct {
	rawMessage       []byte
	message          Message
	acknowledgements []domain.Acknowledgement
	unmarshaler      *Unmarshaler
	hbciVersion      segment.HBCIVersion
}

func (d *DecryptedMessage) MarshalHBCI() ([]byte, error) {
	return d.rawMessage, nil
}

func (d *DecryptedMessage) MessageHeader() *segment.MessageHeaderSegment {
	return d.message.MessageHeader()
}

func (d *DecryptedMessage) MessageEnd() *segment.MessageEndSegment {
	return d.message.MessageEnd()
}

func (d *DecryptedMessage) FindMarshaledSegment(segmentID string) []byte {
	return d.unmarshaler.MarshaledSegmentById(segmentID)
}

func (d *DecryptedMessage) FindMarshaledSegments(segmentID string) [][]byte {
	return d.unmarshaler.MarshaledSegmentsById(segmentID)
}

func (d *DecryptedMessage) MarshaledSegments() [][]byte {
	return d.unmarshaler.MarshaledSegments()
}

func (d *DecryptedMessage) FindSegment(segmentID string) segment.Segment {
	return d.unmarshaler.SegmentById(segmentID)
}

func (d *DecryptedMessage) FindSegments(segmentID string) []segment.Segment {
	return d.unmarshaler.SegmentsById(segmentID)
}

func (d *DecryptedMessage) SegmentNumber(segmentID string) int {
	seg := d.unmarshaler.MarshaledSegmentById(segmentID)
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
	return header.Number.Val()
}

func (d *DecryptedMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{}
}

func (d *DecryptedMessage) HBCIVersion() segment.HBCIVersion {
	return d.hbciVersion
}

func (d *DecryptedMessage) Acknowledgements() []domain.Acknowledgement {
	return d.acknowledgements
}
