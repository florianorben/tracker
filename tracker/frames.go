package tracker

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type (
	Frames      []Frame
	FramesQuery struct {
		Projects map[string]struct{}
		Tags     map[string]struct{}
		From     time.Time
		To       time.Time
	}
)

func (f Frames) Len() int           { return len(f) }
func (f Frames) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f Frames) Less(i, j int) bool { return f[i].Start.Unix() < f[j].Start.Unix() }

func NewFrameQuery(projects, tags []string, from, to string) (FramesQuery, error) {
	defaultStartDate := time.Duration(viper.GetInt("log.defaultStartDate"))
	defaultEndDate := time.Duration(viper.GetInt("log.defaultEndDate"))
	endOfDay, _ := time.ParseDuration("23h59m59s")

	tmpFrom := time.Time{}
	tmpTo := time.Time{}

	if from != "" {
		if tmp, err := time.ParseInLocation(DateFormat, from, DateLocation); err == nil {
			tmpFrom = tmp
		} else {
			return FramesQuery{}, err
		}
	} else {
		defaultFromDate, err := time.ParseInLocation(DateFormat, time.Now().Format(DateFormat), DateLocation)
		if err == nil {
			tmpFrom = defaultFromDate.Add(defaultStartDate * 24 * time.Hour)
		}
	}

	if to != "" {
		if tmp, err := time.ParseInLocation(DateFormat, to, DateLocation); err == nil {
			tmpTo = tmp.Add(endOfDay)
		} else {
			return FramesQuery{}, err
		}
	} else {
		defaultToDate, err := time.ParseInLocation(DateFormat, time.Now().Format(DateFormat), DateLocation)
		if err == nil {
			tmpTo = defaultToDate.Add(defaultEndDate * 24 * time.Hour).Add(endOfDay)
		}
	}

	tmpProjects := make(map[string]struct{})
	if len(projects) > 0 {
		for _, project := range projects {
			tmpProjects[project] = struct{}{}
		}
	}

	tmpTags := make(map[string]struct{})
	if len(tags) > 0 {
		for _, tag := range tags {
			tmpTags[tag] = struct{}{}
		}
	}

	return FramesQuery{
		Projects: tmpProjects,
		Tags:     tmpTags,
		From:     tmpFrom,
		To:       tmpTo,
	}, nil
}

func GetFrames() Frames {
	var file *os.File
	file, err := os.OpenFile(viper.GetString("framesFile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if len(contents) == 0 {
		contents = []byte("[]")
	}

	var f Frames
	err = json.Unmarshal(contents, &f)
	if err != nil {
		panic(err)
	}

	return f
}

func GetFramesFiltered(q FramesQuery) Frames {
	frames := GetFrames()
	filtered := make(Frames, 0, len(frames))

	for _, frame := range frames {
		if !q.From.IsZero() && q.From.Sub(frame.Start) >= 0 {
			continue
		}

		if !q.To.IsZero() && q.To.Sub(frame.End) <= 0 {
			continue
		}

		if len(q.Projects) > 0 {
			if _, exists := q.Projects[frame.Project]; !exists {
				continue
			}
		}

		if len(q.Tags) > 0 {
			contains := false
			for _, tag := range frame.Tags {
				if _, exists := q.Tags[tag]; exists {
					contains = true
				}
			}

			if contains == false {
				continue
			}
		}

		filtered = append(filtered, frame)
	}

	return filtered
}

func (f Frames) Persist() {
	var file *os.File
	framesFile := viper.GetString("framesFile")

	file, err := os.OpenFile(framesFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.Truncate(framesFile, 0)
	if err != nil {
		panic(err)
	}

	file.Write(bytes.TrimSpace(b))
}

func (f Frames) ByUUID(uuid string) (int, Frame) {
	for i, frame := range f {
		if frame.Uuid == uuid {
			return i, frame
		}
	}

	return -1, Frame{}
}

func (f Frames) ByPosition(pos int) (int, Frame) {
	sort.Sort(f)

	length := len(f)
	var framePos int

	switch {
	case pos > length-1:
		framePos = length - 1
	case pos < 0:
		tmp := length + pos
		if tmp < 0 {
			tmp = 0
		}
		framePos = tmp
	}

	return framePos, f[framePos]
}

func (f Frames) MinDate() time.Time {
	var minDate time.Time
	minTimestamp := int64(1<<63 - 1)

	for _, frame := range f {
		timestamp := frame.Start.Unix()
		if timestamp < minTimestamp {
			minDate = frame.Start
			minTimestamp = frame.Start.Unix()
		}
	}

	return minDate
}

func (f Frames) MaxDate() time.Time {
	var maxDate time.Time
	maxTimestamp := int64(-1 << 63)

	for _, frame := range f {
		timestamp := frame.Start.Unix()
		if timestamp > maxTimestamp {
			maxDate = frame.Start
			maxTimestamp = timestamp
		}
	}

	return maxDate
}

func (f Frames) Projects() []string {
	projectsMap := make(map[string]bool)
	for _, frame := range f {
		projectsMap[frame.Project] = true
	}

	projects := make([]string, 0, len(projectsMap))
	for project, _ := range projectsMap {
		projects = append(projects, project)
	}

	sort.Strings(projects)
	return projects
}

func (f Frames) Tags() []string {
	tagsMap := make(map[string]bool)
	for _, frame := range f {
		for _, tag := range frame.Tags {
			tagsMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagsMap))
	for tag, _ := range tagsMap {
		tags = append(tags, tag)
	}

	sort.Strings(tags)
	return tags
}

func (f Frames) Frames() []string {
	frames := make([]string, 0, len(f))
	for _, frame := range f {
		frames = append(frames, frame.Uuid)
	}

	sort.Strings(frames)
	return frames
}
