package scheduler

import (
	"strings"
	"testing"
	"time"

	"git-push-timer/internal/logger"
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

func TestStopWaitsForRunningJobs(t *testing.T) {
	log, err := logger.New()
	if err != nil {
		t.Fatalf("logger.New() error = %v", err)
	}
	defer log.Close()

	s := &Scheduler{
		logger: log,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
	s.runningWG.Add(1)

	stopped := make(chan struct{})
	go func() {
		s.Stop()
		close(stopped)
	}()

	close(s.doneCh)

	select {
	case <-stopped:
		t.Fatal("Stop() returned before running jobs finished")
	case <-time.After(50 * time.Millisecond):
	}

	s.runningWG.Done()

	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("Stop() did not return after running jobs finished")
	}
}
