package domain

const (
	// HBCIVersion220 represents version 2.2.0 of HBCI protocol
	HBCIVersion220 = 220
	// FINTSVersion300 represents version 3.0.0 of FINTS protocol
	FINTSVersion300 = 300
)

// SupportedHBCIVersions provides a list of supported versions
var SupportedHBCIVersions = []int{
	HBCIVersion220,
	FINTSVersion300,
}
