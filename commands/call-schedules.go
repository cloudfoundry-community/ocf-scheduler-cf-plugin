package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	scheduler "github.com/starkandwayne/scheduler-for-ocf/core"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/client"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// cf call-schedules
func CallSchedules(services *core.Services, args []string) {
	if err := callSchedules(services); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func callSchedules(services *core.Services) error {
	space, err := core.MySpace(services)
	if err != nil {
		return fmt.Errorf("Could not get current space.")
	}

	// BEGIN REAL IMPLEMENTATION
	output := make(map[string][]string)
	allApps, err := core.MyApps(services)
	if err != nil {
		return fmt.Errorf("Could not get apps.")
	}

	appCalls := make(map[string][]*scheduler.Call)
	allCalls, err := client.ListCalls(services.Client, space)
	if err != nil {
		return fmt.Errorf("Could not get calls for space %s.\n", space.Name)
	}

	err = core.PrintActionInProgress(services, "Getting scheduled calls")
	if err != nil {
		return err
	}

	for _, call := range allCalls {
		if appCalls[call.AppGUID] == nil {
			appCalls[call.AppGUID] = make([]*scheduler.Call, 0)
		}

		appCalls[call.AppGUID] = append(appCalls[call.AppGUID], call)
	}

	for appGUID, calls := range appCalls {
		if len(calls) == 0 {
			continue
		}

		for _, call := range calls {
			schedules, serr := client.ListCallSchedules(services.Client, call)
			if serr != nil {
				continue
			}

			if len(schedules) == 0 {
				continue
			}

			if output[appGUID] == nil {
				output[appGUID] = []string{
					"Call Name\tURL\tSchedule\tExpression",
				}
			}

			for _, schedule := range schedules {
				output[appGUID] = append(
					output[appGUID],
					fmt.Sprintf(
						"%s\t%s\t%s\t%s",
						call.Name,
						call.URL,
						schedule.GUID,
						schedule.Expression,
					),
				)
			}
		}

		app, err := core.AppByGUID(allApps, appGUID)
		if err != nil {
			continue
		}

		fmt.Printf("\nApp Name: %s\n", app.Name)
		if len(output[appGUID]) == 0 {
			fmt.Println("App", app.Name, "has calls, but none are scheduled.")
			continue
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', uint(0))
		for _, line := range output[appGUID] {
			fmt.Fprintln(writer, line)
		}
		writer.Flush()
	}

	return nil
}
