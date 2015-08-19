package generator

type SegmentIdentifier struct {
	Name          string
	InterfaceName string
	Versions      []SegmentIdentifier
	Version       int
}
