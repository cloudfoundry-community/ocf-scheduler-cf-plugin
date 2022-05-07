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
	/* API Call */
	// ok,err := if HasSpace() {
	//space_guid := "" //GetSpace()???
	//method := "GET"
	//path := fmt.Sprintf("jobs?space_guid=%s", space_guid)
	//headers := "-H 'Accept: application/json'"

	//Loop over response entries and print out.
	//TODO: --json output support
	/* Responses
	 */
	fmt.Println("TODO: Implement OCFScheduler.RunJob().")
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
	space, err := services.CLI.GetCurrentSpace()
	if err != nil {
		fmt.Println("Could not get current Space.")
		return
	}

	fmt.Println("Space GUID:", space.SpaceFields.Guid, ", Space Name:", space.SpaceFields.Name)

	params := hype.Params{}
	params.Set("space_guid", space.SpaceFields.Guid)

	response := services.Client.Get("jobs", params)

	if !response.Okay() {
		fmt.Println("Got a bad response from the Scheduler API.")
		return
	}

	data := struct {
		Resources []*scheduler.Job `json:"resources"`
	}{}

	err := json.Unmarshal(response.Data(), &data)
	if err != nil {
		fmt.Println("Could not decode Scheduler API response.")
		return
	}

	for _, job := range data.Resources {
		fmt.Printf("%s (%s)\n", job.Name, job.GUID)
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
