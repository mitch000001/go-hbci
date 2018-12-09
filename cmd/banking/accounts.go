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

	"github.com/mitch000001/go-hbci/domain"
	"github.com/spf13/cobra"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Lists all accounts associated with the UserID",
	Long: `This command lists all accounts associated with the userID used to
authenticate. For example:

	banking accounts

will list all accounts for the current userID.`,
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := hbciClient.Accounts()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(domain.AccountInfos(accounts))
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
}
