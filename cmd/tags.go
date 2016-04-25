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
}

func tags(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	tags := frames.Tags()

	for _, tag := range tags {
		fmt.Println(helpers.PrintBlue(tag))
	}
}
