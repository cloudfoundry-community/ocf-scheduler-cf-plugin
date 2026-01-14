package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/commands"
	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
)

var Version string = "v0.0.0"

var SemVerMajor string
var SemVerMinor string
var SemVerPatch string
var SemVerPrerelease string
var SemVerBuild string
var BuildDate string
var BuildVcsUrl string
var BuildVcsId string
var BuildVcsIdDate string
var GoArch string
var GoOs string

type OCFScheduler struct{}

// Setup to use the cf cli ui prompts
// CloudFoundry/cli/cf/i18n
//    define the T function which is suppose to translate and process go templates
//    define the teePrinter
//    define the UI interface

var teePrinter *terminal.TeePrinter
var ui terminal.UI

func (c *OCFScheduler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "OCFScheduler",
		Version: plugin.VersionType{
			Major: getVersion("Major", SemVerMajor),
			Minor: getVersion("Minor", SemVerMinor),
			Build: getVersion("Patch", SemVerPatch),
		},
		Commands: []plugin.Command{
			{
				Name:     "create-job",
				HelpText: "Creates a job (task) related to an app.",
				UsageDetails: plugin.Usage{
					Usage: "create-job:\n\tcf create-job [OPTIONS] APP-NAME JOB-NAME COMMAND\n\nWHERE\n\tAPP-NAME is the name of the cf app environment to execute with\n\tJOB-NAME is the name for this job (task)\n\tCOMMAND is the name of the command to execute within the app environment.\n\nOPTIONS\n\t--disk, -k LIMIT set job(task) disk limit (default 1024M)\n\t--memory, -m LIMIT set the job(task) memory limit (default 1024M)\n\n\tNOTE: In both of the above options, LIMIT must be specified as an\n\tinteger with an M or G at the end. This suffix is required to\n\tdifferentiate between megabytes and gigabytes (and to avoid parser\n\terrors).\n",
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
					Usage: "job-history:\n\tcf job-history [OPTIONS] JOB-NAME\n\nWHERE\n\tJOB-NAME is the requested job name for its historical execution data\n\nOPTIONS\n\t--show, -s (scheduled | manual | all)\n\n\tThe show parameter filters job history based on execution type:\n\tscheduled or ad hoc(\"manual\"). The \"all\" parameter shows both\n\texecution types at the same time. The parameter value is prefix-matched,\n\tso you do not need to provide the full value. (default: \"scheduled\")\n",
				},
			},
			{
				Name:     "delete-job",
				HelpText: "Deletes named job.",
				UsageDetails: plugin.Usage{
					Usage: "delete-job:\n\tcf delete-job [OPTIONS] JOB-NAME\n\nWHERE\n\tJOB-NAME is the job (task) name to delete\n\nOPTIONS\n\t--force, -f   Force deletion without confirmation",
				},
			},
			{
				Name:     "delete-job-schedule",
				HelpText: "Deletes the job scheduled with the named GUID.",
				UsageDetails: plugin.Usage{
					Usage: "delete-job-schedule:\n\tcf delete-call-schedule [OPTIONS] JOB-NAME SCHEDULE-GUID\n\nOPTIONS\n\t--force, -f   Force deletion without confirmation",
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
					Usage: "call-history:\n\tcf call-history [OPTIONS] CALL-NAME\n\nWHERE\n\tCALL-NAME is the requested call name for its historical execution data\n\nOPTIONS\n\t--show, -s (scheduled | manual | all)\n\n\tThe show parameter filters call history based on execution type:\n\tscheduled or ad hoc(\"manual\"). The \"all\" parameter shows both\n\texecution types at the same time. The parameter value is prefix-matched,\n\tso you do not need to provide the full value. (default: \"scheduled\")\n",
				},
			},
			{
				Name:     "delete-call",
				HelpText: "Deletes the named call.",
				UsageDetails: plugin.Usage{
					Usage: "delete-call:\n\tcf delete-call [OPTIONS] CALL-NAME\n\nOPTIONS\n\t--force, -f   Force deletion without confirmation",
				},
			},
			{
				Name:     "delete-call-schedule",
				HelpText: "Delete a call scheduled with a given GUID",
				UsageDetails: plugin.Usage{
					Usage: "delete-call-schedule:\n\tcf delete-call-schedule [OPTIONS] CALL-NAME SCHEDULE-GUID\n\nOPTIONS\n\t--force, -f   Force deletion without confirmation",
				},
			},
			{
				Name:     "scheduler-time-zones",
				HelpText: "Lists scheduler time zones",
				UsageDetails: plugin.Usage{
					Usage: "jobs:\ncf scheduler-time-zones",
				},
			},
		},
	}
}

func (c *OCFScheduler) Run(cliConnection plugin.CliConnection, args []string) {

	i18n.T = func(translationID string, args ...interface{}) string {
		var buffer bytes.Buffer

		var keys interface{}
		if len(args) > 0 {
			keys = args[0]
		}

		formattedTemplate := template.Must(template.New("Display Text").Parse(translationID))
		err := formattedTemplate.Execute(&buffer, keys)
		if err != nil {
			return translationID + "\ntemplate processing failed " + err.Error() + "\n"
		}

		return buffer.String()
	}

	if runtime.GOOS == "windows" {
		terminal.UserAskedForColors = "false"
	}
	terminal.InitColorSupport()
	teePrinter = terminal.NewTeePrinter(os.Stdout)
	// Should we really be using NewPluginUI instead TODO
	ui = terminal.NewUI(os.Stdin, os.Stdout, teePrinter, trace.NewWriterPrinter(io.Discard, false))

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

	services := &core.Services{CLI: cliConnection, Client: client, UI: ui}

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

	case "scheduler-time-zones":
		commands.SchedulerTimeZones(services, args)
	}
}

func createBuildMeta(buildOs, buildArch, build string) string {
	p1 := strings.TrimSpace(buildOs)
	p2 := strings.TrimSpace(buildArch)
	p3 := strings.TrimSpace(build)
	if p1 == "" || p2 == "" {
		panic(fmt.Sprintf("Go meta data is missing one of its parts: %s, %s ", p1, p2))
	}
	b := strings.Join([]string{p1, p2}, ".")
	if p3 != "" {
		b += "." + p3
	}
	return b
}

func createSemVer(major, minor, patch, prerelease, build string) string {
	p1 := strings.TrimSpace(major)
	p2 := strings.TrimSpace(minor)
	p3 := strings.TrimSpace(patch)
	p4 := strings.TrimSpace(prerelease)
	p5 := strings.TrimSpace(build)
	if p1 == "" || p2 == "" || p3 == "" {
		panic(fmt.Sprintf("Semanic version is missing one of its parts: %s.%s.%s", p1, p2, p3))
	}

	sv := strings.Join([]string{p1, p2, p3}, ".")
	if p4 != "" {
		sv += "-" + p4
	}
	if p5 != "" {
		sv += "+" + p5
	}
	return sv
}

func getVersion(version, toInt string) int {
	theInt, err := strconv.Atoi(toInt)
	if err != nil {
		theInt = 0
		fmt.Printf("Warning: %v for %v version value.  Defaulting to a zero value\n", err.Error(), version)
	}
	return theInt
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		bm := createBuildMeta(GoOs, GoArch, SemVerBuild)
		sv := createSemVer(SemVerMajor, SemVerMinor, SemVerPatch, SemVerPrerelease, bm)
		f := "%13v %v\n"
		fmt.Printf(f, "Version:", sv)
		fmt.Printf(f, "Build Date:", BuildDate)
		fmt.Printf(f, "VCS Url:", BuildVcsUrl)
		fmt.Printf(f, "VCS Id:", BuildVcsId)
		fmt.Printf(f, "VCS Id Date:", BuildVcsIdDate)
	}

	plugin.Start(new(OCFScheduler))
}
