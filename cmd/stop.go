package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time"
	"tracker/helpers"
	"tracker/tracker"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop monitoring time for the current project.",
	Long: `Stop monitoring time for the current project.

  Example:

  $ tracker stop
  Stopping project apollo11, started a minute ago. (id: b476b37e-a79e-4e0c-8027-4e57772bcaeb)`,
	Run: stop,
}

func init() {
	RootCmd.AddCommand(stopCmd)
}

func stop(cmd *cobra.Command, args []string) {
	var startedFrame tracker.Frame
	frames := tracker.GetFrames()
	for i, frame := range frames {
		if frame.End.IsZero() {
			startedFrame = frame
			frames[i].End = time.Now()
		}
	}

	if startedFrame.Start.IsZero() {
		fmt.Println("Error: " + helpers.PrintRed("No project started"))
	} else {
		frames.Persist()
		fmt.Printf(
			"Stopping project %s %s, started %s. (id: %s)\n",
			startedFrame.FormattedProject(),
			startedFrame.FormattedTags(),
			startedFrame.RelativeTime(),
			startedFrame.Uuid,
		)
	}
}
