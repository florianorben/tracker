package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tracker/tracker"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display when the current project was started and the time spent since.",
	Long: `Display when the current project was started and the time spent since.

	Example:

	$ tracker start test +foo
	$ tracker status
	Project test [foo] started just now (2016-04-25 21:38:00 +0200 CEST)`,
	Run: status,
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	var activeFrame tracker.Frame

	for _, frame := range frames {
		if frame.InProgress() {
			activeFrame = frame
		}
	}

	if activeFrame.Start.IsZero() {
		fmt.Println("No project started")
	} else {
		fmt.Printf(
			"Project %s %s started %s (%s)\n",
			activeFrame.FormattedProject(),
			activeFrame.FormattedTags(),
			activeFrame.RelativeTime(),
			activeFrame.Start,
		)
	}
}
