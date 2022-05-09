package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf delete-job-schedule GUID
func DeleteJobSchedule(services *core.Services, args []string) {
	if len(args) != 4 {
		fmt.Println("cf delete-job-schedule NAME GUID")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[0]
	scheduleGUID := args[1]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	err = client.DeleteJobSchedule(services.Client, job, scheduleGUID)
	if err != nil {
		fmt.Println("Could not delete schedule", scheduleGUID)
		return
	}

	fmt.Println("Schedule", scheduleGUID, "deleted.")
}
