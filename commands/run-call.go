package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf run-call CALL-NAME
func RunCall(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf run-call CALL-NAME")
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

	execution, err := client.ExecuteCall(services.Client, call)
	if err != nil {
		fmt.Println("Could not execute call: " + err.Error())
		return
	}

	fmt.Printf(
		"Executed call %s (%s) [Execution GUID: %s]\n",
		call.Name,
		call.GUID,
		execution.GUID,
	)
}
