package cmd

import (
	"fmt"
	"strconv"
	"github.com/florianorben/tracker/helpers"
	"github.com/florianorben/tracker/tracker"

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
	if len(args) == 0 {
		fmt.Printf("Error: %s\n", helpers.PrintRed("Invalid arguments"))
		return
	}

	frames := tracker.GetFrames()
	newFrames := make(tracker.Frames, 0, len(frames)-1)
	deletedIndexes := make([]int, 0)

	for _, arg := range args {
		index := -1
		var oldFrame tracker.Frame

		if len(arg) == 36 {
			index, oldFrame = frames.ByUUID(arg)
		} else if arg[0] == '@' {
			if pos, err := strconv.Atoi(args[0][1:]); err != nil {
				index = -1
			} else {
				index, oldFrame = frames.ByPosition(pos)
			}
		}

		if index == -1 {
			fmt.Printf("Error: %s %s", helpers.PrintRed("No frame found with id."), args[1])
			continue
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
				continue
			}
		}

		deletedIndexes = append(deletedIndexes, index)
	}

	for i, frame := range frames {
		in := false
		for _, j := range deletedIndexes {
			if i == j {
				in = true
			}
		}

		if !in {
			newFrames = append(newFrames, frame)
		}
	}

	newFrames.Persist()
	fmt.Println("Frame(s) deleted.")
}
