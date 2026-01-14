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

	for _, tz := range *tzs {
		dstStr := "no"
		if tz.HasDst {
			dstStr = "yes"
		}
		aliasStr := ""
		if len(tz.Aliases) > 0 {
			aliasStr = fmt.Sprintf("[ %s ]", strings.Join(tz.Aliases, ", "))
		}
		table.Add(tz.Name, dstStr, aliasStr)
	}

	table.Print()
	return nil
}
