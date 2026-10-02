package commands

import (
	"fmt"
	"io"
	"time"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

// renderAnalysis prints one validated expression after the OK/FAILED line.
// display, when set, is the zone run times are shown in instead of the
// schedule's own; the schedule still runs in its zone.
func renderAnalysis(w io.Writer, a *scheduler.ScheduleAnalysis, display string) {
	if len(a.Errors) > 0 {
		renderFindings(w, a.Expression, a.Errors)
		fmt.Fprintln(w)
	}
	label := func(name, value string) { fmt.Fprintf(w, "%-15s%s\n", name+":", value) }
	if a.ScheduleGUID != "" {
		state := "enabled"
		if a.Enabled != nil && !*a.Enabled {
			state = "disabled"
		}
		label("Schedule", a.ScheduleGUID+" ("+state+")")
	}
	label("Expression", a.Expression)
	label("Description", a.Description)
	if a.Location != "" {
		label("Time zone", a.Location)
	}
	for i, h := range a.HashedFields {
		name := ""
		if i == 0 {
			name = "Hashed:"
		}
		fmt.Fprintf(w, "%-15s%s %s → %s\n", name, h.Field, h.Value, h.Resolved)
	}
	runs := func(title string, times []time.Time) {
		if len(times) == 0 {
			return
		}
		zone := a.Location
		if display != "" {
			title += " in " + display
			zone = display
		}
		if a.Illustrative {
			scope := "the job or call"
			if a.Ref != nil {
				scope = "the " + a.Ref.Type
			}
			title += " (illustrative, H values depend on " + scope + ")"
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w, title+":")
		for _, t := range times {
			fmt.Fprintln(w, "  "+formatRun(t, zone, times))
		}
	}
	runs(countRuns("Next", len(a.NextRuns)), a.NextRuns)
	runs(countRuns("Previous", len(a.PrevRuns)), a.PrevRuns)
	if len(a.Warnings) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Warnings:")
		for _, f := range a.Warnings {
			fmt.Fprintln(w, "  "+f.Message)
		}
	}
}

// formatRun shows t in the schedule's zone when this machine knows it, with
// seconds only if some run in the list has them. "Local" is the server's
// zone, not this machine's, so t keeps the offset the server sent.
func formatRun(t time.Time, location string, all []time.Time) string {
	if loc, err := time.LoadLocation(location); err == nil && location != "" && location != "Local" {
		t = t.In(loc)
	}
	layout := "Mon 2006-01-02 15:04 MST"
	for _, other := range all {
		if other.Second() != 0 {
			layout = "Mon 2006-01-02 15:04:05 MST"
			break
		}
	}
	return t.Format(layout)
}

// countRuns titles a run list: "Next run", "Next 5 runs".
func countRuns(which string, n int) string {
	if n == 1 {
		return which + " run"
	}
	return fmt.Sprintf("%s %d runs", which, n)
}
