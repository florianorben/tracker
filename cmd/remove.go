package cmd

import (
	"fmt"

	"strconv"
	"tracker/helpers"
	"tracker/tracker"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove a frame.",
	Long: `Remove a frame.

  You can specify the frame either by id or by position
  (ex: '@-1' for the last frame).`,
	Run: remove,
}

var force bool

func init() {
	RootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "Don't ask for confirmation")
}

func remove(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Printf("Error: %s\n", helpers.PrintRed("Invalid arguments"))
		return
	}

	frames := tracker.GetFrames()
	index := -1
	var oldFrame tracker.Frame

	if len(args[0]) == 36 {
		index, oldFrame = frames.ByUUID(args[0])
	} else if args[0][0] == '@' {
		if pos, err := strconv.Atoi(args[0][1:]); err != nil {
			index = -1
		} else {
			index, oldFrame = frames.ByPosition(pos)
		}
	}

	if index == -1 {
		fmt.Printf("Error: %s %s", helpers.PrintRed("No frame found with id."), args[1])
		return
	}

	if force == false {
		ok := helpers.AskForConfirmation(
			fmt.Sprintf(
				"You are about to remove frame %s from %s to %s, continue? [y/N]: ",
				oldFrame.Uuid,
				oldFrame.FormattedStartTime(),
				oldFrame.FormattedEndTime(),
			),
		)

		if ok == false {
			fmt.Println("aborted!")
			return
		}
	}

	newFrames := make(tracker.Frames, 0, len(frames)-1)
	for i, frame := range frames {
		if i != index {
			newFrames = append(newFrames, frame)
		}
	}

	newFrames.Persist()
	fmt.Println("Frame deleted.")
}
