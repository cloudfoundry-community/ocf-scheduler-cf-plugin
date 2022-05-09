package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf call-history CALL-NAME
func CallHistory(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf call-history CALL-NAME")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]

	call, err := client.CallNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find call named %s in space %s.\n", name, space.Name)
		return
	}

	executions, _ := client.ListCallExecutions(services.Client, call)
	if len(executions) == 0 {
		fmt.Printf("No executions for call %s.\n" + name)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "State\tExecution Start Time\tExecution End Time")
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
