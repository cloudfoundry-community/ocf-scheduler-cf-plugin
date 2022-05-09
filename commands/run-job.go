package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf run-job NAME
func RunJob(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf run-job NAME")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]

	job, err := core.JobNamed(services, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	if core.ExecuteJob(services, job) != nil {
		fmt.Printf("Could not execute job %s.\n", job.Name)
		return
	}

	fmt.Println("OK")
}
