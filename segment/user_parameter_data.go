package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment CommonUserParameterDataSegment

type CommonUserParameterDataSegment struct {
	Segment
	UserID     *element.IdentificationDataElement
	UPDVersion *element.NumberDataElement
	// Status |￼Beschreibung
	// -----------------------------------------------------------------
	// 0	  | Die nicht aufgeführten Geschäftsvorfälle sind gesperrt
	//		  | (die aufgeführten Geschäftsvorfälle sind zugelassen).
	// 1 ￼ ￼  | Bei den nicht aufgeführten Geschäftsvorfällen ist anhand
	//        | der UPD keine Aussage darüber möglich, ob diese erlaubt
	//        | oder gesperrt sind. Diese Prüfung kann nur online vom
	//        | Kreditinstitutssystem vorgenommen werden.
	UPDUsage *element.NumberDataElement
}

func (c *CommonUserParameterDataSegment) Version() int         { return 2 }
func (c *CommonUserParameterDataSegment) ID() string           { return "HIUPA" }
func (c *CommonUserParameterDataSegment) referencedId() string { return "HKVVB" }
func (c *CommonUserParameterDataSegment) sender() string       { return senderBank }

func (c *CommonUserParameterDataSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.UserID,
		c.UPDVersion,
		c.UPDUsage,
	}
}

func (c *CommonUserParameterDataSegment) UserParameterData() domain.UserParameterData {
	return domain.UserParameterData{
		UserID:  c.UserID.Val(),
		Version: c.UPDVersion.Val(),
		Usage:   c.UPDUsage.Val(),
	}
}
