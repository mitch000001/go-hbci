package domain

import (
	"fmt"
	"strings"
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
	Type                 string
	Code                 int
	ReferenceDataElement string
	Text                 string
	Params               []string
}

func (a Acknowledgement) String() string {
	return fmt.Sprintf("%s: Code: %d, Position: %s, Text: %s, Parameter: %s", a.Type, a.Code, a.ReferenceDataElement, a.Text, strings.Join(a.Params, ", "))
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
