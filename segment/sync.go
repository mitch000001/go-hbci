package segment

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/element"
)

func NewSynchronisationSegment(modus int) *SynchronisationRequestSegment {
	s := &SynchronisationRequestSegment{
		SyncModus: element.NewNumber(modus, 1),
	}
	s.Segment = NewBasicSegment(5, s)
	return s
}

type SynchronisationRequestSegment struct {
	Segment
	// Code | Bedeutung
	// ---------------------------------------------------------
	// 0 ￼ ￼| Neue Kundensystem-ID zurückmelden
	// 1	| Letzte verarbeitete Nachrichtennummer zurückmelden ￼ ￼
	// 2 ￼ ￼| Signatur-ID zurückmelden
	SyncModus *element.NumberDataElement
}

func (s *SynchronisationRequestSegment) version() int         { return 2 }
func (s *SynchronisationRequestSegment) id() string           { return "HKSYN" }
func (s *SynchronisationRequestSegment) referencedId() string { return "" }
func (s *SynchronisationRequestSegment) sender() string       { return senderUser }

func (s *SynchronisationRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.SyncModus,
	}
}

type SynchronisationResponseSegment struct {
	Segment
	ClientSystemID *element.IdentificationDataElement
	MessageNumber  *element.NumberDataElement
	SignatureID    *element.NumberDataElement
}

func (s *SynchronisationResponseSegment) version() int         { return 3 }
func (s *SynchronisationResponseSegment) id() string           { return "HISYN" }
func (s *SynchronisationResponseSegment) referencedId() string { return "HKSYN" }
func (s *SynchronisationResponseSegment) sender() string       { return senderBank }

func (s *SynchronisationResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.ClientSystemID,
		s.MessageNumber,
		s.SignatureID,
	}
}

func (s *SynchronisationResponseSegment) UnmarshalHBCI(value []byte) error {
	value = bytes.TrimSuffix(value, []byte("'"))
	elements := bytes.Split(value, []byte("+"))
	header := elements[0]
	headerElems := bytes.Split(header, []byte(":"))
	num, err := strconv.Atoi(string(headerElems[1]))
	if err != nil {
		return fmt.Errorf("%T: Malformed segment header", s)
	}
	if len(headerElems) == 4 {
		ref, err := strconv.Atoi(string(headerElems[3]))
		if err != nil {
			return fmt.Errorf("%T: Malformed segment header reference: %v", s, err)
		}
		s.Segment = NewReferencingBasicSegment(num, ref, s)
	} else {
		s.Segment = NewBasicSegment(num, s)
	}
	s.ClientSystemID = element.NewIdentification(string(elements[1]))
	if len(elements) >= 3 && len(elements[2]) > 0 {
		messageNum, err := strconv.Atoi(string(elements[2]))
		if err != nil {
			return fmt.Errorf("%T: Malformed message number: %v", s, err)
		}
		s.MessageNumber = element.NewNumber(messageNum, 4)
	}
	if len(elements) >= 4 && len(elements[3]) > 0 {
		signatureID, err := strconv.Atoi(string(elements[3]))
		if err != nil {
			return fmt.Errorf("%T: Malformed signature id: %v", s, err)
		}
		s.SignatureID = element.NewNumber(signatureID, 16)
	}
	return nil
}
