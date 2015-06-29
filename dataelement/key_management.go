package dataelement

import (
	"fmt"
	"reflect"

	"github.com/mitch000001/go-hbci/domain"
)

func NewPublicKeyDataElement(pubKey *domain.PublicKey) *PublicKeyDataElement {
	if !reflect.DeepEqual(pubKey.Exponent, []byte("65537")) {
		panic(fmt.Errorf("Exponent must equal 65537 (% X)", "65537"))
	}
	p := &PublicKeyDataElement{
		Usage:         NewAlphaNumericDataElement(pubKey.Type, 3),
		OperationMode: NewAlphaNumericDataElement("16", 3),
		Cipher:        NewAlphaNumericDataElement("10", 3),
		Modulus:       NewBinaryDataElement(pubKey.Modulus, 512),
		ModulusID:     NewAlphaNumericDataElement("12", 3),
		Exponent:      NewBinaryDataElement(pubKey.Exponent, 512),
		ExponentID:    NewAlphaNumericDataElement("13", 3),
	}
	p.DataElement = NewDataElementGroup(PublicKeyDEG, 7, p)
	return p
}

type PublicKeyDataElement struct {
	DataElement
	// "5" for OCF, Owner Ciphering (Encryption key)
	// "6" for OSG, Owner Signing (Signing key)
	Usage *AlphaNumericDataElement
	// "16" for DSMR (ISO 9796)
	OperationMode *AlphaNumericDataElement
	// "10" for RSA
	Cipher  *AlphaNumericDataElement
	Modulus *BinaryDataElement
	// "12" for MOD, Modulus
	ModulusID *AlphaNumericDataElement
	// "65537"
	Exponent *BinaryDataElement
	// "13" for EXP, Exponent
	ExponentID *AlphaNumericDataElement
}

func (p *PublicKeyDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		p.Usage,
		p.OperationMode,
		p.Cipher,
		p.Modulus,
		p.ModulusID,
		p.Exponent,
		p.ExponentID,
	}
}

func (p *PublicKeyDataElement) Val() *domain.PublicKey {
	return &domain.PublicKey{
		Type:     p.Usage.Val(),
		Modulus:  p.Modulus.Val(),
		Exponent: p.Exponent.Val(),
	}
}
