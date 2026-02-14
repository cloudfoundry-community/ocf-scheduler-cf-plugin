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
				Name:     "cron-expression",
				HelpText: "Display documentation on how to write ocf-scheduler's cron expression.",
				UsageDetails: plugin.Usage{
					Usage: "",
				},
			},
			{
				Name:     "create-job",
				HelpText: "Creates a job (task) related to an app.",
				UsageDetails: plugin.Usage{
					Usage: "cf create-job [OPTIONS] APP-NAME JOB-NAME COMMAND\n\nWHERE\n   APP-NAME   Name of the app whose environment is used to run the job.\n   JOB-NAME   Name of the job.\n   COMMAND    Command to execute in the app environment.\n\nOPTIONS:\n   --disk[=],   -k LIMIT   Disk limit for the job\n   --memory[=], -m LIMIT   Memory limit for the job\n   --log-rate-limit[=], --logs[=], -l LIMIT   Log rate limit for the job\n\n   Limit values are specified in units of bytes and may include unit \n   multipliers.\n\n   | unit | Multiplier | name | disk | log-rate | memory |\n   | ---- | ---------- | ---- | ---- | -------- | ------ |\n   |      |     1      |      |  X   |    X     |   X    |\n   |  k   |    1000    | Kilo |      |    X     |        |\n   |  Ki  |    1024    | Kibi |      |    X     |        |\n   |  M   |    k * k   | Mega |  X   |    X     |   X    |\n   |  Mi  |   Ki * Ki  | Mebi |  X   |    X     |   X    |\n   |  G   |    M * k   | Giga |  X   |          |   X    |\n   |  Gi  |   Mi * Ki  | Gibi |  X   |          |   X    |\n\n   • The log rate limit is specified in units of bytes per second (Bps)\n\n   • System, org, and space quotas are enforced. For example, a job\n     specifying an unlimited log rate will fail if the log rate quota\n     does not allow unlimited rates.",
				},
			},
			{
				Name:     "run-job",
				HelpText: "Runs the job (task) with the given name once.",
				UsageDetails: plugin.Usage{
					Usage: "cf run-job JOB-NAME\n\nWHERE\n   JOB-NAME is the name of the created job",
				},
			},
			{
				Name:     "schedule-job",
				HelpText: "Schedules the named job (task) to run based on the given cron schedule.",
				UsageDetails: plugin.Usage{
					Usage: "cf schedule-job [OPTIONS] JOB-NAME CRON-EXPRESSION\n\nWHERE\n   JOB-NAME is the name of the created job\n   CRON-EXPRESSION is the cron schedule format \"MINUTE HOUR DAY-OF-MONTH MONTH DAY-OF-WEEK\"\n\nOPTIONS:\n   --timezone[=]| -t string   Specifies a timezone to interpret the cron expression relative to.\n\nSEE ALSO:\n   scheduler-time-zones",
				},
			},
			{
				Name:     "jobs",
				HelpText: "Lists created jobs.",
				UsageDetails: plugin.Usage{
					Usage: "cf jobs",
				},
			},
			{
				Name:     "job-schedules",
				HelpText: "Lists created job schedules",
				UsageDetails: plugin.Usage{
					Usage: "cf job-schedules",
				},
			},
			{
				Name:     "job-history",
				HelpText: "Lists execution history for the given job name",
				UsageDetails: plugin.Usage{
					Usage: "cf job-history [OPTIONS] JOB-NAME\n\nWHERE\n   JOB-NAME is the requested job name for its historical execution data\n\nOPTIONS:\n   --show, -s (scheduled | manual | all)\n\n   The show parameter filters job history based on execution type:\n   scheduled or ad hoc(\"manual\"). The \"all\" parameter shows both\n   execution types at the same time. The parameter value is prefix-matched,\n   so you do not need to provide the full value. (default: \"scheduled\")",
				},
			},
			{
				Name:     "delete-job",
				HelpText: "Deletes named job.",
				UsageDetails: plugin.Usage{
					Usage: "cf delete-job [OPTIONS] JOB-NAME\n\nWHERE\n   JOB-NAME is the job (task) name to delete\n\nOPTIONS:\n   --force, -f   Force deletion without confirmation",
				},
			},
			{
				Name:     "delete-job-schedule",
				Alias:    "djs",
				HelpText: "Deletes the job scheduled with the named GUID.",
				UsageDetails: plugin.Usage{
					Usage: "cf delete-job-schedule [OPTIONS] JOB-NAME SCHEDULE-GUID\n\nWHERE\n   JOB-NAME is the job (task) name to delete\n   SCHEDULE-GUID is the GUID from the job-schedules command.\n\nOPTIONS:\n   --force, -f   Force deletion without confirmation",
				},
			},
			{
				Name:     "create-call",
				HelpText: "Creates a web request call",
				UsageDetails: plugin.Usage{
					Usage: "cf create-call APP-NAME CALL-NAME URL\nWHERE\n   APP-NAME is the name of the cf app to create a call for\n   CALL-NAME is a name to refer to the call as\n   URL is the URL to call.",
				},
			},
			{
				Name:     "run-call",
				HelpText: "Execute a named call request once.",
				UsageDetails: plugin.Usage{
					Usage: "cf run-call CALL-NAME\n\nWHERE\n   CALL-NAME is a name for the scheduled call",
				},
			},
			{
				Name:     "schedule-call",
				HelpText: "Schedules a call to be run based on the supplied cron schedule",
				UsageDetails: plugin.Usage{
					Usage: "cf schedule-call [OPTIONS] CALL-NAME CRON-EXPRESSION\n\nWHERE\n   CALL-NAME is a name for the scheduled call\n   CRON-EXPRESSION is a schedule using cron format \"MINUTE HOUR DAY-OF-MONTH DAY-OF-WEEK\"\n\nOPTIONS:\n   --timezone[=], -t string   Specifies a timezone to interpret the cron expression relative to.\n\nSEE ALSO:\n   scheduler-time-zones",
				},
			},
			{
				Name:     "calls",
				HelpText: "List created calls",
				UsageDetails: plugin.Usage{
					Usage: "cf calls",
				},
			},
			{
				Name:     "call-schedules",
				HelpText: "List calls scheduled to be run with app and schedule.",
				UsageDetails: plugin.Usage{
					Usage: "cf call-schedules",
				},
			},
			{
				Name:     "call-history",
				HelpText: "Shows the execution history for the named call.",
				UsageDetails: plugin.Usage{
					Usage: "cf call-history [OPTIONS] CALL-NAME\n\nWHERE\n   CALL-NAME is the requested call name for its historical execution data\n\nOPTIONS:\n   --show, -s (scheduled | manual | all)\n\n   The show parameter filters call history based on execution type:\n   scheduled or ad hoc(\"manual\"). The \"all\" parameter shows both\n   execution types at the same time. The parameter value is prefix-matched,\n   so you do not need to provide the full value. (default: \"scheduled\")",
				},
			},
			{
				Name:     "delete-call",
				HelpText: "Deletes the named call.",
				UsageDetails: plugin.Usage{
					Usage: "cf delete-call [OPTIONS] CALL-NAME\n\nWHERE\n   CALL-NAME is a name for the scheduled call\nOPTIONS:\n\n   --force, -f   Force deletion without confirmation",
				},
			},
			{
				Name:     "delete-call-schedule",
				Alias:    "dcs",
				HelpText: "Delete a call scheduled with a given GUID",
				UsageDetails: plugin.Usage{
					Usage: "cf delete-call-schedule [OPTIONS] CALL-NAME SCHEDULE-GUID\n\nWHERE\n   CALL-NAME is a name for the scheduled call.\n   SCHEDULE-GUID is the GUID from the job-schedules command.\n\nOPTIONS:\n   --force, -f   Force deletion without confirmation",
				},
			},
			{
				Name:     "scheduler-time-zones",
				Alias:    "stz",
				HelpText: "Lists scheduler time zones",
				UsageDetails: plugin.Usage{
					Usage: "cf scheduler-time-zones\n\nSEE ALSO:\n   schedule-job, schedule-call",
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
	case "cron-expression":
		break;
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
