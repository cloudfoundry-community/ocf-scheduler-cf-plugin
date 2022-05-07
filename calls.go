package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

func (c *OCFScheduler) CreateCall(args []string) {
	/* Inputs */
	/* API
	headers := "-H 'Content-Type: application/json' -H 'Accept: application/json'"
	method := "POST"
	body := fmt.Sprintf("{\"auth_header\": \"%s\", \"name\": \"%s\", \"url\": \"%s\"}", auth_header, name, url)
	path := fmt.Sprintf("/calls",app_guid)
	*/
	/* Reponses
	*/
	fmt.Println("TODO: Implement OCFScheduler.RunCall().")
}

func (c *OCFScheduler) RunCall(args []string) {
	/* Inputs
	call_guid := ""
	*/
	/* API
	method := "POST"
	headers := "-H 'Content-Type: application/json' -H 'Accept: application/json'"
	path := fmt.Sprintf("/calls/%s/execute", call_guid)
	*/
	/* Reponses
	201 Response
	{
		"call_guid": "string",
		"execution_end_time": "string",
		"execution_start_time": "string",
		"guid": "string",
		"message": "string",
		"schedule_guid": "string",
		"scheduled_time": "string",
		"state": "string"
	}
	401 - Unauthorized
	404 - Not Found
	*/
	fmt.Println("TODO: Implement OCFScheduler.RunCall().")
}

func (c *OCFScheduler) ScheduleCall(args []string) {
	/* Inputs */
	call_guid := ""
	/* API */
	method := "POST"
	path := fmt.Sprintf("/calls/%s/schedules",call_guid)
	/* Reponses
	201 Response
	{
		"call_guid": "string",
		"created_at": "string",
		"enabled": false,
		"expression": "string",
		"expression_type": "string",
		"guid": "string",
		"updated_at": "string"
	}
	401 - Unauthorized
	404 - Not Found
	*/
	fmt.Println("TODO: Implement OCFScheduler.ScheduleCall().")
}

func (c *OCFScheduler) Calls(args []string) {
	/* Inputs */
	space_guid := ""
	page := "false"
	/* API */
	method := "GET"
	headers := "Accept: application/json"
	path := fmt.Sprintf("/calls", space_guid)

	/* Reponses
	200 Response
	{
		"pagination": {
			"first": {
				"href": "string"
			},
			"last": {
				"href": "string"
			},
			"next": {
				"href": "string"
			},
			"previous": {
				"href": "string"
			},
			"total_pages": 0,
			"total_results": 0
		},
		"resources": [
		{
			"app_guid": "string",
			"auth_header": "string",
			"created_at": "string",
			"guid": "string",
			"name": "string",
			"space_guid": "string",
			"updated_at": "string",
			"url": "string"
		}
		]
	}
	*/
	fmt.Println("TODO: Implement OCFScheduler.Calls().")
}

func (c *OCFScheduler) CallSchedules(args []string) {
	/* Inputs */
	call_guid := ""
	page := "false"
	/* API */
	method := "GET"
	path := fmt.Sprintf("/calls/%s/schedules",call_guid)

	/* Reponses
	200 Response
	{
		"pagination": {
			"first": {
				"href": "string"
			},
			"last": {
				"href": "string"
			},
			"next": {
				"href": "string"
			},
			"previous": {
				"href": "string"
			},
			"total_pages": 0,
			"total_results": 0
		},
		"resources": [
		{
			"call_guid": "string",
			"created_at": "string",
			"enabled": false,
			"expression": "string",
			"expression_type": "string",
			"guid": "string",
			"updated_at": "string"
		}
		]
	}

	401 - Unauthorized
	404 - Not Found
	*/
	fmt.Println("TODO: Implement OCFScheduler.CallSchedules().")
}

func (c *OCFScheduler) CallHistory(args []string) {
	/* Inputs */
	call_guid := ""
	/* API */
	method := "GET"
	headers := "-H 'Accept: application/json'"
	path := fmt.Sprintf("/calls/%s/history", call_guid)
	/* Reponses
	200 Response
	{
		"pagination": {
			"first": {
				"href": "string"
			},
			"last": {
				"href": "string"
			},
			"next": {
				"href": "string"
			},
			"previous": {
				"href": "string"
			},
			"total_pages": 0,
			"total_results": 0
		},
		"resources": [
		{
			"call_guid": "string",
			"execution_end_time": "string",
			"execution_start_time": "string",
			"guid": "string",
			"message": "string",
			"schedule_guid": "string",
			"scheduled_time": "string",
			"state": "string"
		}
		]
	}

	*/
	fmt.Println("TODO: Implement OCFScheduler.CallHistory().")
}

func (c *OCFScheduler) DeleteCall(args []string) {
	/* Inputs */
	call_guid := ""
	/* API */
	method := "DELETE"
	path := fmt.Sprintf("/calls/%s",call_guid)
	/* Reponses
	204 - No Content
	401 - Unauthorized
	404 - Not Found
	*/
	fmt.Println("TODO: Implement OCFScheduler.DeleteCall().")
}

func (c *OCFScheduler) DeleteCallSchedule(args []string) {
	/* Inputs */
	call_guid := ""
	schedule_guid :=""
	/* API */
	path := fmt.Sprintf("/calls/%s/schedules/%s", job_guid, schedule_guid)
	/* Reponses
	204 - No Content
	401 - Unauthorized
	404 - Not Found
	*/
	fmt.Println("TODO: Implement OCFScheduler.DeleteCallSchedule().")
}

