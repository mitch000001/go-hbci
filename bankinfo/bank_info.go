package bankinfo

import (
	"fmt"
	"math"
)

const (
	version220       = "2.2"
	version300       = "3.0"
	version400       = "4.0"
	version410       = "4.1"
	versionString220 = "HBCI 2.2 Erweiterung PIN/TAN V1.01"
	versionString300 = "FinTS V3.0"
	versionString400 = "FinTS V4.0"
	versionString410 = "FinTS V4.1"
)

type BankData []BankInfo

func FindByBankId(bankId string) BankInfo {
	var bankInfo BankInfo
	for _, entry := range data {
		if entry.BankId == bankId {
			bankInfo = entry
		}
	}
	return bankInfo
}

type BankInfo struct {
	BankId        string
	VersionNumber string
	URL           string
	VersionString string
}

func (b BankInfo) HbciVersion() int {
	parsedVersionNumber := b.parseVersionNumber()
	parsedVersionString := b.parseVersionString()
	return int(math.Max(float64(parsedVersionNumber), float64(parsedVersionString)))
}

func (b BankInfo) parseVersionString() int {
	switch b.VersionString {
	case versionString220:
		return 220
	case versionString300:
		return 300
	case versionString400:
		return 400
	case versionString410:
		return 410
	case "":
		return -1
	default:
		panic(fmt.Errorf("Unknown HBCI Version String: %q", b.VersionString))
	}
}

func (b BankInfo) parseVersionNumber() int {
	switch b.VersionNumber {
	case version220:
		return 220
	case version300:
		return 300
	case version400:
		return 400
	case version410:
		return 410
	case "":
		return -1
	default:
		panic(fmt.Errorf("Unknown HBCI Version: %q", b.VersionNumber))
	}
}
