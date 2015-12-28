package iban

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const countryDE = "DE"
const countyCodeDE = 1314

func New(bankId, accountId string) (Iban, error) {
	accountID, err := strconv.ParseInt(accountId, 10, 64)
	if err != nil {
		return "", err
	}
	step1 := fmt.Sprintf("%s%010d%d%02d", bankId, accountID, countyCodeDE, 0)
	step2 := new(big.Int)
	step2, ok := step2.SetString(step1, 10)
	if !ok {
		return "", fmt.Errorf("Malformed iban string: %s", step1)
	}
	step3 := step2.Mod(step2, big.NewInt(97))
	proofNumber := new(big.Int).Sub(big.NewInt(98), step3)
	iban := fmt.Sprintf("%s%02d%s%010d", countryDE, proofNumber.Int64(), bankId, accountID)
	return Iban(iban), nil
}

type Iban string

func (i Iban) BankId() string {
	return string(i[4:12])
}

func (i Iban) AccountId() string {
	accountId := string(i[12:])
	return strings.TrimLeft(accountId, "0")
}

func (i Iban) Country() string {
	return string(i[:2])
}

func (i Iban) ProofNumber() string {
	return string(i[2:4])
}
