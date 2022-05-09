package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf schedule-job JOB-NAME SCHEDULE
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

	jobName := args[1]
	cronExpression := args[2]
	job, err := core.JobNamed(services, space, jobName)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", jobName, space.Name)
		return
	}

	schedule, err := core.ScheduleJob(services, job, cronExpression)
	if err != nil {
		fmt.Printf("Could not schedule job %s with the expression %s.\n", jobName, cronExpression)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "Job Name\tSchedule\tWhen")
	fmt.Fprintln(writer, "========\t========\t====")

	fmt.Fprintf(
		writer,
		"%s\t%s\t%s\n",
		job.Name,
		schedule.GUID,
		schedule.Expression,
	)

	writer.Flush()
}
