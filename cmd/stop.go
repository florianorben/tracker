package cmd

import (
	"fmt"
	"time"
	"tracker/helpers"
	"tracker/tracker"

	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop monitoring time for the current project.",
	Long: `Stop monitoring time for the current project.

  Add -m / --message "string" to add a comment to the frame.
  Omit the comment and your configured editor will be started to add the comment.

  Example:

  $ tracker stop
  Stopping project apollo11, started a minute ago. (id: b476b37e-a79e-4e0c-8027-4e57772bcaeb)

  $ tracker stop -m "foo bar"
  Stopping project apollo11, started a minute ago. (id: b476b37e-a79e-4e0c-8027-4e57772bcaeb)`,
	Run: stop,
}

var addMessage bool
var noLog bool

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&addMessage, "message", "m", false, "Add message")
	stopCmd.Flags().BoolVarP(&noLog, "no-log", "n", false, "Don't automatically add a worklog, even if configured otherwise")
}

func stop(cmd *cobra.Command, args []string) {
	message := ""

	if addMessage == true {
		if len(args) > 0 {
			message = strings.Join(args, " ")
		} else {
			var b []byte
			msg, err := helpers.OpenInEditor(b)
			if err != nil {
				message = ""
			}
			message = string(msg)
		}
	}

	var startedFrame tracker.Frame
	var frameIndex int
	frames := tracker.GetFrames()

	for i, frame := range frames {
		if frame.End.IsZero() {
			frames[i].End = time.Now()
			frames[i].Comment = message
			startedFrame = frames[i]
			frameIndex = i
		}
	}

	if startedFrame.Start.IsZero() {
		fmt.Println("Error: " + helpers.PrintRed("No project started"))
	} else {
		if viper.GetBool("backend.autoAddWorkLog") && !noLog {
			startedFrame.AddWorkLog()
			frames[frameIndex] = startedFrame
		}

		frames.Persist()

		formattedTags := startedFrame.FormattedTags()
		if formattedTags != "" {
			formattedTags = " " + formattedTags
		}

		fmt.Printf(
			"Stopping project %s%s, started %s. (id: %s)\n",
			startedFrame.FormattedProject(),
			formattedTags,
			startedFrame.RelativeTime(),
			startedFrame.Uuid,
		)
	}
}
