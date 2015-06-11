package hbci

func NewKeyNameDataElement() *KeyNameDataElement {
	a := &KeyNameDataElement{}
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
