package hbci

var HKVVBSegmentNumber = -1

func NewCommonBankParameterSegment(
	bpdVersion int,
	countryCode int,
	bankId string,
	bankName string,
	businessTransactionCount int,
	supportedLanguages []int,
	supportedHBCIVersions []int,
	maxMessageSize int) *CommonBankParameterSegment {
	c := &CommonBankParameterSegment{
		BPDVersion:               NewNumberDataElement(bpdVersion, 3),
		BankID:                   NewBankIndentificationDataElementWithBankId(countryCode, bankId),
		BankName:                 NewAlphaNumericDataElement(bankName, 60),
		BusinessTransactionCount: NewNumberDataElement(businessTransactionCount, 3),
		SupportedLanguages:       NewSupportedLanguagesDataElement(supportedLanguages...),
		SupportedHBCIVersions:    NewSupportedHBCIVersionsDataElement(supportedHBCIVersions...),
		MaxMessageSize:           NewNumberDataElement(maxMessageSize, 4),
	}
	header := NewReferencingSegmentHeader("HIBPA", 1, 2, HKVVBSegmentNumber)
	c.basicSegment = NewBasicSegmentWithHeader(header, c)
	return c
}

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

func (c *CommonBankParameterSegment) elements() []DataElement {
	return []DataElement{
		c.BPDVersion,
		c.BankID,
		c.BankName,
		c.BusinessTransactionCount,
		c.SupportedLanguages,
		c.SupportedHBCIVersions,
		c.MaxMessageSize,
	}
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
