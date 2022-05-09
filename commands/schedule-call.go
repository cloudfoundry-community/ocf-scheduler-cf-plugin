package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf schedule-call CALL-NAME SCHEDULE
func ScheduleCall(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf schedule-call CALL-NAME SCHEDULE")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[0]
	cronExpression := args[1]

	call, err := client.CallNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	schedule, err := client.ScheduleCall(services.Client, call, cronExpression)
	if err != nil {
		fmt.Println("Could not schedule call: " + err.Error())
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "Call Name\tSchedule\tWhen")
	fmt.Fprintln(writer, "=========\t========\t====")

	fmt.Fprintf(
		writer,
		"%s\t%s\t%s\nn",
		call.Name,
		schedule.GUID,
		schedule.Expression,
	)

	writer.Flush()
}
