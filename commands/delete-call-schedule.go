package commands

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf delete-call-schedule CALL-NAME SCHEDULE-GUID
func DeleteCallSchedule(services *core.Services, args []string) {
	var forceFlag bool
	var promptFlag bool

	flags := pflag.NewFlagSet("delete-call-schedule", pflag.ExitOnError)
	flags.BoolVarP(&forceFlag, "force", "f", false, "Force call schedule deletion without confirmation")
	flags.BoolVarP(&promptFlag, "prompt", "p", false, "Allow call schedule deletion with confirmation")
	flags.MarkHidden("prompt")
	flags.Parse(args)
	args = flags.Args()

	if len(args) != 3 {
		fmt.Println("cf delete-call-schedule [OPTIONS] CALL-NAME SCHEDULE-GUID")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]
	scheduleGUID := args[2]

	call, err := client.CallNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find call named %s in space %s.\n", name, space.Name)
		return
	}

	if promptFlag && !forceFlag && !services.UI.ConfirmDelete("call schedule", name+" "+scheduleGUID) {
		return
	}

	err = client.DeleteCallSchedule(services.Client, call, scheduleGUID)
	if err != nil {
		fmt.Println("Could not delete schedule", scheduleGUID)
		return
	}

	fmt.Println("Schedule", scheduleGUID, "deleted.")
}
