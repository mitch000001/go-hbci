// Copyright © 2015 Michael Wagner <mitch.wagna@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/iban"
	"github.com/spf13/cobra"
)

var daysToFetch int
var transactionsAccount string
var disableSepa bool

// transactionsCmd represents the transactions command
var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "fetch transactions for an account",
	Long: `This command allows to fetch transactions for a specific account. By
default it will fetch the transactions for the account used to authenticate
for the last ten days. For example:

	banking transactions --accountID=123456789 --daysToFetch=30

will fetch booked transactions for account 123456789 for the last 30 days.`,
	Run: func(cmd *cobra.Command, args []string) {
		if transactionsAccount == "" {
			transactionsAccount = clientConfig.AccountID
		}
		timeframe := domain.Timeframe{
			StartDate: domain.NewShortDate(time.Now().AddDate(0, 0, -daysToFetch)),
		}

		i, err := iban.NewGerman(clientConfig.BankID, transactionsAccount)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		account = domain.InternationalAccountConnection{
			IBAN:      string(i),
			AccountID: transactionsAccount,
			BankID:    domain.BankID{CountryCode: 280, ID: clientConfig.BankID},
		}

		if disableSepa {
			fetchTransactions(account, timeframe)
			return
		}

		fetchSepaTransactions(account, timeframe)
	},
}

func fetchTransactions(account domain.InternationalAccountConnection, timeframe domain.Timeframe) {
	transactions, err := hbciClient.AccountTransactions(account.ToAccountConnection(), timeframe, false, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(domain.AccountTransactions(transactions))
}

func fetchSepaTransactions(account domain.InternationalAccountConnection, timeframe domain.Timeframe) {
	transactions, err := hbciClient.SepaAccountTransactions(account, timeframe, false, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(domain.AccountTransactions(transactions))
}

func init() {
	rootCmd.AddCommand(transactionsCmd)

	transactionsCmd.Flags().BoolVar(
		&disableSepa, "disableSepa", false,
		"whether the library should not handle account data as sepa compliant",
	)
	transactionsCmd.Flags().IntVar(
		&daysToFetch, "daysToFetch", 10,
		"the number of days to fetch transactions for (10)",
	)
	transactionsCmd.Flags().StringVar(
		&transactionsAccount, "accountID", "",
		"the accountID to fetch transactions for (defaults to the UserID)",
	)
}
