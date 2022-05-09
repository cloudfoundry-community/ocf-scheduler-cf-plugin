package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf job-schedules JOB-NAME
func JobSchedules(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf job-schedules JOB-NAME")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[0]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	schedules, _ := client.ListJobSchedules(services.Client, job)
	if len(schedules) == 0 {
		fmt.Printf("No schedules for job %s.\n" + name)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "Job Name\tSchedule\tWhen")
	fmt.Fprintln(writer, "=========\t========\t====")

	for _, schedule := range schedules {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", job.Name, schedule.GUID, schedule.Expression)
	}

	writer.Flush()
}
