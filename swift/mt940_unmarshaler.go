package swift

import (
	"bytes"
	"fmt"
)

func (m *MT940) Unmarshal(value []byte) error {
	tagExtractor := NewTagExtractor(value)
	tags, err := tagExtractor.Extract()
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	transactions := make(map[int]*TransactionTag)
	customTags := make(map[int]*CustomFieldTag)
	closingBalanceIndex := -1
	for i, tag := range tags {
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
				return err
			}
		case bytes.HasPrefix(tag, []byte(":62")):
			m.ClosingBalance = &BalanceTag{}
			err = m.ClosingBalance.Unmarshal(tag)
			if err != nil {
				return err
			}
		case bytes.HasPrefix(tag, []byte(":64:")):
			m.CurrentValutaBalance = &BalanceTag{}
			err = m.CurrentValutaBalance.Unmarshal(tag)
			if err != nil {
				return err
			}
		case bytes.HasPrefix(tag, []byte(":65:")):
			m.FutureValutaBalance = &BalanceTag{}
			err = m.FutureValutaBalance.Unmarshal(tag)
			if err != nil {
				return err
			}
		case bytes.HasPrefix(tag, []byte(":61:")):
			transaction := &TransactionTag{}
			err = transaction.Unmarshal(tag)
			if err != nil {
				return err
			}
			transactions[i] = transaction
		case bytes.HasPrefix(tag, []byte(":86:")):
			customField := &CustomFieldTag{}
			err = customField.Unmarshal(tag)
			if err != nil {
				return err
			}
			customTags[i] = customField
		default:
			return fmt.Errorf("Malformed marshaled value")
		}
	}
	for idx, transaction := range transactions {
		sequence := &TransactionSequence{Transaction: transaction}
		if cust, ok := customTags[idx+1]; ok {
			sequence.Description = cust
		}
		m.Transactions = append(m.Transactions, sequence)
	}
	for k, cust := range customTags {
		if k > closingBalanceIndex {
			m.CustomField = cust
		}
	}
	return nil
}
