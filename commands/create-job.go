package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf create-job APP-NAME JOB-NAME COMMAND
func CreateJob(services *core.Services, args []string) {
	if len(args) != 4 {
		fmt.Println("cf create-job APP-NAME JOB-NAME COMMAND")
		return
	}

	if err := createJob(services, args[1], args[2], args[3]); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println("OK")
}

func createJob(services *core.Services, appName, jobName, command string) error {
	err := core.PrintActionInProgress(services, "Creating job %s for %s with command '%s'", jobName, appName, command)
	if err != nil {
		return err
	}

	app, err := services.CLI.GetApp(appName)
	if err != nil {
		return fmt.Errorf("could not find app with name %s", appName)
	}

	payload, err := client.CreateJob(services.Client, app.Guid, jobName, command)
	if err != nil {
		return err
	}

	// i'm going to refactor this after we make it work. be prepared.
	core.Table(
		[]string{"Job Name", "App Name", "Command"},
		[]string{payload.Name, appName, payload.Command},
	)

	return nil
}
