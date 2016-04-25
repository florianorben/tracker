package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tracker/tracker"
)

var framesCmd = &cobra.Command{
	Use:   "frames",
	Short: "Display a list of all frame IDs.",
	Long: `Display a list of all frame IDs.

  Each frame ID is a unique UUID v4.

  Example:

  $ tracker frames
  8006e400-dbaf-410b-9f83-a75310ebe8f4
  75216856-4a20-4524-97fe-aa462a08a701
  f68b0a8f-034f-4164-aaca-a6c46b20dd80
  [...]`,
	Run: frames,
}

func init() {
	RootCmd.AddCommand(framesCmd)
}

func frames(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	uuids := frames.Frames()

	for _, uuid := range uuids {
		fmt.Printf("%36s\n", uuid)
	}
}
