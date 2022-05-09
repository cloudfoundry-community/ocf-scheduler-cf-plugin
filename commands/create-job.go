package commands

import (
	"encoding/json"
	"fmt"

	"github.com/ess/hype"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
	scheduler "github.com/starkandwayne/scheduler-for-ocf/core"
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

	params := hype.Params{}
	params.Set("app_guid", app.Guid)

	payload := &scheduler.Job{
		Name:    name,
		Command: command,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Could not prepare the request payload")
		return
	}

	response := services.Client.Post("jobs", params, data)

	if !response.Okay() {
		fmt.Println(response.Error())
		return
	}

	err = json.Unmarshal(response.Data(), payload)
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
