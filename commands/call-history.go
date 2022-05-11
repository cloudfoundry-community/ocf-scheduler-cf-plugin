package commands

import (
	"fmt"

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

	rows := [][]string{
		[]string{
			"Execution GUID",
			"Execution State",
			"Scheduled Time",
			"Execution Start Time",
			"Execution End Time",
			"Exit Message",
		},
	}

	for _, execution := range executions {
		rows = append(
			rows,
			[]string{
				execution.GUID,
				execution.State,
				execution.ScheduledTime.String(),
				execution.ExecutionStartTime.String(),
				execution.ExecutionEndTime.String(),
				execution.Message,
			},
		)
	}

	core.Table(rows)

	return nil
}
