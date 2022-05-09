package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf jobs
func Jobs(services *core.Services, args []string) {
	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	jobs, err := core.AllJobs(services, space)
	if err != nil {
		fmt.Printf("Could not get jobs for space %s.\n", space.Name)
		return
	}

	fmt.Printf("Jobs for Space %s:\n\n", space.Name)
	for _, job := range jobs {
		fmt.Printf("\t%s (%s)\n", job.Name, job.GUID)
	}
}
