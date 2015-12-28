package iban

import "testing"

func TestNew(t *testing.T) {
	bankId := "10090044"
	accountId := "532013018"
	var result Iban

	result, err := New(bankId, accountId)

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
	iban := Iban("DE10100900440532013018")

	bankId := iban.BankId()

	expectedBankId := "10090044"

	if bankId != expectedBankId {
		t.Logf("Expected bankId to equal %q, got %q\n", expectedBankId, bankId)
		t.Fail()
	}
}

func TestIbanAccountId(t *testing.T) {
	iban := Iban("DE10100900440532013018")

	accountId := iban.AccountId()

	expectedAccountId := "532013018"

	if accountId != expectedAccountId {
		t.Logf("Expected accountId to equal %q, got %q\n", expectedAccountId, accountId)
		t.Fail()
	}
}

func TestIbanCountry(t *testing.T) {
	iban := Iban("DE10100900440532013018")

	country := iban.Country()

	expectedCountry := "DE"

	if country != expectedCountry {
		t.Logf("Expected country to equal %q, got %q\n", expectedCountry, country)
		t.Fail()
	}
}

func TestIbanProofNumber(t *testing.T) {
	iban := Iban("DE10100900440532013018")

	proofNumber := iban.ProofNumber()

	expectedProofNumber := "10"

	if proofNumber != expectedProofNumber {
		t.Logf("Expected proofNumber to equal %q, got %q\n", expectedProofNumber, proofNumber)
		t.Fail()
	}
}
