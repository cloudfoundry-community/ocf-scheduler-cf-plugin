package commands

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf delete-call CALL-NAME
func DeleteCall(services *core.Services, args []string) {
	var forceFlag bool
	var promptFlag bool

	flags := pflag.NewFlagSet("delete-call", pflag.ExitOnError)
	flags.BoolVarP(&forceFlag, "force", "f", false, "Force call deletion without confirmation")
	flags.BoolVarP(&promptFlag, "prompt", "p", false, "Allow call deletion with confirmation")
	flags.MarkHidden("prompt")
	flags.Parse(args)
	args = flags.Args()

	if len(args) != 2 {
		fmt.Println("cf delete-call [OPTIONS] CALL-NAME")
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

	if promptFlag && !forceFlag && !services.UI.ConfirmDeleteWithAssociations("call", name) {
		return
	}

	err = client.DeleteCall(services.Client, call)
	if err != nil {
		fmt.Println("Could not delete call: " + err.Error())
		return
	}

	fmt.Printf(
		"Deleted call %s (%s)\n",
		call.Name,
		call.GUID,
	)
}
