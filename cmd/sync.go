package cmd

import (
	"fmt"

	"tracker/tracker"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pushes all unsynced workloads to jira.",
	Long: `Pushes all unsynced workloads to jira.

  Pushes all unsynced workloads to jira

  Example:

  $ tracker login
  $ tracker sync
  Sync successul`,
	Run: sync,
}

func init() {
	RootCmd.AddCommand(syncCmd)
}

func sync(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()

	for i, frame := range frames {
		frame.AddWorkLog()
		frames[i] = frame
	}

	frames.Persist()

	fmt.Println("Sync successul")
}
