package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf delete-call-schedule NAME GUID
func DeleteCallSchedule(services *core.Services, args []string) {
	if len(args) != 4 {
		fmt.Println("cf delete-call-schedule NAME GUID")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[0]
	scheduleGUID := args[1]

	call, err := client.CallNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find call named %s in space %s.\n", name, space.Name)
		return
	}

	err = client.DeleteCallSchedule(services.Client, call, scheduleGUID)
	if err != nil {
		fmt.Println("Could not delete schedule", scheduleGUID)
		return
	}

	fmt.Println("Schedule", scheduleGUID, "deleted.")
}
