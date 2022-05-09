package commands

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf create-call APP_NAME NAME URL
func CreateCall(services *core.Services, args []string) {
	if len(args) != 4 {
		fmt.Println("cf create-call APP_NAME NAME URL")
		return
	}

	appName := args[0]
	name := args[1]
	url := args[2]

	app, err := services.CLI.GetApp(appName)
	if err != nil {
		fmt.Println("Could not find app with name", appName)
		return
	}

	payload, err := client.CreateCall(services.Client, app.Guid, name, url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Created call %s\n\tGUID: %s\n\tApp GUID: %s\n\tSpace GUID: %s\n\tURL: %s\n\tAuth Header: %s\n",
		payload.Name,
		payload.GUID,
		payload.AppGUID,
		payload.SpaceGUID,
		payload.URL,
		payload.AuthHeader,
	)
}
