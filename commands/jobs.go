package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf jobs
func Jobs(services *core.Services, args []string) {

	if err := listJobs(services); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func listJobs(services *core.Services) error {
	err := core.PrintActionInProgress(services, "Listing jobs")
	if err != nil {
		return err
	}

	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	apps, err := core.MyApps(services)
	if err != nil {
		return fmt.Errorf("Could not get apps.")
	}

	jobs, err := client.ListJobs(services.Client, space)
	if err != nil {
		return fmt.Errorf("Could not get jobs for space %s.\n", space.Name)
	}

	table := core.NewTable().Add("Job Name", "App Name", "Command")

	for _, job := range jobs {
		appName := "**UNKNOWN**"
		if app, err := core.AppByGUID(apps, job.AppGUID); err == nil {
			appName = app.Name
		}

		table.Add(job.Name, appName, job.Command)
	}

	table.Print()
	return nil
}
