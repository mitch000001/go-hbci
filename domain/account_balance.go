package domain

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

// AccountBalances represents a printable version of a slice of
// account balances
type AccountBalances []AccountBalance

func (ab AccountBalances) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("Product Name\tBankID\tAccountID\t\tBookingDate\tAmount\tLimit")
	buf.WriteString("\n")
	containsEarmarkedBalances := false
	for _, a := range ab {
		buf.WriteString(a.ProductName)
		buf.WriteString("\t")
		buf.WriteString(a.Account.BankID)
		buf.WriteString("\t")
		buf.WriteString(a.Account.AccountID)
		buf.WriteString("\t")
		buf.WriteString("\t")
		buf.WriteString(a.BookedBalance.TransmissionDate.Format("2006-01-02"))
		buf.WriteString("\t")
		buf.WriteString(fmt.Sprintf("%.2f %s", a.BookedBalance.Amount.Amount, a.BookedBalance.Amount.Currency))
		if a.EarmarkedBalance != nil && a.EarmarkedBalance.Amount.Amount != 0.0 {
			fmt.Fprintf(&buf, " (%.2f %s)*", a.EarmarkedBalance.Amount.Amount, a.EarmarkedBalance.Amount.Currency)
			containsEarmarkedBalances = true
		}
		buf.WriteString("\t")
		if a.CreditLimit != nil {
			buf.WriteString(fmt.Sprintf("%.2f %s", a.CreditLimit.Amount, a.CreditLimit.Currency))
		} else {
			buf.WriteString(" - ")
		}
		buf.WriteString("\t")
		buf.WriteString("\n")
	}
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 20, 1, 0, ' ', tabwriter.TabIndent)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	if containsEarmarkedBalances {
		fmt.Fprintln(&out, "* earmarked transactions")
	}
	return out.String()
}

// AccountBalance represents a balance for a specific account
type AccountBalance struct {
	Account          AccountConnection
	ProductName      string
	Currency         string
	BookedBalance    Balance
	EarmarkedBalance *Balance
	CreditLimit      *Amount
	AvailableAmount  *Amount
	UsedAmount       *Amount
	BookingDate      *time.Time
	DueDate          *time.Time
}

func (a AccountBalance) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BankID\tAccountID\t\tBookingDate\tAmount")
	buf.WriteString("\n")
	buf.WriteString(a.Account.BankID)
	buf.WriteString("\t")
	buf.WriteString(a.Account.AccountID)
	buf.WriteString("\t")
	buf.WriteString("\t")
	buf.WriteString(a.BookedBalance.TransmissionDate.Format("2006-01-02"))
	buf.WriteString("\t")
	buf.WriteString(fmt.Sprintf("%.2f %s", a.BookedBalance.Amount.Amount, a.BookedBalance.Amount.Currency))
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 24, 1, 0, ' ', tabwriter.TabIndent)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}
