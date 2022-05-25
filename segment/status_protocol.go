package segment

import (
	"sort"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func StatusProtocolRequestBuilder(versions []int) (func(from, to time.Time, maxEntries int, continuationReference string) StatusProtocolRequest, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	for _, version := range versions {
		switch version {
		case 4:
			return NewStatusProtocolRequestV3, nil
		case 3:
			return NewStatusProtocolRequestV4, nil
		default:
			continue
		}
	}
	return nil, &unsupportedSegmentVersionError{segmentID: "HKPRO", versions: versions}
}

type StatusProtocolRequest interface {
	ClientSegment
}

func NewStatusProtocolRequestV3(from, to time.Time, maxEntries int, continuationReference string) StatusProtocolRequest {
	s := &StatusProtocolRequestSegmentV3{
		From:       element.NewDate(from),
		To:         element.NewDate(to),
		MaxEntries: element.NewNumber(maxEntries, 4),
	}
	if continuationReference != "" {
		s.ContinuationReference = element.NewAlphaNumeric(continuationReference, 35)
	}
	s.ClientSegment = NewBasicSegment(1, s)
	return s
}

type StatusProtocolRequestSegmentV3 struct {
	ClientSegment
	From                  *element.DateDataElement
	To                    *element.DateDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (s *StatusProtocolRequestSegmentV3) ID() string           { return "HKPRO" }
func (s *StatusProtocolRequestSegmentV3) Version() int         { return 3 }
func (s *StatusProtocolRequestSegmentV3) referencedId() string { return "" }
func (s *StatusProtocolRequestSegmentV3) sender() string       { return senderUser }

func (s *StatusProtocolRequestSegmentV3) elements() []element.DataElement {
	return []element.DataElement{
		s.From,
		s.To,
		s.MaxEntries,
		s.ContinuationReference,
	}
}

func NewStatusProtocolRequestV4(from, to time.Time, maxEntries int, continuationReference string) StatusProtocolRequest {
	s := &StatusProtocolRequestSegmentV4{
		From:       element.NewDate(from),
		To:         element.NewDate(to),
		MaxEntries: element.NewNumber(maxEntries, 4),
	}
	if continuationReference != "" {
		s.ContinuationReference = element.NewAlphaNumeric(continuationReference, 35)
	}
	s.ClientSegment = NewBasicSegment(1, s)
	return s
}

type StatusProtocolRequestSegmentV4 struct {
	ClientSegment
	From                  *element.DateDataElement
	To                    *element.DateDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (s *StatusProtocolRequestSegmentV4) ID() string           { return "HKPRO" }
func (s *StatusProtocolRequestSegmentV4) Version() int         { return 4 }
func (s *StatusProtocolRequestSegmentV4) referencedId() string { return "" }
func (s *StatusProtocolRequestSegmentV4) sender() string       { return senderUser }

func (s *StatusProtocolRequestSegmentV4) elements() []element.DataElement {
	return []element.DataElement{
		s.From,
		s.To,
		s.MaxEntries,
		s.ContinuationReference,
	}
}

type StatusProtocolResponse interface {
	BankSegment
	Status() domain.StatusAcknowledgement
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment StatusProtocolResponseSegment -segment_interface StatusProtocolResponse -segment_versions="StatusProtocolResponseSegmentV3:3:Segment"

type StatusProtocolResponseSegment struct {
	StatusProtocolResponse
}

type StatusProtocolResponseSegmentV3 struct {
	Segment
	ReferencingMessage *element.ReferencingMessageDataElement
	ReferencingSegment *element.NumberDataElement
	Date               *element.DateDataElement
	Time               *element.TimeDataElement
	Acknowledgement    *element.AcknowledgementDataElement
}

func (s *StatusProtocolResponseSegmentV3) ID() string           { return "HIPRO" }
func (s *StatusProtocolResponseSegmentV3) Version() int         { return 3 }
func (s *StatusProtocolResponseSegmentV3) referencedId() string { return "HKPRO" }
func (s *StatusProtocolResponseSegmentV3) sender() string       { return senderBank }

func (s *StatusProtocolResponseSegmentV3) elements() []element.DataElement {
	return []element.DataElement{
		s.ReferencingMessage,
		s.ReferencingSegment,
		s.Date,
		s.Time,
		s.Acknowledgement,
	}
}

func (s *StatusProtocolResponseSegmentV3) Status() domain.StatusAcknowledgement {
	ack := s.Acknowledgement.Val()
	if s.ReferencingMessage != nil {
		ack.ReferencingMessage = s.ReferencingMessage.Val()
		ack.ReferencingSegmentNumber = s.ReferencingSegment.Val()
		ack.Type = domain.SegmentAcknowledgement
	} else {
		ack.Type = domain.MessageAcknowledgement
	}
	status := domain.StatusAcknowledgement{
		Acknowledgement: ack,
		TransmittedAt:   s.Date.Val().Add(time.Since(s.Time.Val())),
	}
	return status
}

type StatusProtocolResponseSegmentV4 struct {
	Segment
	ReferencingMessage *element.ReferencingMessageDataElement
	ReferencingSegment *element.NumberDataElement
	Date               *element.DateDataElement
	Time               *element.TimeDataElement
	Acknowledgement    *element.AcknowledgementDataElement
}

func (s *StatusProtocolResponseSegmentV4) ID() string           { return "HIPRO" }
func (s *StatusProtocolResponseSegmentV4) Version() int         { return 4 }
func (s *StatusProtocolResponseSegmentV4) referencedId() string { return "HKPRO" }
func (s *StatusProtocolResponseSegmentV4) sender() string       { return senderBank }

func (s *StatusProtocolResponseSegmentV4) elements() []element.DataElement {
	return []element.DataElement{
		s.ReferencingMessage,
		s.ReferencingSegment,
		s.Date,
		s.Time,
		s.Acknowledgement,
	}
}

func (s *StatusProtocolResponseSegmentV4) Status() domain.StatusAcknowledgement {
	ack := s.Acknowledgement.Val()
	if s.ReferencingMessage != nil {
		ack.ReferencingMessage = s.ReferencingMessage.Val()
		ack.ReferencingSegmentNumber = s.ReferencingSegment.Val()
		ack.Type = domain.SegmentAcknowledgement
	} else {
		ack.Type = domain.MessageAcknowledgement
	}
	status := domain.StatusAcknowledgement{
		Acknowledgement: ack,
		TransmittedAt:   s.Date.Val().Add(time.Since(s.Time.Val())),
	}
	return status
}
