package commands

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf scheduler-time-zones
func SchedulerTimeZones(services *core.Services, args []string) {

	if err := listTimeZones(services); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("OK")
}

func listTimeZones(services *core.Services) error {
	err := core.PrintActionInProgress(services, "Listing time zones")
	if err != nil {
		return err
	}

	tzs, err := client.ListTimeZones(services.Client)
	if err != nil {
		return fmt.Errorf("Could not get scheduler time zones %w\n", err)
	}

	table := core.NewTable().Add("Time Zone", "Dst", "Aliases")

	var serverTimezone []string = make([]string, 0, 2)

	for _, tz := range *tzs {
		dstStr := "no"
		if tz.HasDst {
			dstStr = "yes"
		}
		if tz.IsServerTimeZone != "" {
			serverTimezone = append(serverTimezone, tz.Name)
		}
		aliasStr := ""
		if len(tz.Aliases) > 0 {
			aliasStr = fmt.Sprintf("[ %s ]", strings.Join(tz.Aliases, ", "))
		}
		table.Add(tz.Name, dstStr, aliasStr)
	}

	table.Print()

	switch len(serverTimezone) {
	case 0:
		fmt.Println("\nServer timezone was not discovered")
	case 1:
		fmt.Println("\nServer timezone is", serverTimezone[0])
	default:
		fmt.Println("\nMultiple server timezones was discovered", serverTimezone)
	}
	return nil
}
