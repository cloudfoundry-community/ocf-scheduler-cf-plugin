package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf validate-schedule [OPTIONS] [job|call] [NAME|GUID] [EXPRESSION]
// Exits 1 when an expression is invalid, so scripts can gate on it.
func ValidateSchedule(services *core.Services, args []string) {
	valid, err := validateSchedule(services, args)
	if err != nil {
		fmt.Println("FAILED")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if !valid {
		os.Exit(1)
	}
}

type validateOptions struct {
	timezone string
	next     *int // nil: the scheduler's default
	prev     int
	from     *time.Time // nil: now
}

// parseValidateFlags reads the options; args starts with the command name,
// which rest drops.
func parseValidateFlags(args []string) (validateOptions, []string, error) {
	var opts validateOptions
	var next int
	var from string
	flags := pflag.NewFlagSet("validate-schedule", pflag.ExitOnError)
	flags.StringVarP(&opts.timezone, "timezone", "t", "", "Time zone for the expression; for stored schedules, the zone run times are shown in")
	flags.IntVar(&next, "next", 0, "Number of next runs to show (default 5; 1 per stored schedule)")
	flags.IntVar(&opts.prev, "prev", 0, "Number of previous runs to show")
	flags.StringVar(&from, "from", "", "List runs from this time instead of now")
	flags.Parse(args)

	opts.timezone = strings.TrimSpace(opts.timezone)
	if strings.Contains(opts.timezone, " ") {
		return opts, nil, fmt.Errorf("no spaces allowed in the timezone")
	}
	if flags.Changed("next") {
		if next < 0 {
			return opts, nil, fmt.Errorf("--next must be 0 or more")
		}
		opts.next = &next
	}
	if from != "" {
		t, err := parseFrom(strings.TrimSpace(from), opts.timezone)
		if err != nil {
			return opts, nil, err
		}
		opts.from = &t
	}
	return opts, flags.Args()[1:], nil
}

// parseFrom reads --from: RFC 3339, or a date or a date and time without an
// offset, read in zone (UTC when zone is empty).
func parseFrom(s, zone string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	loc := time.UTC
	if zone != "" {
		l, err := time.LoadLocation(zone)
		if err != nil {
			return time.Time{}, fmt.Errorf("unknown time zone %q", zone)
		}
		loc = l
	}
	for _, layout := range []string{"2006-01-02T15:04", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("--from must be a date (2030-01-01), a date and time (2030-01-01T09:00) or RFC 3339 (2030-01-01T09:00:00Z)")
}

func validateSchedule(services *core.Services, args []string) (bool, error) {
	opts, args, err := parseValidateFlags(args)
	if err != nil {
		return false, err
	}
	timezone := opts.timezone
	if err := checkArgs(args); err != nil {
		return false, err // usage errors need no lookups
	}
	// fail keeps the "Validating ... in org / space as user" line ahead of
	// FAILED for errors found before the target is known.
	fail := func(err error) (bool, error) {
		core.PrintActionInProgress(services, "%s", "Validating schedule")
		return false, err
	}

	space, err := core.MySpace(services)
	if err != nil {
		return fail(fmt.Errorf("Could not get current space."))
	}
	jobs, err := client.ListJobs(services.Client, space)
	if err != nil {
		return fail(fmt.Errorf("Could not get jobs for space %s.", space.Name))
	}
	calls, err := client.ListCalls(services.Client, space)
	if err != nil {
		return fail(fmt.Errorf("Could not get calls for space %s.", space.Name))
	}

	r, err := resolveArgs(args, jobs, calls)
	var amb *AmbiguousError
	if errors.As(err, &amb) {
		return fail(ambiguity(services, amb))
	}
	if err != nil {
		return fail(err)
	}

	req := scheduler.ValidateRequest{Prev: opts.prev, Next: opts.next, From: opts.from}
	var from time.Time
	if opts.from != nil {
		from = *opts.from
	}
	display := "" // stored schedules: -t only changes the zone runs are shown in
	if r.Expression != "" {
		req.Expression = strings.TrimSpace(r.Expression)
		if timezone != "" {
			req.Expression = fmt.Sprintf("CRON_TZ=%s %s", timezone, req.Expression)
		}
	} else if timezone != "" {
		if _, err := time.LoadLocation(timezone); err != nil {
			return fail(fmt.Errorf("unknown time zone %q", timezone))
		}
		display = timezone
	}
	action := "Validating cron expression"
	if r.Target != nil {
		req.RefType, req.RefGUID = r.Target.Type, r.Target.GUID
		action = fmt.Sprintf("Validating schedule for %s %s", r.Target.Type, r.Target.Name)
		if r.Expression == "" {
			action = fmt.Sprintf("Validating stored schedules of %s %s", r.Target.Type, r.Target.Name)
		}
	}
	if err := core.PrintActionInProgress(services, "%s", action); err != nil {
		return false, err
	}

	if req.Expression != "" {
		a, err := client.ValidateExpression(services.Client, req)
		if err != nil {
			return false, err
		}
		fmt.Println(verdict(a.Valid))
		fmt.Println()
		renderAnalysis(os.Stdout, a, display, from)
		return a.Valid, nil
	}

	if req.Next == nil {
		one := 1
		req.Next = &one
	}
	stored, err := client.ValidateStored(services.Client, req)
	if err != nil {
		return false, err
	}
	valid := true
	for _, a := range stored.Resources {
		valid = valid && a.Valid
	}
	fmt.Println(verdict(valid))
	if len(stored.Resources) == 0 {
		fmt.Printf("\n%s %s has no schedules.\n", r.Target.Type, r.Target.Name)
	}
	for _, a := range stored.Resources {
		fmt.Println()
		renderAnalysis(os.Stdout, a, display, from)
	}
	return valid, nil
}

func verdict(valid bool) string {
	if valid {
		return "OK"
	}
	return "FAILED"
}

func ambiguity(services *core.Services, amb *AmbiguousError) error {
	appNames := map[string]string{}
	if apps, err := core.MyApps(services); err == nil {
		for _, app := range apps {
			appNames[app.Guid] = app.Name
		}
	}
	var b strings.Builder
	b.WriteString(amb.Error())
	for _, m := range amb.Matches {
		app := appNames[m.AppGUID]
		if app == "" {
			app = m.AppGUID
		}
		fmt.Fprintf(&b, "\n  %-5s %-20s app %-20s %s", m.Type, m.Name, app, m.GUID)
	}
	return errors.New(b.String())
}
