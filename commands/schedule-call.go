package commands

import (
	"fmt"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf schedule-call CALL-NAME SCHEDULE
func ScheduleCall(services *core.Services, args []string) {
	if len(args) != 3 {
		fmt.Println("cf schedule-call CALL-NAME SCHEDULE")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]
	cronExpression := args[2]

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

	core.
		NewTable().
		Add("Call Name", "Schedule GUID", "Expression").
		Add(call.Name, schedule.GUID, schedule.Expression).
		Print()
}
