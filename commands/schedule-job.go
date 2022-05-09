package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf schedule-job NAME SCHEDULE
func ScheduleJob(services *core.Services, args []string) {
	if len(args) != 3 {
		fmt.Println("cf schedule-job NAME SCHEDULE")
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
