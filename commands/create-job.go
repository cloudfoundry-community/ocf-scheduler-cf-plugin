package commands

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

func quotaInMb(quota string) (int, error) {
	parsed, err := resource.ParseQuantity(quota)
	if err != nil {
		return 0, err
	}

	return int(parsed.ScaledValue(resource.Mega)), nil
}

// cf create-job APP-NAME JOB-NAME COMMAND
func CreateJob(services *core.Services, args []string) {
	var diskQuota string
	var memoryQuota string

	flags := pflag.NewFlagSet("create-job", pflag.ExitOnError)
	flags.StringVarP(&diskQuota, "disk", "k", "1024M", "disk limit")
	flags.StringVarP(&memoryQuota, "memory", "m", "1024M", "memory limit")
	flags.Parse(args)
	args = flags.Args()

	diskInMb, err := quotaInMb(diskQuota)
	if err != nil {
		fmt.Println("Error: Couldn't parse disk limit:", err.Error())
		return
	}

	memoryInMb, err := quotaInMb(memoryQuota)
	if err != nil {
		fmt.Println("Error: Couldn't parse memory limit:", err.Error())
		return
	}

	if len(args) != 4 {
		fmt.Println("cf create-job APP-NAME JOB-NAME COMMAND")
		return
	}

	if err := createJob(services, args[1], args[2], args[3], diskInMb, memoryInMb); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println("OK")
}

func createJob(services *core.Services, appName, jobName, command string, diskInMb, memoryInMb int) error {
	err := core.PrintActionInProgress(services, "Creating job %s for %s with command '%s'", jobName, appName, command)
	if err != nil {
		return err
	}

	app, err := services.CLI.GetApp(appName)
	if err != nil {
		return fmt.Errorf("could not find app with name %s", appName)
	}

	payload, err := client.CreateJob(services.Client, app.Guid, jobName, command, diskInMb, memoryInMb)
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
