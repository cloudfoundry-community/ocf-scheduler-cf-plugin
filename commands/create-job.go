package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf create-job APP_NAME NAME COMMAND
func CreateJob(services *core.Services, args []string) {
	if len(args) != 4 {
		fmt.Println("cf create-job APP_NAME NAME COMMAND")
		return
	}

	appName := args[1]
	name := args[2]
	command := args[3]

	app, err := services.CLI.GetApp(appName)
	if err != nil {
		fmt.Println("Could not find app with name", appName)
		return
	}

	payload, err := client.CreateJob(services.Client, app.Guid, name, command)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Created job %s\n\tGUID: %s\n\tApp Name: %s\n\tApp GUID: %s\n\tSpace GUID: %s\n\tCommand: %s\n",
		payload.Name,
		payload.GUID,
		appName,
		payload.AppGUID,
		payload.SpaceGUID,
		payload.Command,
	)
}
