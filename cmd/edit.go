package cmd

import (
	"fmt"

	"encoding/json"
	"github.com/spf13/cobra"
	"strconv"
	"time"
	"tracker/helpers"
	"tracker/tracker"
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

	frames := tracker.GetFrames()
	index := -1
	var frame tracker.Frame

	if len(args[0]) == 36 {
		index, frame = frames.ByUUID(args[0])
	} else if args[0][0] == '@' {
		if pos, err := strconv.Atoi(args[0][1:]); err != nil {
			index = -1
		} else {
			index, frame = frames.ByPosition(pos)
		}
	}

	if index == -1 {
		fmt.Printf("Error: %s %s.\n", helpers.PrintRed("No frame found with id"), args[0])
		return
	}

	b, err := json.MarshalIndent(&struct {
		Start   string   `json:"start"`
		End     string   `json:"end"`
		Project string   `json:"project"`
		Tags    []string `json:"tags"`
	}{
		Start:   frame.Start.Format(tracker.DateTimeFormat),
		End:     frame.End.Format(tracker.DateTimeFormat),
		Project: frame.Project,
		Tags:    frame.Tags,
	}, "", "  ")
	if err != nil {
		fmt.Printf("Error: Creating temp file failed: %s\n", helpers.PrintRed(err.Error()))
		return
	}

	newFrameContents, err := helpers.OpenInEditor(b)
	if err != nil {
		return
	}

	var newFrame tracker.Frame
	json.Unmarshal(newFrameContents, &newFrame)
	newFrame.Uuid = frame.Uuid
	newFrame.LastEdit = time.Now()

	if frames[index].Equals(newFrame) {
		fmt.Println("No changes made.")
		return
	}

	frames[index] = newFrame
	frames.Persist()

	fmt.Printf(
		"Edited frame for project %s %s, from %s to %s (%s)\n",
		newFrame.FormattedProject(),
		newFrame.FormattedTags(),
		newFrame.FormattedStartTime(),
		newFrame.FormattedEndTime(),
		newFrame.Duration().String(),
	)
}
