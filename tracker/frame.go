package tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/florianorben/tracker/helpers"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"gopkg.in/andygrunwald/go-jira.v1"
)

const (
	LongDateFormat = "Monday 02 Jan 2006"
	DateTimeFormat = "02.01.2006 15:04"
	TimeFormat     = "15:04"
	DateFormat     = "2006-01-02"
)

var DateLocation = time.Now().Location()

type (
	Frame struct {
		Start    time.Time
		End      time.Time
		Project  string
		Tags     []string
		Uuid     string
		LastEdit time.Time
		Comment  string
		Synced   bool
	}
	EditFrameOpts struct {
		UUID     string
		Position int
	}
)

func NewFrame(p string, t []string) Frame {
	return Frame{
		Start:    time.Now(),
		End:      time.Time{},
		Project:  p,
		Uuid:     uuid.NewV4().String(),
		Tags:     t,
		LastEdit: time.Now(),
		Comment:  "",
		Synced:   false,
	}
}

func (f *Frame) MarshalJSON() ([]byte, error) {
	var end = ""

	if !f.End.IsZero() {
		end = f.End.Format(DateTimeFormat)
	}

	return json.Marshal(&struct {
		Start    string   `json:"start"`
		End      string   `json:"end"`
		Project  string   `json:"project"`
		Tags     []string `json:"tags"`
		Uuid     string   `json:"uuid"`
		LastEdit string   `json:"lastEdit"`
		Comment  string   `json:"comment"`
		Synced   bool     `json:"synced"`
	}{
		f.Start.Format(DateTimeFormat),
		end,
		f.Project,
		f.Tags,
		f.Uuid,
		f.LastEdit.Format(DateTimeFormat),
		f.Comment,
		f.Synced,
	})
}

func (f *Frame) UnmarshalJSON(data []byte) error {
	type Alias Frame
	aux := &struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		LastEdit string `json:"lastEdit"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if start, err := time.ParseInLocation(DateTimeFormat, aux.Start, DateLocation); err != nil {
		f.Start = time.Time{}
	} else {
		f.Start = start
	}

	if aux.End == "" {
		f.End = time.Time{}
	} else if end, err := time.ParseInLocation(DateTimeFormat, aux.End, DateLocation); err != nil {
		f.End = time.Time{}
	} else {
		f.End = end
	}

	if lastEdit, err := time.ParseInLocation(DateTimeFormat, aux.LastEdit, DateLocation); err != nil {
		f.LastEdit = time.Now()
	} else {
		f.LastEdit = lastEdit
	}

	return nil
}

func (f Frame) Equals(other Frame) bool {
	if f.Synced != other.Synced {
		return false
	}

	if f.Start != other.Start {
		return false
	}

	if f.End != other.End {
		return false
	}

	if f.Project != other.Project {
		return false
	}

	if f.Uuid != other.Uuid {
		return false
	}

	if len(f.Tags) != len(other.Tags) {
		return false
	}

	if f.Comment != other.Comment {
		return false
	}

	tmpTags := make(map[string]bool)

	for _, tag := range f.Tags {
		tmpTags[tag] = true
	}

	for _, otherTag := range other.Tags {
		tmpTags[otherTag] = true
	}

	if len(tmpTags) != len(f.Tags) {
		return false
	}

	return true
}

func (f Frame) InProgress() bool {
	return f.End.IsZero()
}

func (f Frame) Finished() bool {
	return !f.End.IsZero()
}

func (f Frame) Duration() TrackerDuration {
	return TrackerDuration{f.End.Sub(f.Start)}
}

func (f Frame) FormattedProject() string {
	return helpers.PrintPurple(f.Project)
}

func (f Frame) FormattedStartTime() string {
	return helpers.PrintGreen(f.Start.Format(TimeFormat))
}

func (f Frame) FormattedEndTime() string {
	return helpers.PrintGreen(f.End.Format(TimeFormat))
}

func (f Frame) FormattedTags() string {
	if len(f.Tags) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s]", helpers.PrintBlue(strings.Join(f.Tags, " ")))
}

func (f Frame) RelativeTime() string {
	t := helpers.PrintGreen("just now")
	if time.Since(f.Start).Minutes() > 1 {
		t = f.FormattedStartTime()
	}

	return t
}

func (f *Frame) AddWorkLog() {
	if len(f.Tags) == 0 || f.Synced {
		return
	}

	user := viper.GetString("backend.user")
	if user == "" {
		fmt.Printf("Error: %s\n", helpers.PrintRed("Please use tracker login first"))
	}

	pass, err := keyring.Get("tracker", user)
	if err != nil {
		if err != keyring.ErrNotFound {
			fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
			os.Exit(1)
		} else {
			fmt.Printf("Error: %s\n", helpers.PrintRed("Please use tracker login first"))
			return
		}
	}

	tp := jira.BasicAuthTransport{
		Username: user,
		Password: pass,
	}

	client, err := jira.NewClient(tp.Client(), viper.GetString("backend.url"))
	if err != nil {
		return
	}

	for _, tag := range f.Tags {
		_, _, err = client.Issue.AddWorklogRecord(tag, &jira.WorklogRecord{
			Comment:          f.Comment,
			TimeSpentSeconds: int(f.End.Sub(f.Start).Truncate(time.Minute).Seconds()),
		})

		if err != nil {
			if jiraError, ok := err.(*jira.Error); !ok {
				fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
				return
			} else if jiraError.ErrorMessages[0] == "Issue Does Not Exist" {
				continue
			}
		}

		f.Synced = true
		return
	}
}

func EditFrame(opts EditFrameOpts) (Frame, error) {
	var frame Frame
	frames := GetFrames()
	index := -1
	arg := ""

	if opts.UUID != "" {
		index, frame = frames.ByUUID(opts.UUID)
		arg = opts.UUID
	} else if opts.Position < 0 {
		index, frame = frames.ByPosition(opts.Position)
		arg = fmt.Sprintf("%d", opts.Position)
	}

	if index == -1 {
		return frame, fmt.Errorf("Error: %s %s.\n", helpers.PrintRed("No frame found with id"), arg)
	}

	b, err := json.MarshalIndent(&struct {
		Start   string   `json:"start"`
		End     string   `json:"end"`
		Project string   `json:"project"`
		Tags    []string `json:"tags"`
		Comment string   `json:"comment"`
	}{
		Start:   frame.Start.Format(DateTimeFormat),
		End:     frame.End.Format(DateTimeFormat),
		Project: frame.Project,
		Tags:    frame.Tags,
		Comment: frame.Comment,
	}, "", "  ")
	if err != nil {
		return frame, fmt.Errorf("Error: Creating temp file failed: %s\n", helpers.PrintRed(err.Error()))
	}

	newFrameContents, err := helpers.OpenInEditor(b)
	if err != nil {
		return frame, err
	}

	var newFrame Frame
	json.Unmarshal(newFrameContents, &newFrame)
	newFrame.Uuid = frame.Uuid
	newFrame.LastEdit = time.Now()

	if frames[index].Equals(newFrame) {
		return frame, fmt.Errorf("%s\n", "No changes made.")
	}

	frames[index] = newFrame
	frames.Persist()

	return newFrame, nil
}
