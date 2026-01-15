package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf schedule-job JOB-NAME CRON-EXPRESSION
func ScheduleJob(services *core.Services, args []string) {
	var timezone string
	flags := pflag.NewFlagSet("schedule-call", pflag.ExitOnError)
	flags.StringVarP(&timezone, "timezone", "t", "", "Interpret Cron Expression relative to the given timezone")
	flags.Parse(args)
	args = flags.Args()
	if len(args) != 3 {
		fmt.Println("cf schedule-job [OPTIONS] JOB-NAME CRON-EXPRESSION")
		return
	}

	timezone = strings.TrimSpace(timezone)
	if strings.Contains(timezone, " ") {
		fmt.Println("no spaces allowed in the timezone")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]
	cronExpression := strings.TrimSpace(args[2])
	if timezone != "" {
		cronExpression = fmt.Sprintf("CRON_TZ=%s %s", timezone, cronExpression)
	}

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	schedule, err := client.ScheduleJob(services.Client, job, cronExpression)
	if err != nil {
		fmt.Printf("Could not schedule job %s with the expression %s.\n", name, cronExpression)
		return
	}

	core.
		NewTable().
		Add("Job Name", "Schedule GUID", "Expression").
		Add(job.Name, schedule.GUID, schedule.Expression).
		Print()
}
