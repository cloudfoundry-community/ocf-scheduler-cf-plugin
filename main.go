package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/commands"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

type OCFScheduler struct{}

func (c *OCFScheduler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "OCFScheduler",
		Commands: []plugin.Command{
			{
				Name:     "create-job",
				HelpText: "Creates a job (task) related to an app.",
				UsageDetails: plugin.Usage{
					Usage: "create-job:\n\tcf create-job APP-NAME JOB-NAME COMMAND\n\nWHERE\n\tAPP-NAME is the name of the cf app environment to execute with\n\tJOB-NAME is the name for this job (task)\n\tCOMMAND is the name of the command to execute within the app environment.\n\nTODO: --disk LIMIT/-k LIMIT set job(task) disk limit\n\tTODO: --memory LIMIT//-m LIMIT set the job(task) memory limit\n\tLIMIT is an integer suffixed with MB or GB.",
				},
			},
			{
				Name:     "run-job",
				HelpText: "Runs the job (task) with the given name once.",
				UsageDetails: plugin.Usage{
					Usage: "run-job:\n\tcf run-job JOB-NAME",
				},
			},
			{
				Name:     "schedule-job",
				HelpText: "Schedules the named job (task) to run based on the given cron schedule.",
				UsageDetails: plugin.Usage{
					Usage: "schedule-job:\n\tcf schedule-job JOB-NAME CRON-EXPRESSION\n\nWHERE\n\tJOB-NAME is the name of the created job\n\tCRON-EXPRESSION is the cron schedule format \"MIN HOUR DAY-OF-MONTH MONTH DAY-OF-WEEK\"",
				},
			},
			{
				Name:     "jobs",
				HelpText: "Lists created jobs.",
				UsageDetails: plugin.Usage{
					Usage: "jobs:\ncf jobs",
				},
			},
			{
				Name:     "job-schedules",
				HelpText: "Lists created job schedules",
				UsageDetails: plugin.Usage{
					Usage: "job-schedules:\n\tcf job-schedules",
				},
			},
			{
				Name:     "job-history",
				HelpText: "Lists execution history for the given job name",
				UsageDetails: plugin.Usage{
					Usage: "job-history:\n\tcf job-history JOB-NAME",
				},
			},
			{
				Name:     "delete-job",
				HelpText: "Deletes named job.",
				UsageDetails: plugin.Usage{
					Usage: "delete-job:\n\tcf delete-job JOB-NAME",
				},
			},
			{
				Name:     "delete-job-schedule",
				HelpText: "Deletes the job scheduled with the named GUID.",
				UsageDetails: plugin.Usage{
					Usage: "delete-job-schedule:\n\tcf delete-job-schedule JOB-NAME SCHEDULE-GUID",
				},
			},
			{
				Name:     "create-call",
				HelpText: "Creates a web request call",
				UsageDetails: plugin.Usage{
					Usage: "create-call:\n\tcf create-call APP-NAME CALL-NAME URL\nWHERE\n\tAPP-NAME is the name of the cf app to create a call for\n\tCALL-NAME is a name to refer to the call as\n\tURL is the URL to call.",
				},
			},
			{
				Name:     "run-call",
				HelpText: "Execute a named call request once.",
				UsageDetails: plugin.Usage{
					Usage: "run-call:\n\tcf run-call CALL-NAME",
				},
			},
			{
				Name:     "schedule-call",
				HelpText: "Schedules a call to be run based on the supplied cron schedule",
				UsageDetails: plugin.Usage{
					Usage: "schedule-call:\n\tcf schedule-call CALL-NAME SCHEDULE\n\tCALL-NAME is a name for the scheduled call\n\tSCHEUDLE is a schedule using cron schedule format \"MIN HOUR DAY-OF-MONTH DAY-OF-WEEK\"\n\nEXAMPLE\n\tcf schedule-call hourlyrun \"0 * * * *\"",
				},
			},
			{
				Name:     "calls",
				HelpText: "List created calls",
				UsageDetails: plugin.Usage{
					Usage: "calls:\n\tcf calls",
				},
			},
			{
				Name:     "call-schedules",
				HelpText: "List calls scheduled to be run with app and schedule.",
				UsageDetails: plugin.Usage{
					Usage: "call-schedules:\n\tcf call-schedules",
				},
			},
			{
				Name:     "call-history",
				HelpText: "Shows the execution history for the named call.",
				UsageDetails: plugin.Usage{
					Usage: "call-history:\n\tcf call-history CALL-NAME",
				},
			},
			{
				Name:     "delete-call",
				HelpText: "Deletes the named call.",
				UsageDetails: plugin.Usage{
					Usage: "delete-call:\n\tcf delete-call CALL-NAME",
				},
			},
			{
				Name:     "delete-call-schedule",
				HelpText: "Delete a call scheduled with a given GUID",
				UsageDetails: plugin.Usage{
					Usage: "delete-call-schedule:\n\tcf delete-call-schedule CALL-NAME SCHEDULE-GUID",
				},
			},
		},
	}
}

func (c *OCFScheduler) Run(cliConnection plugin.CliConnection, args []string) {
	scheduler, err := core.GetScheduler(cliConnection)
	if err != nil {
		fmt.Println(err)
		return
	}

	token, err := core.GetBearer(cliConnection)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := core.NewDriver(scheduler, token)
	if err != nil {
		fmt.Println("Could not create a Scheduler API client.")
		return
	}

	services := &core.Services{CLI: cliConnection, Client: client}

	switch args[0] {
	case "create-job":
		commands.CreateJob(services, args)

	case "run-job":
		commands.RunJob(services, args)

	case "schedule-job":
		commands.ScheduleJob(services, args)

	case "jobs":
		commands.Jobs(services, args)

	case "job-schedules":
		commands.JobSchedules(services, args)

	case "job-history":
		commands.JobHistory(services, args)

	case "delete-job":
		commands.DeleteJob(services, args)

	case "delete-job-schedule":
		commands.DeleteJobSchedule(services, args)

	case "create-call":
		commands.CreateCall(services, args)

	case "run-call":
		commands.RunCall(services, args)

	case "schedule-call":
		commands.ScheduleCall(services, args)

	case "calls":
		commands.Calls(services, args)

	case "call-schedules":
		commands.CallSchedules(services, args)

	case "call-history":
		commands.CallHistory(services, args)

	case "delete-call":
		commands.DeleteCall(services, args)

	case "delete-call-schedule":
		commands.DeleteCallSchedule(services, args)
	}
}

func main() {
	plugin.Start(new(OCFScheduler))
}
