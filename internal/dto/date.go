package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

// Date is a strict JSON calendar date encoded as YYYY-MM-DD.
type Date struct {
	t time.Time
}

// DateFromTime returns the UTC calendar date for t at midnight.
func DateFromTime(t time.Time) Date {
	year, month, day := t.UTC().Date()
	return Date{t: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// Time returns the date as UTC midnight.
func (d Date) Time() time.Time {
	return d.t
}

// MarshalJSON encodes the date as YYYY-MM-DD.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.t.Format(dateLayout))
}

// UnmarshalJSON accepts only JSON strings formatted as YYYY-MM-DD.
func (d *Date) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("date must be a JSON string")
	}
	if len(value) != len(dateLayout) {
		return fmt.Errorf("date must use YYYY-MM-DD")
	}
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return fmt.Errorf("date must use YYYY-MM-DD")
	}
	d.t = DateFromTime(parsed).Time()
	return nil
}
