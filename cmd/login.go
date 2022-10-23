package cmd

import (
	"fmt"
	"github.com/florianorben/tracker/helpers"
	"github.com/florianorben/tracker/tracker"
	"os"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save login credentials to backend (jira).",
	Long:  `Save login credentials to backend (jira).`,
	Run:   loginFrame,
}

func init() {
	RootCmd.AddCommand(loginCmd)
}

func loginFrame(cmd *cobra.Command, args []string) {
	service := "tracker"
	user := helpers.GetUserInput("Username: ")

	if user == "" {
		os.Exit(0)
	}

	if _, err := keyring.Get(service, user); err != nil {
		if err != keyring.ErrNotFound {
			fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
			os.Exit(1)
		}

	} else {
		if !helpers.AskForConfirmation("User " + user + " already saved, set new password? [y/N]: ") {
			return
		}
	}

	pass := helpers.GetUserInputPassword("Password: ")
	if pass == "" {
		os.Exit(0)
	}

	tracker.SetConfig("backend.user", user)
	err := tracker.WriteConfigToFile()
	if err != nil {
		fmt.Printf("Error: %s: %s\n", helpers.PrintRed("configuration could not be saved to disk"), err.Error())
		os.Exit(1)
	}

	err = keyring.Set(service, user, pass)
	if err != nil {
		fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
		os.Exit(1)
	}

	fmt.Printf("%s\n", helpers.PrintGreen("Login saved!"))
}
