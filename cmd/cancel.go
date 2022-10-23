package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/florianorben/tracker/helpers"
	"github.com/florianorben/tracker/tracker"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel the last call to the start command.",
	Long:  `Cancel the last call to the start command. The time will not be recorded.`,
	Run:   cancelFrame,
}

func init() {
	RootCmd.AddCommand(cancelCmd)
}

func cancelFrame(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	newFrames := make(tracker.Frames, 0, len(frames)-1)

	for _, frame := range frames {
		if !frame.InProgress() {
			newFrames = append(newFrames, frame)
		} else {
			formattedTags := frame.FormattedTags()
			if formattedTags != "" {
				formattedTags = " " + formattedTags
			}

			fmt.Printf("Canceling the timer for project %s%s\n", frame.FormattedProject(), formattedTags)
		}
	}

	if len(newFrames) == len(frames) {
		fmt.Printf("Error: %s\n", helpers.PrintRed("No project started."))
		return
	}

	newFrames.Persist()

}
