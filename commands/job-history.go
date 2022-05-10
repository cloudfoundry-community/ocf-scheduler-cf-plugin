package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
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
	if len(executions) == 0 {
		return fmt.Errorf("No executions for job %s.\n" + name)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "Execution GUID\tExecution State\tScheduled Time\tExecution Start Time\tExecution End Time\tExit Message")

	for _, execution := range executions {
		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			execution.GUID,
			execution.State,
			execution.ScheduledTime,
			execution.ExecutionStartTime,
			execution.ExecutionEndTime,
			execution.Message,
		)
	}

	writer.Flush()
	return nil
}
