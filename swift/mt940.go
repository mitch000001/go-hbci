package swift

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mitch000001/go-hbci/domain"
)

type MT940 struct {
	JobReference         *AlphaNumericTag
	Reference            *AlphaNumericTag
	Account              *AccountTag
	StatementNumber      *StatementNumberTag
	StartingBalance      *BalanceTag
	Transactions         []*TransactionSequence
	ClosingBalance       *BalanceTag
	CurrentValutaBalance *BalanceTag
	FutureValutaBalance  *BalanceTag
	CustomField          *CustomFieldTag
}

func (m *MT940) AccountTransactions() []domain.AccountTransaction {
	accountConnection := domain.AccountConnection{BankID: m.Account.BankID, AccountID: m.Account.AccountID, CountryCode: 280}
	var transactions []domain.AccountTransaction
	for _, transactionSequence := range m.Transactions {
		tr := transactionSequence.Transaction
		descr := transactionSequence.Description
		var amount float64
		if tr.DebitCreditIndicator == "D" {
			amount = -tr.Amount
		} else {
			amount = tr.Amount
		}
		transaction := domain.AccountTransaction{
			Account:     accountConnection,
			Amount:      domain.Amount{amount, m.StartingBalance.Currency},
			ValutaDate:  tr.ValutaDate.Time,
			BookingDate: tr.BookingDate.Time,
			AccountBalanceBefore: domain.Balance{
				Amount: domain.Amount{
					m.StartingBalance.Amount,
					m.StartingBalance.Currency,
				},
				TransmissionDate: m.StartingBalance.BookingDate.Time,
			},
			AccountBalanceAfter: domain.Balance{
				Amount: domain.Amount{
					m.ClosingBalance.Amount,
					m.ClosingBalance.Currency,
				},
				TransmissionDate: m.ClosingBalance.BookingDate.Time,
			},
		}
		if descr != nil {
			transaction.BankID = descr.BankID
			transaction.AccountID = descr.AccountID
			transaction.Purpose = descr.Purpose
			transaction.Purpose2 = descr.Purpose2
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

type AccountTag struct {
	Tag       string
	BankID    string
	AccountID string
}

func (a *AccountTag) Unmarshal(value []byte) error {
	elements, err := ExtractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", a)
	}
	a.Tag = string(elements[0])
	fields := bytes.Split(elements[1], []byte("/"))
	if len(fields) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", a)
	}
	a.BankID = string(fields[0])
	a.AccountID = string(fields[1])
	return nil
}

type StatementNumberTag struct {
	Tag         string
	Number      int
	SheetNumber int
}

func (s *StatementNumberTag) Unmarshal(value []byte) error {
	elements, err := ExtractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", s)
	}
	s.Tag = string(elements[0])
	var numBytes []byte
	if bytes.IndexByte(elements[1], '/') != -1 {
		buf := bytes.NewBuffer(elements[1])
		numBytes, err = buf.ReadBytes('/')
		if err != nil {
			return err
		}
		numBytes = numBytes[:len(numBytes)-1]
		sheetNum, err := strconv.Atoi(buf.String())
		if err != nil {
			return err
		}
		s.SheetNumber = sheetNum
	} else {
		numBytes = elements[1]
	}
	num, err := strconv.Atoi(string(numBytes))
	if err != nil {
		return err
	}
	s.Number = num
	return nil
}

type BalanceTag struct {
	Tag                  string
	DebitCreditIndicator string
	BookingDate          domain.ShortDate
	Currency             string
	Amount               float64
}

func (b *BalanceTag) Unmarshal(value []byte) error {
	elements, err := ExtractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", b)
	}
	b.Tag = string(elements[0])
	buf := bytes.NewBuffer(elements[1])
	b.DebitCreditIndicator = string(buf.Next(1))
	dateBytes := buf.Next(6)
	date, err := parseDate(dateBytes, time.Now().Year())
	if err != nil {
		return err
	}
	b.BookingDate = domain.NewShortDate(date)
	b.Currency = string(buf.Next(3))
	amountString := strings.Replace(buf.String(), ",", ".", 1)
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return err
	}
	b.Amount = amount
	return nil
}

type TransactionSequence struct {
	Transaction *TransactionTag
	Description *CustomFieldTag
}

type TransactionTag struct {
	Tag                   string
	ValutaDate            domain.ShortDate
	BookingDate           domain.ShortDate
	DebitCreditIndicator  string
	CurrencyKind          string
	Amount                float64
	BookingKey            string
	Reference             string
	BankReference         string
	AdditionalInformation string
}

func (t *TransactionTag) Unmarshal(value []byte) error {
	elements, err := ExtractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", t)
	}
	t.Tag = string(elements[0])
	buf := bytes.NewBuffer(elements[1])
	dateBytes := buf.Next(6)
	date, err := parseDate(dateBytes, time.Now().Year())
	if err != nil {
		return err
	}
	t.ValutaDate = domain.NewShortDate(date)
	r, _, err := buf.ReadRune()
	if err != nil {
		return err
	}
	if unicode.IsDigit(r) {
		buf.UnreadRune()
		dateBytes = buf.Next(4)
		date, err = parseDate(dateBytes, t.ValutaDate.Year())
		if err != nil {
			return err
		}
		t.BookingDate = domain.NewShortDate(date)
		monthDiff := int(math.Abs(float64(t.ValutaDate.Month() - t.BookingDate.Month())))
		if monthDiff > 1 {
			t.BookingDate = domain.NewShortDate(t.BookingDate.AddDate(1, 0, 0))
		}
	}
	var runes []rune
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return err
		}
		runes = append(runes, r)
		if unicode.IsDigit(r) {
			buf.UnreadRune()
			runes = runes[:len(runes)-1]
			if len(runes) == 3 {
				t.DebitCreditIndicator = string(runes[:2])
				t.CurrencyKind = string(runes[2:])
			} else if len(runes) == 2 {
				t.DebitCreditIndicator = string(runes[:1])
				t.CurrencyKind = string(runes[1:])
			} else if len(runes) == 1 {
				t.DebitCreditIndicator = string(runes)
			} else {
				return fmt.Errorf("%T: Malformed marshaled value", t)
			}
			break
		}
	}
	amountBytes, err := buf.ReadBytes('N')
	if err != nil {
		return err
	}
	amountBytes = bytes.Replace(amountBytes[:len(amountBytes)-1], []byte(","), []byte("."), 1)
	amount, err := strconv.ParseFloat(string(amountBytes), 64)
	if err != nil {
		return err
	}
	t.Amount = amount
	t.BookingKey = string(buf.Next(3))
	remaining := buf.String()
	addInfSepIdx := strings.Index(remaining, "\r\n/")
	doubleSlashIdx := strings.Index(remaining, "//")

	if doubleSlashIdx != -1 && addInfSepIdx != -1 {
		t.Reference = remaining[:doubleSlashIdx]
		if doubleSlashIdx < addInfSepIdx {
			t.BankReference = remaining[doubleSlashIdx+2 : addInfSepIdx]
			t.AdditionalInformation = remaining[addInfSepIdx+3:]
		} else {
			// The only valid case in the FINTS30 documenation in the other one, but the data we receive are sometimes
			// formatted like that :(
			t.BankReference = remaining[addInfSepIdx+3 : doubleSlashIdx]
			t.AdditionalInformation = remaining[doubleSlashIdx+2:]
		}
	} else {
		t.Reference = remaining
		if doubleSlashIdx != -1 {
			t.Reference = remaining[:doubleSlashIdx]
			t.BankReference = remaining[doubleSlashIdx+2:]
		}
		if addInfSepIdx != -1 {
			t.Reference = remaining[:addInfSepIdx]
			t.AdditionalInformation = remaining[addInfSepIdx+3:]
		}
	}
	return nil
}

func parseDate(value []byte, referenceYear int) (time.Time, error) {
	var offset int
	if len(value) == 6 {
		offset = 2
	} else {
		offset = 4
	}
	yearBegin := fmt.Sprintf("%d", referenceYear)[:offset]
	dateString := yearBegin + string(value)
	date, err := time.Parse("20060102", dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date.Truncate(24 * time.Hour), nil
}
