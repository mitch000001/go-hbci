package element

import "fmt"

// A DataElementType represents the type of a DataElement
type DataElementType int

const (
	// DataElements

	alphaNumericDE DataElementType = iota + 1
	textDE
	numberDE
	digitDE
	floatDE
	dtausCharsetDE
	binaryDE

	// Derived types

	booleanDE
	codeDE
	dateDE
	virtualDateDE
	timeDE
	identificationDE
	countryCodeDE
	currencyDE
	valueDE

	// Multiple used element

	amountGDEG
	bankIdentificationGDEG
	accountConnectionGDEG
	internationalAccountConnectionGDEG
	balanceGDEG
	addressGDEG
	securityMethodVersionGDEG
	acknowlegdementParamsGDEG
	pinTanBusinessTransactionParameterGDEG

	// DataElementGroups

	segmentHeaderDEG
	referenceMessageDEG
	acknowledgementDEG
	securityIdentificationDEG
	securityDateDEG
	hashAlgorithmDEG
	signatureAlgorithmDEG
	encryptionAlgorithmDEG
	keyNameDEG
	certificateDEG
	publicKeyDEG
	supportedLanguagesDEG
	supportedHBCIVersionDEG
	communicationParameterDEG
	supportedSecurityMethodDEG
	pinTanDEG
	accountLimitDEG
	allowedBusinessTransactionDEG
	disposalEligiblePersonDEG
	securityProfileDEG
	tan2StepSubmissionParameterDEG
	tan2StepSubmissionProcessParameterDEG
)

var typeName = map[DataElementType]string{
	alphaNumericDE: "an",
	textDE:         "txt",
	numberDE:       "num",
	digitDE:        "dig",
	floatDE:        "float",
	dtausCharsetDE: "dta",
	binaryDE:       "bin",
	// Derived types
	booleanDE:        "jn",
	codeDE:           "code",
	dateDE:           "dat",
	virtualDateDE:    "vdat",
	timeDE:           "tim",
	identificationDE: "id",
	countryCodeDE:    "ctr",
	currencyDE:       "cur",
	valueDE:          "wrt",
	// Multiple used elements or GroupDataElementGroups
	amountGDEG:                             "btg",
	bankIdentificationGDEG:                 "kik",
	accountConnectionGDEG:                  "ktv",
	internationalAccountConnectionGDEG:     "kti",
	balanceGDEG:                            "sdo",
	addressGDEG:                            "addr",
	securityMethodVersionGDEG:              "Unterstützte Sicherheitsverfahren",
	acknowlegdementParamsGDEG:              "Rückmeldungsparameter",
	pinTanBusinessTransactionParameterGDEG: "Geschäftsvorfallspezifische PIN-TAN-Informationen",
	// DataElementGroups
	segmentHeaderDEG:                      "Segmentkopf",
	referenceMessageDEG:                   "Bezugsnachricht",
	acknowledgementDEG:                    "Rückmeldung",
	securityIdentificationDEG:             "Sicherheitsidentifikation, Details",
	securityDateDEG:                       "Sicherheitsdatum und -uhrzeit",
	hashAlgorithmDEG:                      "Hashalgorithmus",
	signatureAlgorithmDEG:                 "Signaturalgorithmus",
	encryptionAlgorithmDEG:                "Verschlüsselungsalgorithmus",
	keyNameDEG:                            "Schlüsselname",
	certificateDEG:                        "Zertifikat",
	publicKeyDEG:                          "Öffentlicher Schlüssel",
	supportedLanguagesDEG:                 "Unterstützte Sprachen",
	supportedHBCIVersionDEG:               "Unterstützte HBCI-Versionen",
	communicationParameterDEG:             "Kommunikationsparameter",
	pinTanDEG:                             "PIN-TAN",
	accountLimitDEG:                       "Kontolimit",
	allowedBusinessTransactionDEG:         "Erlaubte Geschäftsvorfälle",
	disposalEligiblePersonDEG:             "Verfügungsberechtigte",
	securityProfileDEG:                    "Sicherheitsprofil",
	tan2StepSubmissionParameterDEG:        "Parameter Zwei-Schritt-TAN-Einreichung",
	tan2StepSubmissionProcessParameterDEG: "Verfahrensparameter Zwei-Schritt-Verfahren",
}

func (d DataElementType) String() string {
	s := typeName[d]
	if s == "" {
		return fmt.Sprintf("DataElementType%d", int(d))
	}
	return s
}
