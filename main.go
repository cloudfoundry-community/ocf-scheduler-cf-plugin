package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

type OCFScheduler struct{}

func (c *OCFScheduler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "OCFScheduler",
		Commands: []plugin.Command {
			{
				Name: "create-job",
				HelpText: "Creates a job (task) related to an app.",
				UsageDetails: plugin.Usage {
					Usage: "create-job:\n\tcf create-job APP_NAME NAME COMMAND\n\nWHERE\n\tAPP_NAME is the name of the cf app environment to execute with\n\tNAME is the name for this job (task)\n\tCOMMAND is the name of the command to execute within the app environment.\n\nTODO: --disk LIMIT/-k LIMIT set job(task) disk limit\n\tTODO: --memory LIMIT//-m LIMIT set the job(task) memory limit\n\tLIMIT is an integer suffixed with MB or GB.",
				},
			},
			{
				Name: "run-job",
				HelpText: "Runs the job (task) with the given name once.",
				UsageDetails: plugin.Usage {
					Usage: "run-job:\n\tcf run-job NAME",
				},
			},
			{
				Name: "schedule-job",
				HelpText: "Schedules the named job (task) to run based on the given cron schedule.",
				UsageDetails: plugin.Usage {
					Usage: "schedule-job:\n\tcf schedule-job GUID SCHEDULE\n\nWHERE\n\tGUID is the guid of the created job\n\tSCHEDULE is the cron schedule format \"MIN HOUR DAY-OF-MONTH DAY-OF-WEEK\"",
				},
			},
			{
				Name: "jobs",
				HelpText: "Lists created jobs.",
				UsageDetails: plugin.Usage {
					Usage: "jobs:\ncf jobs",
				},
			},
			{
				Name: "job-schedules",
				HelpText: "Lists created job schedules",
				UsageDetails: plugin.Usage {
					Usage: "job-schedules:\n\tcf job-schedules SCHEDULE",
				},
			},
			{
				Name: "job-history",
				HelpText: "Lists execution history for the given job name",
				UsageDetails: plugin.Usage {
					Usage: "job-history:\n\tcf job-history NAME",
				},
			},
			{
				Name: "delete-job",
				HelpText: "Deletes named job.",
				UsageDetails: plugin.Usage {
					Usage: "delete-job:\n\tcf delete-job NAME",
				},
			},
			{
				Name: "delete-job-schedule",
				HelpText: "Deletes the job scheduled with the named GUID.",
				UsageDetails: plugin.Usage {
					Usage: "delete-job-schedule:\n\tcf delete-job-schedule SCHEDULE_GUID",
				},
			},
			{
				Name: "create-call",
				HelpText: "Creates a web request call",
				UsageDetails: plugin.Usage {
					Usage: "create-call:\n\tcreate-call APP_NAME NAME URL \nWHERE\n\tAPP_NAME is the name of the cf app to create a call for\n\tNAME is a name to refer to the call as\n\tURL is the URL to call.",
				},
			},
			{
				Name: "run-call ",
				HelpText: "Execute a named call request once.",
				UsageDetails: plugin.Usage {
					Usage: "run-call:\n\trun-call NAME",
				},
			},
			{
				Name: "schedule-call",
				HelpText: "Schedules a call to be run based on the supplied cron schedule",
				UsageDetails: plugin.Usage {
					Usage: "schedule-call:\n\tschedule-call NAME SCHEDULE\n\tNAME is a name for the scheduled call\n\tSCHEUDLE is a schedule using cron schedule format \"MIN HOUR DAY-OF-MONTH DAY-OF-WEEK\"\n\nEXAMPLE\n\tcf schedule-call hourlyrun \"0 * * * *\"",
				},
			},
			{
				Name: "calls",
				HelpText: "List created calls",
				UsageDetails: plugin.Usage {
					Usage: "calls:\n\tcalls",
				},
			},
			{
				Name: "call-schedules",
				HelpText: "List calls scheduled to be run with app and schedule.",
				UsageDetails: plugin.Usage {
					Usage: "call-schedules:\n\tcf call-schedules",
				},
			},
			{
				Name: "call-history",
				HelpText: "Shows the execution history for the named call.",
				UsageDetails: plugin.Usage {
					Usage: "call-history:\n\tcf call-history NAME",
				},
			},
			{
				Name: "delete-call",
				HelpText: "Deletes the named call.",
				UsageDetails: plugin.Usage {
					Usage: "delete-call:\n\tcf delete-call NAME",
				},
			},
			{
				Name: "delete-call-schedule",
				HelpText: "Delete a call scheduled with a given GUID",
				UsageDetails: plugin.Usage {
					Usage: "delete-call-schedule:\n\tcf delete-call-schedule GUID",
				},
			},
		},
	}
}

func (c *OCFScheduler) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "create-job" {
		c.CreateJob()
	} else if args[0] == "run-job" {
		c.RunJob()
	} else if args[0] == "schedule-job" {
		c.ScheduleJob()
	} else if args[0] == "jobs" {
		c.Jobs()
	} else if args[0] == "job-schedules" {
		c.JobSchedules()
	} else if args[0] == "job-history" {
		c.JobHistory()
	} else if args[0] == "delete-job" {
		c.DeleteJob()
	} else if args[0] == "delete-job-schedule" {
		c.DeleteJobSchedule()
	} else if args[0] == "create-call" {
		c.CreateCall()
	} else if args[0] == "run-call" {
		c.RunCall()
	} else if args[0] == "schedule-call" {
		c.ScheduleCall()
	} else if args[0] == "calls" {
		c.Calls()
	} else if args[0] == "call-schedules" {
		c.CallSchedules()
	} else if args[0] == "call-history" {
		c.CallHistory()
	} else if args[0] == "delete-call" {
		c.DeleteCall()
	} else if args[0] == "delete-call-schedule" {
		c.DeleteCallSchedule()
	}
}

func main() {
	plugin.Start(new(OCFScheduler))
}

