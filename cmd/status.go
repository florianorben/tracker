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
	Project test [foo] started just now (2016-04-25 21:38:00 +0200 CEST)
    
    You can display a shorter version of status with the -short flag
    
    Example:
    
    $ tracker start test +foo
    $ tracker status
    Project test started just now`,
	Run: status,
}

func init() {
	RootCmd.AddCommand(statusCmd)

	statusCmd.Flags().StringP("short", "s", "", "Truncated output")
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
		formattedTags := activeFrame.FormattedTags()
		if formattedTags != "" {
			formattedTags = " " + formattedTags
		}

        if cmd.Flag("short").Value.String() != "" {
            fmt.Printf("Project %s started %s\n", activeFrame.FormattedProject(), activeFrame.RelativeTime())
        } else {
		    fmt.Printf(
		        "Project %s%s started %s (%s)\n",
			    activeFrame.FormattedProject(),
		        formattedTags,
		        activeFrame.RelativeTime(),
		        activeFrame.Start,
		    )
        }
	}
}
