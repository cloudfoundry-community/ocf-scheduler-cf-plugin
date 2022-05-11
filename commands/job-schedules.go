package commands

import (
	"fmt"

	scheduler "github.com/starkandwayne/scheduler-for-ocf/core"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf job-schedules
func JobSchedules(services *core.Services, args []string) {
	if err := jobSchedules(services); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func jobSchedules(services *core.Services) error {
	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	// BEGIN REAL IMPLEMENTATION
	output := make(map[string][]string)
	allApps, err := core.MyApps(services)
	if err != nil {
		return fmt.Errorf("Could not get apps.")
	}

	appJobs := make(map[string][]*scheduler.Job)
	allJobs, err := client.ListJobs(services.Client, space)
	if err != nil {
		return fmt.Errorf("Could not get jobs for space %s.\n", space.Name)
	}

	err = core.PrintActionInProgress(services, "Getting scheduled jobs")
	if err != nil {
		return err
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
					"Job Name\tCommand\tSchedule\tExpression",
				}
			}

			for _, schedule := range schedules {
				output[appGUID] = append(
					output[appGUID],
					fmt.Sprintf(
						"%s\t%s\t%s\t%s",
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

		if len(output[appGUID]) == 0 {
			fmt.Println("App", app.Name, "has jobs, but none are scheduled.")
		}

		table := core.NewTable()

		for _, line := range output[appGUID] {
			table.AddRow(line)
		}

		table.Print()
	}

	return nil
}
