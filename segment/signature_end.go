package segment

import "github.com/mitch000001/go-hbci/element"

type SignatureEnd interface {
	SetSignature(signature []byte)
	SetPinTan(pin, tan string)
}

func NewSignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndV1{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.ClientSegment = NewBasicSegment(number, s)

	segment := &SignatureEndSegment{
		signatureEndSegment: s,
	}
	return segment
}

type signatureEndSegment interface {
	ClientSegment
	SignatureEnd
}

type SignatureEndSegment struct {
	signatureEndSegment
}

func (s *SignatureEndSegment) SetSignature(signature []byte) {
	s.signatureEndSegment.SetSignature(signature)
}

func (s *SignatureEndSegment) SetPinTan(pin, tan string) {
	s.signatureEndSegment.SetPinTan(pin, tan)
}

type SignatureEndV1 struct {
	ClientSegment
	SecurityControlRef *element.AlphaNumericDataElement
	Signature          *element.BinaryDataElement
	PinTan             *element.PinTanDataElement
}

func (s *SignatureEndV1) Version() int         { return 1 }
func (s *SignatureEndV1) ID() string           { return "HNSHA" }
func (s *SignatureEndV1) referencedId() string { return "" }
func (s *SignatureEndV1) sender() string       { return senderBoth }

func (s *SignatureEndV1) elements() []element.DataElement {
	return []element.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndV1) SetSignature(signature []byte) {
	s.Signature = element.NewBinary(signature, 512)
}

func (s *SignatureEndV1) SetPinTan(pin, tan string) {
	s.PinTan = element.NewPinTan(pin, tan)
}
func NewSignatureEndSegmentV2(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndV2{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.ClientSegment = NewBasicSegment(number, s)

	segment := &SignatureEndSegment{
		signatureEndSegment: s,
	}
	return segment
}

type SignatureEndV2 struct {
	ClientSegment
	SecurityControlRef *element.AlphaNumericDataElement
	Signature          *element.BinaryDataElement
	CustomSignature    *element.CustomSignatureDataElement
}

func (s *SignatureEndV2) Version() int         { return 2 }
func (s *SignatureEndV2) ID() string           { return "HNSHA" }
func (s *SignatureEndV2) referencedId() string { return "" }
func (s *SignatureEndV2) sender() string       { return senderBoth }

func (s *SignatureEndV2) elements() []element.DataElement {
	return []element.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.CustomSignature,
	}
}

func (s *SignatureEndV2) SetSignature(signature []byte) {
	s.Signature = element.NewBinary(signature, 512)
}

func (s *SignatureEndV2) SetPinTan(pin, tan string) {
	s.CustomSignature = element.NewCustomSignature(pin, tan)
}
