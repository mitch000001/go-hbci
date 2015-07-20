package dialog

import (
	"net/http"
	"reflect"
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

	syncResponseMessage := encryptedTestMessage(
		"HIRMG:2:2:1+0100::Dialog beendet'",
		"HISYN:2:3:8+newClientSystemID'",
		"HIUPD:3:4:8+12345::280:1000000+54321+EUR+Muster+Max+++HKTAN:1+HKKAZ:1'",
	)

	dialogEndResponseMessage := encryptedTestMessage("HIRMG:2:2:1+0020::Auftrag entgegengenommen'")

	transport.SetResponsePayloads([][]byte{
		syncResponseMessage,
		dialogEndResponseMessage,
	})

	if len(dialog.Accounts) != 0 {
		t.Logf("Expected no accounts, got %d\n", len(dialog.Accounts))
		t.Fail()
	}

	dialogID := dialog.dialogID
	if dialogID != initialDialogID {
		t.Logf("Expected dialogID to equal\n%q\n\tgot\n%q\n", initialDialogID, dialogID)
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

	// message errors
	syncResponseMessage = encryptedTestMessage("HIRMG:2:2:1+9000::Nachricht enthält Fehler'")

	transport.SetResponsePayloads([][]byte{
		syncResponseMessage,
		[]byte(""),
	})

	res, err = dialog.SyncClientSystemID()

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	}

	if err != nil {
		errMessage := err.Error()
		expectedMessage := "Institute returned errors:\nMessageAcknowledgement: Code: 9000, Position: , Text: Nachricht enthält Fehler, Parameter: "
		if expectedMessage != errMessage {
			t.Logf("Expected error to equal\n%q\n\tgot\n%q\n", expectedMessage, errMessage)
			t.Fail()
		}
	}
}
