package element

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
)

func NewSupportedSecurityMethod(methodCode string, versions ...int) *SupportedSecurityMethodDataElement {
	s := &SupportedSecurityMethodDataElement{
		MethodCode: NewAlphaNumeric(methodCode, 3),
		Versions:   NewSecurityMethodVersions(1, 9, versions...),
	}
	s.DataElement = NewDataElementGroup(SupportedSecurityMethodDEG, 2, s)
	return s
}

type SupportedSecurityMethodDataElement struct {
	DataElement
	// Code | Bedeutung
	// ------------------------------
	// DDV  | DES-DES-Verfahren
	// RDH  | RSA-DES-Hybridverfahren
	MethodCode *AlphaNumericDataElement
	// At the moment only "1" is allowed
	Versions *SecurityMethodVersionsDataElement
}

func (s *SupportedSecurityMethodDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.MethodCode,
		s.Versions,
	}
}

func (s *SupportedSecurityMethodDataElement) UnmarshalHBCI(value []byte) error {
	s.DataElement = NewDataElementGroup(SupportedSecurityMethodDEG, 2, s)
	return s.DataElement.UnmarshalHBCI(value)
}

func NewSecurityMethodVersions(min, max int, versions ...int) *SecurityMethodVersionsDataElement {
	versionDEs := make([]DataElement, len(versions))
	for i, version := range versions {
		versionDEs[i] = NewNumber(version, 3)
	}
	s := &SecurityMethodVersionsDataElement{}
	s.arrayElementGroup = NewArrayElementGroup(SecurityMethodVersionGDEG, min, max, versionDEs)
	return s
}

type SecurityMethodVersionsDataElement struct {
	*arrayElementGroup
}

func (s *SecurityMethodVersionsDataElement) Elements() []DataElement {
	return s.arrayElementGroup.array
}

func (s *SecurityMethodVersionsDataElement) UnmarshalHBCI(value []byte) error {
	dataElements := make([]DataElement, 9)
	s.arrayElementGroup = NewArrayElementGroup(SecurityMethodVersionGDEG, 1, 9, dataElements)
	return s.arrayElementGroup.UnmarshalHBCI(value)
}

func (s *SecurityMethodVersionsDataElement) Versions() []*NumberDataElement {
	versions := make([]*NumberDataElement, len(s.arrayElementGroup.array))
	for i, version := range s.arrayElementGroup.array {
		versions[i] = version.(*NumberDataElement)
	}
	return versions
}

func NewSupportedHBCIVersions(versions ...int) *SupportedHBCIVersionsDataElement {
	versionDEs := make([]DataElement, len(versions))
	for i, version := range versions {
		versionDEs[i] = NewNumber(version, 3)
	}
	s := &SupportedHBCIVersionsDataElement{}
	s.arrayElementGroup = NewArrayElementGroup(SupportedHBCIVersionDEG, 1, 9, versionDEs)
	return s
}

var validHBCIVersions = []int{201, 210, 220, 300}

type SupportedHBCIVersionsDataElement struct {
	*arrayElementGroup
}

func (s *SupportedHBCIVersionsDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 || len(elements) > 9 {
		return fmt.Errorf("Malformed marshaled value")
	}
	versions := make([]DataElement, len(elements))
	for i, elem := range elements {
		version, err := strconv.Atoi(charset.ToUTF8(elem))
		if err != nil {
			return err
		}
		versions[i] = NewNumber(version, 3)
	}
	s.arrayElementGroup = NewArrayElementGroup(SupportedHBCIVersionDEG, 1, 9, versions)
	return nil
}

func NewSupportedLanguages(languages ...int) *SupportedLanguagesDataElement {
	languageDEs := make([]DataElement, len(languages))
	for i, lang := range languages {
		languageDEs[i] = NewNumber(lang, 3)
	}
	s := &SupportedLanguagesDataElement{}
	s.arrayElementGroup = NewArrayElementGroup(SupportedLanguagesDEG, 1, 9, languageDEs)
	return s
}

var validLanguages = []int{1, 2, 3}

type SupportedLanguagesDataElement struct {
	*arrayElementGroup
}

func (s *SupportedLanguagesDataElement) Languages() []*NumberDataElement {
	languages := make([]*NumberDataElement, len(s.arrayElementGroup.array))
	for i, lang := range s.arrayElementGroup.array {
		languages[i] = lang.(*NumberDataElement)
	}
	return languages
}

func (s *SupportedLanguagesDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 || len(elements) > 9 {
		return fmt.Errorf("Malformed marshaled value")
	}
	languages := make([]DataElement, len(elements))
	for i, elem := range elements {
		lang, err := strconv.Atoi(charset.ToUTF8(elem))
		if err != nil {
			return err
		}
		if sort.SearchInts(validLanguages, lang) >= len(validLanguages) {
			return fmt.Errorf("Unsupported language code: %d", lang)
		}
		languages[i] = NewNumber(lang, 3)
	}
	s.arrayElementGroup = NewArrayElementGroup(SupportedLanguagesDEG, 1, 9, languages)
	return nil
}

type SupportedCompressionMethodsDataElement struct {
	*arrayElementGroup
}

type BusinessTransactionParameter struct {
	DataElement
	DataElements []DataElement
}

func (b *BusinessTransactionParameter) GroupDataElements() []DataElement {
	return b.DataElements
}

func NewPinTanBusinessTransactionParameters(pinTanTransactions []domain.PinTanBusinessTransaction) *PinTanBusinessTransactionParameters {
	transactionsDEs := make([]DataElement, len(pinTanTransactions))
	for i, transaction := range pinTanTransactions {
		pinTanBusinessTransaction := &PinTanBusinessTransactionParameter{
			SegmentID: NewAlphaNumeric(transaction.SegmentID, 6),
			NeedsTAN:  NewBoolean(transaction.NeedsTan),
		}
		pinTanBusinessTransaction.DataElement = NewGroupDataElementGroup(PinTanBusinessTransactionParameterGDEG, 2, pinTanBusinessTransaction)
		transactionsDEs[i] = pinTanBusinessTransaction
	}
	p := &PinTanBusinessTransactionParameters{}
	p.arrayElementGroup = NewArrayElementGroup(PinTanBusinessTransactionParameterGDEG, len(transactionsDEs), len(transactionsDEs), transactionsDEs)
	return p
}

type PinTanBusinessTransactionParameters struct {
	*arrayElementGroup
}

func (p *PinTanBusinessTransactionParameters) Val() []domain.PinTanBusinessTransaction {
	transactions := make([]domain.PinTanBusinessTransaction, len(p.array))
	for i, elem := range p.array {
		transactions[i] = elem.(*PinTanBusinessTransactionParameter).Val()
	}
	return transactions
}

func (p *PinTanBusinessTransactionParameters) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements)%2 != 0 {
		return fmt.Errorf("Malformed marshaled value: value pairs not even")
	}
	dataElements := make([]DataElement, len(elements)/2)
	for i := 0; i < len(elements); i += 2 {
		elem := bytes.Join(elements[i:i+2], []byte(":"))
		pinTanTransaction := &PinTanBusinessTransactionParameter{}
		err := pinTanTransaction.UnmarshalHBCI(elem)
		if err != nil {
			return err
		}
		dataElements[i/2] = pinTanTransaction
	}
	p.arrayElementGroup = NewArrayElementGroup(PinTanBusinessTransactionParameterGDEG, len(dataElements), len(dataElements), dataElements)
	return nil
}

type PinTanBusinessTransactionParameter struct {
	DataElement
	SegmentID *AlphaNumericDataElement
	NeedsTAN  *BooleanDataElement
}

func (p *PinTanBusinessTransactionParameter) Elements() []DataElement {
	return []DataElement{
		p.SegmentID,
		p.NeedsTAN,
	}
}

func (p *PinTanBusinessTransactionParameter) Val() domain.PinTanBusinessTransaction {
	return domain.PinTanBusinessTransaction{
		SegmentID: p.SegmentID.Val(),
		NeedsTan:  p.NeedsTAN.Val(),
	}
}

func (p *PinTanBusinessTransactionParameter) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	p.SegmentID = NewAlphaNumeric(charset.ToUTF8(elements[0]), 6)
	needsTan := &BooleanDataElement{}
	err = needsTan.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	p.NeedsTAN = needsTan
	p.DataElement = NewGroupDataElementGroup(PinTanBusinessTransactionParameterGDEG, 2, p)
	return nil
}
