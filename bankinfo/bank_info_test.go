package bankinfo

import "testing"

func TestFindByBankId(t *testing.T) {
	data = BankData{
		BankInfo{BankId: "1000000", URL: "1"},
		BankInfo{BankId: "2000000", URL: "2"},
		BankInfo{BankId: "3000000", URL: "3"},
	}

	var url string

	url = FindByBankId("1000000").URL

	if url != "1" {
		t.Logf("Expected url to equal %q, got %q\n", "1", url)
		t.Fail()
	}

	url = FindByBankId("2000000").URL

	if url != "2" {
		t.Logf("Expected url to equal %q, got %q\n", "2", url)
		t.Fail()
	}

	url = FindByBankId("3000000").URL

	if url != "3" {
		t.Logf("Expected url to equal %q, got %q\n", "3", url)
		t.Fail()
	}
}

func TestHbciVersion(t *testing.T) {
	tests := []struct {
		versionNumber string
		versionString string
		result        int
	}{
		{
			versionNumber: "3.0",
			versionString: "FinTS V3.0",
			result:        300,
		},
		{
			versionNumber: "2.2",
			versionString: "HBCI 2.2 Erweiterung PIN/TAN V1.01",
			result:        220,
		},
		{
			versionNumber: "",
			versionString: "FinTS V3.0",
			result:        300,
		},
		{
			versionNumber: "",
			versionString: "HBCI 2.2 Erweiterung PIN/TAN V1.01",
			result:        220,
		},
		{
			versionNumber: "2.2",
			versionString: "FinTS V3.0",
			result:        300,
		},
	}

	for _, test := range tests {

		version, err := hbciVersion(test.versionString, test.versionNumber)

		if err != nil {
			t.Logf("Expected no error, got %q\n", err)
			t.Fail()
		}

		expectedVersion := test.result

		if expectedVersion != version {
			t.Logf("Expected version to equal %d, got %d\n", expectedVersion, version)
			t.Fail()
		}
	}
}
