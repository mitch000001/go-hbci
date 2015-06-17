package hbci

type CommonBankParameterSegment struct {
	*basicSegment
	BPDVersion               *NumberDataElement
	BankID                   *BankIdentificationDataElement
	BankName                 *AlphaNumericDataElement
	BusinessTransactionCount *NumberDataElement
	SupportedLanguages       *SupportedLanguagesDataElement
	SupportedHBCIVersions    *SupportedHBCIVersionsDataElement
	MaxMessageSize           *NumberDataElement
}

func NewSupportedHBCIVersionsDataElement(versions ...int) *SupportedHBCIVersionsDataElement {
	versionDEs := make([]DataElement, len(versions))
	for i, version := range versions {
		versionDEs[i] = NewNumberDataElement(version, 3)
	}
	s := &SupportedHBCIVersionsDataElement{}
	s.arrayElementGroup = NewArrayElementGroup(SupportedHBCIVersionDEG, 1, 9, versionDEs...)
	return s
}

type SupportedHBCIVersionsDataElement struct {
	*arrayElementGroup
}

func NewSupportedLanguagesDataElement(languages ...int) *SupportedLanguagesDataElement {
	languageDEs := make([]DataElement, len(languages))
	for i, lang := range languages {
		languageDEs[i] = NewNumberDataElement(lang, 3)
	}
	s := &SupportedLanguagesDataElement{}
	s.arrayElementGroup = NewArrayElementGroup(SupportedLanguagesDEG, 1, 9, languageDEs...)
	return s
}

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
