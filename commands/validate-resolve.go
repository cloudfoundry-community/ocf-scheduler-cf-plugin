package commands

import (
	"fmt"
	"strings"

	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
)

// target is the job or call an argument named.
type target struct {
	Type, GUID, Name, AppGUID string
}

// resolution is what validate-schedule's arguments mean.
type resolution struct {
	Target     *target // nil for a bare expression
	Expression string  // "" to re-check the target's stored schedules
}

// AmbiguousError is a name that matches more than one job or call.
type AmbiguousError struct {
	Name    string
	Matches []target
}

func (e *AmbiguousError) Error() string {
	for _, m := range e.Matches[1:] {
		if m.Type != e.Matches[0].Type {
			return fmt.Sprintf("%q matches %d jobs or calls; put job or call before it, or use the GUID", e.Name, len(e.Matches))
		}
	}
	// Same type on different apps: only the GUID tells them apart.
	return fmt.Sprintf("%q matches %d %ss; use the GUID", e.Name, len(e.Matches), e.Matches[0].Type)
}

// resolveArgs reads [job|call] [NAME|GUID] [EXPRESSION]. A name or GUID is
// looked up before anything is treated as an expression: names may contain
// spaces, and are unique per app, not per space.
func resolveArgs(args []string, jobs []*scheduler.Job, calls []*scheduler.Call) (resolution, error) {
	if err := checkArgs(args); err != nil {
		return resolution{}, err
	}
	kinds := []string{"job", "call"}
	if len(args) >= 2 && (args[0] == "job" || args[0] == "call") {
		kinds, args = args[:1], args[1:]
	}
	var matches []target
	for _, kind := range kinds {
		if kind == "job" {
			for _, j := range jobs {
				if j.GUID == args[0] || j.Name == args[0] {
					matches = append(matches, target{"job", j.GUID, j.Name, j.AppGUID})
				}
			}
		} else {
			for _, c := range calls {
				if c.GUID == args[0] || c.Name == args[0] {
					matches = append(matches, target{"call", c.GUID, c.Name, c.AppGUID})
				}
			}
		}
	}
	expression := ""
	if len(args) == 2 {
		expression = args[1]
	}
	switch {
	case len(matches) == 1:
		return resolution{Target: &matches[0], Expression: expression}, nil
	case len(matches) > 1:
		return resolution{}, &AmbiguousError{Name: args[0], Matches: matches}
	case len(kinds) == 1:
		return resolution{}, fmt.Errorf("no %s named %q in this space", kinds[0], args[0])
	case len(args) == 1:
		return resolution{Expression: args[0]}, nil
	}
	return resolution{}, fmt.Errorf("no job or call named %q in this space; quote a cron expression: cf validate-schedule \"0 2 * * *\"", args[0])
}

// validateUsage is shown when the arguments cannot be read.
const validateUsage = `USAGE:
   cf validate-schedule [OPTIONS] CRON-EXPRESSION
   cf validate-schedule [OPTIONS] [job|call] NAME-OR-GUID [CRON-EXPRESSION]

OPTIONS:
   --timezone, -t ZONE   Zone for CRON-EXPRESSION; for stored schedules, the zone run times are shown in.
   --next N              Number of next runs to show (default 5; 1 per stored schedule).
   --prev N              Number of previous runs to show (default 0).`

// checkArgs rejects argument shapes resolveArgs cannot read, before any
// request is made.
func checkArgs(args []string) error {
	if len(args) >= 2 && (args[0] == "job" || args[0] == "call") {
		args = args[1:]
	}
	switch {
	case len(args) == 0:
		return fmt.Errorf("give a cron expression, a job or call, or both\n\n%s", validateUsage)
	case len(args) > 2:
		return fmt.Errorf("too many arguments; quote a cron expression: cf validate-schedule \"0 2 * * *\"\n\n%s", validateUsage)
	}
	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			return fmt.Errorf("empty argument; quote a cron expression or name a job or call")
		}
	}
	return nil
}
