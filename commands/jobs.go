package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

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

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "Job Name\tApp Name\tCommand")

	for _, job := range jobs {
		appName := "**UNKNOWN**"
		if app, err := core.AppByGUID(apps, job.AppGUID); err != nil {
			appName = app.Name
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\n", job.Name, appName, job.Command)
	}

	writer.Flush()
	return nil
}
