package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/florianorben/tracker/helpers"
	"github.com/florianorben/tracker/tracker"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Display a list of all existing projects.",
	Long: `Display a list of all existing projects.

  Example:

  $ tracker projects
  ARGUSII
  SN
  VAN
  private
  [...]
`,
	Run: projects,
}

func init() {
	RootCmd.AddCommand(projectsCmd)

	projectsCmd.Flags().BoolP("no-color", "b", false, "No color mode")
}

func projects(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	projects := frames.Projects()

	var noColor bool
	var err error
	if noColor, err = cmd.Flags().GetBool("no-color"); err != nil {
		noColor = false
	}

	for _, project := range projects {
		if noColor {
			fmt.Println(project)
		} else {
			fmt.Println(helpers.PrintPurple(project))
		}
	}
}
