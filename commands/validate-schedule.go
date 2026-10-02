package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
}

// parseValidateFlags reads the options; args starts with the command name,
// which rest drops.
func parseValidateFlags(args []string) (validateOptions, []string, error) {
	var opts validateOptions
	var next int
	flags := pflag.NewFlagSet("validate-schedule", pflag.ExitOnError)
	flags.StringVarP(&opts.timezone, "timezone", "t", "", "Interpret Cron Expression relative to the given timezone")
	flags.IntVar(&next, "next", 0, "Number of next runs to show (default 5; 1 per stored schedule)")
	flags.IntVar(&opts.prev, "prev", 0, "Number of previous runs to show")
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
	return opts, flags.Args()[1:], nil
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

	req := scheduler.ValidateRequest{Prev: opts.prev, Next: opts.next}
	if r.Expression != "" {
		req.Expression = strings.TrimSpace(r.Expression)
		if timezone != "" {
			req.Expression = fmt.Sprintf("CRON_TZ=%s %s", timezone, req.Expression)
		}
	} else if timezone != "" {
		return fail(fmt.Errorf("--timezone needs an expression"))
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
		renderAnalysis(os.Stdout, a)
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
		renderAnalysis(os.Stdout, a)
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
