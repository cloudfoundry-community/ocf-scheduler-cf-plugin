package commands

import (
	"fmt"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf run-job JOB-NAME
func RunJob(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf run-job JOB-NAME")
		return
	}

	if err := runJob(services, args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func runJob(services *core.Services, args []string) error {
	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	apps, err := core.MyApps(services)
	if err != nil {
		return fmt.Errorf("Could not get apps.")
	}

	name := args[1]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		return fmt.Errorf("Could not find job named %s in space %s.\n", name, space.Name)
	}

	appName := "**UNKNOWN**"
	if app, err := core.AppByGUID(apps, job.AppGUID); err == nil {
		appName = app.Name
	}

	err = core.PrintActionInProgress(services, "Enqueuing job %s for app %s", job.Name, appName)
	if err != nil {
		return err
	}

	_, err = client.ExecuteJob(services.Client, job)
	if err != nil {
		return err
	}

	return nil
}
