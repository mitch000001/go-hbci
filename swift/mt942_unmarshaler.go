package swift

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

// Unmarshal unmarshals value into m
func (m *MT942) Unmarshal(value []byte) error {
	tagExtractor := newTagExtractor(value)
	tags, err := tagExtractor.Extract()
	if err != nil {
		return fmt.Errorf("error extracting tags: %w", err)
	}
	if len(tags) == 0 {
		return fmt.Errorf("malformed marshaled value: no tags found")
	}
	for _, tag := range tags {
		tagID, err := extractTagID(tag)
		if err != nil {
			return fmt.Errorf("error extracting tag ID: %w", err)
		}
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
		case bytes.HasPrefix(tag, []byte(":34F:")):
			minAmount := &MinAmountTag{}
			err = minAmount.Unmarshal(tag)
			if err != nil {
				return fmt.Errorf("error unmarshaling min amount tag: %w", err)
			}
			if m.MinAmount != nil {
				m.MinAmount2 = minAmount
				continue
			}
			m.MinAmount = minAmount
		case bytes.HasPrefix(tag, []byte(":13D:")):
			m.CreationDate = &CreationDateTag{}
			err = m.CreationDate.Unmarshal(tag)
			if err != nil {
				return fmt.Errorf("error unmarshaling creation date tag: %w", err)
			}
		case bytes.HasPrefix(tag, []byte(":13:")):
			creationDateTag := &CreationDateWithoutTimezoneTag{}
			err = creationDateTag.Unmarshal(tag)
			if err != nil {
				return fmt.Errorf("error unmarshaling creation date tag: %w", err)
			}
			m.CreationDate = creationDateTag.CreationDateTag
		case bytes.HasPrefix(tag, []byte(":90D:")):
			m.DebitCountAndBalance = &DebitCountAndBalanceTag{}
			err = m.DebitCountAndBalance.Unmarshal(tag)
			if err != nil {
				return fmt.Errorf("error unmarshaling debit count and balance tag: %w", err)
			}
		case bytes.HasPrefix(tag, []byte(":90C:")):
			m.CreditCountAndBalance = &CreditCountAndBalanceTag{}
			err = m.CreditCountAndBalance.Unmarshal(tag)
			if err != nil {
				return fmt.Errorf("error unmarshaling credit count and balance tag: %w", err)
			}
		case bytes.HasPrefix(tag, []byte(":61:")):
			transaction := &TransactionTag{}
			err = transaction.Unmarshal(tag)
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
			indexLastSliceitem := len(m.Transactions) - 1
			if indexLastSliceitem < 0 {
				return errors.New("Uuexpected CustomTag before first TransactionTag")
			}
			if m.Transactions[indexLastSliceitem].Description != nil {
				return errors.Errorf("unexpected CustomTag: CustomTag would replace Description of %v", m.Transactions[indexLastSliceitem])
			}
			m.Transactions[indexLastSliceitem].Description = customField
		default:
			return fmt.Errorf("malformed marshaled value: unknown tag %q: %s", string(tagID), string(tag))
		}
	}

	return nil
}
