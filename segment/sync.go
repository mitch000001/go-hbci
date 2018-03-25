package segment

import (
	"strconv"

	"github.com/mitch000001/go-hbci/element"
)

// Possible sync modes
var (
	SyncModeAquireClientID               = SyncMode{mode: 0}
	SyncModeAquireLastProcessedMessageID = SyncMode{mode: 1}
	SyncModeAquireSignatureID            = SyncMode{mode: 2}
)

type SyncMode struct {
	mode int
}

func NewSynchronisationSegmentV2(modus SyncMode) *SynchronisationRequestSegment {
	s := &SynchronisationRequestV2{
		SyncModus: element.NewNumber(modus.mode, 1),
	}
	s.ClientSegment = NewBasicSegment(5, s)

	segment := &SynchronisationRequestSegment{
		ClientSegment: s,
	}
	return segment
}

type SynchronisationRequestSegment struct {
	ClientSegment
}

type SynchronisationRequestV2 struct {
	ClientSegment
	// Code | Bedeutung
	// ---------------------------------------------------------
	// 0 ￼ ￼| Neue Kundensystem-ID zurückmelden
	// 1	| Letzte verarbeitete Nachrichtennummer zurückmelden ￼ ￼
	// 2 ￼ ￼| Signatur-ID zurückmelden
	SyncModus *element.NumberDataElement
}

func (s *SynchronisationRequestV2) Version() int         { return 2 }
func (s *SynchronisationRequestV2) ID() string           { return "HKSYN" }
func (s *SynchronisationRequestV2) referencedId() string { return "" }
func (s *SynchronisationRequestV2) sender() string       { return senderUser }

func (s *SynchronisationRequestV2) elements() []element.DataElement {
	return []element.DataElement{
		s.SyncModus,
	}
}

func NewSynchronisationSegmentV3(modus SyncMode) *SynchronisationRequestSegment {
	s := &SynchronisationRequestV3{
		SyncModus: element.NewCode(strconv.Itoa(modus.mode), 1, []string{"0", "1", "2"}),
	}
	s.ClientSegment = NewBasicSegment(5, s)

	segment := &SynchronisationRequestSegment{
		ClientSegment: s,
	}
	return segment
}

type SynchronisationRequestV3 struct {
	ClientSegment
	// Code | Bedeutung
	// ---------------------------------------------------------
	// 0 ￼ ￼| Neue Kundensystem-ID zurückmelden
	// 1	| Letzte verarbeitete Nachrichtennummer zurückmelden ￼ ￼
	// 2 ￼ ￼| Signatur-ID zurückmelden
	SyncModus *element.CodeDataElement
}

func (s *SynchronisationRequestV3) Version() int         { return 3 }
func (s *SynchronisationRequestV3) ID() string           { return "HKSYN" }
func (s *SynchronisationRequestV3) referencedId() string { return "" }
func (s *SynchronisationRequestV3) sender() string       { return senderUser }

func (s *SynchronisationRequestV3) elements() []element.DataElement {
	return []element.DataElement{
		s.SyncModus,
	}
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment SynchronisationResponseSegment -segment_interface SynchronisationResponse -segment_versions="SynchronisationResponseSegmentV3:3,SynchronisationResponseSegmentV4:4"

type SynchronisationResponse interface {
	BankSegment
	ClientSystemID() string
	LastMessageNumber() int
	SignatureID() int
}

type SynchronisationResponseSegment struct {
	SynchronisationResponse
}

type SynchronisationResponseSegmentV3 struct {
	Segment
	ClientSystemIDResponse *element.IdentificationDataElement
	MessageNumberResponse  *element.NumberDataElement
	SignatureIDResponse    *element.NumberDataElement
}

func (s *SynchronisationResponseSegmentV3) Version() int         { return 3 }
func (s *SynchronisationResponseSegmentV3) ID() string           { return "HISYN" }
func (s *SynchronisationResponseSegmentV3) referencedId() string { return "HKSYN" }
func (s *SynchronisationResponseSegmentV3) sender() string       { return senderBank }

func (s *SynchronisationResponseSegmentV3) elements() []element.DataElement {
	return []element.DataElement{
		s.ClientSystemIDResponse,
		s.MessageNumberResponse,
		s.SignatureIDResponse,
	}
}

func (s *SynchronisationResponseSegmentV3) ClientSystemID() string {
	return s.ClientSystemIDResponse.Val()
}

func (s *SynchronisationResponseSegmentV3) LastMessageNumber() int {
	return s.MessageNumberResponse.Val()
}

func (s *SynchronisationResponseSegmentV3) SignatureID() int {
	return s.SignatureIDResponse.Val()
}

type SynchronisationResponseSegmentV4 struct {
	Segment
	ClientSystemIDResponse *element.IdentificationDataElement
	MessageNumberResponse  *element.NumberDataElement
	SignatureIDResponse    *element.NumberDataElement
}

func (s *SynchronisationResponseSegmentV4) Version() int         { return 4 }
func (s *SynchronisationResponseSegmentV4) ID() string           { return "HISYN" }
func (s *SynchronisationResponseSegmentV4) referencedId() string { return "HKSYN" }
func (s *SynchronisationResponseSegmentV4) sender() string       { return senderBank }

func (s *SynchronisationResponseSegmentV4) elements() []element.DataElement {
	return []element.DataElement{
		s.ClientSystemIDResponse,
		s.MessageNumberResponse,
		s.SignatureIDResponse,
	}
}

func (s *SynchronisationResponseSegmentV4) ClientSystemID() string {
	return s.ClientSystemIDResponse.Val()
}

func (s *SynchronisationResponseSegmentV4) LastMessageNumber() int {
	return s.MessageNumberResponse.Val()
}

func (s *SynchronisationResponseSegmentV4) SignatureID() int {
	return s.SignatureIDResponse.Val()
}
