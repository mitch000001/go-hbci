package element

import "fmt"

// NewPinTanEncryptionAlgorithm returns an EncryptionAlgorithmDataElement ready
// to use in pin/tan flow
func NewPinTanEncryptionAlgorithm() *EncryptionAlgorithmDataElement {
	e := &EncryptionAlgorithmDataElement{
		Usage:                      NewAlphaNumeric("2", 3),
		OperationMode:              NewAlphaNumeric("2", 3),
		Algorithm:                  NewAlphaNumeric("13", 3),
		Key:                        NewBinary([]byte(defaultPinTan), 512),
		KeyParamID:                 NewAlphaNumeric("5", 3),
		InitializationValueParamID: NewAlphaNumeric("1", 3),
	}
	e.DataElement = NewDataElementGroup(encryptionAlgorithmDEG, 7, e)
	return e
}

// NewRDHEncryptionAlgorithm returns an EncryptionAlgorithmDataElement ready to
// use in CardReader flow
func NewRDHEncryptionAlgorithm(pubKey []byte) *EncryptionAlgorithmDataElement {
	e := &EncryptionAlgorithmDataElement{
		Usage:                      NewAlphaNumeric("2", 3),
		OperationMode:              NewAlphaNumeric("2", 3),
		Algorithm:                  NewAlphaNumeric("13", 3),
		Key:                        NewBinary(pubKey, 512),
		KeyParamID:                 NewAlphaNumeric("6", 3),
		InitializationValueParamID: NewAlphaNumeric("1", 3),
	}
	e.DataElement = NewDataElementGroup(encryptionAlgorithmDEG, 7, e)
	return e
}

// EncryptionAlgorithmDataElement represents an encryption algorithm
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

// GroupDataElements returns the grouped DataElements
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

// UnmarshalHBCI unmarshals value into the DataElement
func (e *EncryptionAlgorithmDataElement) UnmarshalHBCI(value []byte) error {
	e = &EncryptionAlgorithmDataElement{}
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 7 {
		return fmt.Errorf("Malformed marshaled value")
	}
	e.DataElement = NewDataElementGroup(encryptionAlgorithmDEG, 5, e)
	e.Usage = &AlphaNumericDataElement{}
	err = e.Usage.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	e.OperationMode = &AlphaNumericDataElement{}
	err = e.OperationMode.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		e.Algorithm = &AlphaNumericDataElement{}
		err = e.Algorithm.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		e.Key = &BinaryDataElement{}
		err = e.Key.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		e.KeyParamID = &AlphaNumericDataElement{}
		err = e.KeyParamID.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		e.InitializationValueParamID = &AlphaNumericDataElement{}
		err = e.InitializationValueParamID.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		e.InitializationValue = &BinaryDataElement{}
		err = e.InitializationValue.UnmarshalHBCI(elements[6])
		if err != nil {
			return err
		}
	}
	return nil
}
