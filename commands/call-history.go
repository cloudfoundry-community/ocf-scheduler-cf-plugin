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
	flags.FuncP("show", "s", "display scheduled, manual or complete call execution history", func(value string) error {
		if strings.HasPrefix("scheduled", value) {
			filterOutput = "scheduled"
			return nil
		} else if strings.HasPrefix("manual", value) {
			filterOutput = "manual"
			return nil
		} else if strings.HasPrefix("all", value) {
			filterOutput = "complete"
			return nil
		} else {
			return errors.New("The show parameter value must be a prefix of one of theses words, \"scheduled\", \"manual\" or \"all\".")
		}
	})
	flags.Parse(args)
	args = flags.Args()

	if len(args) != 2 {
		fmt.Println("cf call-history [OPTIONS] CALL-NAME")
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

	err = core.PrintActionInProgress(services, "Getting %s call history for %s", filterOutput, name)
	if err != nil {
		return err
	}

	executions, _ := client.ListCallExecutions(services.Client, call)
	totalCount := len(executions)
	var manualCount, scheduledCount int

	filterCount := func() int {
		switch filterOutput {
		case "scheduled":
			return scheduledCount
		case "manual":
			return manualCount
		}
		return totalCount
	}

	filterDisplayName := func() string {
		switch filterOutput {
		case "scheduled":
			return filterOutput
		case "manual":
			return "ad hoc"
		}
		return ""
	}

	table := core.NewTable().Add(
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
			manualCount++
		} else {
			scheduledTime = execution.ScheduledTime.String()
			scheduledCount++
		}

		if filterOutput == "complete" || filterOutput == scheduledTime || (filterOutput == "scheduled" && scheduledTime != "manual") {
			table.Add(
				execution.GUID,
				execution.State,
				scheduledTime,
				execution.ExecutionStartTime.String(),
				execution.ExecutionEndTime.String(),
				execution.Message,
			)
		}
	}

	if filterCount() == 0 {
		fmt.Printf("No%s executions for call %s\n", core.AddSpace(filterDisplayName()), name)
		return nil
	}

	makeExecutionPlural := func(v int) string {
		s := "execution"
		if v != 0 {
			s += "s"
		}
		return s
	}

	fmt.Printf("1 - %v out of %v call %s\n", filterCount(), totalCount, makeExecutionPlural(totalCount))

	table.Print()

	return nil
}
