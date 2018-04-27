package cmd

import (
	"fmt"

	"strconv"
	"tracker/helpers"
	"tracker/tracker"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a frame.",
	Long: `Edit a frame.


  You can specify the frame to edit by its position (prefixed with an '@') or by its frame id.
  For example, to edit the second-to-last frame, pass @-2 as the frame index.
  You can get the id of a frame with the 'tracker log' command.

  If no id or index is given, the frame defaults to the current frame or the
  last recorded frame, if no project is currently running.

  You can configure the used editor via
    $ tracker config core.editor.`,
	Run: edit,
}

func init() {
	RootCmd.AddCommand(editCmd)
}

func edit(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Printf("Error: %s\n", helpers.PrintRed("Invalid arguments"))
		return
	}

	opts := tracker.EditFrameOpts{}

	if len(args[0]) == 36 {
		opts.UUID = args[0]
	} else if args[0][0] == '@' {
		if pos, err := strconv.Atoi(args[0][1:]); err != nil {
			opts.Position = -1
		} else {
			opts.Position = pos
		}
	}

	newFrame, err := tracker.EditFrame(opts)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	formattedTags := newFrame.FormattedTags()
	if formattedTags != "" {
		formattedTags = " " + formattedTags
	}

	fmt.Printf(
		"Edited frame for project %s%s, from %s to %s (%s)\n",
		newFrame.FormattedProject(),
		formattedTags,
		newFrame.FormattedStartTime(),
		newFrame.FormattedEndTime(),
		newFrame.Duration().String(),
	)
}
