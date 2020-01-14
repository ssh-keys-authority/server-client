package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all ssh keys with SSH Keys Authority",
	Long:  `This command will sync the authorized_keys file of each system account you have configured with SSH Keys Authority.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.ReadInConfig()

		var accounts []Account
		configErr := viper.UnmarshalKey("accounts", &accounts)

		if configErr != nil {
			color.Red("There was a problem setting up the user account.\nPlease try again or contact SSH Keys Authority for assistance.")
			os.Exit(1)
		}

		var baseDomain string
		viper.UnmarshalKey("basedomain", &baseDomain)

		if len(baseDomain) <= 0 {
			color.Red("The base domain is missing.\nPlease check you've correctly configured SSH Keys Authority on this server and try again.")
			os.Exit(1)
		}

		var serverAPIKey string
		viper.UnmarshalKey("serverkey", &serverAPIKey)

		if len(serverAPIKey) <= 0 {
			color.Red("The server API key is missing.\nPlease check you've correctly configured SSH Keys Authority on this server and try again.")
			os.Exit(1)
		}

		baseURL := baseDomain + "/api/v1/keys/" + serverAPIKey + "/"

		for _, account := range accounts {
			u, userErr := user.Lookup(account.Username)
			if userErr != nil {
				color.Red("Unable to find user `%s`. Please check the username, and re-create the user on SSH Keys Authority.", username)
				return
			}

			accountAPIURL := baseURL + account.ApiKey
			color.Green("Loading API Key for %s from %s", account.Username, accountAPIURL)

			httpClient := http.Client{
				Timeout: time.Second * 10,
			}

			req, err := http.NewRequest(http.MethodPost, accountAPIURL, nil)
			if err != nil {
				log.Fatal(err)
			}

			res, getErr := httpClient.Do(req)
			if getErr != nil {
				log.Fatal(getErr)
			}

			body, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {
				log.Fatal(readErr)
			}

			keys := string(body)

			validStart := strings.Contains(keys, "START SSH Keys Authority Managed Keys File")
			validEnd := strings.Contains(keys, "END SSH Keys Authority Managed Keys File")

			if !validStart || !validEnd {
				color.Red("The response from the SSH Keys Authority api was invalid.")
				os.Exit(1)
			}

			uid, _ := strconv.Atoi(u.Uid)
			gid, _ := strconv.Atoi(u.Gid)

			homeDir := u.HomeDir
			keysDir := homeDir + "/.ssh"
			keysFile := keysDir + "/authorized_keys"
			color.Green("Writing to " + keysFile)

			if _, keysDirErr := os.Stat(keysDir); os.IsNotExist(keysDirErr) {
				color.Yellow("It looks like " + keysDir + " does not yet exist. Lets create it now.")
				os.MkdirAll(keysDir, 0700)
				os.Chown(keysDir, uid, gid)
			}

			ioutil.WriteFile(keysFile, []byte(keys), 0600)

			os.Chown(keysFile, uid, gid)

			color.Green("Done!")
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
