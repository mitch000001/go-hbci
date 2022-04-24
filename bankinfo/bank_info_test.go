package bankinfo

import (
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func Test_ParseBankInfos(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		wantResult []BankInfo
		wantError  error
	}{
		{
			name: "success",
			data: `Nr.;BLZ;Institut;Ort;RZ;Organisation;HBCI-Zugang DNS;HBCI- Zugang     IP-Adresse;HBCI-Version;DDV;RDH-1;RDH-2;RDH-3;RDH-4;RDH-5;RDH-6;RDH-7;RDH-8;RDH-9;RDH-10;RAH-7;RAH-9;RAH-10;PIN/TAN-Zugang URL;Version;Datum letzte Änderung;;;;;
2;10010010;Postbank;Berlin;eigenes Rechenzentrum;BdB;;;;;;;;;;;;;;;;;;https://hbci.postbank.de/banking/hbci.do;FinTS V3.0;30.04.2015;;;;;
3;10020200;BHF-Bank AG;Berlin;Bank-Verlag GmbH;BdB;hbciserver.bankverlag.de;nicht unterstützt;3.0;;ja;;ja;;ja;;;;;;;;;https://www.bv-activebanking.de/hbciTunnel/hbciTransfer.jsp;FinTS V3.0;;;;;;`,
			wantResult: []BankInfo{
				{
					BankID:        "10010010",
					VersionNumber: "",
					URL:           "https://hbci.postbank.de/banking/hbci.do",
					VersionName:   "FinTS V3.0",
					Institute:     "Postbank",
					City:          "Berlin",
					LastChanged:   "30.04.2015",
				},
				{
					BankID:        "10020200",
					VersionNumber: "3.0",
					URL:           "https://www.bv-activebanking.de/hbciTunnel/hbciTransfer.jsp",
					VersionName:   "FinTS V3.0",
					Institute:     "BHF-Bank AG",
					City:          "Berlin",
					LastChanged:   "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseBankInfos(strings.NewReader(tt.data))

			if !reflect.DeepEqual(tt.wantError, err) {
				t.Errorf("Expected error to equal\n%v\n\tgot\n%v", tt.wantError, err)
			}

			if !reflect.DeepEqual(tt.wantResult, info) {
				t.Errorf("Results differ:\n%v", pretty.Diff(tt.wantResult, info))
			}
		})
	}
}

func TestFindByBankId(t *testing.T) {
	data = []BankInfo{
		{BankID: "1000000", URL: "1"},
		{BankID: "2000000", URL: "2"},
		{BankID: "3000000", URL: "3"},
	}

	var url string

	url = FindByBankID("1000000").URL

	if url != "1" {
		t.Logf("Expected url to equal %q, got %q\n", "1", url)
		t.Fail()
	}

	url = FindByBankID("2000000").URL

	if url != "2" {
		t.Logf("Expected url to equal %q, got %q\n", "2", url)
		t.Fail()
	}

	url = FindByBankID("3000000").URL

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
		BankInfo{BankID: "30000000"},
		BankInfo{BankID: "10000000"},
		BankInfo{BankID: "20000000"},
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

	if sorted[0].BankID != "10000000" {
		t.Logf("Expected first entry to have BankId '10000000', but was %q\n", sorted[0].BankID)
		t.Fail()
	}

	if sorted[1].BankID != "30000000" {
		t.Logf("Expected first entry to have BankId '30000000', but was %q\n", sorted[1].BankID)
		t.Fail()
	}
}
