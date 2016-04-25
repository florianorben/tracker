package tracker

import (
	"strconv"
	"strings"
	"time"
)

type TrackerDuration struct {
	time.Duration
}

func (t TrackerDuration) String() string {
	h := int(t.Hours())
	m := int(t.Minutes()) - h*60

	if h == 0 && m == 0 {
		return "< 1m"
	}

	parts := make([]string, 0, 2)

	if h > 0 {
		parts = append(parts, strconv.Itoa(h)+"h")
	}

	if m > 0 {
		parts = append(parts, strconv.Itoa(m)+"m")
	}

	return strings.Join(parts, " ")
}
