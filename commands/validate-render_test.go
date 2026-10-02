package commands

import (
	"bytes"
	"strings"
	"testing"
	"time"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

func ptr[T any](v T) *T { return &v }

func TestRenderAnalysis(t *testing.T) {
	la, _ := time.LoadLocation("America/Los_Angeles")
	a := &scheduler.ScheduleAnalysis{
		Expression:   "CRON_TZ=America/Los_Angeles H(15-45) 2 * * MON-FRI",
		Ref:          &scheduler.Ref{Type: "job", GUID: "j1", Name: "backup"},
		Valid:        true,
		Description:  "at a hashed minute between 15 and 45 of hour 2, on Monday to Friday, America/Los_Angeles time",
		Location:     "America/Los_Angeles",
		Illustrative: true,
		HashedFields: []scheduler.HashedField{{Field: "minute", Value: "H(15-45)", Resolved: "18"}},
		NextRuns:     []time.Time{time.Date(2026, 10, 5, 2, 18, 0, 0, la), time.Date(2026, 10, 6, 2, 18, 0, 0, la)},
		Warnings:     []scheduler.Finding{{Code: "dst_skipped", Message: "02:18 on 2027-03-14 does not exist in America/Los_Angeles; it runs at 03:18"}},
	}
	var b bytes.Buffer
	renderAnalysis(&b, a)
	want := `Expression:    CRON_TZ=America/Los_Angeles H(15-45) 2 * * MON-FRI
Description:   at a hashed minute between 15 and 45 of hour 2, on Monday to Friday, America/Los_Angeles time
Time zone:     America/Los_Angeles
Hashed:        minute H(15-45) → 18

Next 2 runs (illustrative, H values depend on the job):
  Mon 2026-10-05 02:18 PDT
  Tue 2026-10-06 02:18 PDT

Warnings:
  02:18 on 2027-03-14 does not exist in America/Los_Angeles; it runs at 03:18
`
	if b.String() != want {
		t.Errorf("got\n%s\nwant\n%s", b.String(), want)
	}
}

func TestRenderParseError(t *testing.T) {
	a := &scheduler.ScheduleAnalysis{
		Expression:  "60 * * * *",
		Description: "minute 60, hour *, day-of-month *, month *, day-of-week *",
		Errors: []scheduler.Finding{{Code: "parse_error", Message: "minute field value 60 is out of range (valid: 0-59)",
			Field: "minute", Value: "60", Offset: ptr(0)}},
	}
	var b bytes.Buffer
	renderAnalysis(&b, a)
	if !strings.HasPrefix(b.String(), "minute field value 60 is out of range (valid: 0-59)\n  60 * * * *\n  ^^\n") {
		t.Errorf("got\n%s", b.String())
	}
}

func TestRenderStoredSchedule(t *testing.T) {
	a := &scheduler.ScheduleAnalysis{ScheduleGUID: "s1", Enabled: ptr(false), Expression: "0 0 30 2 *",
		Description: "at 00:00, on day 30 of the month, in February",
		Errors:      []scheduler.Finding{{Code: "never_fires", Message: "this schedule never runs"}}}
	var b bytes.Buffer
	renderAnalysis(&b, a)
	if !strings.HasPrefix(b.String(), "this schedule never runs\n\nSchedule:      s1 (disabled)\nExpression:    0 0 30 2 *\n") {
		t.Errorf("got\n%s", b.String())
	}
}

func TestFormatRunKeepsLocalUnconverted(t *testing.T) {
	run := time.Date(2026, 10, 5, 2, 18, 0, 0, time.FixedZone("SRV", 3600))
	if got := formatRun(run, "Local", []time.Time{run}); got != "Mon 2026-10-05 02:18 SRV" {
		t.Errorf("got %q", got)
	}
}

func TestRenderOneRun(t *testing.T) {
	run := time.Date(2027, 1, 1, 3, 0, 0, 0, time.UTC)
	a := &scheduler.ScheduleAnalysis{Expression: "0 3 1 1 *", Valid: true, Location: "Etc/UTC",
		NextRuns: []time.Time{run}, PrevRuns: []time.Time{run.AddDate(-1, 0, 0)}}
	var b bytes.Buffer
	renderAnalysis(&b, a)
	if !strings.Contains(b.String(), "\nNext run:\n") || !strings.Contains(b.String(), "\nPrevious run:\n") {
		t.Errorf("got\n%s", b.String())
	}
}
