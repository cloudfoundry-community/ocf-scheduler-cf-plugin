package commands

import (
	"fmt"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf create-call APP-NAME CALL-NAME URL
func CreateCall(services *core.Services, args []string) {
	if len(args) != 4 {
		fmt.Println("cf create-call APP-NAME CALL-NAME URL")
		return
	}

	if err := createCall(services, args[1], args[2], args[3]); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println("OK")
}

func createCall(services *core.Services, appName, callName, url string) error {
	err := core.PrintActionInProgress(services, "Creating call %s for %s with url '%s'", callName, appName, url)
	if err != nil {
		return err
	}

	app, err := services.CLI.GetApp(appName)
	if err != nil {
		return fmt.Errorf("could not find app with name %s", appName)
	}

	payload, err := client.CreateCall(services.Client, app.Guid, callName, url)
	if err != nil {
		return err
	}

	core.
		NewTable().
		Add("Call Name", "App Name", "URL").
		Add(payload.Name, appName, payload.URL).
		Print()

	return nil
}
