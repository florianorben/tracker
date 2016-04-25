package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tracker/helpers"
	"tracker/tracker"
)

var configCmd = &cobra.Command{
	Use:   "config SECTION.OPTION [VALUE]",
	Short: "Get and set configuration options.",
	Long: `Get and set configuration options.

  If value is not provided, the content of the key is displayed. Else, the
  given value is set.

  You can edit the config file with an editor with the '--edit' option.

  Example:

  $ tracker config backend.token 7e329263e329
  $ tracker config backend.token
  7e329263e329`,
	Run: config,
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolP("edit", "e", false, "Edit the configuration file with an editor.")
}

func config(cmd *cobra.Command, args []string) {
	if e, err := cmd.Flags().GetBool("edit"); err == nil && e == true {
		if err := tracker.EditConfig(); err != nil {
			fmt.Printf("Error: %s: %s\n", helpers.PrintRed("configuration file could not be opened/saved"), err.Error())
		}
		return
	}

	if len(args) == 0 || len(args) > 2 {
		fmt.Printf(cmd.UsageString())
		return
	}

	if len(args) == 1 {
		printConfig(args[0])
		return
	}

	tracker.SetConfig(args[0], args[1])
	err := tracker.WriteConfigToFile()
	if err != nil {
		fmt.Printf("Error: %s: %s\n", helpers.PrintRed("configuration could not be saved to disk"), err.Error())
	}
}

func printConfig(configKey string) {
	configValue := viper.Get(configKey)
	switch configValue.(type) {
	case map[string]interface{}:
		for key, val := range configValue.(map[string]interface{}) {
			fmt.Printf("%s: %s\n", helpers.PrintBold(configKey+"."+key), val)
		}
	case nil:
		fmt.Println(helpers.PrintBold(helpers.PrintBold("Configuration value " + configKey + " does not exist")))
	default:
		fmt.Printf("%s: %s\n", helpers.PrintBold(configKey), configValue)
	}
}
