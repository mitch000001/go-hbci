package dataelement

func NewPinTanEncryptionAlgorithm() *EncryptionAlgorithmDataElement {
	e := &EncryptionAlgorithmDataElement{
		Usage:                      NewAlphaNumeric("2", 3),
		OperationMode:              NewAlphaNumeric("2", 3),
		Algorithm:                  NewAlphaNumeric("13", 3),
		Key:                        NewBinary([]byte(defaultPinTan), 512),
		KeyParamID:                 NewAlphaNumeric("5", 3),
		InitializationValueParamID: NewAlphaNumeric("1", 3),
	}
	e.DataElement = NewDataElementGroup(EncryptionAlgorithmDEG, 7, e)
	return e
}

func NewRDHEncryptionAlgorithm(pubKey []byte) *EncryptionAlgorithmDataElement {
	e := &EncryptionAlgorithmDataElement{
		Usage:                      NewAlphaNumeric("2", 3),
		OperationMode:              NewAlphaNumeric("2", 3),
		Algorithm:                  NewAlphaNumeric("13", 3),
		Key:                        NewBinary(pubKey, 512),
		KeyParamID:                 NewAlphaNumeric("6", 3),
		InitializationValueParamID: NewAlphaNumeric("1", 3),
	}
	e.DataElement = NewDataElementGroup(EncryptionAlgorithmDEG, 7, e)
	return e
}

type EncryptionAlgorithmDataElement struct {
	DataElement
	// "2" for OSY, Owner Symmetric
	Usage *AlphaNumericDataElement
	// "2" for CBC, Cipher Block Chaining.
	OperationMode *AlphaNumericDataElement
	// "13" for 2-Key-Triple-DES
	Algorithm *AlphaNumericDataElement
	Key       *BinaryDataElement
	// "5" for KYE, Symmetric key, en-/decryption with a symmetric key (DDV)
	// "6" for KYP, Symmetric key, encryption with a public key (RDH).
	KeyParamID                 *AlphaNumericDataElement
	InitializationValueParamID *AlphaNumericDataElement
	InitializationValue        *BinaryDataElement
}

func (e *EncryptionAlgorithmDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		e.Usage,
		e.OperationMode,
		e.Algorithm,
		e.Key,
		e.KeyParamID,
		e.InitializationValueParamID,
		e.InitializationValue,
	}
}
