package bankinfo

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParserParse(t *testing.T) {
	bankData := fmt.Sprintf(
		`%s;BLA;%s;%s;XYZ;%s
		1000000;xxx;3.0;https://foo.example.com;FOO;FinTS V3.0`,
		BANK_IDENTIFIER, VERSION_NUMBER_HEADER, URL_HEADER, VERSION_NAME_HEADER,
	)
	parser := Parser{}
	var result []BankInfo

	result, err := parser.Parse(strings.NewReader(bankData))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	var expectedResult = []BankInfo{
		BankInfo{
			BankId:        "1000000",
			VersionNumber: "3.0",
			URL:           "https://foo.example.com",
			VersionName:   "FinTS V3.0",
		},
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Logf("Expected result to equal\n%q\n\tgot:\n%q\n", expectedResult, result)
		t.Fail()
	}
}
