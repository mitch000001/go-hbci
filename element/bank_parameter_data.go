package element

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
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

func NewSecurityMethodVersions(min, max int, versions ...int) *SecurityMethodVersionsDataElement {
	versionDEs := make([]DataElement, len(versions))
	for i, version := range versions {
		versionDEs[i] = NewNumber(version, 3)
	}
	s := &SecurityMethodVersionsDataElement{}
	s.arrayElementGroup = NewArrayElementGroup(SecurityMethodVersionGDEG, min, max, versionDEs...)
	return s
}

type SecurityMethodVersionsDataElement struct {
	*arrayElementGroup
}

func (s *SecurityMethodVersionsDataElement) Elements() []DataElement {
	return s.arrayElementGroup.array
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
	s.arrayElementGroup = NewArrayElementGroup(SupportedHBCIVersionDEG, 1, 9, versionDEs...)
	return s
}

var validHBCIVersions = []int{201, 210, 220}

type SupportedHBCIVersionsDataElement struct {
	*arrayElementGroup
}

func (s *SupportedHBCIVersionsDataElement) UnmarshalHBCI(value []byte) error {
	elements := bytes.Split(value, []byte(":"))
	if len(elements) == 0 || len(elements) > 9 {
		return fmt.Errorf("Malformed marshaled value")
	}
	versions := make([]DataElement, len(elements))
	for i, elem := range elements {
		version, err := strconv.Atoi(string(elem))
		if err != nil {
			return err
		}
		if sort.SearchInts(validHBCIVersions, version) >= len(validHBCIVersions) {
			return fmt.Errorf("Unsupported HBCI version: %d", version)
		}
		versions[i] = NewNumber(version, 3)
	}
	s.arrayElementGroup = NewArrayElementGroup(SupportedHBCIVersionDEG, 1, 9, versions...)
	return nil
}

func NewSupportedLanguages(languages ...int) *SupportedLanguagesDataElement {
	languageDEs := make([]DataElement, len(languages))
	for i, lang := range languages {
		languageDEs[i] = NewNumber(lang, 3)
	}
	s := &SupportedLanguagesDataElement{}
	s.arrayElementGroup = NewArrayElementGroup(SupportedLanguagesDEG, 1, 9, languageDEs...)
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
	elements := bytes.Split(value, []byte(":"))
	if len(elements) == 0 || len(elements) > 9 {
		return fmt.Errorf("Malformed marshaled value")
	}
	languages := make([]DataElement, len(elements))
	for i, elem := range elements {
		lang, err := strconv.Atoi(string(elem))
		if err != nil {
			return err
		}
		if sort.SearchInts(validLanguages, lang) >= len(validLanguages) {
			return fmt.Errorf("Unsupported language code: %d", lang)
		}
		languages[i] = NewNumber(lang, 3)
	}
	s.arrayElementGroup = NewArrayElementGroup(SupportedLanguagesDEG, 1, 9, languages...)
	return nil
}

type SupportedCompressionMethodsDataElement struct {
	*arrayElementGroup
}

type BusinessTransactionParameter struct {
	*elementGroup
	DataElements []DataElement
}

func (b *BusinessTransactionParameter) GroupDataElements() []DataElement {
	return b.DataElements
}
