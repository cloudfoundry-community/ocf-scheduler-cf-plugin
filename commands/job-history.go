package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf job-history NAME
func JobHistory(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf job-history NAME")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[0]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	executions, _ := client.ListJobExecutions(core.Client, job)
	if len(executions) == 0 {
		fmt.Printf("No executions for job %s.\n" + name)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "State\tStart Time\tEnd Time")
	fmt.Fprintln(writer, "=====\t====================\t==================")

	for _, execution := range executions {
		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\n",
			execution.State,
			execution.ExecutionStartTime,
			execution.ExecutionEndTime,
		)
	}

	writer.Flush()
}
