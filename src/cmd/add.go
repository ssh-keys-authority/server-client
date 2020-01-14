package cmd

import (
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/spf13/viper"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var username string
var apikey string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a system account to SSH Keys Authority",
	Long:  `Add a new system account (e.g root) to SSH Keys Authority to have it's SSH Keys automatically managed.`,
	Run: func(cmd *cobra.Command, args []string) {
		u, err := user.Lookup(username)
		if err != nil {
			color.Red("Unable to find user `%s`. Please check the username, and re-create the user on SSH Keys Authority.", username)
			os.Exit(1)
		}

		color.Green("Found system user: %s\nSetting up SSH Keys Authority for the account.", u.Username)

		viper.ReadInConfig()

		var accounts []Account
		configErr := viper.UnmarshalKey("accounts", &accounts)

		if configErr != nil {
			panic("There was a problem setting up the user account. Please try again.")
		}

		accounts = append(accounts, Account{username, apikey})
		viper.Set("accounts", accounts)
		viper.WriteConfig()

		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)

		homeDir := u.HomeDir
		keysDir := homeDir + "/.ssh"
		keysFile := keysDir + "/authorized_keys"
		backupKeysFile := keysDir + "/authorized_keys.bak"

		if _, keysDirErr := os.Stat(keysDir); os.IsNotExist(keysDirErr) {
			color.Yellow("It looks like " + keysDir + " does not yet exist. Lets create it now.")
			os.MkdirAll(keysDir, 0700)
			os.Chown(keysDir, uid, gid)
		}

		if _, keysFileErr := os.Stat(keysFile); !os.IsNotExist(keysFileErr) {
			backupFileErr := os.Rename(keysFile, backupKeysFile)
			if backupFileErr != nil {
				log.Fatal(backupFileErr)
			}
			color.Yellow("An existing authorized_keys file was found. This has been moved to %s", backupKeysFile)

		}

		file, fileError := os.Create(keysFile)
		if fileError != nil {
			panic("Unable to create authorized_keys file. Please check that the user you are running the agent as has the correct privileges.")
		}

		file.Write(keysFileTemplate)

		os.Chown(keysFile, uid, gid)

		color.Green("The user was successfully configured and is now managed by SSH Keys Authority.")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&username, "username", "u", "", "The username of the system account to add to SSH Keys Authority")
	addCmd.MarkFlagRequired("user")

	addCmd.Flags().StringVarP(&apikey, "apikey", "k", "", "The unique API Key for the system account, provided when adding the account via your SSH Keys Authority control panel.")
	addCmd.MarkFlagRequired("api-key")
}
