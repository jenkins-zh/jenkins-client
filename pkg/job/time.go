package job

import (
	"encoding/json"
	"time"
)

const jenkinsTimeLayout = "2006-01-02T15:04:05.000-0700"

// Time wraps time.Time for more flexible operations.
type Time struct {
	time.Time
}

// IsZero returns the  true if the value is nil or time is zero.
func (t *Time) IsZero() bool {
	return t == nil || t.Time.IsZero()
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	parsedTime, err := time.Parse(jenkinsTimeLayout, timeStr)
	if err != nil {
		return err
	}
	t.Time = parsedTime.Local()
	return nil
}
