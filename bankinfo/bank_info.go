package bankinfo

import (
	"fmt"
	"math"
	"strings"
)

const (
	version210       = "2.1"
	version220       = "2.2"
	version300       = "3.0"
	version400       = "4.0"
	version410       = "4.1"
	versionString220 = "HBCI 2.2 Erweiterung PIN/TAN V1.01"
	versionString300 = "FinTS V3.0"
	versionString400 = "FinTS V4.0"
	versionString410 = "FinTS V4.1"
)

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
	VersionName   string
	Institute     string
	City          string
}

func (b BankInfo) HbciVersion() int {
	version, err := hbciVersion(b.VersionName, b.VersionNumber)
	if err != nil {
		panic(err)
	}
	return version
}

type SortableBankInfos []BankInfo

func (s SortableBankInfos) Len() int           { return len(s) }
func (s SortableBankInfos) Swap(a, b int)      { s[a], s[b] = s[b], s[a] }
func (s SortableBankInfos) Less(a, b int) bool { return s[a].BankId < s[b].BankId }

func hbciVersion(versionName, versionNumber string) (int, error) {
	var errs []string
	parsedVersionName, err := parseVersionName(versionName)
	if err != nil {
		errs = append(errs, err.Error())
	}
	parsedVersionNumber, err := parseVersionNumber(versionNumber)
	if err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) != 0 {
		return 0, fmt.Errorf(strings.Join(errs, "\n"))
	}
	return int(math.Max(float64(parsedVersionNumber), float64(parsedVersionName))), nil
}

func parseVersionName(versionName string) (int, error) {
	switch versionName {
	case versionString220:
		return 220, nil
	case versionString300:
		return 300, nil
	case versionString400:
		return 400, nil
	case versionString410:
		return 410, nil
	case "":
		return -1, nil
	default:
		return 0, fmt.Errorf("Unknown HBCI Version Name: %q", versionName)
	}
}

func parseVersionNumber(versionNumber string) (int, error) {
	switch versionNumber {
	case version210:
		return 210, nil
	case version220:
		return 220, nil
	case version300:
		return 300, nil
	case version400:
		return 400, nil
	case version410:
		return 410, nil
	case "":
		return -1, nil
	default:
		return 0, fmt.Errorf("Unknown HBCI Version Number: %q", versionNumber)
	}
}
