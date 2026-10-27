package snowflake

import (
	"testing"
	"time"
)

func TestComposeAndTimeRoundTrip(t *testing.T) {
	at := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	id, err := Compose(at, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got := Time(id); !got.Equal(at) {
		t.Fatalf("Time = %s, want %s", got, at)
	}
	start, end := MonthBounds(id)
	if !start.Equal(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("month start = %s", start)
	}
	if !end.Equal(time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("month end = %s", end)
	}
}

func TestNextMonotonic(t *testing.T) {
	ms := time.Date(2026, 10, 27, 10, 0, 0, 0, time.UTC).UnixMilli()
	now := func() time.Time {
		cur := ms
		ms++
		return time.UnixMilli(cur).UTC()
	}
	gen, err := New(0, now)
	if err != nil {
		t.Fatal(err)
	}
	a, _, err := gen.Next()
	if err != nil {
		t.Fatal(err)
	}
	b, _, err := gen.Next()
	if err != nil {
		t.Fatal(err)
	}
	if a >= b {
		t.Fatalf("ids not increasing: %d then %d", a, b)
	}
}
