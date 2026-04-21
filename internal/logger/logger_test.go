package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoggerRotatesLogFileOnDateChange(t *testing.T) {
	logDir := t.TempDir()
	current := time.Date(2026, 4, 17, 23, 59, 58, 0, time.Local)

	logger, err := newWithOptions(logDir, func() time.Time {
		return current
	}, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("newWithOptions() error = %v", err)
	}
	defer logger.Close()

	logger.Info("day one entry")

	current = time.Date(2026, 4, 18, 0, 0, 1, 0, time.Local)
	logger.Info("day two entry")

	dayOneLog, err := os.ReadFile(filepath.Join(logDir, "2026-04-17.log"))
	if err != nil {
		t.Fatalf("ReadFile(day one) error = %v", err)
	}

	dayTwoLog, err := os.ReadFile(filepath.Join(logDir, "2026-04-18.log"))
	if err != nil {
		t.Fatalf("ReadFile(day two) error = %v", err)
	}

	if !strings.Contains(string(dayOneLog), "day one entry") {
		t.Fatalf("day one log = %q, want first entry", string(dayOneLog))
	}

	if strings.Contains(string(dayOneLog), "day two entry") {
		t.Fatalf("day one log = %q, should not contain second entry", string(dayOneLog))
	}

	if !strings.Contains(string(dayTwoLog), "day two entry") {
		t.Fatalf("day two log = %q, want second entry", string(dayTwoLog))
	}
}
