package commands

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/client"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

// cf delete-job JOB-NAME
func DeleteJob(services *core.Services, args []string) {
	var forceFlag bool
	var promptFlag bool

	flags := pflag.NewFlagSet("delete-job", pflag.ExitOnError)
	flags.BoolVarP(&forceFlag, "force", "f", false, "Force job deletion without confirmation")
	flags.BoolVarP(&promptFlag, "prompt", "p", false, "Allow job deletion with confirmation")
	flags.MarkHidden("prompt")
	flags.Parse(args)
	args = flags.Args()

	if len(args) != 2 {
		fmt.Println("cf delete-job [OPTIONS] JOB-NAME")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]

	job, err := client.JobNamed(services.Client, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	if promptFlag && !forceFlag && !services.UI.ConfirmDeleteWithAssociations("job", name) {
		return
	}

	err = client.DeleteJob(services.Client, job)
	if err != nil {
		fmt.Println("Could not delete job: " + err.Error())
		return
	}

	fmt.Printf(
		"Deleted job %s (%s)\n",
		job.Name,
		job.GUID,
	)
}
