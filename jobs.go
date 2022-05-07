package main

import (
	"fmt"

	"github.com/starkandwayne/ocf-scheduler-cf-plugin/core"
)

// FIRST GOAL!!!
// cf create-job APP_NAME NAME COMMAND
func (c *OCFScheduler) CreateJob(services *core.Services, args []string) {

	//name := ""
	//command := ""
	//disk := ""   // 1024MB default
	//memory := "" // 1024MB default
	/* API */
	//method := "POST"
	//headers = ["Content-Type: application/json", "Accept: application/json"]
	//path := "/jobs?app_guid=GUID "
	//body := fmt.Sprintf("{\"name\":\"%s\", \"command\":\"%s\", \"disk_in_mb\":%s, \"memory_in_mb\": %s}", name, command, disk, memory)
	/* Responses
	 */
	fmt.Println("TODO: Implement OCFScheduler.CreateJob().")
	/* hype method:

	  appGUID := args[0]
	  jobName := args[1]
	  jobCommand := args[2]

	  params := hype.Params{}
	  params.Set("app_guid", appGUID)

	  payload := &scheduler.Job{
	    Name:    jobName,
	    Command: jobCommand,
	  }

	  data, err := json.Marshal(payload)

	  response := core.Client.Post("jobs", params, data)

	  if !response.Okay() {
	    return response.Error()
	  }

	  err = json.Unmarshal(response.Data(), payload)
	  if err != nil {
	    return err
	  }

	  fmt.Printf(
	    "Created job %s\n\tGUID: %s\n\tApp GUID: %s\n\tSpace GUID: %s\n\tCommand: %s\n",
	    payload.Name,
	    payload.GUID,
	    payload.AppGUID,
	    payload.SpaceGUID,
	    payload.Command,
	  )

	  return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
	*/
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
		panic("in the far reaches of space")
	}

	fmt.Println("Space GUID:", space.SpaceFields.Guid, ", Space Name:", space.SpaceFields.Name)
	/* Inputs
	detailed := "false"
	page := "false"
	space_guid := GUID
	*/
	/* API
	method := "GET"
	path := fmt.Sprintf("/jobs?space_guid=%s", space_guid)
	*/
	/* Reponses
	 */

	fmt.Println("TODO: Implement OCFScheduler.Jobs().")
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
