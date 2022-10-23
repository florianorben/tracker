package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/florianorben/tracker/helpers"
	"github.com/florianorben/tracker/tracker"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Display a report of the time spent on each project.",
	Long: `A longer description that spans multiple lines and likely contains examples

  If a project is given, the time spent on this project is printed. Else,
  print the total for each root project.

  By default, the time spent the last 7 days is printed. This timespan can
  be controlled with the '--from' and '--to' arguments. The dates must have
  the format 'YYYY-MM-DD', like: '2014-05-19'.

  You can limit the report to a project or a tag using the '--project' and
  '--tag' options. They can be specified several times each to add multiple
  projects or tags to the report.

  Example:

  $ tracker report
  Mon 05 May 2014 -> Mon 12 May 2014

  apollo11 - 13h 22m 20s
          [brakes    7h 53m 18s]
          [module    7h 41m 41s]
          [reactor   8h 35m 50s]
          [steering 10h 33m 37s]
          [wheels   10h 11m 35s]

  hubble - 8h 54m 46s
          [camera        8h 38m 17s]
          [lens          5h 56m 22s]
          [transmission  6h 27m 07s]

  voyager1 - 11h 45m 13s
          [antenna     5h 53m 57s]
          [generators  9h 04m 58s]
          [probe      10h 14m 29s]
          [sensors    10h 30m 26s]

  voyager2 - 16h 16m 09s
          [antenna     7h 05m 50s]
          [generators 12h 20m 29s]
          [probe      12h 20m 29s]
          [sensors    11h 23m 17s]

  Total: 43h 42m 20s

  $ tracker report --from 2014-04-01 --to 2014-04-30 --project apollo11
  Tue 01 April 2014 -> Wed 30 April 2014

  apollo11 - 13h 22m 20s
          [brakes    7h 53m 18s]
          [module    7h 41m 41s]
          [reactor   8h 35m 50s]
          [steering 10h 33m 37s]
          [wheels   10h 11m 35s]`,
	Run: report,
}

func init() {
	RootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringP("from", "f", "", "Start date (YYYY-MM-DD)")
	reportCmd.Flags().StringP("to", "t", "", "End date (YYYY-MM-DD)")
	reportCmd.Flags().StringSliceP("project", "p", []string{}, `Reports activity only for the given project.
                        You can add other projects by using this option several times.`)
	reportCmd.Flags().StringSliceP("tag", "T", []string{}, `Reports activity only for frames containing the given tag.
                        You can add several tags by using this option multiple times`)

}

func report(cmd *cobra.Command, args []string) {
	projects, err := cmd.Flags().GetStringSlice("project")
	if err != nil {
		projects = make([]string, 0)
	}

	tags, err := cmd.Flags().GetStringSlice("tag")
	if err != nil {
		tags = make([]string, 0)
	}

	query, err := tracker.NewFrameQuery(
		projects,
		tags,
		cmd.Flag("from").Value.String(),
		cmd.Flag("to").Value.String(),
	)

	if err != nil {
		fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
		return
	}

	tracker.Report(query)
}
