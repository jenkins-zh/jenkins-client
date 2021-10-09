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
	// parse using Jenkins time layout first
	parsedTime, err := time.Parse(jenkinsTimeLayout, timeStr)
	if err != nil {
		// parse using RFC3339 layout corresponding to MarshalJSON the timeStr isn't sufficient Jenkins time layout
		parsedTime, err = time.Parse(time.RFC3339, timeStr)
	}
	if err != nil {
		// return error if we tried the above methods
		return err
	}
	t.Time = parsedTime.Local()
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t *Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	buf := make([]byte, 0, len(time.RFC3339)+2)
	buf = append(buf, '"')
	buf = t.UTC().AppendFormat(buf, time.RFC3339)
	buf = append(buf, '"')
	return buf, nil
}
