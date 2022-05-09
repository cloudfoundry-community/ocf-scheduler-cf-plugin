package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
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

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
	fmt.Fprintln(writer, "Call Name\tApp Name\tURL")

	for _, call := range calls {
		appName := "**UNKNOWN**"
		if app, err := core.AppByGUID(apps, call.AppGUID); err == nil {
			appName = app.Name
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\n", call.Name, appName, call.URL)
	}

	writer.Flush()
	return nil
}
