package swift

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

// Unmarshal unmarshals value into m
func (m *MT940) Unmarshal(value []byte) error {
	tagExtractor := newTagExtractor(value)
	tags, err := tagExtractor.Extract()
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	balanceTagOpen := false
	for _, tag := range tags {

		switch {
		case bytes.HasPrefix(tag, []byte(":20:")):
			m.JobReference = &AlphaNumericTag{}
			err = m.JobReference.Unmarshal(tag)
			if err != nil {
				return err
			}
		case bytes.HasPrefix(tag, []byte(":21:")):
			m.Reference = &AlphaNumericTag{}
			err = m.Reference.Unmarshal(tag)
			if err != nil {
				return err
			}
		case bytes.HasPrefix(tag, []byte(":25:")):
			m.Account = &AccountTag{}
			err = m.Account.Unmarshal(tag)
			if err != nil {
				return err
			}
		case bytes.HasPrefix(tag, []byte(":28C:")):
			m.StatementNumber = &StatementNumberTag{}
			err = m.StatementNumber.Unmarshal(tag)
			if err != nil {
				return err
			}
		case bytes.HasPrefix(tag, []byte(":60")):
			m.StartingBalance = &BalanceTag{}
			err = m.StartingBalance.Unmarshal(tag)
			if err != nil {
				return errors.WithMessage(err, "unmarshal starting balance tag")
			}
			balanceTagOpen = true
		case bytes.HasPrefix(tag, []byte(":62")):

			m.ClosingBalance = &BalanceTag{}
			err = m.ClosingBalance.Unmarshal(tag)
			if err != nil {
				return errors.WithMessage(err, "unmarshal closing balance tag")
			}

			balanceTagOpen = false
		case bytes.HasPrefix(tag, []byte(":64:")):
			m.CurrentValutaBalance = &BalanceTag{}
			err = m.CurrentValutaBalance.Unmarshal(tag)
			if err != nil {
				return errors.WithMessage(err, "unmarshal current valuta balance tag")
			}
		case bytes.HasPrefix(tag, []byte(":65:")):
			m.FutureValutaBalance = &BalanceTag{}
			err = m.FutureValutaBalance.Unmarshal(tag)
			if err != nil {
				return errors.WithMessage(err, "unmarshal future valuta balance tag")
			}
		case bytes.HasPrefix(tag, []byte(":61:")):

			transaction := &TransactionTag{}
			err = transaction.Unmarshal(tag, m.StartingBalance.BookingDate.Year())
			if err != nil {
				return err
			}
			m.Transactions = append(m.Transactions, &TransactionSequence{Transaction: transaction})
		case bytes.HasPrefix(tag, []byte(":86:")):
			customField := &CustomFieldTag{}
			err = customField.Unmarshal(tag)
			if err != nil {
				return err
			}
			if balanceTagOpen {
				indexLastSliceitem := len(m.Transactions) - 1
				if indexLastSliceitem < 0 {
					return errors.New("Unexpected CustomTag before first TransactionTag")
				}
				if m.Transactions[indexLastSliceitem].Description != nil {
					return errors.Errorf("Unexpected CustomTag: CustomTag would replace Description of %v", m.Transactions[indexLastSliceitem])
				}
				m.Transactions[indexLastSliceitem].Description = customField
			} else {
				m.CustomField = customField
			}
		default:
			return fmt.Errorf("Malformed marshaled value")
		}
	}

	return nil
}
