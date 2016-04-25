package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tracker/helpers"
	"tracker/tracker"
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
}

func projects(cmd *cobra.Command, args []string) {
	frames := tracker.GetFrames()
	projects := frames.Projects()

	for _, project := range projects {
		fmt.Println(helpers.PrintPurple(project))
	}
}
