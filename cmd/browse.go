package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toqueteos/webbrowser"
	"github.com/florianorben/tracker/helpers"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Open a browser window to your configured backend.",
	Long: `Open a browser window to your configured backend.

  Backend url and token have to be configured via
  $ tracker config backend.url URL
  $ tracker config backend.token TOKEN
  first`,
	Run: browse,
}

func init() {
	RootCmd.AddCommand(browseCmd)
}

func browse(c *cobra.Command, args []string) {
	url := viper.GetString("backend.url")
	token := viper.GetString("backend.token")

	if url == "" || token == "" {
		fmt.Printf("Error: %s\n", helpers.PrintRed("You need to set backend url and token before being able to use 'browse'."))
		fmt.Println("        tracker config backend.url http://some.url")
		fmt.Println("        tracker config backend.token mytoken")
		return
	}

	webbrowser.Open(url)
}
