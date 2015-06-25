package hbci

func NewDialogInitializationClientMessage() *DialogInitializationClientMessage {
	d := &DialogInitializationClientMessage{}
	d.basicClientMessage = newBasicClientMessage(d)
	return d
}

type DialogInitializationClientMessage struct {
	*basicClientMessage
	Identification             *IdentificationSegment
	ProcessingPreparation      *ProcessingPreparationSegment
	PublicSigningKeyRequest    *PublicKeyRequestSegment
	PublicEncryptionKeyRequest *PublicKeyRequestSegment
}

func (d *DialogInitializationClientMessage) Jobs() SegmentSequence {
	return SegmentSequence{
		d.Identification,
		d.ProcessingPreparation,
		d.PublicSigningKeyRequest,
		d.PublicEncryptionKeyRequest,
	}
}

type DialogInitializationBankMessage struct {
	*basicBankMessage
	BankParams            SegmentSequence
	UserParams            SegmentSequence
	PublicKeyTransmission *PublicKeyTransmissionSegment
	Announcement          *BankAnnouncementSegment
}

type DialogFinishingMessage struct {
	*basicClientMessage
	DialogEnd *DialogEndSegment
}

func (d *DialogFinishingMessage) Jobs() SegmentSequence {
	return SegmentSequence{
		d.DialogEnd,
	}
}

func NewDialogCancellationMessage(messageAcknowledgement *MessageAcknowledgement) *DialogCancellationMessage {
	d := &DialogCancellationMessage{
		MessageAcknowledgements: messageAcknowledgement,
	}
	return d
}

type DialogCancellationMessage struct {
	*basicMessage
	MessageAcknowledgements *MessageAcknowledgement
}

type AnonymousDialogMessage struct {
	*basicMessage
	Identification        *IdentificationSegment
	ProcessingPreparation *ProcessingPreparationSegment
}

func NewDialogEndSegment(dialogId string) *DialogEndSegment {
	d := &DialogEndSegment{
		DialogID: NewIdentificationDataElement(dialogId),
	}
	d.basicSegment = NewBasicSegment("HKEND", 3, 1, d)
	return d
}

type DialogEndSegment struct {
	*basicSegment
	DialogID *IdentificationDataElement
}

func (d *DialogEndSegment) elements() []DataElement {
	return []DataElement{
		d.DialogID,
	}
}

func NewProcessingPreparationSegment(bdpVersion int, udpVersion int, language int) *ProcessingPreparationSegment {
	p := &ProcessingPreparationSegment{
		BPDVersion:     NewNumberDataElement(bdpVersion, 3),
		UPDVersion:     NewNumberDataElement(udpVersion, 3),
		DialogLanguage: NewNumberDataElement(language, 3),
		ProductName:    NewAlphaNumericDataElement(productName, 25),
		ProductVersion: NewAlphaNumericDataElement(productVersion, 5),
	}
	p.basicSegment = NewBasicSegment("HKVVB", 4, 2, p)
	return p
}

type ProcessingPreparationSegment struct {
	*basicSegment
	BPDVersion *NumberDataElement
	UPDVersion *NumberDataElement
	// 0 for undefined
	// Sprachkennzeichen | Bedeutung   | Sprachencode ISO 639 | ISO 8859 Subset | ISO 8859- Codeset
	// --------------------------------------------------------------------------------------------
	// 1				 | Deutsch	   | de (German) ￼	      | Deutsch ￼ ￼		| 1 (Latin 1)
	// 2				 | Englisch	   | en (English)		  | Englisch		| 1 (Latin 1)
	// 3 				 | Französisch | fr (French)  		  | Französisch ￼	| 1 (Latin 1)
	DialogLanguage *NumberDataElement
	ProductName    *AlphaNumericDataElement
	ProductVersion *AlphaNumericDataElement
}

func (p *ProcessingPreparationSegment) elements() []DataElement {
	return []DataElement{
		p.BPDVersion,
		p.UPDVersion,
		p.DialogLanguage,
		p.ProductName,
		p.ProductVersion,
	}
}

func NewBankAnnouncementSegment(subject, body string) *BankAnnouncementSegment {
	b := &BankAnnouncementSegment{
		Subject: NewAlphaNumericDataElement(subject, 35),
		Body:    NewTextDataElement(body, 2048),
	}
	b.basicSegment = NewBasicSegment("HIKIM", 8, 2, b)
	return b
}

type BankAnnouncementSegment struct {
	*basicSegment
	Subject *AlphaNumericDataElement
	Body    *TextDataElement
}

func (b *BankAnnouncementSegment) elements() []DataElement {
	return []DataElement{
		b.Subject,
		b.Body,
	}
}
