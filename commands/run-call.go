package commands

import (
	"fmt"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf run-call CALL-NAME
func RunCall(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf run-call CALL-NAME")
		return
	}

	if err := runCall(services, args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func runCall(services *core.Services, args []string) error {
	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	apps, err := core.MyApps(services)
	if err != nil {
		return fmt.Errorf("Could not get apps.")
	}

	name := args[1]

	call, err := client.CallNamed(services.Client, space, name)
	if err != nil {
		return fmt.Errorf("Could not find call named %s in space %s.\n", name, space.Name)
	}

	appName := "**UNKNOWN**"
	if app, err := core.AppByGUID(apps, call.AppGUID); err == nil {
		appName = app.Name
	}

	err = core.PrintActionInProgress(services, "Enqueuing call %s for app %s", call.Name, appName)
	if err != nil {
		return err
	}

	_, err = client.ExecuteCall(services.Client, call)
	if err != nil {
		return err
	}

	return nil
}
