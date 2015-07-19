package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewIdentificationSegment(bankId domain.BankId, clientId string, clientSystemId string, systemIdRequired bool) *IdentificationSegment {
	var clientSystemStatus *element.NumberDataElement
	if systemIdRequired {
		clientSystemStatus = element.NewNumber(1, 1)
	} else {
		clientSystemStatus = element.NewNumber(0, 1)
	}
	id := &IdentificationSegment{
		BankId:             element.NewBankIndentification(bankId),
		ClientId:           element.NewIdentification(clientId),
		ClientSystemId:     element.NewIdentification(clientSystemId),
		ClientSystemStatus: clientSystemStatus,
	}
	id.Segment = NewBasicSegment(3, id)
	return id
}

type IdentificationSegment struct {
	Segment
	BankId             *element.BankIdentificationDataElement
	ClientId           *element.IdentificationDataElement
	ClientSystemId     *element.IdentificationDataElement
	ClientSystemStatus *element.NumberDataElement
}

func (i *IdentificationSegment) version() int         { return 2 }
func (i *IdentificationSegment) id() string           { return "HKIDN" }
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
