package bankinfo

import "testing"

func TestFindByBankId(t *testing.T) {
	data = []BankInfo{
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
		versionName   string
		result        int
	}{
		{
			versionNumber: "3.0",
			versionName:   "FinTS V3.0",
			result:        300,
		},
		{
			versionNumber: "2.2",
			versionName:   "HBCI 2.2 Erweiterung PIN/TAN V1.01",
			result:        220,
		},
		{
			versionNumber: "",
			versionName:   "FinTS V3.0",
			result:        300,
		},
		{
			versionNumber: "",
			versionName:   "HBCI 2.2 Erweiterung PIN/TAN V1.01",
			result:        220,
		},
		{
			versionNumber: "2.2",
			versionName:   "FinTS V3.0",
			result:        300,
		},
	}

	for _, test := range tests {

		version, err := hbciVersion(test.versionName, test.versionNumber)

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

func TestSortableBankInfosSortInterface(t *testing.T) {
	sorted := SortableBankInfos{
		BankInfo{BankId: "30000000"},
		BankInfo{BankId: "10000000"},
		BankInfo{BankId: "20000000"},
	}

	length := sorted.Len()

	if len(sorted) != length {
		t.Logf("Expected length to equal %d, got %d\n", len(sorted), length)
		t.Fail()
	}

	less := sorted.Less(0, 1)

	if less {
		t.Logf("Expected first entry not to be less than second, but was\n")
		t.Fail()
	}

	sorted.Swap(0, 1)

	if sorted[0].BankId != "10000000" {
		t.Logf("Expected first entry to have BankId '10000000', but was %q\n", sorted[0].BankId)
		t.Fail()
	}

	if sorted[1].BankId != "30000000" {
		t.Logf("Expected first entry to have BankId '30000000', but was %q\n", sorted[1].BankId)
		t.Fail()
	}
}
