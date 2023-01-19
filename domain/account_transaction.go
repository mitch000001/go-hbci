package domain

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

// AccountTransactions represents a printable version of a slice of
// account transactions
type AccountTransactions []AccountTransaction

func (at AccountTransactions) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	if len(at) != 0 {
		first := at[0]
		fmt.Fprintf(
			&buf, "Transactions for account %s/%s\n",
			first.Account.BankID, first.Account.AccountID,
		)
		fmt.Fprintf(
			&buf, "Balance at %s: %.2f %s\n",
			first.AccountBalanceBefore.TransmissionDate.Format("2006-01-02"),
			first.AccountBalanceBefore.Amount.Amount,
			first.AccountBalanceBefore.Amount.Currency,
		)
	}
	buf.WriteString("BookingDate\tBooking Text\tAmount\tBankID\tAccountID\tName\tPurpose")
	buf.WriteString("\n")
	for _, a := range at {
		buf.WriteString(a.BookingDate.Format("2006-01-02"))
		buf.WriteString("\t")
		buf.WriteString(a.BookingText)
		buf.WriteString("\t")
		buf.WriteString(fmt.Sprintf("%.2f %s", a.Amount.Amount, a.Amount.Currency))
		buf.WriteString("\t")
		buf.WriteString(a.BankID)
		buf.WriteString("\t")
		buf.WriteString(a.AccountID)
		buf.WriteString("\t")
		buf.WriteString(a.Name)
		buf.WriteString("\t")
		buf.WriteString(a.Purpose)
		buf.WriteString("\n")
	}
	if len(at) != 0 {
		last := at[len(at)-1]
		fmt.Fprintf(
			&buf, "Balance at %s: %.2f %s\n",
			last.AccountBalanceAfter.TransmissionDate.Format("2006-01-02"),
			last.AccountBalanceAfter.Amount.Amount,
			last.AccountBalanceAfter.Amount.Currency,
		)
	}
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 24, 1, 0, ' ', tabwriter.TabIndent)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}

// AccountTransaction represents one transaction entry for a given account
type AccountTransaction struct {
	Account              AccountConnection
	Amount               Amount
	ValueDate            time.Time
	BookingDate          time.Time
	BookingText          string
	BankID               string
	AccountID            string
	Name                 string
	Purpose              string
	Purpose2             string
	TransactionID        int
	AccountBalanceBefore Balance
	AccountBalanceAfter  Balance
}

func (a AccountTransaction) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BookingDate\tAmount\tBankID\tAccountID\tPurpose")
	buf.WriteString("\n")
	buf.WriteString(a.BookingDate.Format("2006-01-02"))
	buf.WriteString("\t")
	buf.WriteString(fmt.Sprintf("%.2f %s", a.Amount.Amount, a.Amount.Currency))
	buf.WriteString("\t")
	buf.WriteString(a.BankID)
	buf.WriteString("\t")
	buf.WriteString(a.AccountID)
	buf.WriteString("\t")
	buf.WriteString(a.Purpose)
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 0, 8, 0, '\t', 0)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}
