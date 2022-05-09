package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf call-schedules CALL-NAME
func CallSchedules(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf call-schedules CALL-NAME")
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

	schedules, _ := client.ListCallSchedules(services.Client, call)
	if len(schedules) == 0 {
		fmt.Printf("No schedules for call %s.\n" + name)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "Call Name\tSchedule\tWhen")
	fmt.Fprintln(writer, "=========\t========\t====")

	for _, schedule := range schedules {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", call.Name, schedule.GUID, schedule.Expression)
	}

	writer.Flush()
}
