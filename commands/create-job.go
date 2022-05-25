package commands

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf create-job APP-NAME JOB-NAME COMMAND
func CreateJob(services *core.Services, args []string) {
	var diskInMb int
	var memoryInMb int

	flags := pflag.NewFlagSet("create-job", pflag.ExitOnError)
	flags.IntVarP(&diskInMb, "disk", "k", 1024, "disk limit in MB")
	flags.IntVarP(&memoryInMb, "memory", "m", 1024, "memory limit in MB")
	flags.Parse(args)

	fmt.Println("disk:", diskInMb, "memory:", memoryInMb)
	fmt.Println("args:", args)

	os.Exit(0)

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

	core.
		NewTable().
		Add("Job Name", "App Name", "Command").
		Add(payload.Name, appName, payload.Command).
		Print()

	return nil
}
