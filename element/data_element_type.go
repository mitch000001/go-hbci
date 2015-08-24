package element

import "fmt"

type DataElementType int

const (
	// DataElements
	AlphaNumericDE DataElementType = iota + 1
	TextDE
	NumberDE
	DigitDE
	FloatDE
	DTAUSCharsetDE
	BinaryDE
	// Derived types
	BooleanDE
	CodeDE
	DateDE
	VirtualDateDE
	TimeDE
	IdentificationDE
	CountryCodeDE
	CurrencyDE
	ValueDE
	// Multiple used element
	AmountGDEG
	BankIdentificationGDEG
	AccountConnectionGDEG
	InternationalAccountConnectionGDEG
	BalanceGDEG
	AddressGDEG
	SecurityMethodVersionGDEG
	AcknowlegdementParamsGDEG
	PinTanBusinessTransactionParameterGDEG
	// DataElementGroups
	SegmentHeaderDEG
	ReferenceMessageDEG
	AcknowledgementDEG
	SecurityIdentificationDEG
	SecurityDateDEG
	HashAlgorithmDEG
	SignatureAlgorithmDEG
	EncryptionAlgorithmDEG
	KeyNameDEG
	CertificateDEG
	PublicKeyDEG
	SupportedLanguagesDEG
	SupportedHBCIVersionDEG
	CommunicationParameterDEG
	SupportedSecurityMethodDEG
	PinTanDEG
	AccountLimitDEG
	AllowedBusinessTransactionDEG
	DisposalEligiblePersonDEG
	SecurityProfileDEG
)

var typeName = map[DataElementType]string{
	AlphaNumericDE: "an",
	TextDE:         "txt",
	NumberDE:       "num",
	DigitDE:        "dig",
	FloatDE:        "float",
	DTAUSCharsetDE: "dta",
	BinaryDE:       "bin",
	// Derived types
	BooleanDE:        "jn",
	CodeDE:           "code",
	DateDE:           "dat",
	VirtualDateDE:    "vdat",
	TimeDE:           "tim",
	IdentificationDE: "id",
	CountryCodeDE:    "ctr",
	CurrencyDE:       "cur",
	ValueDE:          "wrt",
	// Multiple used elements or GroupDataElementGroups
	AmountGDEG:                             "btg",
	BankIdentificationGDEG:                 "kik",
	AccountConnectionGDEG:                  "ktv",
	InternationalAccountConnectionGDEG:     "kti",
	BalanceGDEG:                            "sdo",
	AddressGDEG:                            "addr",
	SecurityMethodVersionGDEG:              "Unterstützte Sicherheitsverfahren",
	AcknowlegdementParamsGDEG:              "Rückmeldungsparameter",
	PinTanBusinessTransactionParameterGDEG: "Geschäftsvorfallspezifische PIN-TAN-Informationen",
	// DataElementGroups
	SegmentHeaderDEG:              "Segmentkopf",
	ReferenceMessageDEG:           "Bezugsnachricht",
	AcknowledgementDEG:            "Rückmeldung",
	SecurityIdentificationDEG:     "Sicherheitsidentifikation, Details",
	SecurityDateDEG:               "Sicherheitsdatum und -uhrzeit",
	HashAlgorithmDEG:              "Hashalgorithmus",
	SignatureAlgorithmDEG:         "Signaturalgorithmus",
	EncryptionAlgorithmDEG:        "Verschlüsselungsalgorithmus",
	KeyNameDEG:                    "Schlüsselname",
	CertificateDEG:                "Zertifikat",
	PublicKeyDEG:                  "Öffentlicher Schlüssel",
	SupportedLanguagesDEG:         "Unterstützte Sprachen",
	SupportedHBCIVersionDEG:       "Unterstützte HBCI-Versionen",
	CommunicationParameterDEG:     "Kommunikationsparameter",
	PinTanDEG:                     "PIN-TAN",
	AccountLimitDEG:               "Kontolimit",
	AllowedBusinessTransactionDEG: "Erlaubte Geschäftsvorfälle",
	DisposalEligiblePersonDEG:     "Verfügungsberechtigte",
	SecurityProfileDEG:            "Sicherheitsprofil",
}

func (d DataElementType) String() string {
	s := typeName[d]
	if s == "" {
		return fmt.Sprintf("DataElementType%d", int(d))
	}
	return s
}
