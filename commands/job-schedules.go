package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	scheduler "github.com/starkandwayne/scheduler-for-ocf/core"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf job-schedules
func JobSchedules(services *core.Services, args []string) {
	// TODO: redesign this command such that it shows schedules for all jobs
	// so as to match the upstream UX
	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	// BEGIN REAL IMPLEMENTATION
	output := make(map[string][]string)
	allApps, err := core.MyApps(services)
	if err != nil {
		fmt.Print("Could not get apps.")
		return
	}

	appJobs := make(map[string][]*scheduler.Job)
	allJobs, err := client.ListJobs(services.Client, space)
	if err != nil {
		fmt.Printf("Could not get jobs for space %s.\n", space.Name)
		return
	}

	for _, job := range allJobs {
		if appJobs[job.AppGUID] == nil {
			appJobs[job.AppGUID] = make([]*scheduler.Job, 0)
		}

		appJobs[job.AppGUID] = append(appJobs[job.AppGUID], job)
	}

	for appGUID, jobs := range appJobs {
		if len(jobs) == 0 {
			continue
		}

		for _, job := range jobs {
			schedules, serr := client.ListJobSchedules(services.Client, job)
			if serr != nil {
				continue
			}

			if len(schedules) == 0 {
				continue
			}

			if output[appGUID] == nil {
				output[appGUID] = []string{
					"Job Name\tCommand\tSchedule\tExpression\n",
				}
			}

			for _, schedule := range schedules {
				output[appGUID] = append(
					output[appGUID],
					fmt.Sprintf(
						"%s\t%s\t%s\t%s\n",
						job.Name,
						job.Command,
						schedule.GUID,
						schedule.Expression,
					),
				)
			}
		}

		app, err := core.AppByGUID(allApps, appGUID)
		if err != nil {
			continue
		}

		fmt.Printf("\nApp Name: %s\n", app.Name)

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
		for _, line := range output[appGUID] {
			fmt.Fprintln(writer, line)
		}
		writer.Flush()
	}

	fmt.Println("OK")
	// END REAL IMPLEMENTATION

	//name := args[1]

	//job, err := client.JobNamed(services.Client, space, name)
	//if err != nil {
	//fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
	//return
	//}

	//schedules, _ := client.ListJobSchedules(services.Client, job)
	//if len(schedules) == 0 {
	//fmt.Printf("No schedules for job %s.\n", name)
	//return
	//}

	//writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	//fmt.Fprintln(writer, "Job Name\tSchedule\tWhen")
	//fmt.Fprintln(writer, "=========\t========\t====")

	//for _, schedule := range schedules {
	//fmt.Fprintf(writer, "%s\t%s\t%s\n", job.Name, schedule.GUID, schedule.Expression)
	//}

	//writer.Flush()
}
