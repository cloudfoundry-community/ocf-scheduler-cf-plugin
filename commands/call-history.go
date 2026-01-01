package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf call-history CALL-NAME
func CallHistory(services *core.Services, args []string) {
	filterOutput := "scheduled"

	flags := pflag.NewFlagSet("call-history", pflag.ExitOnError)
	flags.FuncP("display","d", "display scheduled, manual or complete execution histories", func(value string) error {
		if strings.HasPrefix("scheduled", value) {
			filterOutput="scheduled"
			return nil
		} else if strings.HasPrefix("manual", value) {
			filterOutput = "manual"
			return nil
		} else if strings.HasPrefix("all", value) {
			filterOutput = "all"
			return nil
		} else {
			return errors.New("The display parameter value must be a prefix of one of theses words, \"scheduled\", \"manual\" or \"all\".")
		}
	})
	flags.Parse(args)
	args = flags.Args()

	if len(args) != 2 {
		fmt.Println("cf call-history CALL-NAME")
		return
	}

	if err := callHistory(services, filterOutput, args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func callHistory(services *core.Services, filterOutput string, args []string) error {
	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	name := args[1]

	call, err := client.CallNamed(services.Client, space, name)
	if err != nil {
		return fmt.Errorf("Could not find call named %s in space %s.\n", name, space.Name)
	}

	err = core.PrintActionInProgress(services, "Getting call history for %s", name)
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

	output := core.NewTable().Add(
		"Execution GUID",
		"Execution State",
		"Scheduled Time",
		"Execution Start Time",
		"Execution End Time",
		"Exit Message",
	)

	for _, execution := range executions {
		var scheduledTime string

		if execution.ScheduledTime.IsZero() {
			scheduledTime = "manual"
		} else {
			scheduledTime =execution.ScheduledTime.String()
		}

		if filterOutput == "all" || filterOutput == scheduledTime || (filterOutput == "scheduled" && scheduledTime != "manual") {
			output.Add(
				execution.GUID,
				execution.State,
				scheduledTime,
				execution.ExecutionStartTime.String(),
				execution.ExecutionEndTime.String(),
				execution.Message,
			)
		}
	}

	output.Print()

	return nil
}
