package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf schedule-job JOB-NAME CRON-EXPRESSION
func ScheduleJob(services *core.Services, args []string) {
	if len(args) != 3 {
		fmt.Println("cf schedule-job JOB-NAME CRON-EXPRESSION")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]
	cronExpression := args[2]

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
