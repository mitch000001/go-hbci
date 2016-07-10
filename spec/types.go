package spec

type Type string

const (
	DE  Type = "DataElement"
	DEG Type = "DataElementGroup"
	GD  Type = "GroupDataElement"
	GDG Type = "GroupDataElementGroup"
)

type Format string

const (
	an    Format = "AlphaNumeric"
	txt   Format = "Text"
	num   Format = "Number"
	dig   Format = "Digit"
	float Format = "Float"
	dta   Format = "DTAUSCharset"
	bin   Format = "Binary"
	// Derived types
	jn   Format = "Boolean"
	code Format = "Code"
	dat  Format = "Date"
	vdat Format = "VirtualDate"
	tim  Format = "Time"
	id   Format = "Identification"
	ctr  Format = "CountryCode"
	cur  Format = "Currency"
	wrt  Format = "Value"
	// Multiple used elements or GroupDataElementGroups
	btg  Format = "Amount"
	kik  Format = "BankIdentification"
	ktv  Format = "AccountConnection"
	kti  Format = "InternationalAccountConnection"
	sdo  Format = "Balance"
	addr Format = "Address"

	composed Format = "Composed"
)

type Length string

type Status string

const (
	M Status = "Mandatory"
	C Status = "Conditional"
	N Status = "Not allowed"
	O Status = "Optional"
)

type Sender string

const (
	Kunde          Sender = "Client"
	Kreditinstitut Sender = "Bank"
)

type SegmentSpec struct {
	TypeName    string
	Name        string
	Kind        string
	Sender      Sender
	Id          string
	ReferenceId string
	Version     int
	Elements    []DataElementUsageSpec
}

type DataElementUsageSpec struct {
	Number     int
	FieldName  string
	Name       string
	Format     Format
	Version    int
	Type       Type
	Length     Length
	Status     Status
	Count      int
	Limitation string
}

type DataElementSpec struct {
	TypeName string
	Name     string
	Version  int
	Type     Type
	Format   Format
	Length   Length
	Elements []DataElementUsageSpec
}
