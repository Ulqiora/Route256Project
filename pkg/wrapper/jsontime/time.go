package jtime

import (
	"strings"
	"time"
)

type TimeWrap time.Time

const JFormatTime = "2006-01-02 15:04:05.999999999 -07:00"

func (c *TimeWrap) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}
	t, err := time.Parse(JFormatTime, value)
	if err != nil {
		return err
	}
	*c = TimeWrap(t)
	return nil
}

func (c *TimeWrap) MarshalJSON() ([]byte, error) {
	if c.Time().IsZero() {
		return []byte("null"), nil
	}
	return []byte(time.Time(*c).Format(JFormatTime)), nil
}

func (c *TimeWrap) Time() *time.Time {
	return (*time.Time)(c)
}
func (c *TimeWrap) String() string {
	return time.Time(*c).Format(JFormatTime)
}

func (c *TimeWrap) IsZero() bool { return time.Time(*c).IsZero() }
