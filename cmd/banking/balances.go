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

	"github.com/mitch000001/go-hbci/bankinfo"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/iban"
	"github.com/spf13/cobra"
)

var balanceAccount string

// balancesCmd represents the balances command
var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "Fetches balances for a specific account",
	Long: `This command allows to fetch balances for a specific account. By
default it will fetch the balance for the account used to authenticate. For example:

	banking balances --accountID=123456789

will fetch the balance for account 123456789.`,
	Run: func(cmd *cobra.Command, args []string) {
		if balanceAccount == "" {
			balanceAccount = clientConfig.AccountID
		}
		i, err := iban.NewGerman(clientConfig.BankID, balanceAccount)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		account = domain.InternationalAccountConnection{
			IBAN:      string(i),
			BIC:       bankinfo.FindByBankID(clientConfig.BankID).BIC,
			AccountID: balanceAccount,
			BankID:    domain.BankID{CountryCode: 280, ID: clientConfig.BankID},
		}
		if disableSepa {
			balances, err := hbciClient.AccountBalances(account.ToAccountConnection(), true)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(domain.AccountBalances(balances).String())
			return
		}

		balances, err := hbciClient.SepaAccountBalances(account, true, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(domain.SepaAccountBalances(balances).String())
	},
}

func init() {
	rootCmd.AddCommand(balancesCmd)
	balancesCmd.Flags().BoolVar(
		&disableSepa, "disableSepa", false,
		"whether the library should not handle account data as sepa compliant",
	)
	balancesCmd.Flags().StringVar(
		&balanceAccount, "accountID", "",
		"the accountID to fetch balance for (defaults to the UserID)",
	)
}
