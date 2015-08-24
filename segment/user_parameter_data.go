package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment CommonUserParameterDataSegment -segment_interface commonUserParameterDataSegment -segment_versions="CommonUserParameterDataV2:2:Segment,CommonUserParameterDataV3:3:Segment,CommonUserParameterDataV4:4:Segment"

type CommonUserParameterData interface {
	BankSegment
	UserParameterData() domain.UserParameterData
}

type CommonUserParameterDataSegment struct {
	commonUserParameterDataSegment
}

type commonUserParameterDataSegment interface {
	BankSegment
	UserParameterData() domain.UserParameterData
}

type CommonUserParameterDataV2 struct {
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

func (c *CommonUserParameterDataV2) Version() int         { return 2 }
func (c *CommonUserParameterDataV2) ID() string           { return "HIUPA" }
func (c *CommonUserParameterDataV2) referencedId() string { return "HKVVB" }
func (c *CommonUserParameterDataV2) sender() string       { return senderBank }

func (c *CommonUserParameterDataV2) elements() []element.DataElement {
	return []element.DataElement{
		c.UserID,
		c.UPDVersion,
		c.UPDUsage,
	}
}

func (c *CommonUserParameterDataV2) UserParameterData() domain.UserParameterData {
	return domain.UserParameterData{
		UserID:  c.UserID.Val(),
		Version: c.UPDVersion.Val(),
		Usage:   c.UPDUsage.Val(),
	}
}

type CommonUserParameterDataV3 struct {
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
	UPDUsage         *element.NumberDataElement
	UserName         *element.AlphaNumericDataElement
	CommonExtensions *element.AlphaNumericDataElement
}

func (c *CommonUserParameterDataV3) Version() int         { return 3 }
func (c *CommonUserParameterDataV3) ID() string           { return "HIUPA" }
func (c *CommonUserParameterDataV3) referencedId() string { return "HKVVB" }
func (c *CommonUserParameterDataV3) sender() string       { return senderBank }

func (c *CommonUserParameterDataV3) elements() []element.DataElement {
	return []element.DataElement{
		c.UserID,
		c.UPDVersion,
		c.UPDUsage,
		c.UserName,
		c.CommonExtensions,
	}
}

func (c *CommonUserParameterDataV3) UserParameterData() domain.UserParameterData {
	return domain.UserParameterData{
		UserID:  c.UserID.Val(),
		Version: c.UPDVersion.Val(),
		Usage:   c.UPDUsage.Val(),
	}
}

type CommonUserParameterDataV4 struct {
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
	UPDUsage         *element.NumberDataElement
	UserName         *element.AlphaNumericDataElement
	CommonExtensions *element.AlphaNumericDataElement
}

func (c *CommonUserParameterDataV4) Version() int         { return 4 }
func (c *CommonUserParameterDataV4) ID() string           { return "HIUPA" }
func (c *CommonUserParameterDataV4) referencedId() string { return "HKVVB" }
func (c *CommonUserParameterDataV4) sender() string       { return senderBank }

func (c *CommonUserParameterDataV4) elements() []element.DataElement {
	return []element.DataElement{
		c.UserID,
		c.UPDVersion,
		c.UPDUsage,
		c.UserName,
		c.CommonExtensions,
	}
}

func (c *CommonUserParameterDataV4) UserParameterData() domain.UserParameterData {
	return domain.UserParameterData{
		UserID:  c.UserID.Val(),
		Version: c.UPDVersion.Val(),
		Usage:   c.UPDUsage.Val(),
	}
}
