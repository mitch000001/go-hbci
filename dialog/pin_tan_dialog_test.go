package dialog

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
)

func TestPinTanDialogSyncClientSystemID(t *testing.T) {
	transport := &MockHttpsTransport{}

	url := "http://localhost"
	clientID := "12345"
	bankID := domain.BankId{280, "10000000"}
	dialog := NewPinTanDialog(bankID, url, clientID)
	dialog.SetPin("abcde")
	dialog.transport = transport

	syncResponseMessage := encryptedTestMessage(
		"newDialogID",
		"HIRMG:2:2:1+0100::Dialog beendet'",
		"HIBPA:3:2:+12+280:10000000+Bank Name+3+1+201:210:220+0'",
		"HISYN:4:3:8+newClientSystemID'",
		"HIUPA:5:2:7+12345+4+0'",
		"HIUPD:6:4:8+12345::280:1000000+54321+EUR+Muster+Max+++HKTAN:1+HKKAZ:1'",
	)

	dialogEndResponseMessage := encryptedTestMessage("newDialogID", "HIRMG:2:2:1+0020::Auftrag entgegengenommen'")

	transport.SetResponseMessages([][]byte{
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
		t.Fail()
	}

	bankParamData := dialog.BankParameterData
	if bankParamData.Version != 0 {
		t.Logf("Expected BPD version to equal 0, was %d\n", bankParamData.Version)
		t.Fail()
	}

	userParamData := dialog.UserParameterData
	if userParamData.Version != 0 {
		t.Logf("Expected UPD version to equal 0, was %d\n", userParamData.Version)
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

	bankParamData = dialog.BankParameterData
	if bankParamData.Version != 12 {
		t.Logf("Expected BankParameterData version to equal 12, was %d\n", bankParamData.Version)
		t.Fail()
	}

	userParamData = dialog.UserParameterData
	if userParamData.Version != 4 {
		t.Logf("Expected UPD version to equal 4, was %d\n", userParamData.Version)
		t.Fail()
	}

	dialogID = dialog.dialogID
	if dialogID != "newDialogID" {
		t.Logf("Expected dialogID to equal\n%q\n\tgot\n%q\n", "newDialogID", dialogID)
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
	syncResponseMessage = encryptedTestMessage("ABCDE", "HIRMG:2:2:1+9000::Nachricht enthält Fehler'")

	transport.SetResponseMessages([][]byte{
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

func TestPinTanDialogInit(t *testing.T) {
	transport := &MockHttpsTransport{}

	url := "http://localhost"
	clientID := "12345"
	bankID := domain.BankId{280, "10000000"}
	dialog := NewPinTanDialog(bankID, url, clientID)
	dialog.SetPin("abcde")
	dialog.transport = transport

	dialogID := dialog.dialogID
	if dialogID != initialDialogID {
		t.Logf("Expected dialogID to equal\n%q\n\tgot\n%q\n", initialDialogID, dialogID)
		t.Fail()
	}

	initResponse := encryptedTestMessage(
		"newDialogID",
		"HIRMG:2:2:1+0020::Auftrag entgegengenommen'",
		"HIKIM:10:2+ec-Karte+Ihre neue ec-Karte liegt zur Abholung bereit.'",
	)
	transport.SetResponseMessage(initResponse)

	err := dialog.Init()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	dialogID = dialog.dialogID
	if dialogID != "newDialogID" {
		t.Logf("Expected dialogID to equal\n%q\n\tgot\n%q\n", "newDialogID", dialogID)
		t.Fail()
	}
}
