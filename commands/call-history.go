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

	if err := callHistory(services, args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func callHistory(services *core.Services, args []string) error {
	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	name := args[1]

	call, err := client.CallNamed(services.Client, space, name)
	if err != nil {
		return fmt.Errorf("Could not find call named %s in space %s.\n", name, space.Name)
	}

	err = core.PrintActionInProgress(services, "Getting scheduled call history for %s", name)
	if err != nil {
		return err
	}

	executions, _ := client.ListCallExecutions(services.Client, call)
	count := len(executions)
	if count == 0 {
		fmt.Printf("No executions for call %s.\n", name)
		return nil
	}

	fmt.Println("1 -", count, "of", count, "Total Results")

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
