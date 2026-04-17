package scheduler

import (
	"strings"
	"testing"
	"time"
)

func TestAdvanceNextRunSkipsMissedSchedules(t *testing.T) {
	schedule, err := parseCronSpec("*/5 * * * *")
	if err != nil {
		t.Fatalf("parse cron: %v", err)
	}

	nextRun := time.Date(2026, 4, 17, 18, 5, 0, 0, time.Local)
	now := time.Date(2026, 4, 17, 18, 17, 30, 0, time.Local)

	got := advanceNextRun(schedule, nextRun, now)
	want := time.Date(2026, 4, 17, 18, 20, 0, 0, time.Local)

	if !got.Equal(want) {
		t.Fatalf("advanceNextRun() = %v, want %v", got, want)
	}
}

func TestTimeUntilNextTickAlignsToMinuteBoundary(t *testing.T) {
	now := time.Date(2026, 4, 17, 18, 3, 51, 0, time.Local)

	got := timeUntilNextTick(now)
	want := 9 * time.Second

	if got != want {
		t.Fatalf("timeUntilNextTick() = %v, want %v", got, want)
	}
}

func TestParseCronSpecAcceptsStandardCron(t *testing.T) {
	if _, err := parseCronSpec("0 18 * * *"); err != nil {
		t.Fatalf("parseCronSpec() unexpected error: %v", err)
	}
}

func TestParseCronSpecRejectsDescriptors(t *testing.T) {
	_, err := parseCronSpec("@every 30s")
	if err == nil {
		t.Fatal("parseCronSpec() expected error for descriptor")
	}

	if !strings.Contains(err.Error(), "does not accept descriptors") {
		t.Fatalf("parseCronSpec() error = %q, want descriptor rejection", err)
	}
}
