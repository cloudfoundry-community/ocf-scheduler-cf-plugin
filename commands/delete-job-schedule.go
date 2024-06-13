package commands

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf delete-job-schedule JOB-NAME SCHEDULE-GUID
func DeleteJobSchedule(services *core.Services, args []string) {
    var forceFlag bool
    var promptFlag bool

	flags := pflag.NewFlagSet("delete-job-schedule", pflag.ExitOnError)
    flags.BoolVarP(&forceFlag, "force", "f", false, "Force job schedule deletion without confirmation")
    flags.BoolVarP(&promptFlag, "prompt", "p", false, "Allow job schedule deletion with confirmation")
    flags.MarkHidden("prompt")
	flags.Parse(args)

	args = flags.Args()
	if len(args) != 3 {
		fmt.Println("cf delete-job-schedule JOB-NAME SCHEDULE-GUID [OPTIONS]")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]
	scheduleGUID := args[2]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

    if promptFlag && !forceFlag && !services.UI.ConfirmDelete("job schedule", name + " " + scheduleGUID) {
        return;
    }

	err = client.DeleteJobSchedule(services.Client, job, scheduleGUID)
	if err != nil {
		fmt.Println("Could not delete schedule", scheduleGUID)
		return
	}

	fmt.Println("Schedule", scheduleGUID, "deleted.")
}
