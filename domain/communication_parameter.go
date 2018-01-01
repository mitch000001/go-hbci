package domain

// CommunicationParameter provides information about the access point to a bank
// institute
type CommunicationParameter struct {
	Protocol              int
	Address               string
	AddressAddition       string
	FilterFunction        string
	FilterFunctionVersion int
}
