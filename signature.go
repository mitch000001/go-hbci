package hbci

func NewKeyNameDataElement(countryCode int, bankId string, userId string, keyType string, keyNumber, keyVersion int) *KeyNameDataElement {
	a := &KeyNameDataElement{
		Bank:       NewBankIndentificationDataElementWithBankId(countryCode, bankId),
		UserID:     NewIdentificationDataElement(userId),
		KeyType:    NewAlphaNumericDataElement(keyType, 1),
		KeyNumber:  NewNumberDataElement(keyNumber, 3),
		KeyVersion: NewNumberDataElement(keyVersion, 3),
	}
	a.elementGroup = NewDataElementGroup(KeyNameDEG, 5, a)
	return a
}

type KeyNameDataElement struct {
	*elementGroup
	Bank       *BankIdentificationDataElement
	UserID     *IdentificationDataElement
	KeyType    *AlphaNumericDataElement
	KeyNumber  *NumberDataElement
	KeyVersion *NumberDataElement
}

func (k *KeyNameDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		k.Bank,
		k.UserID,
		k.KeyType,
		k.KeyNumber,
		k.KeyVersion,
	}
}
