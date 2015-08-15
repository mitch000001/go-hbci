package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const productName = "go-hbci library"
const productVersion = "0.0.1"

func NewDialogEndSegment(dialogId string) *DialogEndSegment {
	d := &DialogEndSegment{
		DialogID: element.NewIdentification(dialogId),
	}
	d.Segment = NewBasicSegment(3, d)
	return d
}

type DialogEndSegment struct {
	Segment
	DialogID *element.IdentificationDataElement
}

func (d *DialogEndSegment) Version() int         { return 1 }
func (d *DialogEndSegment) ID() string           { return "HKEND" }
func (d *DialogEndSegment) referencedId() string { return "" }
func (d *DialogEndSegment) sender() string       { return senderUser }

func (d *DialogEndSegment) elements() []element.DataElement {
	return []element.DataElement{
		d.DialogID,
	}
}

func NewProcessingPreparationSegment(bdpVersion int, udpVersion int, language domain.Language) *ProcessingPreparationSegment {
	p := &ProcessingPreparationSegment{
		BPDVersion:     element.NewNumber(bdpVersion, 3),
		UPDVersion:     element.NewNumber(udpVersion, 3),
		DialogLanguage: element.NewNumber(int(language), 3),
		ProductName:    element.NewAlphaNumeric(productName, 25),
		ProductVersion: element.NewAlphaNumeric(productVersion, 5),
	}
	p.Segment = NewBasicSegment(4, p)
	return p
}

type ProcessingPreparationSegment struct {
	Segment
	BPDVersion *element.NumberDataElement
	UPDVersion *element.NumberDataElement
	// 0 for undefined
	// Sprachkennzeichen | Bedeutung   | Sprachencode ISO 639 | ISO 8859 Subset | ISO 8859- Codeset
	// --------------------------------------------------------------------------------------------
	// 1				 | Deutsch	   | de (German) ￼	      | Deutsch ￼ ￼		| 1 (Latin 1)
	// 2				 | Englisch	   | en (English)		  | Englisch		| 1 (Latin 1)
	// 3 				 | Französisch | fr (French)  		  | Französisch ￼	| 1 (Latin 1)
	DialogLanguage *element.NumberDataElement
	ProductName    *element.AlphaNumericDataElement
	ProductVersion *element.AlphaNumericDataElement
}

func (p *ProcessingPreparationSegment) Version() int         { return 2 }
func (p *ProcessingPreparationSegment) ID() string           { return "HKVVB" }
func (p *ProcessingPreparationSegment) referencedId() string { return "" }
func (p *ProcessingPreparationSegment) sender() string       { return senderUser }

func (p *ProcessingPreparationSegment) elements() []element.DataElement {
	return []element.DataElement{
		p.BPDVersion,
		p.UPDVersion,
		p.DialogLanguage,
		p.ProductName,
		p.ProductVersion,
	}
}

func NewBankAnnouncementSegment(subject, body string) *BankAnnouncementSegment {
	b := &BankAnnouncementSegment{
		Subject: element.NewAlphaNumeric(subject, 35),
		Body:    element.NewText(body, 2048),
	}
	b.Segment = NewBasicSegment(8, b)
	return b
}

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

func (b *BankAnnouncementSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 3 {
		return fmt.Errorf("Malformed marshaled value")
	}
	segment, err := SegmentFromHeaderBytes(elements[0], b)
	if err != nil {
		return err
	}
	b.Segment = segment
	b.Subject = element.NewAlphaNumeric(charset.ToUtf8(elements[1]), 35)
	b.Body = element.NewText(charset.ToUtf8(elements[2]), 2048)
	return nil
}
