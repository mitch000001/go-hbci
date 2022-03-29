package domain

// MessageReference represents a reference to another message within a given
// dialog
type MessageReference struct {
	DialogID      string
	MessageNumber int
}
