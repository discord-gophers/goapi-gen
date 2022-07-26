package types

import (
	"encoding/json"
	"time"
)

// DateFormat represents a Go time format that matches the OpenAPI date.
const DateFormat = "2006-01-02"

// Date represents a time that must conform DateFormat.
//
// Date implements the JSON.Marshaler and JSON.Unmarshaler interfaces, and
// validates that date matches DateFormat.
type Date struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(DateFormat))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}
	parsed, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

func (d Date) String() string {
	return d.Time.Format(DateFormat)
}
