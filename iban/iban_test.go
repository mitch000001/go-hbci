package iban

import "testing"

func TestNew(t *testing.T) {
	bankID := "10090044"
	accountID := "532013018"
	var result IBAN

	result, err := NewGerman(bankID, accountID)

	if err != nil {
		t.Logf("Expected no error, got %T:%v", err, err)
		t.Fail()
	}

	expectedResult := "DE10100900440532013018"

	if string(result) != expectedResult {
		t.Logf("Expected result to equal %q, got %q", expectedResult, result)
		t.Fail()
	}
}

func TestIbanBankId(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	bankID := iban.BankID()

	expectedBankID := "10090044"

	if bankID != expectedBankID {
		t.Logf("Expected bankId to equal %q, got %q\n", expectedBankID, bankID)
		t.Fail()
	}
}

func TestIbanAccountId(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	accountID := iban.AccountID()

	expectedAccountID := "532013018"

	if accountID != expectedAccountID {
		t.Logf("Expected accountId to equal %q, got %q\n", expectedAccountID, accountID)
		t.Fail()
	}
}

func TestIbanCountry(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	country := iban.CountryCode()

	expectedCountry := "DE"

	if country != expectedCountry {
		t.Logf("Expected country to equal %q, got %q\n", expectedCountry, country)
		t.Fail()
	}
}

func TestIbanProofNumber(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	proofNumber := iban.ProofNumber()

	expectedProofNumber := "10"

	if proofNumber != expectedProofNumber {
		t.Logf("Expected proofNumber to equal %q, got %q\n", expectedProofNumber, proofNumber)
		t.Fail()
	}
}
