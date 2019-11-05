package segment

import (
	"github.com/mitch000001/go-hbci"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const defaultProductName = "go-hbci library"
const defaultProductVersion = hbci.Version

func NewDialogEndSegment(dialogId string) *DialogEndSegment {
	d := &DialogEndSegment{
		DialogID: element.NewIdentification(dialogId),
	}
	d.ClientSegment = NewBasicSegment(3, d)
	return d
}

type DialogEndSegment struct {
	ClientSegment
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

func NewProcessingPreparationSegment(bdpVersion int, udpVersion int, language domain.Language, productName string) *ProcessingPreparationSegment {
	if productName == "" {
		productName = defaultProductName
	}
	p := &ProcessingPreparationSegment{
		BPDVersion:     element.NewNumber(bdpVersion, 3),
		UPDVersion:     element.NewNumber(udpVersion, 3),
		DialogLanguage: element.NewNumber(int(language), 3),
		ProductName:    element.NewAlphaNumeric(productName, 25),
		ProductVersion: element.NewAlphaNumeric(defaultProductVersion, 5),
	}
	p.ClientSegment = NewBasicSegment(4, p)
	return p
}

type ProcessingPreparationSegment struct {
	ClientSegment
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
