package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf calls
func Calls(services *core.Services, args []string) {
	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	calls, err := client.ListCalls(services.Client, space)
	if err != nil {
		fmt.Printf("Could not get calls for space %s.\n", space.Name)
		return
	}

	fmt.Printf("Calls for Space %s:\n\n", space.Name)
	for _, call := range calls {
		fmt.Printf("\t%s (%s)\n", call.Name, call.GUID)
	}
}
