package dto

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateJSONRoundTrip(t *testing.T) {
	var payload struct {
		Published Date `json:"published"`
	}

	if err := json.Unmarshal([]byte(`{"published":"2026-06-18"}`), &payload); err != nil {
		t.Fatalf("unmarshal date: %v", err)
	}

	want := time.Date(2026, 6, 18, 0, 0, 0, 0, time.UTC)
	if !payload.Published.Time().Equal(want) {
		t.Fatalf("date time = %s, want %s", payload.Published.Time(), want)
	}

	got, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal date: %v", err)
	}
	if string(got) != `{"published":"2026-06-18"}` {
		t.Fatalf("marshaled payload = %s", got)
	}
}

func TestDateZeroValueMarshalsNull(t *testing.T) {
	got, err := json.Marshal(Date{})
	if err != nil {
		t.Fatalf("marshal zero date: %v", err)
	}
	if string(got) != `null` {
		t.Fatalf("zero date = %s, want null", got)
	}
}

func TestDateRejectsDatetimeString(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(`"2026-06-18T00:00:00Z"`), &d); err == nil {
		t.Fatal("expected datetime string to be rejected")
	}
}

func TestDateRejectsInvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "non-string", body: `123`},
		{name: "bad format", body: `"2026/06/18"`},
		{name: "missing zero padding", body: `"2026-6-18"`},
		{name: "invalid calendar date", body: `"2026-02-30"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			if err := json.Unmarshal([]byte(tt.body), &d); err == nil {
				t.Fatalf("expected %s to be rejected", tt.body)
			}
		})
	}
}

func TestDateFromTimeUsesUTCMidnight(t *testing.T) {
	source := time.Date(2026, 6, 18, 1, 30, 0, 0, time.FixedZone("ahead", 2*60*60))

	got := DateFromTime(source).Time()
	want := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("date from time = %s, want %s", got, want)
	}
	if got.Location() != time.UTC {
		t.Fatalf("date location = %s, want UTC", got.Location())
	}
}
