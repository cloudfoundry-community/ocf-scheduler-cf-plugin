package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf job-history JOB-NAME
func JobHistory(services *core.Services, args []string) {
	filterOutput := "scheduled"

	flags := pflag.NewFlagSet("job-history", pflag.ContinueOnError)
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
			return errors.New("The display parameter value must be a prefix of the words, scheduled, manual or all.")
		}
	})
	flags.Parse(args)
	args = flags.Args()

	if len(args) != 2 {
		fmt.Println("cf job-history JOB-NAME")
		return
	}

	if err := jobHistory(services, filterOutput, args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func jobHistory(services *core.Services, filterOutput string, args []string) error {
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

		var scheduledTime string

		if execution.ScheduledTime.IsZero() {
			scheduledTime = "manual"
		} else {
			scheduledTime =execution.ScheduledTime.String()
		}

		if filterOutput == "all" || filterOutput == scheduledTime || (filterOutput == "scheduled" && scheduledTime != "manual") {
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

	table.Print()
	return nil
}
