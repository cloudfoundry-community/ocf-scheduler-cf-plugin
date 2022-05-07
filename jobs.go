package main

import (
	"fmt"
	"net/http"
	"time"
	//   "github.com/ess/hype"
	"code.cloudfoundry.org/cli/plugin"
)

// httpReq("GET", "/...",)
func httpReq(method string, url string, body = "", headers [][]string = [[]], timeout = 10) (http.Reponse, error) {
	api := os.GetEnv("API_URL")
	token := os.GetEnv("BEARER_TOKEN")

	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	req.Header.Set("user-agent", "ocf-scheduler")
	for header := range headers {
    h := strings.Split(header,': ')
		req.Header.Add(h[0], h[1])
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err.Error())
	}
	defer response.Body.Close()
	// TODO: Return response for processing
	return response, nil
}

// FIRST GOAL!!!
// cf create-job ....
func (c *OCFScheduler) CreateJob(args []string) {
	client := &http.Client{
	}

	name := ""
	command := ""
	disk := "" // 1024MB default
	memory := "" // 1024MB default
	/* API */
	method := "POST"
	headers = ["Content-Type: application/json", "Accept: application/json"]
	path := "/jobs?app_guid=GUID "
	body := fmt.Sprintf("{\"name\":\"%s\", \"command\":\"%s\", \"disk_in_mb\":%s, \"memory_in_mb\": %s}", name, command, disk, memory)
	*/
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

func (c *OCFScheduler) RunJob(args []string) {
	/* API Call */
	// ok,err := if HasSpace() { 
	space_guid := "" //GetSpace()???
	method := "GET"
	path := fmt.Sprintf("jobs?space_guid=%s",space_guid)
	headers := "-H 'Accept: application/json'"

	Loop over response entries and print out.
	TODO: --json output support
	/* Responses
	*/
	fmt.Println("TODO: Implement OCFScheduler.RunJob().")
}

func (c *OCFScheduler) ScheduleJob(args []string) {
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

func (c *OCFScheduler) Jobs(args []string) {
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

func (c *OCFScheduler) JobSchedules(args []string) {
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

func (c *OCFScheduler) JobHistory(args []string) {
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

func (c *OCFScheduler) DeleteJob(args []string) {
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

func (c *OCFScheduler) DeleteJobSchedule(args []string) {
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

