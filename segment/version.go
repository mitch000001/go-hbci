package segment

import "github.com/mitch000001/go-hbci/element"

type version interface {
	version() int
	versionedElements() []element.DataElement
}
