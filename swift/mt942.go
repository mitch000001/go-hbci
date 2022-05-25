package swift

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/internal"
)

type MT942 struct {
	JobReference          *AlphaNumericTag    `swift:"20"`
	Reference             *AlphaNumericTag    `swift:"21"`
	Account               *AccountTag         `swift:"25"`
	StatementNumber       *StatementNumberTag `swift:"28C"`
	MinAmount             *MinAmountTag       `swift:"34F"`
	MinAmount2            *MinAmountTag       `swift:"34F"`
	CreationDate          *CreationDateTag    `swift:"13D"`
	Transactions          []*TransactionSequence
	DebitCountAndBalance  *DebitCountAndBalanceTag  `swift:"90D"`
	CreditCountAndBalance *CreditCountAndBalanceTag `swift:"90C"`
}

// AccountTransactions returns a slice of account transactions created from m
func (m *MT942) UnbookedAccountTransactions() domain.UnbookedAccountTransactions {
	accountConnection := domain.AccountConnection{BankID: m.Account.BankID, AccountID: m.Account.AccountID, CountryCode: 280}
	var transactions []domain.Transaction
	for _, transactionSequence := range m.Transactions {
		tr := transactionSequence.Transaction
		descr := transactionSequence.Description
		var amount float64
		if tr.DebitCreditIndicator == "D" {
			amount = -tr.Amount
		} else {
			amount = tr.Amount
		}
		transaction := domain.Transaction{
			Amount:      domain.Amount{Amount: amount, Currency: m.MinAmount.Currency},
			ValutaDate:  tr.ValutaDate.Time,
			BookingDate: tr.BookingDate.Time,
		}
		if descr != nil {
			transaction.BookingText = descr.BookingText
			transaction.BankID = descr.BankID
			transaction.AccountID = descr.AccountID
			transaction.Name = descr.Name
			transaction.Purpose = strings.Join(descr.Purpose, " ")
			transaction.Purpose2 = strings.Join(descr.Purpose2, " ")
			transaction.TransactionID = descr.TransactionID
		}
		transactions = append(transactions, transaction)
	}
	internal.Info.Printf("MT942: %v", m)
	unbookedTransactions := domain.UnbookedAccountTransactions{
		Account:      accountConnection,
		Transactions: transactions,
	}
	if m.CreationDate != nil {
		unbookedTransactions.CreationDate = m.CreationDate.Date
	}
	if m.DebitCountAndBalance != nil {
		unbookedTransactions.DebitAmount = domain.Amount{Amount: m.DebitCountAndBalance.Amount, Currency: m.DebitCountAndBalance.Currency}
		unbookedTransactions.DebitTransactions = m.DebitCountAndBalance.Count
	}
	if m.CreditCountAndBalance != nil {
		unbookedTransactions.CreditAmount = domain.Amount{Amount: m.CreditCountAndBalance.Amount, Currency: m.CreditCountAndBalance.Currency}
		unbookedTransactions.CreditTransactions = m.CreditCountAndBalance.Count
	}
	return unbookedTransactions
}

type MinAmountTag struct {
	Tag                  string
	Currency             string
	DebitCreditIndicator string
	Amount               float64
}

// Unmarshal unmarshals value into m
func (m *MinAmountTag) Unmarshal(value []byte) error {
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", m)
	}
	m.Tag = string(elements[0])
	buf := bytes.NewBuffer(elements[1])
	m.Currency = string(buf.Next(3))
	m.DebitCreditIndicator = string(buf.Next(1))
	amountString := strings.Replace(buf.String(), ",", ".", 1)
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return fmt.Errorf("MT942 Balance tag: error unmarshaling amount: %w", err)
	}
	m.Amount = amount
	return nil
}

type CreationDateTag struct {
	Tag  string
	Date time.Time
}

// Unmarshal unmarshals value into c
func (c *CreationDateTag) Unmarshal(value []byte) error {
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", c)
	}
	c.Tag = string(elements[0])
	buf := bytes.NewBuffer(elements[1])
	dateBytes := buf.Next(6)
	timeStr := string(buf.Next(4))
	timezoneSign := string(buf.Next(1))
	timezoneDifference := string(buf.Next(4))
	date, err := time.Parse(
		"200601021504-0700",
		fmt.Sprintf("20%s%s%s%s", string(dateBytes), timeStr, timezoneSign, timezoneDifference),
	)
	if err != nil {
		return fmt.Errorf("unmarshal creation date tag: parsing date: %w", err)
	}
	c.Date = date
	return nil
}

type CreationDateWithoutTimezoneTag struct {
	*CreationDateTag
}

// Unmarshal unmarshals value into c
func (c *CreationDateWithoutTimezoneTag) Unmarshal(value []byte) error {
	c.CreationDateTag = &CreationDateTag{}
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", c)
	}
	c.Tag = string(elements[0])
	buf := bytes.NewBuffer(elements[1])
	dateBytes := buf.Next(6)
	timeStr := string(buf.Next(4))
	date, err := time.Parse(
		"200601021504",
		fmt.Sprintf("20%s%s", string(dateBytes), timeStr),
	)
	if err != nil {
		return fmt.Errorf("unmarshal creation date tag: parsing date: %w", err)
	}
	c.Date = date
	return nil
}

type DebitCountAndBalanceTag struct {
	Tag      string
	Count    int
	Currency string
	Amount   float64
}

// Unmarshal unmarshals value into d
func (d *DebitCountAndBalanceTag) Unmarshal(value []byte) error {
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", d)
	}
	d.Tag = string(elements[0])
	buf := bytes.NewBuffer(elements[1])
	var countRunes []rune
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return err
		}
		countRunes = append(countRunes, r)
		if !unicode.IsDigit(r) {
			buf.UnreadRune()
			countRunes = countRunes[:len(countRunes)-1]
			break
		}
	}
	count, err := strconv.Atoi(string(countRunes))
	if err != nil {
		return fmt.Errorf("error parsing count: %w", err)
	}
	d.Count = count
	d.Currency = string(buf.Next(3))
	amountString := strings.Replace(buf.String(), ",", ".", 1)
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return fmt.Errorf("MT940 Balance tag: error unmarshaling amount: %w", err)
	}
	d.Amount = amount
	return nil
}

type CreditCountAndBalanceTag struct {
	Tag      string
	Count    int
	Currency string
	Amount   float64
}

// Unmarshal unmarshals value into d
func (d *CreditCountAndBalanceTag) Unmarshal(value []byte) error {
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", d)
	}
	d.Tag = string(elements[0])
	buf := bytes.NewBuffer(elements[1])
	var countRunes []rune
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return err
		}
		countRunes = append(countRunes, r)
		if !unicode.IsDigit(r) {
			buf.UnreadRune()
			countRunes = countRunes[:len(countRunes)-1]
			break
		}
	}
	count, err := strconv.Atoi(string(countRunes))
	if err != nil {
		return fmt.Errorf("error parsing count: %w", err)
	}
	d.Count = count
	d.Currency = string(buf.Next(3))
	amountString := strings.Replace(buf.String(), ",", ".", 1)
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return fmt.Errorf("MT940 Balance tag: error unmarshaling amount: %w", err)
	}
	d.Amount = amount
	return nil
}
