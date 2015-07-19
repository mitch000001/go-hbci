package segment

import "github.com/mitch000001/go-hbci/element"

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

func (d *DialogEndSegment) version() int         { return 1 }
func (d *DialogEndSegment) id() string           { return "HKEND" }
func (d *DialogEndSegment) referencedId() string { return "" }
func (d *DialogEndSegment) sender() string       { return senderUser }

func (d *DialogEndSegment) elements() []element.DataElement {
	return []element.DataElement{
		d.DialogID,
	}
}

func NewProcessingPreparationSegment(bdpVersion int, udpVersion int, language int) *ProcessingPreparationSegment {
	p := &ProcessingPreparationSegment{
		BPDVersion:     element.NewNumber(bdpVersion, 3),
		UPDVersion:     element.NewNumber(udpVersion, 3),
		DialogLanguage: element.NewNumber(language, 3),
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

func (p *ProcessingPreparationSegment) version() int         { return 2 }
func (p *ProcessingPreparationSegment) id() string           { return "HKVVB" }
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

func (b *BankAnnouncementSegment) version() int         { return 2 }
func (b *BankAnnouncementSegment) id() string           { return "HIKIM" }
func (b *BankAnnouncementSegment) referencedId() string { return "" }
func (b *BankAnnouncementSegment) sender() string       { return senderBank }

func (b *BankAnnouncementSegment) elements() []element.DataElement {
	return []element.DataElement{
		b.Subject,
		b.Body,
	}
}
