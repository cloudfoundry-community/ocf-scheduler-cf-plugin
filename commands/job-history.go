package commands

import (
	"fmt"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf job-history JOB-NAME
func JobHistory(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf job-history JOB-NAME")
		return
	}

	if err := jobHistory(services, args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func jobHistory(services *core.Services, args []string) error {
	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	name := args[1]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		return fmt.Errorf("Could not find job named %s in space %s.\n", name, space.Name)
	}

	err = core.PrintActionInProgress(services, "Getting scheduled job history for %s", name)
	if err != nil {
		return err
	}

	executions, _ := client.ListJobExecutions(services.Client, job)
	count := len(executions)
	if count == 0 {
		fmt.Printf("No executions for job %s.\n", name)
		return nil
	}

	fmt.Println("1 -", count, "of", count, "Total Results")

	table := core.NewTable().Add(
		"Execution GUID",
		"Execution State",
		"Scheduled Time",
		"Execution Start Time",
		"Execution End Time",
		"Exit Message",
	)

	for _, execution := range executions {
		table.Add(
			execution.GUID,
			execution.State,
			execution.ScheduledTime.String(),
			execution.ExecutionStartTime.String(),
			execution.ExecutionEndTime.String(),
			execution.Message,
		)
	}

	table.Print()
	return nil
}
