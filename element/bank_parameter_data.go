package element

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

type SupportedHBCIVersionsDataElement struct {
	*arrayElementGroup
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

type SupportedCompressionMethodsDataElement struct {
	*arrayElementGroup
}
