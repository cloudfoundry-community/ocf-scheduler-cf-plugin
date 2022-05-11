package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf delete-job JOB-NAME
func DeleteJob(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf delete-job JOB-NAME")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	err = client.DeleteJob(services.Client, job)
	if err != nil {
		fmt.Println("Could not delete job: " + err.Error())
		return
	}

	fmt.Printf(
		"Deleted job %s (%s)\n",
		job.Name,
		job.GUID,
	)
}
