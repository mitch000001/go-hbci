package segment

import (
	"github.com/mitch000001/go-hbci"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const productName = "5A624F86A785F4024DD914404"
const productVersion = hbci.Version

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

const ProcessingPreparationID = "HKVVB"

func NewProcessingPreparationSegmentV2(bdpVersion int, udpVersion int, language domain.Language) *ProcessingPreparationSegmentV2 {
	p := &ProcessingPreparationSegmentV2{
		BPDVersion:     element.NewNumber(bdpVersion, 3),
		UPDVersion:     element.NewNumber(udpVersion, 3),
		DialogLanguage: element.NewNumber(int(language), 3),
		ProductName:    element.NewAlphaNumeric(productName, 25),
		ProductVersion: element.NewAlphaNumeric(productVersion, 5),
	}
	p.ClientSegment = NewBasicSegment(4, p)
	return p
}

func NewProcessingPreparationSegmentV3(bdpVersion int, udpVersion int, language domain.Language) *ProcessingPreparationSegmentV3 {
	p := &ProcessingPreparationSegmentV3{
		BPDVersion:     element.NewNumber(bdpVersion, 3),
		UPDVersion:     element.NewNumber(udpVersion, 3),
		DialogLanguage: element.NewNumber(int(language), 3),
		ProductName:    element.NewAlphaNumeric(productName, 25),
		ProductVersion: element.NewAlphaNumeric(productVersion, 5),
	}
	p.ClientSegment = NewBasicSegment(4, p)
	return p
}

type ProcessingPreparationSegmentV2 struct {
	ClientSegment
	BPDVersion *element.NumberDataElement
	UPDVersion *element.NumberDataElement
	// 0 for undefined
	// Sprachkennzeichen | Bedeutung   | Sprachencode ISO 639 | ISO 8859 Subset | ISO 8859- Codeset
	// --------------------------------------------------------------------------------------------
	// 1				 | Deutsch	   | de (German) ￼	      | Deutsch ￼ ￼		| 1 (Latin 1)
	// 2				 | Englisch	   | en (English)		  | Englisch		| 1 (Latin 1)
	// 3 				 | Französisch | fr (French)  		   | Französisch ￼	  | 1 (Latin 1)
	DialogLanguage *element.NumberDataElement
	ProductName    *element.AlphaNumericDataElement
	ProductVersion *element.AlphaNumericDataElement
}

func (p *ProcessingPreparationSegmentV2) Version() int         { return 2 }
func (p *ProcessingPreparationSegmentV2) ID() string           { return ProcessingPreparationID }
func (p *ProcessingPreparationSegmentV2) referencedId() string { return "" }
func (p *ProcessingPreparationSegmentV2) sender() string       { return senderUser }

func (p *ProcessingPreparationSegmentV2) elements() []element.DataElement {
	return []element.DataElement{
		p.BPDVersion,
		p.UPDVersion,
		p.DialogLanguage,
		p.ProductName,
		p.ProductVersion,
	}
}

type ProcessingPreparationSegmentV3 struct {
	ClientSegment
	BPDVersion *element.NumberDataElement
	UPDVersion *element.NumberDataElement
	// 0 for Standard = Institute language
	// Sprachkennzeichen | Bedeutung   | Sprachencode ISO 639 | ISO 8859 Subset | ISO 8859- Codeset
	// --------------------------------------------------------------------------------------------
	// 1				 | Deutsch	   | de (German) ￼	      | Deutsch ￼ ￼		| 1 (Latin 1)
	// 2				 | Englisch	   | en (English)		  | Englisch		| 1 (Latin 1)
	// 3 				 | Französisch | fr (French)  		   | Französisch ￼	  | 1 (Latin 1)
	DialogLanguage *element.NumberDataElement
	ProductName    *element.AlphaNumericDataElement
	ProductVersion *element.AlphaNumericDataElement
}

func (p *ProcessingPreparationSegmentV3) Version() int         { return 3 }
func (p *ProcessingPreparationSegmentV3) ID() string           { return ProcessingPreparationID }
func (p *ProcessingPreparationSegmentV3) referencedId() string { return "" }
func (p *ProcessingPreparationSegmentV3) sender() string       { return senderUser }

func (p *ProcessingPreparationSegmentV3) elements() []element.DataElement {
	return []element.DataElement{
		p.BPDVersion,
		p.UPDVersion,
		p.DialogLanguage,
		p.ProductName,
		p.ProductVersion,
	}
}
