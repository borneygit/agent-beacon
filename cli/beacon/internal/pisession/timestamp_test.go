package pisession

import (
	"strconv"
	"testing"
	"time"
)

// A Pi session header carries its time as an RFC3339 string. timestampValue used to read
// the leading number of any string first, so that header became 2026 seconds after the
// epoch, its trace started in 1970, and the local history deleted the session as expired.
func TestTimestampValueReadsEveryShape(t *testing.T) {
	want := time.Date(2026, 10, 31, 2, 41, 17, 937000000, time.UTC)
	second := want.Truncate(time.Second)
	cases := []struct {
		name  string
		value interface{}
		want  time.Time
	}{
		{"rfc3339 string", "2026-10-31T02:41:17.937Z", want},
		{"nanosecond string", "2026-10-31T02:41:17.937000000Z", want},
		{"epoch millis string", strconv.FormatInt(want.UnixMilli(), 10), want},
		{"epoch millis number", float64(want.UnixMilli()), want},
		{"time value", want, want},
		{"epoch seconds string", strconv.FormatInt(want.Unix(), 10), second},
		{"epoch seconds number", float64(want.Unix()), second},
	}
	for _, c := range cases {
		got := timestampValue(c.value)
		if !got.Equal(c.want) {
			t.Errorf("%s: timestampValue(%v) = %s, want %s", c.name, c.value, got.Format(time.RFC3339Nano), c.want.Format(time.RFC3339Nano))
		}
	}
}

// A string that is not a time at all must not be read as its leading digits.
func TestTimestampValueIgnoresNonNumericStrings(t *testing.T) {
	got := timestampValue("2026-not-a-time")
	if got.Year() < 2020 {
		t.Fatalf("timestampValue(\"2026-not-a-time\") = %s, want the current time", got.Format(time.RFC3339Nano))
	}
}
