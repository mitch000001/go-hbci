package domain

import (
	"fmt"
	"strings"
)

func NewAcknowledgement(code int, referenceDataElement, text string, params []string) Acknowledgement {
	return Acknowledgement{
		Code:                 code,
		ReferenceDataElement: referenceDataElement,
		Text:                 text,
		Params:               params,
	}
}

type Acknowledgement struct {
	Code                 int
	ReferenceDataElement string
	Text                 string
	Params               []string
}

func (a Acknowledgement) IsError() bool {
	return a.Code >= 9000
}

func (a Acknowledgement) IsWarning() bool {
	return a.Code >= 3000 && a.Code < 4000
}
