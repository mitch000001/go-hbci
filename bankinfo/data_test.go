package bankinfo

import (
	"strings"
	"testing"
)

func TestBankDataConsistency(t *testing.T) {
	var errs []string
	for _, info := range data {
		_, err := hbciVersion(info.VersionName, info.VersionNumber)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) != 0 {
		t.Logf("Expected no errors, got:\n\t%s", strings.Join(errs, "\n\t"))
		t.Fail()
	}
}
