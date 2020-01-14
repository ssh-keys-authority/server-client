package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

type Account struct {
	Username string
	ApiKey   string
}

var keysFileTemplate = []byte(`# This file is managed by SSH Keys Authority.\n
# Any changes made will be overwritten.\n
# If you do not have an account please contact the server owner for assistance.`)

var rootCmd = &cobra.Command{
	Use:   "ssh-server-client",
	Short: "SSH Keys Authority Client",
	Long:  `The SSH Keys Authority Client is an easy to use command line application, allowing your server to automatically sync your teams SSH keys.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AddConfigPath("/etc/ssh-authority-manager")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		color.Red("Your server is not configured to use SSH Keys Authority!")
		fmt.Println("Please follow the instructions on your server details page inside your SSH Keys Authority account.", viper.ConfigFileUsed())
	}
}
