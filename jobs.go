package main

import (
	"encoding/json"
	"fmt"

	"github.com/ess/hype"
	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
	scheduler "github.com/starkandwayne/scheduler-for-ocf/core"
)

// FIRST GOAL!!!
// cf create-job APP_NAME NAME COMMAND
func (c *OCFScheduler) CreateJob(services *core.Services, args []string) {

	if len(args) != 4 {
		fmt.Println("cf create-job APP_NAME NAME COMMAND")
		return
	}

	appName := args[1]
	name := args[2]
	command := args[3]

	app, err := services.CLI.GetApp(appName)
	if err != nil {
		fmt.Println("Could not find app with name", appName)
		return
	}

	params := hype.Params{}
	params.Set("app_guid", app.Guid)

	payload := &scheduler.Job{
		Name:    name,
		Command: command,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Could not prepare the request payload")
		return
	}

	response := services.Client.Post("jobs", params, data)

	if !response.Okay() {
		fmt.Println(response.Error())
		return
	}

	err = json.Unmarshal(response.Data(), payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Created job %s\n\tGUID: %s\n\tApp Name: %s\n\tApp GUID: %s\n\tSpace GUID: %s\n\tCommand: %s\n",
		payload.Name,
		payload.GUID,
		appName,
		payload.AppGUID,
		payload.SpaceGUID,
		payload.Command,
	)
}

// cf run-job NAME
func (c *OCFScheduler) RunJob(services *core.Services, args []string) {
	if len(args) != 2 {
		fmt.Println("cf run-job NAME")
		return
	}

	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	name := args[1]

	job, err := core.JobNamed(services, space, name)
	if err != nil {
		fmt.Printf("Could not find job named %s in space %s.\n", name, space.Name)
		return
	}

	if core.ExecuteJob(services, job) != nil {
		fmt.Printf("Could not execute job %s.\n", job.Name)
		return
	}

	fmt.Println("OK")
}

// cf schedule-job GUID SCHEDULE
func (c *OCFScheduler) ScheduleJob(services *core.Services, args []string) {
	/* Inputs
	enabled := "false"
	expression :=
	type :=
	*/

	/* API Call
	method := "POST"
	path := fmt.Sprintf("/jobs/%s/schedules", guid)
	body := fmt.Sprintf("{\"enabled\": %s, \"expression\": \"%s\", \"expression_type\": \"%s\" }", enabled, expression, type)
	*/

	/* Responses
	 */
	fmt.Println("TODO: Implement OCFScheduler.ScheduleJob().")
}

// cf jobs
func (c *OCFScheduler) Jobs(services *core.Services, args []string) {
	space, err := core.MySpace(services)
	if err != nil {
		fmt.Println("Could not get current space.")
		return
	}

	jobs, err := core.AllJobs(services, space)
	if err != nil {
		fmt.Printf("Could not get jobs for space %s.\n", space.Name)
		return
	}

	fmt.Printf("Jobs for Space %s:\n\n", space.Name)
	for _, job := range jobs {
		fmt.Printf("\t%s (%s)\n", job.Name, job.GUID)
	}
}

// cf job-schedules SCHEDULE
func (c *OCFScheduler) JobSchedules(services *core.Services, args []string) {
	/* Inputs
	job_guid := ""
	page := "false"
	*/
	/* API
	method := "GET"
	path := fmt.Sprintf("/jobs/%s/schedules",)
	headers := "Accept: application/json"
	*/
	/* Reponses
	 */
	fmt.Println("TODO: Implement OCFScheduler.JobSchedules().")
}

// cf job-history NAME
func (c *OCFScheduler) JobHistory(services *core.Services, args []string) {
	/* Inputs
	job_guid := ""
	page := "false"
	*/
	/* API
	method := "GET"
	path := fmt.Sprintf("/jobs/%s/history",job_guid)
	*/
	/* Reponses
	 */
	fmt.Println("TODO: Implement OCFScheduler.JobHistory().")
}

// cf delete-job NAME
func (c *OCFScheduler) DeleteJob(services *core.Services, args []string) {
	/* Inputs
	job_guid := ""
	*/
	/* API
	method := "DELETE"
	path := fmt.Sprintf("/jobs/%s",job_guid)
	*/
	/* Reponses
	 */
	fmt.Println("TODO: Implement OCFScheduler.DeleteJob().")
}

// cf delete-job-schedule SCHEDULE_GUID
func (c *OCFScheduler) DeleteJobSchedule(services *core.Services, args []string) {
	/* Inputs
	job_guid := ""
	schedule_guid :=  ""
	*/
	/* API
	method := "DELETE"
	path := fmt.Sprintf("/jobs/%s/schedules/%s",job_guid, schedule_guid)
	*/
	/* Reponses
	 */
	fmt.Println("TODO: Implement OCFScheduler.DeleteJobSchedule().")
}
