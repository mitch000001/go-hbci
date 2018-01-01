package bankinfo

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseBankInfos(t *testing.T) {
	header := []string{
		bankIdentifierHeader,
		"Another header",
		bankInstituteHeader,
		versionNumberHeader,
		urlHeader,
		"header of no interest",
		versionNameHeader,
		cityHeader,
	}
	content := []string{
		"1000000",
		"xxx",
		"Bank Institute",
		"3.0",
		"https://foo.example.com",
		"FOO",
		"FinTS V3.0",
		"Hamburg",
	}

	bankData := fmt.Sprintf(`%s
		%s`,
		strings.Join(header, ";"),
		strings.Join(content, ";"),
	)
	var result []BankInfo

	result, err := ParseBankInfos(strings.NewReader(bankData))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	var expectedResult = []BankInfo{
		BankInfo{
			BankID:        "1000000",
			VersionNumber: "3.0",
			URL:           "https://foo.example.com",
			VersionName:   "FinTS V3.0",
			Institute:     "Bank Institute",
			City:          "Hamburg",
		},
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Logf("Expected result to equal\n%#v\n\tgot:\n%#v\n", expectedResult, result)
		t.Fail()
	}
}

func TestParseBicData(t *testing.T) {
	bicData := fmt.Sprintf(
		`%s;BLA;%s;XYZ
		1000000;xxx;MARKDEF1100;abc`,
		bicBankIdentifier, bicIdentifier,
	)
	var result []BicInfo

	result, err := ParseBicData(strings.NewReader(bicData))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	var expectedResult = []BicInfo{
		{
			BankID: "1000000",
			BIC:    "MARKDEF1100",
		},
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Logf("Expected result to equal\n%q\n\tgot:\n%q\n", expectedResult, result)
		t.Fail()
	}
}
