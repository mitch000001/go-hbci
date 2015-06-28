package hbci

type version interface {
	version() int
	versionedElements() []DataElement
}
