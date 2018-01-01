package domain

// ReferencingMessage represents a reference to another message within a given
// dialog
type ReferencingMessage struct {
	DialogID      string
	MessageNumber int
}
