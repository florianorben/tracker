package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tracker/helpers"
	"tracker/tracker"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Display a list of all tags.",
	Long: `Display a list of all tags.

  Example:

  $ tracker tags
  ARGUSII-7100
  SN-4
  meetings
  private
  QA-1498
  [...]`,
	Run: tags,
}

func init() {
	RootCmd.AddCommand(tagsCmd)

	tagsCmd.Flags().BoolP("no-color", "b", false, "No color mode")
}

func tags(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	tags := frames.Tags()

	var noColor bool
	var err error
	if noColor, err = cmd.Flags().GetBool("no-color"); err != nil {
		noColor = false
	}

	for _, tag := range tags {
		if noColor {
			fmt.Println(tag)
		} else {
			fmt.Println(helpers.PrintBlue(tag))
		}
	}
}
