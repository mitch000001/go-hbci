package hbci

const defaultPinTan = "\x00\x00\x00\x00\x00\x00\x00\x00"

func NewPinKey(pin string, keyName *KeyName) *PinKey {
	return &PinKey{pin: pin, keyName: keyName}
}

type PinKey struct {
	pin     string
	keyName *KeyName
}

func (p *PinKey) KeyName() KeyName {
	return *p.keyName
}

func (p *PinKey) SetKeyNumber(number int) {
	p.keyName.KeyNumber = number
}

func (p *PinKey) SetKeyVersion(version int) {
	p.keyName.KeyVersion = version
}

func (p *PinKey) CanSign() bool {
	return true
}

func (p *PinKey) CanEncrypt() bool {
	return true
}

func (p *PinKey) Pin() string {
	return p.pin
}

func (p *PinKey) Sign(message []byte) ([]byte, error) {
	return []byte(p.pin), nil
}

func (p *PinKey) Encrypt(message []byte) ([]byte, error) {
	encMessage := make([]byte, len(message))
	// Make a deep copy, just in case
	copy(encMessage, message)
	return encMessage, nil
}

func NewPinTanEncryptionProvider(key *PinKey, clientSystemId string) *PinTanEncryptionProvider {
	return &PinTanEncryptionProvider{
		key:            key,
		clientSystemId: clientSystemId,
	}
}

type PinTanEncryptionProvider struct {
	key            *PinKey
	clientSystemId string
}

func (p *PinTanEncryptionProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanEncryptionProvider) Encrypt(message []byte) (*EncryptedMessage, error) {
	encryptedBytes, _ := p.key.Encrypt(message)
	encryptedMessage := NewEncryptedPinTanMessage(p.clientSystemId, p.key.KeyName(), encryptedBytes)
	return encryptedMessage, nil
}

func (p *PinTanEncryptionProvider) EncryptWithInitialKeyName(message []byte) (*EncryptedMessage, error) {
	keyName := p.key.KeyName()
	keyName.SetInitial()
	encryptedBytes, _ := p.key.Encrypt(message)
	encryptedMessage := NewEncryptedPinTanMessage(p.clientSystemId, keyName, encryptedBytes)
	return encryptedMessage, nil
}

func NewPinTanSignatureProvider(key *PinKey, clientSystemId string) SignatureProvider {
	return &PinTanSignatureProvider{key: key, clientSystemId: clientSystemId}
}

type PinTanSignatureProvider struct {
	key            *PinKey
	clientSystemId string
}

func (p *PinTanSignatureProvider) SetClientSystemID(clientSystemId string) {
	p.clientSystemId = clientSystemId
}

func (p *PinTanSignatureProvider) SignMessage(signedMessage SignedHBCIMessage) error {
	signedMessage.SignatureEndSegment().SetPinTan(p.key.Pin(), "")
	return nil
}

func (p *PinTanSignatureProvider) NewSignatureHeader(controlReference string, signatureId int) *SignatureHeaderSegment {
	return NewPinTanSignatureHeaderSegment(controlReference, p.clientSystemId, p.key.KeyName())
}

func NewPinTanDataElement(pin, tan string) *PinTanDataElement {
	p := &PinTanDataElement{
		PIN: NewAlphaNumericDataElement(pin, 6),
	}
	if tan != "" {
		p.TAN = NewAlphaNumericDataElement(tan, 35)
	}
	p.DataElement = NewDataElementGroup(PinTanDEG, 2, p)
	return p
}

type PinTanDataElement struct {
	DataElement
	PIN *AlphaNumericDataElement
	TAN *AlphaNumericDataElement
}

func (p *PinTanDataElement) GroupDataElements() []DataElement {
	return []DataElement{
		p.PIN,
		p.TAN,
	}
}
