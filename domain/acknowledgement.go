package domain

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
