package iban

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
)

const countryDE = "DE"
const germanAccountIDLength = 10
const maxAllowedIBANLength = 34

// IsValid validates the IBAN for its proof number
func IsValid(input string) bool {
	iban, err := From(input)
	if err != nil {
		return false
	}
	return iban.Valid()
}

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
	if len(countryCode) != 2 {
		return "", fmt.Errorf("malformed countryCode: must have two characters")
	}
	countryCode = strings.ToUpper(countryCode)

	ibanLen, ok := ibanMaxLengthForCountry[countryCode]
	if !ok {
		ibanLen = maxAllowedIBANLength
	}
	if len(BBAN) == 0 {
		return "", fmt.Errorf("missing BBAN")
	}
	if len(BBAN) > (ibanLen - 4) {
		return "", fmt.Errorf("malformed BBAN: must have at max %d characters", ibanLen)
	}

	BBAN = strings.ToUpper(BBAN)

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

// From creates an IBAN from the provided input. If the input is not a valid IBAN, it
// returns an error.
func From(input string) (IBAN, error) {
	countryCode := input[:2]
	proofNumber := input[2:4]
	bban := input[4:]
	iban, err := New(countryCode, bban)
	if err != nil {
		return "", err
	}
	if proofNumber != iban.ProofNumber() {
		return "", fmt.Errorf("proof number invalid")
	}
	return iban, nil
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

// Valid returns true if the IBAN is valid, false otherwise
func (i IBAN) Valid() bool {
	if len(i) > maxAllowedIBANLength {
		return false
	}
	countryCode := i.CountryCode()
	countryCode = strings.ToUpper(countryCode)
	bban := i.BBAN()
	bban = strings.ToUpper(bban)

	digitCountryCode := transformLettersToDigits(countryCode)
	digitBBAN := transformLettersToDigits(bban)
	checkSumString := fmt.Sprintf("%s%s%s", digitBBAN, digitCountryCode, i.ProofNumber())
	checkSum, ok := new(big.Int).SetString(checkSumString, 10)
	if !ok {
		return false
	}
	mod := checkSum.Mod(checkSum, big.NewInt(97))
	return mod.Cmp(big.NewInt(1)) == 0
}

// String returns the string representation of i
func (i IBAN) String() string {
	return string(i)
}

// Print returns the IBAN in paper format, i.e. with spaces after every fourth
// character
func Print(iban IBAN) string {
	ibanStr := iban.String()
	for len(ibanStr)%4 != 0 {
		ibanStr += " "
	}
	var out bytes.Buffer
	for i := 4; i <= len(ibanStr); i += 4 {
		out.WriteString(ibanStr[i-4 : i])
		out.WriteString(" ")
	}
	return strings.TrimSpace(out.String())
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

var ibanMaxLengthForCountry = map[string]int{
	"AL": 28,
	"AD": 24,
	"AZ": 28,
	"BH": 22,
	"BE": 16,
	"BA": 20,
	"BR": 29,
	"VG": 24,
	"BG": 22,
	"CR": 22,
	"DK": 18,
	"DE": 22,
	"DO": 28,
	"SV": 28,
	"EE": 20,
	"FO": 18,
	"FI": 18,
	"FR": 27,
	"GE": 22,
	"GI": 23,
	"GR": 27,
	"GL": 18,
	"GB": 22,
	"GT": 28,
	"IQ": 23,
	"IE": 22,
	"IS": 26,
	"IL": 23,
	"IT": 27,
	"JO": 30,
	"KZ": 20,
	"QA": 29,
	"XK": 20,
	"HR": 21,
	"KW": 30,
	"LV": 21,
	"LB": 28,
	"LI": 21,
	"LT": 20,
	"LU": 20,
	"MT": 31,
	"MR": 27,
	"MU": 30,
	"MK": 19,
	"MD": 24,
	"MC": 27,
	"ME": 22,
	"NL": 18,
	"NO": 15,
	"AT": 20,
	"PK": 24,
	"PS": 29,
	"PL": 28,
	"PT": 25,
	"RO": 24,
	"LC": 32,
	"SM": 27,
	"ST": 25,
	"SA": 24,
	"SE": 24,
	"CH": 21,
	"RS": 22,
	"SC": 31,
	"SK": 24,
	"SI": 19,
	"ES": 24,
	"TL": 23,
	"TR": 26,
	"CZ": 24,
	"TN": 24,
	"UA": 29,
	"HU": 28,
	"AE": 23,
	"BY": 28,
	"CY": 28,
}
var sepaSupportForCountry = map[string]bool{
	"AL": false,
	"AD": false,
	"AZ": false,
	"BH": false,
	"BE": true,
	"BA": false,
	"BR": false,
	"VG": false,
	"BG": true,
	"CR": false,
	"DK": true,
	"DE": true,
	"DO": false,
	"SV": false,
	"EE": true,
	"FO": true,
	"FI": true,
	"FR": true,
	"GE": false,
	"GI": true,
	"GR": true,
	"GL": true,
	"GB": true,
	"GT": false,
	"IQ": false,
	"IE": true,
	"IS": true,
	"IL": false,
	"IT": true,
	"JO": false,
	"KZ": false,
	"QA": false,
	"XK": false,
	"HR": true,
	"KW": false,
	"LV": true,
	"LB": false,
	"LI": true,
	"LT": true,
	"LU": true,
	"MT": true,
	"MR": false,
	"MU": false,
	"MK": false,
	"MD": false,
	"MC": true,
	"ME": false,
	"NL": true,
	"NO": true,
	"AT": true,
	"PK": false,
	"PS": false,
	"PL": true,
	"PT": true,
	"RO": true,
	"LC": false,
	"SM": true,
	"ST": false,
	"SA": false,
	"SE": true,
	"CH": true,
	"RS": false,
	"SC": false,
	"SK": true,
	"SI": true,
	"ES": true,
	"TL": false,
	"TR": false,
	"CZ": true,
	"TN": false,
	"UA": false,
	"HU": true,
	"AE": false,
	"BY": false,
	"CY": true,
}
