package commands

import (
	"fmt"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf calls
func Calls(services *core.Services, args []string) {
	if err := listCalls(services); err != nil {
		fmt.Println("Error:", err.Error())
	}

	fmt.Println("OK")
}

func listCalls(services *core.Services) error {
	err := core.PrintActionInProgress(services, "Listing calls")
	if err != nil {
		return err
	}

	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	apps, err := core.MyApps(services)
	if err != nil {
		return fmt.Errorf("Could not get apps.")
	}

	calls, err := client.ListCalls(services.Client, space)
	if err != nil {
		return fmt.Errorf("Could not get calls for space %s.\n", space.Name)
	}

	table := core.NewTable().Add("Call Name", "App Name", "URL")

	for _, call := range calls {
		appName := "**UNKNOWN**"
		if app, err := core.AppByGUID(apps, call.AppGUID); err == nil {
			appName = app.Name
		}

		table.Add(call.Name, appName, call.URL)
	}

	table.Print()
	return nil
}
