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
	renderAnalysis(&b, a, "", time.Time{})
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

func TestRenderDescriptionNote(t *testing.T) {
	a := &scheduler.ScheduleAnalysis{Expression: "H/15 * * * *", Valid: true,
		Description:     "every 15 minutes from a hashed start",
		DescriptionNote: "at minutes h, h+15, h+30 and h+45 of every hour, where h is a hashed minute from 0 to 14"}
	var b bytes.Buffer
	renderAnalysis(&b, a, "", time.Time{})
	want := `Expression:    H/15 * * * *
Description:   every 15 minutes from a hashed start
Note:          at minutes h, h+15, h+30 and h+45 of every hour, where h is a hashed minute from 0 to 14
`
	if b.String() != want {
		t.Errorf("got\n%s\nwant\n%s", b.String(), want)
	}
}

func TestRenderFrom(t *testing.T) {
	a := &scheduler.ScheduleAnalysis{Expression: "CRON_TZ=UTC 0 9 * * *", Valid: true, Description: "at 09:00", Location: "UTC",
		NextRuns: []time.Time{time.Date(2030, 1, 1, 9, 0, 0, 0, time.UTC)},
		PrevRuns: []time.Time{time.Date(2029, 12, 31, 9, 0, 0, 0, time.UTC)}}
	var b bytes.Buffer
	renderAnalysis(&b, a, "", time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	for _, want := range []string{"Next run after Tue 2030-01-01 00:00 UTC:\n", "Previous run before Tue 2030-01-01 00:00 UTC:\n"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("missing %q in\n%s", want, b.String())
		}
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
	renderAnalysis(&b, a, "", time.Time{})
	if !strings.HasPrefix(b.String(), "minute field value 60 is out of range (valid: 0-59)\n  60 * * * *\n  ^^\n") {
		t.Errorf("got\n%s", b.String())
	}
}

func TestRenderStoredSchedule(t *testing.T) {
	a := &scheduler.ScheduleAnalysis{ScheduleGUID: "s1", Enabled: ptr(false), Expression: "0 0 30 2 *",
		Description: "at 00:00, on day 30 of the month, in February",
		Errors:      []scheduler.Finding{{Code: "never_fires", Message: "this schedule never runs"}}}
	var b bytes.Buffer
	renderAnalysis(&b, a, "", time.Time{})
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
	renderAnalysis(&b, a, "", time.Time{})
	if !strings.Contains(b.String(), "\nNext run:\n") || !strings.Contains(b.String(), "\nPrevious run:\n") {
		t.Errorf("got\n%s", b.String())
	}
}

func TestRenderStoredInAnotherZone(t *testing.T) {
	la, _ := time.LoadLocation("America/Los_Angeles")
	run := time.Date(2027, 1, 1, 3, 20, 0, 0, la)
	a := &scheduler.ScheduleAnalysis{ScheduleGUID: "s1", Expression: "CRON_TZ=America/Los_Angeles 20 3 1 1 *",
		Valid: true, Location: "America/Los_Angeles", NextRuns: []time.Time{run}}
	var b bytes.Buffer
	renderAnalysis(&b, a, "America/New_York", time.Time{})
	for _, want := range []string{"Time zone:     America/Los_Angeles\n", "\nNext run in America/New_York:\n  Fri 2027-01-01 06:20 EST\n"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("missing %q in\n%s", want, b.String())
		}
	}
}
