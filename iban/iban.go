package iban

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const countryDE = "DE"
const countyCodeDE = 1314

// NewGerman calculates the Iban for the provided bankID and accountID.
// It will return an error if the accountID can not be parsed as int.
// It returns only valid german IBANs, as it is hard coded to use german
// settings.
func NewGerman(bankID, accountID string) (IBAN, error) {
	accountIDAsInt, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		return "", err
	}
	step1 := fmt.Sprintf("%s%010d%d%02d", bankID, accountIDAsInt, countyCodeDE, 0)
	step2 := new(big.Int)
	step2, ok := step2.SetString(step1, 10)
	if !ok {
		return "", fmt.Errorf("Malformed iban string: %s", step1)
	}
	step3 := step2.Mod(step2, big.NewInt(97))
	proofNumber := new(big.Int).Sub(big.NewInt(98), step3)
	iban := fmt.Sprintf("%s%02d%s%010d", countryDE, proofNumber.Int64(), bankID, accountIDAsInt)
	return IBAN(iban), nil
}

// IBAN represents an International Bank Account Number.
// It is defined by ISO 13616:2007
type IBAN string

// BankID returns the parts of the IBAN which refer to the bank institute ID
func (i IBAN) BankID() string {
	return string(i[4:12])
}

// AccountID returns the parts of the IBAN which refer to the AccountID
func (i IBAN) AccountID() string {
	accountID := string(i[12:])
	return strings.TrimLeft(accountID, "0")
}

// CountryCode returns the country code used by the IBAN.
// The country code is defined by ISO 3166-1 alpha-2
func (i IBAN) CountryCode() string {
	return string(i[:2])
}

// ProofNumber returns the used number to make a sanity check whether the IBAN is valid or not.
func (i IBAN) ProofNumber() string {
	return string(i[2:4])
}
