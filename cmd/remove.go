package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"tracker/helpers"
	"tracker/tracker"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a frame.",
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
		ok := askForConfirmation(
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

func askForConfirmation(q string) bool {
	fmt.Print(q)

	var response string
	read, err := fmt.Scanln(&response)
	if err != nil && read != 0 {
		fmt.Println("Error: " + helpers.PrintRed(err.Error()))
		return false
	} else if read == 0 {
		return false
	}

	response = strings.ToLower(response)

	if response == "y" || response == "yes" {
		return true
	} else if response == "n" || response == "no" || response == "" {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation(q)
	}
}
