package cmd

import (
	"os/user"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a system account from SSH Keys Authority",
	Long:  `Remove a system account from SSH Keys Authority control panel.`,
	Run: func(cmd *cobra.Command, args []string) {
		u, err := user.Lookup(username)
		if err != nil || u == nil {
			color.Red("Unable to find user `%s`. Please check the username.", username)
			return
		}

		var accounts []Account
		configErr := viper.UnmarshalKey("accounts", &accounts)

		if configErr != nil {
			panic("There was a problem setting up the user account. Please try again.")
		}

		var updatedAccounts []Account
		for _, data := range accounts {
			if data.Username != username {
				updatedAccounts = append(updatedAccounts, Account{data.Username, data.ApiKey})
			}
		}

		viper.Set("accounts", updatedAccounts)
		viper.WriteConfig()
		color.Green("\nThe selected account has been removed from SSH Keys Authority.")
		color.Green("\nThe authorized_keys file has been left in tact to allow you to manually update it.")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVarP(&username, "username", "u", "", "The username of the system account to add to SSH Keys Authority")
	removeCmd.MarkFlagRequired("user")
}
