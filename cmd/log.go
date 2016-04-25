package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tracker/helpers"
	"tracker/tracker"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Display each recorded session during the given timespan.",
	Long: `Display each recorded session during the given timespan.

  By default, the sessions from the last 14 days are printed. This timespan
  can be controlled with the '--from' and '--to' arguments. The dates must
  have the format 'YYYY-MM-DD', like: '2014-05-19'.
  Default values can be set via
    $ tracker config log.defaultStartDate
    $ tracker config log.defaultEndDate

  You can limit the log to a project or a tag using the '--project' and
  '--tag' options. They can be specified several times each to add multiple
  projects or tags to the log.

  Example:

  $ tracker log --project voyager2 --project apollo11
  Thursday 08 May 2015 (56m 33s)
          01406ee8-210a-4eb9-b498-00851592c541  09:26 to 10:22      56m 33s  apollo11  [reactor, brakes, steering, wheels, module]

  Wednesday 07 May 2015 (27m 29s)
          1aadf9eb-e0bd-4cbe-aa47-1f2ea805d5e7  09:48 to 10:15      27m 29s  voyager2  [sensors, generators, probe]

  Tuesday 06 May 2015 (1h 47m 22s)
          2ee999b1-23ad-49a8-aa5f-d6766bee39b8  12:40 to 14:16   1h 35m 45s  apollo11  [wheels]
          85b17756-677b-4495-9bda-3b6a4ab1aedd  14:23 to 14:35      11m 37s  apollo11  [brakes, steering]

  Monday 05 May 2015 (8h 18m 26s)
          9df1b32e-3dc8-4a17-bff7-8495dcee8781  09:05 to 10:03      57m 12s  voyager2  [probe, generators]
          c0774766-77aa-4224-8b34-b61bf314c78a  10:51 to 14:47   3h 55m 40s  apollo11
          d4c0157f-d399-43c8-a402-15bda600e919  15:12 to 18:38   3h 25m 34s  voyager2  [probe, generators, sensors, antenna]

  $ tracker log --from 2014-04-16 --to 2014-04-17
  Thursday 17 April 2014 (4h 19m 13s)
          ddf12a4c-5690-4d99-a7d7-cdbbaf2d0563  09:15 to 09:43      28m 11s    hubble  [lens, camera, transmission]
          fd7dd185-1e48-48f3-aeeb-2e3bf45af584  10:19 to 12:59   2h 39m 15s    hubble  [camera, transmission]
          c0774766-77aa-4224-8b34-b61bf314c78a  14:42 to 15:54   1h 11m 47s  voyager1  [antenna]

  Wednesday 16 April 2014 (5h 19m 18s)
          b476b37e-a79e-4e0c-8027-4e57772bcaeb  09:53 to 12:43   2h 50m 07s  apollo11  [wheels]
          514a9690-e734-42f7-9588-4f925d1c986a  13:48 to 16:17   2h 29m 11s  voyager1  [antenna, sensors]
`,
	Run: log,
}

func init() {
	RootCmd.AddCommand(logCmd)

	logCmd.Flags().StringP("from", "f", "", "Start date (YYYY-MM-DD")
	logCmd.Flags().StringP("to", "t", "", "End date (YYYY-MM-DD)")
	logCmd.Flags().StringSliceP("project", "p", []string{}, `Reports activity only for the given project.
                        You can add other projects by using this option several times.`)
	logCmd.Flags().StringSliceP("tag", "T", []string{}, `Reports activity only for frames containing the given tag.
                        You can add several tags by using this option multiple times`)
	logCmd.Flags().Bool("oneline", false, "Compact output")

}

func log(cmd *cobra.Command, args []string) {
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

	oneline, err := cmd.Flags().GetBool("oneline")
	if err != nil {
		oneline = false
	}

	tracker.Log(query, oneline)
}
