package domain

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// AccountInfos represent a printable version of a slice of AccountInformation
type AccountInfos []AccountInformation

func (ai AccountInfos) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BankID\tAccountID\tUserID\tCurrency\tName\tProductID\tLimit")
	buf.WriteString("\n")
	for _, a := range ai {
		buf.WriteString(a.AccountConnection.BankID)
		buf.WriteString("\t")
		buf.WriteString(a.AccountConnection.AccountID)
		buf.WriteString("\t")
		buf.WriteString(a.UserID)
		buf.WriteString("\t")
		buf.WriteString(a.Currency)
		buf.WriteString("\t")
		fmt.Fprintf(&buf, "%s, %s", a.Name1, a.Name2)
		buf.WriteString("\t")
		buf.WriteString(a.ProductID)
		if a.Limit != nil {
			buf.WriteString("\t")
			fmt.Fprintf(&buf, "%s: %.2f %s", a.Limit.Kind, a.Limit.Amount.Amount, a.Limit.Amount.Currency)
		} else {
			buf.WriteString("\t - ")
		}
		buf.WriteString("\n")
	}
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 20, 1, 0, ' ', tabwriter.TabIndent)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}

// AccountInformation represents bank specific information about an account
type AccountInformation struct {
	AccountConnection           AccountConnection
	UserID                      string
	Currency                    string
	Name1                       string
	Name2                       string
	ProductID                   string
	Limit                       *AccountLimit
	AllowedBusinessTransactions []BusinessTransaction
}

func (a AccountInformation) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BankID\tAccountID\tUserID\tCurrency\tName\tProductID\tLimit")
	buf.WriteString("\n")
	buf.WriteString(a.AccountConnection.BankID)
	buf.WriteString("\t")
	buf.WriteString(a.AccountConnection.AccountID)
	buf.WriteString("\t")
	buf.WriteString(a.UserID)
	buf.WriteString("\t")
	buf.WriteString(a.Currency)
	buf.WriteString("\t")
	fmt.Fprintf(&buf, "%s, %s", a.Name1, a.Name2)
	buf.WriteString("\t")
	buf.WriteString(a.ProductID)
	if a.Limit != nil {
		buf.WriteString("\t")
		fmt.Fprintf(&buf, "%s: %.2f %s", a.Limit.Kind, a.Limit.Amount.Amount, a.Limit.Amount.Currency)
	} else {
		buf.WriteString("\t - ")
	}
	buf.WriteString("\n")
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 20, 1, 0, ' ', tabwriter.Debug)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}
