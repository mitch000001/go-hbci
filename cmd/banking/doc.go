// banking is a CLI to access bank data via HBCI
//
// Usage:
//   banking [command]
//
// Available Commands:
//   accounts     Lists all accounts associated with the UserID
//   balances     Fetches balances for a specific account
//   help         Help about any command
//   transactions fetch transactions for an account
//
// Flags:
//       --blz string        the identifier for the bank institute
//       --config string     config file (default is $HOME/.banking.yaml)
//   -d, --debug             enable debug logging (very verbose)
//       --hbci.url string   the URL to the bank institute
//   -h, --help              help for banking
//       --pin string        the pin for the provided account
//       --userID string     the account ID to authenticate with
//
// Use "banking [command] --help" for more information about a command.
package main
