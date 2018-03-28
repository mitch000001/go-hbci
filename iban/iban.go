package iban

import (
	"fmt"
	"math/big"
	"strings"
)

const countryDE = "DE"
const germanAccountIDLength = 10
const maxAllowedIBANLength = 34

// NewGerman calculates the Iban for the provided bankID and accountID.
// It will return an error if the accountID can not be parsed as int.
// It returns only valid german IBANs, as it is hard coded to use german
// settings.
func NewGerman(bankID, accountID string) (IBAN, error) {
	if len(accountID) == (germanAccountIDLength - 1) {
		accountID = "0" + accountID
	}

	bban := fmt.Sprintf("%s%s", bankID, accountID)

	return New(countryDE, bban)
}

// New returns a new IBAN for the provided countryCode and BBAN
func New(countryCode string, BBAN string) (IBAN, error) {
	step1 := fmt.Sprintf("%s%s00", transformLettersToDigits(BBAN), transformLettersToDigits(countryCode))
	step2, ok := new(big.Int).SetString(step1, 10)
	if !ok {
		return "", fmt.Errorf("Malformed iban string: %s", step1)
	}
	step3 := step2.Mod(step2, big.NewInt(97))
	proofNumber := new(big.Int).Sub(big.NewInt(98), step3)

	iban := fmt.Sprintf("%s%02d%s", countryCode, proofNumber.Int64(), BBAN)
	return IBAN(iban), nil
}

// IsValid validates the IBAN for its proof number
func IsValid(iban IBAN) bool {
	if len(iban) > maxAllowedIBANLength {
		return false
	}
	countryCode := iban.CountryCode()
	digitCountryCode := transformLettersToDigits(countryCode)
	checkSumString := fmt.Sprintf("%s%s%s", iban.BBAN(), digitCountryCode, iban.ProofNumber())
	checkSum, ok := new(big.Int).SetString(checkSumString, 10)
	if !ok {
		return false
	}
	mod := checkSum.Mod(checkSum, big.NewInt(97))
	return mod.Cmp(big.NewInt(1)) == 0
}

// IBAN represents an International Bank Account Number.
// It is defined by ISO 13616:2007
type IBAN string

// BBAN returns the BBAN, that is, the bank account identifier
func (i IBAN) BBAN() string {
	return string(i[4:])
}

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

func transformLettersToDigits(letters string) string {
	var replaced []string
	for k, v := range alphaToDigit {
		replaced = append(replaced, k)
		replaced = append(replaced, v)
	}
	return strings.NewReplacer(replaced...).Replace(letters)
}

var alphaToDigit = map[string]string{
	"A": "10",
	"B": "11",
	"C": "12",
	"D": "13",
	"E": "14",
	"F": "15",
	"G": "16",
	"H": "17",
	"I": "18",
	"J": "19",
	"K": "20",
	"L": "21",
	"M": "22",
	"N": "23",
	"O": "24",
	"P": "25",
	"Q": "26",
	"R": "27",
	"S": "28",
	"T": "29",
	"U": "30",
	"V": "31",
	"W": "32",
	"X": "33",
	"Y": "34",
	"Z": "35",
}
