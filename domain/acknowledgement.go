package domain

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// NewMessageAcknowledgement creactes a new message acknowledgement
func NewMessageAcknowledgement(code int, referenceDataElement, text string, params []string) Acknowledgement {
	return Acknowledgement{
		Type:                         MessageAcknowledgement,
		Code:                         code,
		ReferenceDataElementPosition: referenceDataElement,
		Text:                         text,
		Params:                       params,
	}
}

// NewSegmentAcknowledgement creactes a new segment acknowledgement
func NewSegmentAcknowledgement(code int, referenceDataElement, text string, params []string) Acknowledgement {
	return Acknowledgement{
		Type:                         SegmentAcknowledgement,
		Code:                         code,
		ReferenceDataElementPosition: referenceDataElement,
		Text:                         text,
		Params:                       params,
	}
}

const (
	// MessageAcknowledgement defines the message ack type
	MessageAcknowledgement = "MessageAcknowledgement"
	// SegmentAcknowledgement defines the segemtn ack type
	SegmentAcknowledgement = "SegmentAcknowledgement"
)

// Acknowledgement represents an acknowledgement from the bank institute
type Acknowledgement struct {
	Type                         string
	Code                         int
	ReferenceDataElementPosition string
	ReferencingDataElement       string
	Text                         string
	Params                       []string
	ReferencingMessage           MessageReference
	ReferencingSegmentNumber     int
	ReferencingSegmentID         string
}

func (a Acknowledgement) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s for message %d (%s)", a.Type, a.ReferencingMessage.MessageNumber, a.ReferencingMessage.DialogID)
	if a.ReferencingSegmentNumber > 0 {
		fmt.Fprintf(&buf, ", segment %d", a.ReferencingSegmentNumber)
	}
	if a.ReferencingSegmentID != "" {
		fmt.Fprintf(&buf, " (%s)", a.ReferencingSegmentID)
	}
	fmt.Fprintf(&buf, ": ")
	fmt.Fprintf(&buf, "Code: %d, ", a.Code)
	if a.ReferenceDataElementPosition != "" {
		fmt.Fprintf(&buf, "Position: %s", a.ReferenceDataElementPosition)
	} else {
		fmt.Fprintf(&buf, "Position: none")
	}
	if a.ReferencingDataElement != "" {
		fmt.Fprintf(&buf, " (%s)", a.ReferencingDataElement)
	}
	fmt.Fprint(&buf, ", ")
	fmt.Fprintf(&buf, "Text: '%s'", a.Text)
	if len(a.Params) != 0 {
		fmt.Fprintf(&buf, ", Parameters: %q", strings.Join(a.Params, ", "))
	}
	return buf.String()
}

// IsMessageAcknowledgement returns true if the type is MessageAcknowledgement, false otherwise
func (a Acknowledgement) IsMessageAcknowledgement() bool {
	return a.Type == MessageAcknowledgement
}

// IsSegmentAcknowledgement returns true if type is SegmentAcknowledgement, false otherwise
func (a Acknowledgement) IsSegmentAcknowledgement() bool {
	return a.Type == SegmentAcknowledgement
}

// IsError returns true if the acknowledgement represents an error
func (a Acknowledgement) IsError() bool {
	return a.Code >= 9000
}

// IsWarning returns true if the acknowledgement represents a warning
func (a Acknowledgement) IsWarning() bool {
	return a.Code >= 3000 && a.Code < 4000
}

// IsSuccess returns true if the acknowledgement represents a success
func (a Acknowledgement) IsSuccess() bool {
	return a.Code > 0 && a.Code < 1000
}

// StatusAcknowledgement represents an Acknowledgement with a transmission date
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
	if s.ReferenceDataElementPosition != "" {
		fmt.Fprintf(&buf, "Position: %s, ", s.ReferenceDataElementPosition)
	} else {
		fmt.Fprintf(&buf, "Position: none, ")
	}
	fmt.Fprintf(&buf, "Text: '%s'", s.Text)
	if len(s.Params) != 0 {
		fmt.Fprintf(&buf, ", Parameters: %s", strings.Join(s.Params, ", "))
	}
	return buf.String()
}
