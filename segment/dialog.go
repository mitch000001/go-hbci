package segment

import "github.com/mitch000001/go-hbci/dataelement"

const productName = "go-hbci library"
const productVersion = "0.0.1"

func NewDialogEndSegment(dialogId string) *DialogEndSegment {
	d := &DialogEndSegment{
		DialogID: dataelement.NewIdentificationDataElement(dialogId),
	}
	d.Segment = NewBasicSegment("HKEND", 3, 1, d)
	return d
}

type DialogEndSegment struct {
	Segment
	DialogID *dataelement.IdentificationDataElement
}

func (d *DialogEndSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		d.DialogID,
	}
}

func NewProcessingPreparationSegment(bdpVersion int, udpVersion int, language int) *ProcessingPreparationSegment {
	p := &ProcessingPreparationSegment{
		BPDVersion:     dataelement.NewNumberDataElement(bdpVersion, 3),
		UPDVersion:     dataelement.NewNumberDataElement(udpVersion, 3),
		DialogLanguage: dataelement.NewNumberDataElement(language, 3),
		ProductName:    dataelement.NewAlphaNumericDataElement(productName, 25),
		ProductVersion: dataelement.NewAlphaNumericDataElement(productVersion, 5),
	}
	p.Segment = NewBasicSegment("HKVVB", 4, 2, p)
	return p
}

type ProcessingPreparationSegment struct {
	Segment
	BPDVersion *dataelement.NumberDataElement
	UPDVersion *dataelement.NumberDataElement
	// 0 for undefined
	// Sprachkennzeichen | Bedeutung   | Sprachencode ISO 639 | ISO 8859 Subset | ISO 8859- Codeset
	// --------------------------------------------------------------------------------------------
	// 1				 | Deutsch	   | de (German) ￼	      | Deutsch ￼ ￼		| 1 (Latin 1)
	// 2				 | Englisch	   | en (English)		  | Englisch		| 1 (Latin 1)
	// 3 				 | Französisch | fr (French)  		  | Französisch ￼	| 1 (Latin 1)
	DialogLanguage *dataelement.NumberDataElement
	ProductName    *dataelement.AlphaNumericDataElement
	ProductVersion *dataelement.AlphaNumericDataElement
}

func (p *ProcessingPreparationSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		p.BPDVersion,
		p.UPDVersion,
		p.DialogLanguage,
		p.ProductName,
		p.ProductVersion,
	}
}

func NewBankAnnouncementSegment(subject, body string) *BankAnnouncementSegment {
	b := &BankAnnouncementSegment{
		Subject: dataelement.NewAlphaNumericDataElement(subject, 35),
		Body:    dataelement.NewTextDataElement(body, 2048),
	}
	b.Segment = NewBasicSegment("HIKIM", 8, 2, b)
	return b
}

type BankAnnouncementSegment struct {
	Segment
	Subject *dataelement.AlphaNumericDataElement
	Body    *dataelement.TextDataElement
}

func (b *BankAnnouncementSegment) elements() []dataelement.DataElement {
	return []dataelement.DataElement{
		b.Subject,
		b.Body,
	}
}
