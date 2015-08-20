package segment

import "github.com/mitch000001/go-hbci/element"

func NewBankAnnouncementSegment(subject, body string) *BankAnnouncementSegment {
	b := &BankAnnouncementSegment{
		Subject: element.NewAlphaNumeric(subject, 35),
		Body:    element.NewText(body, 2048),
	}
	b.Segment = NewBasicSegment(8, b)
	return b
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment BankAnnouncementSegment

type BankAnnouncementSegment struct {
	Segment
	Subject *element.AlphaNumericDataElement
	Body    *element.TextDataElement
}

func (b *BankAnnouncementSegment) Version() int         { return 2 }
func (b *BankAnnouncementSegment) ID() string           { return "HIKIM" }
func (b *BankAnnouncementSegment) referencedId() string { return "" }
func (b *BankAnnouncementSegment) sender() string       { return senderBank }

func (b *BankAnnouncementSegment) elements() []element.DataElement {
	return []element.DataElement{
		b.Subject,
		b.Body,
	}
}
