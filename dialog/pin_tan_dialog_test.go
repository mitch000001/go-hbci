package dialog

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
)

func TestPinTanDialogSyncClientSystemID(t *testing.T) {
	transport := &MockHttpTransport{}
	httpClient := &http.Client{Transport: transport}

	url := "http://localhost"
	clientID := "12345"
	bankID := domain.BankId{280, "10000000"}
	dialog := NewPinTanDialog(bankID, url, clientID)
	dialog.SetPin("abcde")
	dialog.httpClient = httpClient

	encryptedData := []string{
		"HISYN:2:3:8+newClientSystemID'",
		"HIUPD:3:4:8+12345::280:1000000+54321+EUR+Muster+Max+++HKTAN:1+HKKAZ:1'",
	}
	syncResponseMessage := []string{
		"HNHBK:1:3+000000000123+220+abcde+1+'",
		"HNVSK:998:2:+998+1+1::0+1:20150713:173634+2:2:13:@8@\x00\x00\x00\x00\x00\x00\x00\x00:5:1:+280:10000000:12345:V:0:0+0+'",
		fmt.Sprintf("HNVSD:999:1:+@%d@%s'", len(strings.Join(encryptedData, "")), strings.Join(encryptedData, "")),
		"HNHBS:3:1:+1'",
	}

	transport.SetResponsePayloads([][]byte{
		[]byte(strings.Join(syncResponseMessage, "")),
		[]byte(""),
	})

	if len(dialog.Accounts) != 0 {
		t.Logf("Expected no accounts, got %d\n", len(dialog.Accounts))
		t.Fail()
	}

	res, err := dialog.SyncClientSystemID()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := "newClientSystemID"

	if res != expected {
		t.Logf("Expected response to equal\n%q\n\tgot\n%q\n", expected, res)
		t.Fail()
	}

	if dialog.ClientSystemID != expected {
		t.Logf("Expected ClientSystemID to equal %q, got %q\n", expected, dialog.ClientSystemID)
		t.Fail()
	}

	accounts := dialog.Accounts

	if len(accounts) == 0 {
		t.Logf("Expected %d accounts, got 0\n", 1)
		t.Fail()
	}

	if len(accounts) > 0 {
		account := accounts[0]
		expected := domain.AccountInformation{
			AccountConnection: &domain.AccountConnection{AccountID: "12345", CountryCode: 280, BankID: "1000000"},
			UserID:            "54321",
			Currency:          "EUR",
			Name1:             "Muster",
			Name2:             "Max",
			AllowedBusinessTransactions: []domain.BusinessTransaction{
				domain.BusinessTransaction{ID: "HKTAN", NeededSignatures: 1},
				domain.BusinessTransaction{ID: "HKKAZ", NeededSignatures: 1},
			},
		}
		if !reflect.DeepEqual(expected, account) {
			t.Logf("Expected account to eqaul\n%+#v\n\tgot\n%+#v\n", expected, account)
			t.Fail()
		}
	}
}
