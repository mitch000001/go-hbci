package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

const IdentificationID = "HKIDN"

func NewIdentificationSegment(bankId domain.BankID, clientId string, clientSystemId string, systemIdRequired bool) *IdentificationSegment {
	var clientSystemStatus *element.NumberDataElement
	if systemIdRequired {
		clientSystemStatus = element.NewNumber(1, 1)
	} else {
		clientSystemStatus = element.NewNumber(0, 1)
	}
	id := &IdentificationSegment{
		BankId:             element.NewBankIdentification(bankId),
		ClientId:           element.NewIdentification(clientId),
		ClientSystemId:     element.NewIdentification(clientSystemId),
		ClientSystemStatus: clientSystemStatus,
	}
	id.ClientSegment = NewBasicSegment(3, id)
	return id
}

type IdentificationSegment struct {
	ClientSegment
	BankId             *element.BankIdentificationDataElement
	ClientId           *element.IdentificationDataElement
	ClientSystemId     *element.IdentificationDataElement
	ClientSystemStatus *element.NumberDataElement
}

func (i *IdentificationSegment) Version() int         { return 2 }
func (i *IdentificationSegment) ID() string           { return IdentificationID }
func (i *IdentificationSegment) referencedId() string { return "" }
func (i *IdentificationSegment) sender() string       { return senderUser }

func (i *IdentificationSegment) elements() []element.DataElement {
	return []element.DataElement{
		i.BankId,
		i.ClientId,
		i.ClientSystemId,
		i.ClientSystemStatus,
	}
}
