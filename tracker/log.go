package tracker

import (
	"fmt"
	"github.com/florianorben/tracker/helpers"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Log(q FramesQuery, oneline bool, verbose bool, quiet bool) {
	frames := GetFramesFiltered(q)
	groupedFrames := make(map[int][]Frame)
	order := make([]int, 0, len(frames))

	sort.Sort(frames)

	for _, frame := range frames {
		tmpTime, err := time.ParseInLocation(DateFormat, frame.Start.Format(DateFormat), DateLocation)
		if err != nil {
			tmpTime = time.Now()
		}
		date := int(tmpTime.Unix())

		if _, exists := groupedFrames[date]; !exists {
			groupedFrames[date] = make([]Frame, 0)
			order = append(order, date)
		}

		groupedFrames[date] = append(groupedFrames[date], frame)
	}

	maxProjectNameLen := 0
	maxDurationLen := 0
	durationPerDay := make(map[int]float64)

	for date, framesPerDate := range groupedFrames {
		for _, frame := range framesPerDate {

			if frame.InProgress() {
				continue
			}

			duration := frame.Duration()
			durationPerDay[date] += duration.Seconds()

			if len(frame.Project) > maxProjectNameLen {
				maxProjectNameLen = len(frame.Project)
			}

			if len(duration.String()) > maxDurationLen {
				maxDurationLen = len(duration.String())
			}
		}
	}

	for _, timestamp := range order {
		if !oneline && !quiet {
			date := time.Unix(int64(timestamp), 0).Format(LongDateFormat)
			fmt.Printf(
				"%s (%s)\n",
				helpers.PrintBold(helpers.PrintGreen(date)),
				TrackerDuration{time.Duration(durationPerDay[timestamp]) * time.Second}.String(),
			)
		}

		for _, frame := range groupedFrames[timestamp] {
			if frame.InProgress() {
				continue
			}

			if quiet {
				fmt.Println(frame.Uuid)
				continue
			}

			indent := ""
			if !oneline {
				indent = "        "
			}

			shortDate := ""
			if oneline {
				shortDate = helpers.PrintTeal(time.Unix(int64(timestamp), 0).Format(DateFormat) + " ")
			}

			comment := ""
			if !oneline && verbose && frame.Comment != "" {
				commentIndent := indent + "  "
				commentLines := strings.Split(frame.Comment, "\n")
				for i := range commentLines {
					commentLines[i] = commentIndent + commentLines[i]
				}
				comment = strings.Join(commentLines, "\n") + "\n"
			}

			fmt.Printf(
				indent+"%36s %s%s to %s     %"+strconv.Itoa(maxDurationLen)+"s %s %s\n%s",
				frame.Uuid,
				shortDate,
				frame.FormattedStartTime(),
				frame.FormattedEndTime(),
				frame.Duration().String(),
				helpers.PrintPurple(fmt.Sprintf("%"+strconv.Itoa(maxProjectNameLen)+"s", frame.Project)),
				frame.FormattedTags(),
				comment,
			)
		}

		if !oneline && !quiet {
			fmt.Print("\n")
		}
	}
}

func Report(q FramesQuery) {
	frames := GetFramesFiltered(q)

	durationPerProject := make(map[string]float64)
	durationPerTagPerProject := make(map[string]map[string]float64)
	projects := frames.Projects()
	tags := frames.Tags()

	tmpMaxProjectNameLength := 0
	tmpMaxTagNameLength := 0
	tmpMaxDurationLength := 0

	for _, frame := range frames {
		if len(frame.Project) > tmpMaxProjectNameLength {
			tmpMaxProjectNameLength = len(frame.Project)
		}

		if _, exists := durationPerProject[frame.Project]; !exists {
			durationPerProject[frame.Project] = 0
			durationPerTagPerProject[frame.Project] = make(map[string]float64)
		}

		dur := frame.Duration().Seconds()
		durationPerProject[frame.Project] += dur
		for _, tag := range frame.Tags {
			if len(tag) > tmpMaxTagNameLength {
				tmpMaxTagNameLength = len(tag)
			}

			if _, exists := durationPerTagPerProject[frame.Project][tag]; !exists {
				durationPerTagPerProject[frame.Project][tag] = 0
			}
			durationPerTagPerProject[frame.Project][tag] += dur
		}
	}

	for _, tags := range durationPerTagPerProject {
		for _, dur := range tags {
			if l := len(TrackerDuration{time.Duration(dur) * time.Second}.String()); l > tmpMaxDurationLength {
				tmpMaxDurationLength = l
			}
		}
	}

	maxProjectNameLength := strconv.Itoa(tmpMaxProjectNameLength)
	maxTagNameLength := strconv.Itoa(tmpMaxTagNameLength)
	maxDurationLength := strconv.Itoa(tmpMaxDurationLength)

	fmt.Println(
		helpers.PrintTeal(
			fmt.Sprintf("%s -> %s\n",
				frames.MinDate().Format(LongDateFormat),
				frames.MaxDate().Format(LongDateFormat),
			),
		),
	)

	for _, project := range projects {
		if _, exists := durationPerProject[project]; !exists {
			continue
		}

		fmt.Printf(
			"%s - %s\n",
			helpers.PrintPurple(fmt.Sprintf("%-"+maxProjectNameLength+"s", project)),
			helpers.PrintGreen(TrackerDuration{(time.Duration(durationPerProject[project]) * time.Second)}.String()),
		)

		for _, tag := range tags {
			if _, exists := durationPerTagPerProject[project][tag]; !exists {
				continue
			}
			if len(q.Tags) > 0 {
				if _, exists := q.Tags[tag]; !exists {
					continue
				}
			}

			dur := TrackerDuration{(time.Duration(durationPerTagPerProject[project][tag]) * time.Second)}.String()
			fmt.Printf(
				"        [%s    %s]\n",
				helpers.PrintBlue(fmt.Sprintf("%-"+maxTagNameLength+"s", tag)),
				helpers.PrintGreen(fmt.Sprintf("%"+maxDurationLength+"s", dur)),
			)
		}

		fmt.Println("")

	}
}
