package hbci

import "github.com/mitch000001/go-hbci/dataelement"

type version interface {
	version() int
	versionedElements() []dataelement.DataElement
}
