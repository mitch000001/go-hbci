package swift

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mitch000001/go-hbci/domain"
)

type MT940 struct {
	JobReference         *AlphaNumericTag
	Reference            *AlphaNumericTag
	AccountID            *AlphaNumericTag
	StatementNumber      *StatementNumberTag
	StartingBalance      *BalanceTag
	Transactions         []*TransactionSequence
	ClosingBalance       *BalanceTag
	CurrentValutaBalance *BalanceTag
	FutureValutaBalance  *BalanceTag
	CustomField          *AlphaNumericTag
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
	date, err := parseDate(dateBytes)
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
	CustomTag   *AlphaNumericTag
}

type TransactionTag struct {
	Tag                   string
	Date                  domain.ShortDate
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
	date, err := parseDate(dateBytes)
	if err != nil {
		return err
	}
	t.Date = domain.NewShortDate(date)
	r, _, err := buf.ReadRune()
	if err != nil {
		return err
	}
	if unicode.IsDigit(r) {
		buf.UnreadRune()
		dateBytes = buf.Next(4)
		date, err = parseDate(dateBytes)
		if err != nil {
			return err
		}
		t.BookingDate = domain.NewShortDate(date)
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
	if bytes.IndexByte(buf.Bytes(), '/') == -1 {
		t.Reference = buf.String()
	} else {
		if bytes.Index(buf.Bytes(), []byte("//")) != -1 {
			ref, err := buf.ReadString('/')
			if err != nil {
				return err
			}
			if len(ref) > 1 {
				t.Reference = ref[:len(ref)-1]
			}
		}
	}
	return nil
}

func parseDate(value []byte) (time.Time, error) {
	var offset int
	if len(value) == 6 {
		offset = 2
	} else {
		offset = 4
	}
	yearBegin := fmt.Sprintf("%d", time.Now().Year())[:offset]
	dateString := yearBegin + string(value)
	date, err := time.Parse("20060102", dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date.Truncate(24 * time.Hour), nil
}
