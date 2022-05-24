package element

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
)

// NewSupportedSecurityMethod returns a new SupportedSecurityMethodDataElement
func NewSupportedSecurityMethod(methodCode string, versions ...int) *SupportedSecurityMethodDataElement {
	s := &SupportedSecurityMethodDataElement{
		MethodCode: NewAlphaNumeric(methodCode, 3),
		Versions:   NewSecurityMethodVersions(1, 9, versions...),
	}
	s.DataElement = NewDataElementGroup(supportedSecurityMethodDEG, 2, s)
	return s
}

// SupportedSecurityMethodDataElement defines a DataElement for supported
// security methods
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

// GroupDataElements returns the grouped DataElements within the
// SupportedSecurityMethodDataElement
func (s *SupportedSecurityMethodDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		s.MethodCode,
		s.Versions,
	}
}

// UnmarshalHBCI unmarshals the value to a SupportedSecurityMethodDataElement
func (s *SupportedSecurityMethodDataElement) UnmarshalHBCI(value []byte) error {
	s.DataElement = NewDataElementGroup(supportedSecurityMethodDEG, 2, s)
	return s.DataElement.UnmarshalHBCI(value)
}

// NewSecurityMethodVersions returns a new SecurityMethodVersionsDataElement
func NewSecurityMethodVersions(min, max int, versions ...int) *SecurityMethodVersionsDataElement {
	versionDEs := make([]DataElement, len(versions))
	for i, version := range versions {
		versionDEs[i] = NewNumber(version, 3)
	}
	s := &SecurityMethodVersionsDataElement{}
	s.arrayElementGroup = newArrayElementGroup(securityMethodVersionGDEG, min, max, versionDEs)
	return s
}

// SecurityMethodVersionsDataElement represents the possible versions of a
// security method.
type SecurityMethodVersionsDataElement struct {
	*arrayElementGroup
}

// Elements returns the elements within the SecurityMethodVersionsDataElement
func (s *SecurityMethodVersionsDataElement) Elements() []DataElement {
	return s.arrayElementGroup.array
}

// UnmarshalHBCI unmarshals the value to a SecurityMethodVersionsDataElement
func (s *SecurityMethodVersionsDataElement) UnmarshalHBCI(value []byte) error {
	dataElements := make([]DataElement, 9)
	s.arrayElementGroup = newArrayElementGroup(securityMethodVersionGDEG, 1, 9, dataElements)
	return s.arrayElementGroup.UnmarshalHBCI(value)
}

// Versions returns the possible versions packaged in a slice of
// NumberDataElements
func (s *SecurityMethodVersionsDataElement) Versions() []*NumberDataElement {
	versions := make([]*NumberDataElement, len(s.arrayElementGroup.array))
	for i, version := range s.arrayElementGroup.array {
		versions[i] = version.(*NumberDataElement)
	}
	return versions
}

// NewSupportedHBCIVersions returns a new SupportedHBCIVersionsDataElement
func NewSupportedHBCIVersions(versions ...int) *SupportedHBCIVersionsDataElement {
	versionDEs := make([]DataElement, len(versions))
	for i, version := range versions {
		versionDEs[i] = NewNumber(version, 3)
	}
	s := &SupportedHBCIVersionsDataElement{}
	s.arrayElementGroup = newArrayElementGroup(supportedHBCIVersionDEG, 1, 9, versionDEs)
	return s
}

var validHBCIVersions = []int{201, 210, 220, 300}

// SupportedHBCIVersionsDataElement represents a DataElement for supported HBCI
// versions
type SupportedHBCIVersionsDataElement struct {
	*arrayElementGroup
}

// UnmarshalHBCI unmarshals the value to a SupportedHBCIVersionsDataElement
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
	s.arrayElementGroup = newArrayElementGroup(supportedHBCIVersionDEG, 1, 9, versions)
	return nil
}

// NewSupportedLanguages returns a new SupportedLanguagesDataElement
func NewSupportedLanguages(languages ...int) *SupportedLanguagesDataElement {
	languageDEs := make([]DataElement, len(languages))
	for i, lang := range languages {
		languageDEs[i] = NewNumber(lang, 3)
	}
	s := &SupportedLanguagesDataElement{}
	s.arrayElementGroup = newArrayElementGroup(supportedLanguagesDEG, 1, 9, languageDEs)
	return s
}

var validLanguages = []int{1, 2, 3}

// SupportedLanguagesDataElement represents the supported languages by the bank
// institute
type SupportedLanguagesDataElement struct {
	*arrayElementGroup
}

// Languages returns the supported languages packaged within a slice of
// NumberDataElements
func (s *SupportedLanguagesDataElement) Languages() []*NumberDataElement {
	languages := make([]*NumberDataElement, len(s.arrayElementGroup.array))
	for i, lang := range s.arrayElementGroup.array {
		languages[i] = lang.(*NumberDataElement)
	}
	return languages
}

// UnmarshalHBCI unmarshals the value into the SupportedLanguagesDataElement
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
	s.arrayElementGroup = newArrayElementGroup(supportedLanguagesDEG, 1, 9, languages)
	return nil
}

// SupportedCompressionMethodsDataElement represents the compression methods
// supported by the bank institute
type SupportedCompressionMethodsDataElement struct {
	*arrayElementGroup
}

// A BusinessTransactionParameter defines parameters for a specific business
// transaction.
type BusinessTransactionParameter struct {
	DataElement
	DataElements []DataElement
}

// GroupDataElements returns the grouped DataElements.
func (b *BusinessTransactionParameter) GroupDataElements() []DataElement {
	return b.DataElements
}

// NewPinTanBusinessTransactionParameters returns a new
// PinTanBusinessTransactionParameters DataElement
func NewPinTanBusinessTransactionParameters(pinTanTransactions []domain.PinTanBusinessTransaction) *PinTanBusinessTransactionParameters {
	transactionsDEs := make([]DataElement, len(pinTanTransactions))
	for i, transaction := range pinTanTransactions {
		pinTanBusinessTransaction := &PinTanBusinessTransactionParameter{
			SegmentID: NewAlphaNumeric(transaction.SegmentID, 6),
			NeedsTAN:  NewBoolean(transaction.NeedsTan),
		}
		pinTanBusinessTransaction.DataElement = NewGroupDataElementGroup(pinTanBusinessTransactionParameterGDEG, 2, pinTanBusinessTransaction)
		transactionsDEs[i] = pinTanBusinessTransaction
	}
	p := &PinTanBusinessTransactionParameters{}
	p.arrayElementGroup = newArrayElementGroup(pinTanBusinessTransactionParameterGDEG, len(transactionsDEs), len(transactionsDEs), transactionsDEs)
	return p
}

// PinTanBusinessTransactionParameters represents a slice of
// PinTanBusinessTransactionParameter DataElements
type PinTanBusinessTransactionParameters struct {
	*arrayElementGroup
}

// Val returns the underlying PinTanBusinessTransactions
func (p *PinTanBusinessTransactionParameters) Val() []domain.PinTanBusinessTransaction {
	transactions := make([]domain.PinTanBusinessTransaction, len(p.array))
	for i, elem := range p.array {
		transactions[i] = elem.(*PinTanBusinessTransactionParameter).Val()
	}
	return transactions
}

// UnmarshalHBCI unmarshals value into the PinTanBusinessTransactionParameters
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
	p.arrayElementGroup = newArrayElementGroup(pinTanBusinessTransactionParameterGDEG, len(dataElements), len(dataElements), dataElements)
	return nil
}

// PinTanBusinessTransactionParameter defines a specific
// PinTanBusinessTransactionParameter DataElement
type PinTanBusinessTransactionParameter struct {
	DataElement `yaml:"-"`
	SegmentID   *AlphaNumericDataElement `yaml:"segmentID"`
	NeedsTAN    *BooleanDataElement      `yaml:"needsTan"`
}

// Elements returns the elements of this DataElement.
func (p *PinTanBusinessTransactionParameter) Elements() []DataElement {
	return []DataElement{
		p.SegmentID,
		p.NeedsTAN,
	}
}

// Val returns the underlying PinTanBusinessTransaction
func (p *PinTanBusinessTransactionParameter) Val() domain.PinTanBusinessTransaction {
	return domain.PinTanBusinessTransaction{
		SegmentID: p.SegmentID.Val(),
		NeedsTan:  p.NeedsTAN.Val(),
	}
}

// UnmarshalHBCI unmarshals value into p
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
	p.DataElement = NewGroupDataElementGroup(pinTanBusinessTransactionParameterGDEG, 2, p)
	return nil
}
