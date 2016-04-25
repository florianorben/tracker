package tracker

import (
	"time"

	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
	"tracker/helpers"
)

const (
	LongDateFormat = "Monday 02 Jan 2006"
	DateTimeFormat = "02.01.2006 15:04"
	TimeFormat     = "15:04"
	DateFormat     = "2006-01-02"
)

var DateLocation = time.Now().Location()

type Frame struct {
	Start    time.Time
	End      time.Time
	Project  string
	Tags     []string
	Uuid     string
	LastEdit time.Time
}

func NewFrame(p string, t []string) Frame {
	return Frame{
		Start:    time.Now(),
		End:      time.Time{},
		Project:  p,
		Uuid:     uuid.NewV4().String(),
		Tags:     t,
		LastEdit: time.Now(),
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
	}{
		f.Start.Format(DateTimeFormat),
		end,
		f.Project,
		f.Tags,
		f.Uuid,
		f.LastEdit.Format(DateTimeFormat),
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
	if f.Start != other.Start {
		return false
	} else if f.End != other.End {
		return false
	} else if f.Project != other.Project {
		return false
	} else if f.Uuid != other.Uuid {
		return false
	} else if len(f.Tags) != len(other.Tags) {
		return false
	} else {
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
