package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tracker/tracker"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restarts the last stopped time tracking.",
	Long: `Restarts the last stopped time tracking.

  Example:

  $ tracker start someproject +sometag
  // ... wait
  $ tracker stop
  // ... wait
  $ tracker restart

  In this case tracker restart will now start a new time tracking with project "someproject" and tags "sometag".
  i.e. tracker restart will repeat the last not cancelled tracker start command.`,
	Run: restart,
}

func init() {
	RootCmd.AddCommand(restartCmd)
}

func restart(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()

	last := frames[len(frames)-1]
	newFrame := tracker.NewFrame(last.Project, last.Tags)

	frames = append(frames, newFrame)
	frames.Persist()

	fmt.Printf("Starting project %s %s at %s\n", newFrame.FormattedProject(), newFrame.FormattedTags(), newFrame.FormattedStartTime())
}
