package domain

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

func NewMessageAcknowledgement(code int, referenceDataElement, text string, params []string) Acknowledgement {
	return Acknowledgement{
		Type:                 MessageAcknowledgement,
		Code:                 code,
		ReferenceDataElement: referenceDataElement,
		Text:                 text,
		Params:               params,
	}
}

func NewSegmentAcknowledgement(code int, referenceDataElement, text string, params []string) Acknowledgement {
	return Acknowledgement{
		Type:                 SegmentAcknowledgement,
		Code:                 code,
		ReferenceDataElement: referenceDataElement,
		Text:                 text,
		Params:               params,
	}
}

const (
	MessageAcknowledgement = "MessageAcknowledgement"
	SegmentAcknowledgement = "SegmentAcknowledgement"
)

type Acknowledgement struct {
	Type                     string
	Code                     int
	ReferenceDataElement     string
	Text                     string
	Params                   []string
	ReferencingMessage       ReferencingMessage
	ReferencingSegmentNumber int
}

func (a Acknowledgement) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s for message %d (%s)", a.Type, a.ReferencingMessage.MessageNumber, a.ReferencingMessage.DialogID)
	if a.ReferencingSegmentNumber > 0 {
		fmt.Fprintf(&buf, ", segment %d: ", a.ReferencingSegmentNumber)
	} else {
		fmt.Fprintf(&buf, ": ")
	}
	fmt.Fprintf(&buf, "Code: %d, ", a.Code)
	if a.ReferenceDataElement != "" {
		fmt.Fprintf(&buf, "Position: %s, ", a.ReferenceDataElement)
	} else {
		fmt.Fprintf(&buf, "Position: none, ")
	}
	fmt.Fprintf(&buf, "Text: '%s'", a.Text)
	if len(a.Params) != 0 {
		fmt.Fprintf(&buf, ", Parameters: %s", strings.Join(a.Params, ", "))
	}
	return buf.String()
}

func (a Acknowledgement) IsMessageAcknowledgement() bool {
	return a.Type == MessageAcknowledgement
}

func (a Acknowledgement) IsSegmentAcknowledgement() bool {
	return a.Type == SegmentAcknowledgement
}

func (a Acknowledgement) IsError() bool {
	return a.Code >= 9000
}

func (a Acknowledgement) IsWarning() bool {
	return a.Code >= 3000 && a.Code < 4000
}

func (a Acknowledgement) IsSuccess() bool {
	return a.Code > 0 && a.Code < 1000
}

type StatusAcknowledgement struct {
	Acknowledgement
	TransmittedAt time.Time
}

func (s StatusAcknowledgement) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s for message %d (%s)", s.Type, s.ReferencingMessage.MessageNumber, s.ReferencingMessage.DialogID)
	if s.ReferencingSegmentNumber > 0 {
		fmt.Fprintf(&buf, ", segment %d: ", s.ReferencingSegmentNumber)
	} else {
		fmt.Fprintf(&buf, ": ")
	}
	fmt.Fprintf(&buf, "Transmitted at: %s, ", s.TransmittedAt)
	fmt.Fprintf(&buf, "Code: %d, ", s.Code)
	if s.ReferenceDataElement != "" {
		fmt.Fprintf(&buf, "Position: %s, ", s.ReferenceDataElement)
	} else {
		fmt.Fprintf(&buf, "Position: none, ")
	}
	fmt.Fprintf(&buf, "Text: '%s'", s.Text)
	if len(s.Params) != 0 {
		fmt.Fprintf(&buf, ", Parameters: %s", strings.Join(s.Params, ", "))
	}
	return buf.String()
}
