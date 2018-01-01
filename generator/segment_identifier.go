package generator

// SegmentIdentifier represent a segment definition for the generator
type SegmentIdentifier struct {
	Name          string
	InterfaceName string
	Versions      []SegmentIdentifier
	Version       int
}
